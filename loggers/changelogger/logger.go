package changelogger

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// New creates a new ChangeMergeLogger.
func New() *ChangeMergeLogger {
	return &ChangeMergeLogger{
		PlainLog:   PlainLogMap{},
		ManagedLog: ManagedLogMap{},
	}
}

// New implements the `model.CreatableMergeLogger` interface.
func (sl *ChangeMergeLogger) New() model.MergeLogger {
	return New()
}

// Find finds a `ChangeMergeLogger` in a slice of `model.MergeLogger`.
// If not found, it returns `nil`.
func Find(loggers []model.MergeLogger) *ChangeMergeLogger {
	for _, logger := range loggers {
		if sl, ok := logger.(*ChangeMergeLogger); ok {
			return sl
		}
	}

	return nil
}

func (sl *ChangeMergeLogger) Plain(path string, operation model.MergeOperation, oldValue, newValue any) {
	sl.PlainLog[operation] = append(sl.PlainLog[operation], PlainValue{
		Path:     path,
		OldValue: oldValue,
		NewValue: newValue,
	})
}

func (sl *ChangeMergeLogger) Managed(
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	sl.ManagedLog[operation] = append(sl.ManagedLog[operation], ManagedValue{
		Path:         path,
		OldValue:     oldValue,
		NewValue:     newValue,
		OldTimeStamp: oldTimeStamp,
		NewTimeStamp: newTimeStamp,
	})
}
