package merge

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// MergeLoggers is a slice of MergeLogger.
type MergeLoggers []model.MergeLogger

func (ml MergeLoggers) NotifyProcessed(
	path string,
	operation model.MergeOperation,
	oldValue, newValue any,
	oldTimeStamp, newTimeStamp time.Time,
) {
	for _, l := range ml {
		l.Processed(path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
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
