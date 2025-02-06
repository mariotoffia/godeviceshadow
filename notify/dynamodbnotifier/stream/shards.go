package stream

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

type ShardsOptions struct {
	// NumPollBeforeShardDiscovery specifies how many `Poll` invocations needs
	// to be done before calling _DynamoDB_ to query for new shards.
	//
	// Default is 5.
	NumPollsBeforeShardDiscovery int
	// LastEvaluatedShardId will make sure to start iterating shards from last discovery.
	LastEvaluatedShardId string
}

// Shards manages the active shards, including discovery,
// round-robin selection, polling, and commit handling.
type Shards struct {
	// streamArn is the _ARN_ to the table's stream.
	streamArn string
	// iteratorType is determines what the iterator sees when open a new shard
	iteratorType streamTypes.ShardIteratorType
	// shardManager handles discovery and management of shards.
	shardManager *ShardManager
}

// NewShards initializes a Shards instance and discovers initial shards.
func NewShards(
	ctx context.Context,
	streamsClient *dynamodbstreams.Client,
	streamArn string,
	iteratorType streamTypes.ShardIteratorType,
	opts ...ShardsOptions,
) (*Shards, error) {
	var opt ShardsOptions

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.NumPollsBeforeShardDiscovery == 0 {
		opt.NumPollsBeforeShardDiscovery = 5
	}

	d := &Shards{
		streamArn:    streamArn,
		iteratorType: iteratorType,
		shardManager: NewShardManager(streamArn, iteratorType, opt.NumPollsBeforeShardDiscovery),
	}

	return d, nil
}

// Poll retrieves records from one shard using round-robin selection.
//
// It always uses the shard's CommittedIterator for GetRecords,
// and updates its WorkingIterator with the NextShardIterator returned.
// If the NextShardIterator is nil, the shard is considered closed and removed.
func (d *Shards) Poll(
	ctx context.Context,
	streamsClient *dynamodbstreams.Client,
) (string, []streamTypes.Record, error) {
	var lastShardID string

	for {
		shard, err := d.shardManager.NextShard(ctx, streamsClient)

		if err != nil {
			return "", nil, fmt.Errorf("failed to get next shard: %w", err)
		}

		if shard == nil {
			// No more shards to poll.
			return "", nil, nil
		}

		if lastShardID != "" && lastShardID == shard.ShardID {
			// All shards have been polled.
			return "", nil, nil
		}

		lastShardID = shard.ShardID

		// Poll records from the shard.
		out, err := streamsClient.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: aws.String(shard.CommittedIterator),
		})

		if err != nil {
			return "", nil, fmt.Errorf("failed to get records for shard %s: %w", shard.ShardID, err)
		}

		// NextShardIterator == nil -> shard is closed.
		if out.NextShardIterator == nil {
			records := out.Records

			// Will be delete when client `Commit` the shard.
			d.shardManager.MarkedForDelete(shard.ShardID)

			if len(records) > 0 {
				return shard.ShardID, records, nil
			}

			continue
		}

		d.shardManager.UpdateWorkingIterator(shard.ShardID, aws.ToString(out.NextShardIterator))

		if len(out.Records) > 0 {
			return shard.ShardID, out.Records, nil
		}

		// No records -> make sure to use next shard iterator.
		d.shardManager.Commit(shard.ShardID)
	}
}

// Commit advances the committed iterator for a given shard to the working iterator.
// This is called after successful processing of the records.
func (d *Shards) Commit(shardID string) {
	d.shardManager.Commit(shardID)
}
