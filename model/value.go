package model

import "time"

// ValueAndTimestamp is the interface that fields must implement if they
// support timestamp-based merging.
type ValueAndTimestamp interface {
	// GetTimestamp will return the timestamp associated with the value. This is
	// used to determine which value is newer when a merge is commenced.
	GetTimestamp() time.Time
	// GetValue will return the value that the timestamp is associated with.
	//
	// If multiple values, the instance itself is the value and this method
	// will return the _"default"_ value. If the value is a map[string]any
	// it will return all values where the key is the name of the value.
	//
	// The latter gives the caller a way of knowing what values are relevant
	// to e.g. log instead of iterate the whole struct.
	GetValue() any
}

// IdValueAndTimestamp extends ValueAndTimestamp by adding a GetID method.
// This is used for slices/arrays to identify and merge elements by ID rather than position.
// When items in a slice/array implement this interface, the merge algorithm will match
// items by ID rather than by index.
type IdValueAndTimestamp interface {
	ValueAndTimestamp
	// GetID returns a unique identifier for this value, used to match items
	// in slices/arrays during merges.
	GetID() string
}

// Merger is an interface that can be implemented by types that want to
// provide custom merge logic. When a type implements this interface, the
// merge algorithm will defer to the type's Merge method instead of using
// the default algorithm.
type Merger interface {
	// Merge takes another instance of the same type and merges it into this instance
	// according to the provided options. Returns the merged result.
	Merge(other any, mode MergeMode) (any, error)
}

// MergeMode indicates how merging is done regarding deletions.
type MergeMode int

const (
	// ClientIsMaster is when a client is considered the master
	// and deletions are propagated.
	ClientIsMaster MergeMode = 1
	// ServerIsMaster, only updates and additions are propagated.
	ServerIsMaster MergeMode = 2
)

// ValueAndTimestampImpl is a standard implementation of the `ValueAndTimestamp` interface
// mostly used by unit test. Use your own custom logic in production.
type ValueAndTimestampImpl struct {
	Timestamp time.Time
	Value     any
}

func (v *ValueAndTimestampImpl) GetTimestamp() time.Time {
	return v.Timestamp
}

func (v *ValueAndTimestampImpl) GetValue() any {
	return v.Value
}
