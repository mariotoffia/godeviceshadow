package dblogger

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// New creates a new DsqlLogger.
func New(client DbLogger, batchSize int) *DsqlLogger {
	return &DsqlLogger{
		client:    client,
		batchSize: batchSize,
	}
}

// New implements the `model.CreatableMergeLogger` interface.
func (sl *DsqlLogger) New() model.MergeLogger {
	return &DsqlLogger{client: sl.client, batchSize: sl.batchSize}
}

// Find finds a `DsqlLogger` in a slice of `model.MergeLogger`.
// If not found, it returns `nil`.
func Find(loggers []model.MergeLogger) *DsqlLogger {
	for _, logger := range loggers {
		if sl, ok := logger.(*DsqlLogger); ok {
			return sl
		}
	}

	return nil
}

func (sl *DsqlLogger) Plain(ctx context.Context, path string, operation model.MergeOperation, oldValue, newValue any) {
	// NOOP - We're only logging managed values in this implementation
}

func (sl *DsqlLogger) Managed(
	ctx context.Context,
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	// Only process add and update operations, skip remove and unchanged operations
	if !operation.In(model.MergeOperationAdd, model.MergeOperationUpdate) {
		return
	}

	// Skip if we don't have a newValue (shouldn't happen for add/update but just in case)
	if newValue == nil {
		return
	}

	sl.currentBatch = append(sl.currentBatch, LogValue{
		Operation: operation,
		Value:     newValue,
	})

	if len(sl.currentBatch) >= sl.batchSize {
		if err := sl.client.Upsert(ctx, sl.currentBatch); err != nil {
			// Handle error (e.g., log it, return it, etc.)
			return
		}
		sl.currentBatch = nil // Reset the batch after successful upsert
	}
}

// Prepare implements the model.MergeLoggerPrepare interface
func (sl *DsqlLogger) Prepare(ctx context.Context) error {
	if sl.client == nil {
		return fmt.Errorf("client is not initialized (nil)")
	}

	return sl.client.Begin(ctx)
}

// Post implements the model.MergeLoggerPost interface
func (sl *DsqlLogger) Post(ctx context.Context, err error) error {
	if err != nil {
		return err
	}

	if len(sl.currentBatch) > 0 {
		if err := sl.client.Upsert(ctx, sl.currentBatch); err != nil {
			return err
		}

		sl.currentBatch = nil
	}

	// Commit the transaction if supported
	return sl.client.Commit(ctx)
}
