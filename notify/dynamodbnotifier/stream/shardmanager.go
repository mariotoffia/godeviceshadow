package stream

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"

	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

// ShardManager handles discovery, state management, and iterator updates for shards.
type ShardManager struct {
	streamArn                    string
	iteratorType                 streamTypes.ShardIteratorType
	activeShards                 *ShardCollection
	lastShardIndex               int    // Index to track round-robin position.
	lastEvaluatedShardID         string // Paginating DescribeStream calls.
	numPollsBeforeShardDiscovery int    // Number of calls to `NextShard` before calling `discoverShards`.
	cntNextShards                int    // Counter to track number of calls to `NextShard`.
}

// ShardManagerState is used for persisting and resuming shard states.
type ShardManagerState struct {
	LastEvaluatedShardID string       `json:"last_evaluated_shard_id"`
	Shards               []ShardState `json:"shards"`
}

// NewShardManager creates a new instance of ShardManager.
func NewShardManager(streamArn string, iteratorType streamTypes.ShardIteratorType, numPollsBeforeShardDiscovery int) *ShardManager {
	return &ShardManager{
		streamArn:                    streamArn,
		iteratorType:                 iteratorType,
		numPollsBeforeShardDiscovery: numPollsBeforeShardDiscovery,
		activeShards:                 NewShardCollection(),
		lastShardIndex:               0,
	}
}

// discoverShards queries DynamoDB Streams to discover new shards,
// update existing shard states (adding parent shard info if shard split),
// and remove closed shards.
func (sm *ShardManager) discoverShards(ctx context.Context, client *dynamodbstreams.Client) error {
	var lastShardID *string

	if sm.lastEvaluatedShardID != "" {
		lastShardID = aws.String(sm.lastEvaluatedShardID)
	}

	// Collect all shards discovered across paginated requests.
	var discoveredShards []streamTypes.Shard

	for {
		input := &dynamodbstreams.DescribeStreamInput{
			StreamArn:             aws.String(sm.streamArn),
			ExclusiveStartShardId: lastShardID,
		}

		resp, err := client.DescribeStream(ctx, input)

		if err != nil {
			return fmt.Errorf("DescribeStream error: %w", err)
		}

		if resp.StreamDescription == nil {
			return fmt.Errorf("no stream description available")
		}

		discoveredShards = append(discoveredShards, resp.StreamDescription.Shards...)

		// If there are no more pages, break out.
		if resp.StreamDescription.LastEvaluatedShardId == nil {
			break
		}

		lastShardID = resp.StreamDescription.LastEvaluatedShardId
		sm.lastEvaluatedShardID = aws.ToString(lastShardID)
	}

	//
	// Process
	//
	for _, shard := range discoveredShards {
		shardID := aws.ToString(shard.ShardId)

		// already tracked?
		if existing := sm.activeShards.Get(shardID); existing != nil {
			// Update parent shard info if shard split.
			if existing.ParentShardID == "" && shard.ParentShardId != nil {
				existing.ParentShardID = aws.ToString(shard.ParentShardId)
			}

			continue
		}

		// new shard -> get iterator.
		iterInput := &dynamodbstreams.GetShardIteratorInput{
			StreamArn:         aws.String(sm.streamArn),
			ShardId:           shard.ShardId,
			ShardIteratorType: sm.iteratorType,
		}

		iterResp, err := client.GetShardIterator(ctx, iterInput)

		if err != nil {
			return fmt.Errorf("failed to get shard iterator for shard %s: %w", shardID, err)
		}

		initialIterator := aws.ToString(iterResp.ShardIterator)

		if initialIterator == "" {
			// no valid iterator -> skip
			continue
		}

		// Add the new shard state.
		sm.activeShards.Upsert(&ShardState{
			ShardID:           shardID,
			ParentShardID:     aws.ToString(shard.ParentShardId),
			CommittedIterator: initialIterator,
			WorkingIterator:   initialIterator,
		})
	}

	return nil
}

// NextShard returns a shard to be used.
//
// If shard is detected as closed, it is removed from the active shards.
//
// If no more active shards, this will return `nil`. It will automatically
// call `discoverShards` to refresh the shard list.
func (sm *ShardManager) NextShard(ctx context.Context, client *dynamodbstreams.Client) (*ShardState, error) {
	sm.cntNextShards++

	size := sm.activeShards.Size()

	//
	// Shard discovery.
	//
	if size == 0 {
		if err := sm.discoverShards(ctx, client); err != nil {
			return nil, err
		}
	} else if sm.cntNextShards%sm.numPollsBeforeShardDiscovery == 0 {
		if err := sm.discoverShards(ctx, client); err != nil {
			return nil, err
		}
	}

	if size = sm.activeShards.Size(); size == 0 {
		return nil, nil
	}

	// Select shards using round-robin starting from lastShardIndex.
	var shard *ShardState

	for i := 0; i < size; i++ {
		idx := (sm.lastShardIndex + i) % size

		if shard = sm.activeShards.GetByIndex(idx); shard == nil {
			return nil, nil
		}

		// Set last index + 1 -> next call will start from the next shard.
		sm.lastShardIndex = (sm.lastShardIndex + i + 1) % size

		return shard, nil
	}

	return nil, nil
}

// CloseShard is when a shard is detected as closed and should be removed.
func (sm *ShardManager) CloseShard(shardID string) {
	sm.activeShards.Delete(shardID)
}

// MarkedForDelete will mark a shard for deletion.
//
// The delete is done when doing `Commit`. This allows the poller
// return records and when the client is done processing, it will
// commit and thus will delete the shard.
//
// This is not to miss any records that are in the shard but it is
// finished otherwise.
func (sm *ShardManager) MarkedForDelete(shardID string) {
	if state := sm.activeShards.Get(shardID); state != nil {
		state.MarkedForDelete = true
	}
}

// Commit advances the committed iterator for the shard with shardID
// to the current working iterator. This should be called after successful processing.
func (sm *ShardManager) Commit(shardID string) {
	if state := sm.activeShards.Get(shardID); state != nil {
		if state.MarkedForDelete {
			sm.CloseShard(shardID)
		} else {
			state.CommittedIterator = state.WorkingIterator
		}
	}
}

func (sm *ShardManager) UpdateWorkingIterator(shardID, iterator string) {
	if state := sm.activeShards.Get(shardID); state != nil {
		state.WorkingIterator = iterator
	}
}

// GetStates returns a snapshot of all shard states for persistence.
func (sm *ShardManager) GetStates() ShardManagerState {
	list := sm.activeShards.List()

	states := make([]ShardState, 0, len(list))
	for _, state := range list {
		states = append(states, *state)
	}
	return ShardManagerState{
		LastEvaluatedShardID: sm.lastEvaluatedShardID,
		Shards:               states,
	}
}

// SetStates loads shard states from a persisted snapshot,
// allowing the manager to resume from a known state.
func (sm *ShardManager) SetStates(state ShardManagerState) {
	sm.activeShards = NewShardCollection()

	for _, s := range state.Shards {
		sm.activeShards.Upsert(&s)
	}

	sm.lastEvaluatedShardID = state.LastEvaluatedShardID
	sm.lastShardIndex = 0
}
