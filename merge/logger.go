package merge

import (
	"context"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// MergeLoggers is a slice of MergeLogger.
type MergeLoggers []model.MergeLogger

// DesiredLoggers is a slice of DesiredLogger.
type DesiredLoggers []model.DesiredLogger

func (dl DesiredLoggers) NotifyAcknowledge(ctx context.Context, path string, value model.ValueAndTimestamp) {
	for _, l := range dl {
		l.Acknowledge(ctx, path, value)
	}
}

func (ml MergeLoggers) NotifyPrepare(ctx context.Context) error {
	for _, l := range ml {
		if p, ok := l.(model.MergeLoggerPrepare); ok {
			if err := p.Prepare(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ml MergeLoggers) NotifyPost(ctx context.Context, err error) error {
	for _, l := range ml {
		if p, ok := l.(model.MergeLoggerPost); ok {
			if err := p.Post(ctx, err); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ml MergeLoggers) NotifyManaged(
	ctx context.Context,
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	for _, l := range ml {
		l.Managed(ctx, path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
	}
}

func (ml MergeLoggers) NotifyPlain(
	ctx context.Context,
	path string,
	operation model.MergeOperation,
	oldValue, newValue any,
) {
	for _, l := range ml {
		l.Plain(ctx, path, operation, oldValue, newValue)
	}
}
