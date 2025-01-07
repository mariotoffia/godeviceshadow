package merge

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// MergeLoggers is a slice of MergeLogger.
type MergeLoggers []model.MergeLogger

func (ml MergeLoggers) NotifyPrepare() error {
	for _, l := range ml {
		if p, ok := l.(model.MergeLoggerPrepare); ok {
			if err := p.Prepare(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ml MergeLoggers) NotifyPost(err error) error {
	for _, l := range ml {
		if p, ok := l.(model.MergeLoggerPost); ok {
			if err := p.Post(err); err != nil {
				return err
			}
		}
	}

	return nil
}

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
