package stdmgr

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// DesiredAckLogger purpose is to set a desired document as dirty if any
// desired "managed" properties where acknowledged and thus needs to be persisted.
type DesiredAckLogger struct {
	Dirty bool
}

// MergeDirtyLogger will detect all operations except unchanged as dirty.
type MergeDirtyLogger struct {
	Dirty bool
	// IgnorePlain instructs the dirty detection to ignore plain values and
	// only react to managed values.
	IgnorePlain bool
}

// New implements the `model.CreatableDesiredLogger`
func (dal *DesiredAckLogger) New() model.DesiredLogger {
	return &DesiredAckLogger{}
}

// Acknowledge is called when a desired value has been acknowledged (`model.DesiredLogger` interface).
func (dal *DesiredAckLogger) Acknowledge(path string, value model.ValueAndTimestamp) {
	dal.Dirty = true
}

// New implements the `model.CreatableMergeLogger`
func (mdl *MergeDirtyLogger) New() model.MergeLogger {
	return &MergeDirtyLogger{}
}

func (mdl *MergeDirtyLogger) Managed(
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	if operation != model.MergeOperationNotChanged {
		mdl.Dirty = true
	}
}

func (mdl *MergeDirtyLogger) Plain(path string, operation model.MergeOperation, oldValue, newValue any) {
	if !mdl.IgnorePlain && operation != model.MergeOperationNotChanged {
		mdl.Dirty = true
	}
}

// HasDesiredAckLoggerCreator checks if there is any `DesiredAckLogger` inside the _loggers_ slice.
func HasDesiredAckLoggerCreator(loggers []model.CreatableDesiredLogger) bool {
	for _, logger := range loggers {
		if _, ok := logger.(*DesiredAckLogger); ok {
			return true
		}
	}
	return false
}

func HasMergeDirtyLoggerCreator(loggers []model.CreatableMergeLogger) bool {
	for _, logger := range loggers {
		if _, ok := logger.(*MergeDirtyLogger); ok {
			return true
		}
	}
	return false
}

// FindDesiredAckLogger finds the `DesiredAckLogger` instance inside the _loggers_ slice.
func FindDesiredAckLogger(loggers []model.DesiredLogger) (*DesiredAckLogger, bool) {
	for _, logger := range loggers {
		if dal, ok := logger.(*DesiredAckLogger); ok {
			return dal, true
		}
	}
	return nil, false
}

// FindMergeDirtyLogger finds the `MergeDirtyLogger` instance inside the _loggers_ slice.
func FindMergeDirtyLogger(loggers []model.MergeLogger) (*MergeDirtyLogger, bool) {
	for _, logger := range loggers {
		if mdl, ok := logger.(*MergeDirtyLogger); ok {
			return mdl, true
		}
	}
	return nil, false
}
