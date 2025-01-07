package model

import (
	"time"
)

// MergeLogger is a interface that will be called in the different merge
// operations that has been performed.
type MergeLogger interface {
	// Managed is called when a managed value (`ValueAndTimestamp`) has been processed in a
	// merge operation. It is even called when a value has not been changed.
	//
	// The path is a a _JSON_ path to the value that has been processed. It is extracted from the
	// field names, map keys, and slice indexes. If field names do have a JSON tag, the tag is used instead.
	Managed(
		path string,
		operation MergeOperation,
		oldValue, newValue ValueAndTimestamp,
		oldTimeStamp, newTimeStamp time.Time,
	)

	// Plain is called when a value has been processed in a merge operation. It is even called when a value
	// has not been changed. This is called when a "plain" value has been processed and not a "managed" value.
	Plain(path string, operation MergeOperation, oldValue, newValue any)
}

// MergeLoggerPrepare will be called before any merge operation occurs.
type MergeLoggerPrepare interface {
	// Prepare will be called just before any merge operation is taking place.
	// If it returns an error, the merge operation _may_ be aborted.
	Prepare() error
}

// MergeLoggerPost is called after all merge operations have been performed.
type MergeLoggerPost interface {
	// Post is invoked when finished (either successfully or erroneously) and
	// returns an error if the post operation failed.
	Post(err error) error
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
	// from a struct, map or slice.
	MergeOperationRemove MergeOperation = 3
	// MergeOperationNotChanged is the operation that indicates that a value has not been
	// changed and thus left as is.
	MergeOperationNotChanged MergeOperation = 4
)

func (op MergeOperation) String() string {
	switch op {
	case MergeOperationAdd:
		return "add"
	case MergeOperationUpdate:
		return "update"
	case MergeOperationRemove:
		return "remove"
	case MergeOperationNotChanged:
		return "noop"
	default:
		return "unknown"
	}
}
