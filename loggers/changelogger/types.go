package changelogger

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

type ManagedLogMap map[model.MergeOperation][]ManagedValue
type PlainLogMap map[model.MergeOperation][]PlainValue

// ChangeMergeLogger stores the changes that can be queried later.
type ChangeMergeLogger struct {
	PlainLog   PlainLogMap
	ManagedLog ManagedLogMap
}

type PlainValue struct {
	Path     string
	OldValue any
	NewValue any
}

type ManagedValue struct {
	Path         string
	OldValue     model.ValueAndTimestamp
	NewValue     model.ValueAndTimestamp
	OldTimeStamp time.Time
	NewTimeStamp time.Time
}
