package stream

// ShardState holds the state for a shard and its iterators.
type ShardState struct {
	ShardID           string // Unique shard identifier.
	ParentShardID     string // Parent shard identifier (if split).
	CommittedIterator string // Last committed iterator.
	WorkingIterator   string // Iterator returned from the latest GetRecords call.
	MarkedForDelete   bool   // Flag to indicate the shard is closed and should be removed
}

// ShardCollection manages shard states using a map for quick lookup
// and an ordered slice for round-robin selection.
type ShardCollection struct {
	shards map[string]*ShardState
	order  []string
}

// NewShardCollection creates a new instance of ShardCollection.
func NewShardCollection() *ShardCollection {
	return &ShardCollection{
		shards: make(map[string]*ShardState),
		order:  []string{},
	}
}

// Upsert inserts or updates a shard _state_ in the collection.
// If the shard is new, it appends the shardID to the ordered list.
func (sc *ShardCollection) Upsert(state *ShardState) {
	shardID := state.ShardID

	if _, exists := sc.shards[shardID]; !exists {
		sc.order = append(sc.order, shardID)
	}

	sc.shards[shardID] = state
}

// Delete removes a shard from the collection using its shardID.
func (sc *ShardCollection) Delete(shardID string) {
	if _, exists := sc.shards[shardID]; !exists {
		return
	}

	delete(sc.shards, shardID)

	for i, id := range sc.order {
		if id == shardID {
			sc.order = append(sc.order[:i], sc.order[i+1:]...)
			break
		}
	}
}

// Get retrieves the shard state associated with the given _shardID_.
// It returns `nil` if the shard is not found.
func (sc *ShardCollection) Get(shardID string) *ShardState {
	if state, exists := sc.shards[shardID]; exists {
		return state
	}

	return nil
}

func (sc *ShardCollection) GetByIndex(index int) *ShardState {
	if index < 0 || index >= len(sc.order) {
		return nil
	}

	return sc.Get(sc.order[index])
}

func (sc *ShardCollection) Size() int {
	return len(sc.order)
}

// List returns a slice of all shard states in the order maintained.
func (sc *ShardCollection) List() []*ShardState {
	result := make([]*ShardState, 0, len(sc.order))

	for _, id := range sc.order {
		if state, exists := sc.shards[id]; exists {
			result = append(result, state)
		}
	}

	return result
}
