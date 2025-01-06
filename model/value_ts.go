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
