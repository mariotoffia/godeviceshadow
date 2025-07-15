package loggers

import (
	"context"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/desirelogger"
	"github.com/mariotoffia/godeviceshadow/model"
)

// DynamoDbDesireLogger is a special logger that is a merge logger that
// feeds all merge deletes
type DynamoDbDesireLogger struct {
	*desirelogger.DesireLogger
}

// Find finds a `DynamoDbDesireLogger` in a slice of `model.MergeLogger`.
// If not found, it returns `nil`.
func Find(loggers []model.MergeLogger) *DynamoDbDesireLogger {
	for _, logger := range loggers {
		if sl, ok := logger.(*DynamoDbDesireLogger); ok {
			return sl
		}
	}

	return nil
}

// New implements the `model.CreatableMergeLogger` interface.
func (l *DynamoDbDesireLogger) New() model.MergeLogger {
	return &DynamoDbDesireLogger{DesireLogger: desirelogger.New()}
}

func (sl *DynamoDbDesireLogger) Plain(ctx context.Context, path string, operation model.MergeOperation, oldValue, newValue any) {
	// NOOP
}

// Managed is called in a merge operation, it only handles removed -> add to acknowledged.
func (sl *DynamoDbDesireLogger) Managed(
	ctx context.Context,
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	if operation == model.MergeOperationRemove {
		sl.Acknowledge(ctx, path, newValue)
	}
}
