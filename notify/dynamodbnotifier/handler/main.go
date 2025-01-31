package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mariotoffia/godeviceshadow/model"
	gtr "github.com/mariotoffia/godeviceshadow/types"
)

func main() {
	processor := &Processor{}

	// TODO: Initialization callback?

	if processor.tr == nil {
		processor.tr = gtr.NewRegistry()
	}

	processor.StartLambda()
}

type Processor struct {
	tr model.TypeRegistry
}

// StartLambda will start the lambda loop and freeze.
func (p *Processor) StartLambda() {
	lambda.Start(p.HandleRequest)
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
