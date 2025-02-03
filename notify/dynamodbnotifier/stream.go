package dynamodbnotifier

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
	// IteratorType is which iterator type to start polling for records. Default is `TRIM_HORIZON`.
	IteratorType streamTypes.ShardIteratorType
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
		opt.IteratorType = streamTypes.ShardIteratorTypeTrimHorizon
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

	return &DynamoDBStream{
		client:        opt.Client,
		streamsClient: opt.StreamsClient,
		table:         table,
		iteratorType:  opt.IteratorType,
		restoreState:  opt.RestoreState,
		maxWaitTime:   opt.MaxWaitTime,
	}, nil
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
func (s *DynamoDBStream) Start(ctx context.Context, cb func(ctx context.Context, event events.DynamoDBEvent) error) error {
	enabled, err := s.EnableStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to enable stream: %w", err)
	}
	// Ensure the stream state is restored upon exit.
	defer s.Close(ctx)

	if !enabled {
		return fmt.Errorf("stream is not enabled")
	}

	for {
		select {
		case <-ctx.Done():
			// When the context is cancelled, exit gracefully.
			return ctx.Err()
		default:
			// Poll the stream for new records.
			evt, err := s.PollAsDynamoDBEvent(ctx)

			if err != nil {
				return fmt.Errorf("error polling stream: %w", err)
			}

			if len(evt.Records) > 0 {
				if err := cb(ctx, evt); err != nil {
					_ = fmt.Errorf("error handling event: %w", err)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}
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
// It will return `true` if the stream was enabled and the stream is published by this invocation.
func (s *DynamoDBStream) EnableStream(ctx context.Context) (bool, error) {
	//
	enabled, err := s.IsStreamEnabled(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check if stream is enabled: %w", err)
	}

	if enabled {
		s.wasStreamEnabled = true

		return false, nil
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
func (s *DynamoDBStream) Poll(ctx context.Context) ([]streamTypes.Record, error) {
	//
	if s.shardIterator == "" {
		if err := s.initializeShardIterator(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize shard iterator: %w", err)
		}
	}

	// Use the current shard iterator to retrieve records.
	out, err := s.streamsClient.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
		ShardIterator: &s.shardIterator,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	// Update the shard iterator for the next poll.
	s.shardIterator = aws.ToString(out.NextShardIterator)
	return out.Records, nil
}

// initializeShardIterator sets up the shard iterator by retrieving the stream ARN,
// describing the stream to find a shard, and then getting a shard iterator for that shard.
func (s *DynamoDBStream) initializeShardIterator(ctx context.Context) error {
	// Get the stream ARN from the table description.
	streamArn, err := s.getStreamArn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get stream ARN: %w", err)
	}

	// Describe the stream to get shard information.
	descStreamOut, err := s.streamsClient.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: &streamArn,
	})
	if err != nil {
		return fmt.Errorf("failed to describe stream: %w", err)
	}

	if descStreamOut.StreamDescription == nil || len(descStreamOut.StreamDescription.Shards) == 0 {
		return fmt.Errorf("no shards found in stream")
	}

	// For simplicity, pick the first shard.
	shardID := aws.ToString(descStreamOut.StreamDescription.Shards[0].ShardId)

	iterOut, err := s.streamsClient.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
		StreamArn: &streamArn, ShardId: &shardID, ShardIteratorType: s.iteratorType,
	})

	if err != nil {
		return fmt.Errorf("failed to get shard iterator: %w", err)
	}

	s.shardIterator = aws.ToString(iterOut.ShardIterator)

	return nil
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
