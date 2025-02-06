package stream

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"

	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

// StreamPollerCallback is used in the `Start` function to handle DynamoDB stream events.
type StreamPollerCallback func(ctx context.Context, event events.DynamoDBEvent) error

// StreamPollerDoneCallback is called when the poll loop in `Start` is done.
type StreamPollerDoneCallback func(ctx context.Context, err error)

// PollErrorCallback is invoked when an error did occur during polling. If the callback
// returns an error, it will exit the poll loop with that error.
type PollErrorCallback func(ctx context.Context, err error) error

type DynamoDBStream struct {
	// client is the DynamoDB client to use.
	client *dynamodb.Client
	// streamsClient is the DynamoDB Streams client to use.
	streamsClient *dynamodbstreams.Client
	// table is the DynamoDB table name.
	table string
	// wasStreamEnabled is true if the stream was enabled by this instance.
	wasStreamEnabled bool
	// restoreState will make sure to restore the stream state when done. If
	// set to false, it will not.
	restoreState bool
	// maxWaitTime is the maximum time to wait for the stream to be enabled.
	maxWaitTime time.Duration
	// shardIterator is stored between polls.
	shardIterator string
	// iteratorType is the iterator type to start polling for records.
	iteratorType streamTypes.ShardIteratorType
	// callback is the callback function to handle DynamoDB stream events when calling `Start`.
	callback StreamPollerCallback
	// startDone is called when the `Start` function has finished processing.
	startDone StreamPollerDoneCallback
	// shards keeps track on shards.
	shards *Shards
	// logPollErrors will log poll errors if set to `true`.
	logPollErrors bool
	// pollErrorCallback is invoked when an error did occur during polling.
	pollErrorCallback PollErrorCallback
}

type DynamoDBStreamOptions struct {
	// Client is the DynamoDB client to use. If nil, a new client will be created.
	Client *dynamodb.Client
	// StreamsClient is the DynamoDB Streams client to use. If nil, a new client will be created.
	StreamsClient *dynamodbstreams.Client
	// If _Client_ is nil, the AWS region to use for the DynamoDB client. If not set,
	// it will use the profile default.
	Region string
	// RestoreState will make sure to restore the stream state when done. If set to false,
	// it will not. Default is `false`.
	RestoreState bool
	// MaxWaitTime is the maximum time to wait for the stream to be enabled. Default is 2 minutes.
	MaxWaitTime time.Duration
	// IteratorType is which iterator type to start polling for records. Default is `LATEST`.
	IteratorType streamTypes.ShardIteratorType
	// Callback is the callback function to handle DynamoDB stream events when calling `Start`.
	Callback StreamPollerCallback
	// StartDone is called when the `Start` function has finished processing.
	StartDone StreamPollerDoneCallback
	// PollErrorCallback is invoked when an error did occur during polling.
	PollErrorCallback PollErrorCallback
	// LogPollErrors will log poll errors if set to `true`. Default is `false`.
	LogPollErrors bool
}

func NewDynamoDBStream(table string, opts ...DynamoDBStreamOptions) (*DynamoDBStream, error) {
	if table == "" {
		return nil, fmt.Errorf("table name cannot be empty")
	}

	var opt DynamoDBStreamOptions

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.MaxWaitTime == 0 {
		opt.MaxWaitTime = 2 * time.Minute
	}

	if opt.IteratorType == "" {
		opt.IteratorType = streamTypes.ShardIteratorTypeLatest
	}

	var (
		cfg aws.Config
		err error
	)

	if opt.Client == nil || opt.StreamsClient == nil {
		if opt.Region == "" {
			cfg, err = config.LoadDefaultConfig(context.Background())
		} else {
			cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(opt.Region))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
		}
	}

	if opt.Client == nil {
		opt.Client = dynamodb.NewFromConfig(cfg)
	}

	if opt.StreamsClient == nil {
		opt.StreamsClient = dynamodbstreams.NewFromConfig(cfg)
	}

	ds := &DynamoDBStream{
		client:            opt.Client,
		streamsClient:     opt.StreamsClient,
		table:             table,
		iteratorType:      opt.IteratorType,
		restoreState:      opt.RestoreState,
		maxWaitTime:       opt.MaxWaitTime,
		callback:          opt.Callback,
		startDone:         opt.StartDone,
		pollErrorCallback: opt.PollErrorCallback,
		logPollErrors:     opt.LogPollErrors,
	}

	ctx := context.Background()

	arn, err := ds.getStreamArn(ctx)
	if err != nil {
		return nil, err
	}

	ds.shards, err = NewShards(ctx, opt.StreamsClient, arn, opt.IteratorType)

	if err != nil {
		return nil, err
	}

	return ds, nil
}

// Close will release the stream if _restoreState_ is set to `true` and the stream was enabled by this instance.
//
// This instance shall *not* be used after this has been executed.
func (s *DynamoDBStream) Close(ctx context.Context) error {
	s.shardIterator = ""

	if s.restoreState && s.wasStreamEnabled {
		if err := s.ReleaseDbStream(ctx, true); err != nil {
			return err
		}

		s.restoreState = false
	}

	return nil
}

// Start will start manual polling of the DynamoDB stream. It will enable the stream if not already enabled.
//
// This function will block until the context is cancelled or an error occurs. It is automatically closed upon
// cancellation or error.
//
// When _async_ is set to `true`, it will ensure the table is enabled and then fork out the polling in a go routine.
// If not, it will block until the context is cancelled or an error occurs.
//
// Even if _async_ is `true`, it will invoke `Close` when the context is cancelled or an error occurs.
func (s *DynamoDBStream) Start(ctx context.Context, async bool) error {
	cb := s.callback

	if cb == nil {
		return fmt.Errorf("no callback function set")
	}

	enabled, err := s.EnableStream(ctx)

	if err != nil {
		return fmt.Errorf("failed to enable stream: %w", err)
	}

	if !enabled {
		return fmt.Errorf("stream is not enabled")
	}

	poll := func() error {
		var (
			streamID string
			err      error
			evt      events.DynamoDBEvent
		)

		defer func() {
			s.Close(ctx)

			if s.startDone != nil {
				s.startDone(ctx, err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				// When the context is cancelled, exit gracefully.
				err = ctx.Err()
				return err
			default:
				// Poll the stream for new records.
				streamID, evt, err = s.PollAsDynamoDBEvent(ctx)

				if err != nil {
					err = fmt.Errorf("error polling stream: %w", err)
				} else if len(evt.Records) > 0 {
					if err = cb(ctx, evt); err != nil {
						err = fmt.Errorf("error handling event on shard %s: %w", streamID, err)
					} else {
						s.Commit(streamID)
					}
				}

				if err != nil {
					if s.logPollErrors {
						fmt.Println(err)
					}

					if s.pollErrorCallback != nil {
						if err = s.pollErrorCallback(ctx, err); err != nil {
							return err
						}
					}
				}

				time.Sleep(5 * time.Second)
			}
		}
	}

	if async {
		go poll()
	} else {
		return poll()
	}

	return nil
}

// IsStreamEnabled checks if the stream is enabled on the table. If it is enabled,
// it returns `true`. If there is an error checking the stream status, it returns an error.
func (s *DynamoDBStream) IsStreamEnabled(ctx context.Context) (bool, error) {
	descOutput, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &s.table,
	})

	if err != nil {
		return false, fmt.Errorf("failed to describe table: %w", err)
	}

	return isStreamEnabled(descOutput), nil
}

// EnableStream will enable the stream on the table if not already enabled and store the
// state it was before. This state can be used to restore the stream state when done when
// calling`Close`.
//
// It will return `true` if the stream is enabled and the stream is published by this invocation.
func (s *DynamoDBStream) EnableStream(ctx context.Context) (bool, error) {
	//
	enabled, err := s.IsStreamEnabled(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check if stream is enabled: %w", err)
	}

	if enabled {
		s.wasStreamEnabled = true

		return true, nil // already enabled
	}

	_, err = s.client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: &s.table,
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	})

	if err != nil {
		return false, fmt.Errorf("failed to update table to enable stream: %w", err)
	}

	// Wait until the stream ARN becomes available.
	now := time.Now()

	for {
		if time.Since(now) > s.maxWaitTime {
			return false, fmt.Errorf("timeout waiting for stream to become available")
		}

		time.Sleep(5 * time.Second)

		descOutput, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: &s.table,
		})

		if err != nil {
			return false, fmt.Errorf("failed to describe table after update: %w", err)
		}

		if descOutput.Table.LatestStreamArn != nil && *descOutput.Table.LatestStreamArn != "" {
			return true, nil
		}
	}
}

// ReleaseDbStream will release the DynamoDB stream independent if it was enabled by this instance
// or not. If _wait_ is set to `true`, it will wait until the stream is fully disabled (max wait time).
func (s *DynamoDBStream) ReleaseDbStream(ctx context.Context, wait bool) error {
	_, err := s.client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: &s.table,
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled: aws.Bool(false),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to disable stream on table %s: %w", s.table, err)
	}

	if !wait {
		return nil
	}

	now := time.Now()

	for {
		if time.Since(now) > s.maxWaitTime {
			return fmt.Errorf("timeout waiting for stream to be disabled")
		}

		time.Sleep(5 * time.Second)

		descOutput, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: &s.table,
		})

		if err != nil {
			return fmt.Errorf("failed to describe table after disabling stream: %w", err)
		}

		if !isStreamEnabled(descOutput) {
			return nil
		}
	}
}

// Poll polls the DynamoDB stream for new records.
// It returns a slice of records and updates the shard iterator for subsequent polls.
//
// Even if this function do not return an `error`, the records may be empty.
// Poll delegates polling to the DynamoDbShards instance.
//
// The first parameter is the shardID that shall be used in the `Commit` function when
// all records where successfully processed.
func (s *DynamoDBStream) Poll(ctx context.Context) (string, []streamTypes.Record, error) {
	return s.shards.Poll(ctx, s.streamsClient)
}

// Commit delegates the commit operation to DynamoDbShards for a specific shard.
func (s *DynamoDBStream) Commit(shardID string) {
	s.shards.Commit(shardID)
}

// getStreamArn retrieves the LatestStreamArn from the table description.
func (s *DynamoDBStream) getStreamArn(ctx context.Context) (string, error) {
	descOutput, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &s.table,
	})

	if err != nil {
		return "", fmt.Errorf("failed to describe table: %w", err)
	}

	if descOutput.Table.LatestStreamArn == nil || *descOutput.Table.LatestStreamArn == "" {
		return "", fmt.Errorf("stream not enabled on table %s", s.table)
	}

	return *descOutput.Table.LatestStreamArn, nil
}

func isStreamEnabled(descOutput *dynamodb.DescribeTableOutput) bool {
	return descOutput.Table.StreamSpecification != nil &&
		descOutput.Table.StreamSpecification.StreamEnabled != nil &&
		*descOutput.Table.StreamSpecification.StreamEnabled
}
