package dynamodbnotifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mariotoffia/godeviceshadow/model"
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

// Processor is the one that listens to DynamoDB stream events and
// does the report/desire merge and then feed it to a notification
// handler that may send in memory, SQS etc. based on plugins.
type Processor struct {
	tr model.TypeRegistry
	// eventModel is the event model that the processor is running in. Default is `ProcessorEventModelLambda`.
	eventModel ProcessorEventModel
	// stream is used if the event model is `ProcessorEventModelManualAttach`.
	stream *DynamoDBStream
}

// StartLambda will start the lambda loop and freeze.
func (p *Processor) Start(ctx context.Context) {
	if p.eventModel == ProcessorEventModelManualAttach {
		p.stream.Start(ctx, p.HandleRequest)
	} else {
		lambda.StartWithOptions(p.HandleRequest, lambda.WithContext(ctx))
	}
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

		// TODO:
		_ = oldImage
		_ = newImage
	}

	return nil
}
