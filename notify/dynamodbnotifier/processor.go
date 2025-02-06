package dynamodbnotifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
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
type ProcessImagesFunc func(ctx context.Context, oldImage, newImage any, err error) error

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

		if p.cbProcessImage != nil {
			if err := p.cbProcessImage(ctx, oldImage, newImage, err); err != nil {
				return err
			}
		}

		// TODO: Diff old, new (reported or desired) and notify.
	}

	return nil
}
