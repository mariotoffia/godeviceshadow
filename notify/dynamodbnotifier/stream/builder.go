package stream

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"

	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

type StreamBuilder struct {
	opts  DynamoDBStreamOptions
	table string
}

func NewStreamBuilder() *StreamBuilder {
	return &StreamBuilder{}
}

func (b *StreamBuilder) WithTable(table string) *StreamBuilder {
	b.table = table
	return b
}

func (b *StreamBuilder) WithCallback(cb StreamPollerCallback) *StreamBuilder {
	b.opts.Callback = cb
	return b
}

func (b *StreamBuilder) UseClient(client *dynamodb.Client) *StreamBuilder {
	b.opts.Client = client
	return b
}

func (b *StreamBuilder) UseStreamClient(client *dynamodbstreams.Client) *StreamBuilder {
	b.opts.StreamsClient = client
	return b
}

func (b *StreamBuilder) WithRegion(region string) *StreamBuilder {
	b.opts.Region = region
	return b
}

func (b *StreamBuilder) WithShardIteratorType(t streamTypes.ShardIteratorType) *StreamBuilder {
	b.opts.IteratorType = t
	return b
}

func (b *StreamBuilder) WithMaxWaitTime(waitTime time.Duration) *StreamBuilder {
	b.opts.MaxWaitTime = waitTime
	return b
}

// WithStartDone is called when the `Start` function has finished processing.
//
// TIP: This function can be called many times to have multiple callbacks.
func (b *StreamBuilder) WithStartDone(cb StreamPollerDoneCallback) *StreamBuilder {
	if b.opts.StartDone != nil {
		old := b.opts.StartDone

		b.opts.StartDone = func(ctx context.Context, err error) {
			old(ctx, err)
			cb(ctx, err)
		}
	} else {
		b.opts.StartDone = cb
	}

	return b
}

// WithPollErrorCallback will be called every time an error did occur during polling in the `Start` function.
//
// When the function returns an error, the `Start` function will stop and return that error.
func (b *StreamBuilder) WithPollErrorCallback(cb PollErrorCallback) *StreamBuilder {
	b.opts.PollErrorCallback = cb
	return b
}

func (b *StreamBuilder) LogPollErrors() *StreamBuilder {
	b.opts.LogPollErrors = true
	return b
}

func (b *StreamBuilder) Build(ctx context.Context) (*DynamoDBStream, error) {
	return NewDynamoDBStream(ctx, b.table, b.opts)
}
