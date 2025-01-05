package model

import "time"

// ValueAndTimestamp is the interface that fields must implement if they
// support timestamp-based merging.
type ValueAndTimestamp interface {
	GetTimestamp() time.Time
	SetTimestamp(t time.Time)
	// GetValue will return the value that the timestamp is associated with.
	GetValue() any
}
