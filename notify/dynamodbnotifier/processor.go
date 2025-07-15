package dynamodbnotifier

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/loggers"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
	"github.com/mariotoffia/godeviceshadow/utils/loggerutils"
)

type ProcessorEventModel int

const (
	// ProcessorEventModelLambda is the lambda event model where it is hosted as a _AWS_
	// lambda function that get's invoked for each _DynamoDB_ stream event.
	ProcessorEventModelLambda ProcessorEventModel = iota
	// ProcessorEventModelManualAttach is when the processor is manually attached to a
	// _DynamoDB_ stream and polls for events. It may even, enable the stream on a
	// table if needed.
	ProcessorEventModelManualAttach
)

// ProcessImagesFunc is called just after the record has been processed and the images (possibly)
// extracted. The `oldImage` and `newImage` is the images extracted from the record. If the record
// processing failed, the `err` is set.
type ProcessImagesFunc func(ctx context.Context, oldImage, newImage *PersistenceObject, err error) error

// ProcessorDoneFunc is called when the `Start` function has finished processing.
//
// NOTE: When lambda and non async, this will *not* be called.
type ProcessorDoneFunc func(ctx context.Context, err error)

// Processor is the one that listens to DynamoDB stream events and
// does the report/desire merge and then feed it to a notification
// handler that may send in memory, SQS etc. based on plugins.
type Processor struct {
	tr model.TypeRegistry
	// eventModel is the event model that the processor is running in. Default is `ProcessorEventModelLambda`.
	eventModel ProcessorEventModel
	// stream is used if the event model is `ProcessorEventModelManualAttach`.
	stream *stream.DynamoDBStream
	// cbProcessImage is a optional callback function that may cancel any further processing but nothing else.
	cbProcessImage ProcessImagesFunc
	// startDoneFunc is called when the `Start` function has finished processing.
	startDoneFunc ProcessorDoneFunc
	// cbNotifyProcess is invoked before the notification is sent. It will ignore the return values.
	cbNotifyProcess notifiermodel.ProcessFunc
	// reportedLoggers are `model.MergeLogger` that are used to log reported values.
	//
	// It will automatically add the `changelogger.ChangeMergeLogger`. It will capture the reported model merges.
	reportedLoggers []model.CreatableMergeLogger
	// desiredLoggers are `model.MergeLogger` that are used to log desired values.
	//
	// It will automatically add the `changelogger.ChangeMergeLogger` and `loggers.DynamoDbDesireLogger` to detect
	// when a desire has been acknowledged. The former will log a plain desire document merge.
	desiredLoggers []model.CreatableMergeLogger
	// notificationManager will be invoked for each event that has been processed.
	notificationManager notifiermodel.Notifier
}

// StartLambda will start the lambda loop and freeze.
//
// If _async_ is set to `true` it will start the processing in a go routine.
//
// NOTE: Do only call this function *once* - Create a new processor instead.
func (p *Processor) Start(ctx context.Context, async bool) error {
	if p.eventModel == ProcessorEventModelManualAttach {
		return p.stream.Start(ctx, async)
	} else {
		// Lambda
		if async {
			go func() {
				lambda.StartWithOptions(p.HandleRequest, lambda.WithContext(ctx))

				if p.startDoneFunc != nil {
					p.startDoneFunc(ctx, nil)
				}
			}()
		} else {
			lambda.StartWithOptions(p.HandleRequest, lambda.WithContext(ctx))

			if p.startDoneFunc != nil {
				p.startDoneFunc(ctx, nil)
			}
		}
	}

	return fmt.Errorf("lambda returned (it should not)")
}

// HandleRequest processes DynamoDB stream events and does the report/desire
// merge (again) then feed it to a notification handler that may send
// in memory, SQS etc. based on plugins.
//
// This ensures that the change event in the model persisted by `dynamodbpersistence`
// is at least notified once (this is the DynamoDB stream guarantee).
//
// Since it is just doing the merge part with loggers, it is much more lightweight
// than the upsert operation in the `dynamodbpersistence` package.
func (p *Processor) HandleRequest(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		oldImage, newImage, err := processRecord(record, p.tr)

		if err != nil {
			fmt.Printf("Error processing record: %v\n", err)

			continue
		}

		var id persistencemodel.ID

		if oldImage != nil {
			oldImage.Meta["record"] = &record
			id = oldImage.ID()
		}

		if newImage != nil {
			newImage.Meta["record"] = &record

			if id.ID == "" {
				id = newImage.ID()
			}
		}

		if p.cbProcessImage != nil {
			if err := p.cbProcessImage(ctx, oldImage, newImage, err); err != nil {
				return err
			}
		}

		// Merge the reported and desired values
		reportMergeOpts := merge.MergeOptions{
			Mode:    merge.ClientIsMaster,
			Loggers: loggerutils.CreateMergeLoggers(p.reportedLoggers...),
		}

		desiredMergeOpts := merge.MergeOptions{
			Mode:    merge.ClientIsMaster,
			Loggers: loggerutils.CreateMergeLoggers(p.desiredLoggers...),
		}

		var desired, reported bool

		switch DynamoDbEventType(record.EventName) {
		case DynamoDbEventTypeInsert:
			if newImage.Reported != nil {
				merge.MergeAny(ctx, createInstance(newImage.Reported), newImage.Reported, reportMergeOpts)
				reported = true
			}

			if newImage.Desired != nil {
				merge.MergeAny(ctx, createInstance(newImage.Desired), newImage.Desired, desiredMergeOpts)
				desired = true
			}
		case DynamoDbEventTypeModify:
			if newImage.Reported != nil && oldImage.Reported != nil {
				merge.MergeAny(ctx, oldImage.Reported, newImage.Reported, reportMergeOpts)
				reported = true
			}

			if newImage.Desired != nil && oldImage.Desired != nil {
				merge.MergeAny(ctx, oldImage.Desired, newImage.Desired, desiredMergeOpts)
				desired = true
			}
		case DynamoDbEventTypeRemove:
			if oldImage.Reported != nil {
				merge.MergeAny(ctx, oldImage.Reported, createInstance(oldImage.Reported), reportMergeOpts)
				reported = true
			}

			if oldImage.Desired != nil {
				merge.MergeAny(ctx, oldImage.Desired, createInstance(oldImage.Desired), desiredMergeOpts)
				desired = true
			}
		}

		oper := make([]notifiermodel.NotifierOperation, 0, 2)

		if reported {
			oper = append(oper, notifiermodel.NotifierOperation{
				ID:           id.ToPersistenceID(persistencemodel.ModelTypeReported),
				MergeLogger:  *loggerutils.FindMerge[*changelogger.ChangeMergeLogger](reportMergeOpts.Loggers),
				DesireLogger: *loggerutils.FindMerge[*loggers.DynamoDbDesireLogger](desiredMergeOpts.Loggers).DesireLogger,
				Operation:    notifiermodel.OperationTypeReport,
				Reported:     newImage.Reported,
				Desired:      newImage.Desired,
			})
		}

		if desired {
			oper = append(oper, notifiermodel.NotifierOperation{
				ID:          id.ToPersistenceID(persistencemodel.ModelTypeDesired),
				MergeLogger: loggerutils.FindMerge[changelogger.ChangeMergeLogger](desiredMergeOpts.Loggers),
				Operation:   notifiermodel.OperationTypeDesired,
				Reported:    newImage.Reported,
				Desired:     newImage.Desired,
			})
		}

		if p.cbNotifyProcess != nil {
			p.cbNotifyProcess(ctx, nil /*tx*/, oper...)
		}

		if p.notificationManager != nil {
			res := p.notificationManager.Process(ctx, nil /*tx*/, oper...)
			_ = res // TODO: If we get an error, we need to handle it (and return an error)
		}
	}

	return nil
}

// createInstance creates a copy of the instance `v`.
func createInstance(v any) any {
	if v == nil {
		return nil
	}

	return reflect.New(reflect.TypeOf(v)).Interface()
}
