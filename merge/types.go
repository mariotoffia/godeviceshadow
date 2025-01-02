package merge

import "time"

// ValueAndTimestamp is the interface that fields must implement if they
// support timestamp-based merging.
type ValueAndTimestamp interface {
	GetTimestamp() time.Time
	SetTimestamp(t time.Time)
}

type MergeOperation int

const (
	// MergeOperationAdd is the operation that indicates that a value has been added
	// to a map or slice.
	MergeOperationAdd MergeOperation = 1
	// MergeOperationUpdate is the operation that indicates that a value has been updated
	// in a map or slice.
	MergeOperationUpdate MergeOperation = 2
	// MergeOperationRemove is the operation that indicates that a value has been removed
	// from a map or slice.
	MergeOperationRemove MergeOperation = 3
	// MergeOperationNotChanged is the operation that indicates that a value has not been
	// changed and thus left as is.
	MergeOperationNotChanged MergeOperation = 4
)

// MergeLogger is a interface that will be called in the different merge
// operations that has been performed.
type MergeLogger interface {
	// Processed is called when a value has been processed in a merge operation. It is even called
	// when a value has not been changed.
	//
	// The path is a a _JSON_ path to the value that has been processed. It is extracted from the
	// field names, map keys, and slice indexes. If field names do have a JSON tag, the tag is used instead.
	Processed(
		path string,
		operation MergeOperation,
		oldValue, newValue any,
		oldTimeStamp, newTimeStamp time.Time,
	)
}

// MergeLoggers is a slice of MergeLogger.
type MergeLoggers []MergeLogger
