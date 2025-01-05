package merge

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// MergeLoggers is a slice of MergeLogger.
type MergeLoggers []model.MergeLogger

func (ml MergeLoggers) NotifyManaged(
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	for _, l := range ml {
		l.Managed(path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
	}
}

func (ml MergeLoggers) NotifyPlain(
	path string,
	operation model.MergeOperation,
	oldValue, newValue any,
) {
	for _, l := range ml {
		l.Plain(path, operation, oldValue, newValue)
	}
}
