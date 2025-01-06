package str

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// lf is the log format
const lf = "%-10s %-26s %-26s %-30s %-40s %-40s\n"

// StringLogger is a `model.MergeLogger` that logs the operations to a bytes buffer
// where it is possible to retrieve the log as a string.
type StringLogger struct {
	// log is the buffer that contains the log.
	log bytes.Buffer
	// header is set to true when it has printed the header.
	header bool
}

// NewStringLogger creates a new StringLogger.
func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

// String returns the log as a string.
func (sl *StringLogger) String() string {
	return sl.log.String()
}

func (sl *StringLogger) printHeader() {
	if !sl.header {
		fmt.Fprintf(&sl.log, lf,
			"Operation", "Old Timestamp", "New Timestamp", "Path", "OldValue", "NewValue")
		sl.header = true
	}
}
func (sl *StringLogger) Plain(path string, operation model.MergeOperation, oldValue, newValue any) {
	sl.printHeader()

	fmt.Fprintf(&sl.log, lf,
		operation.String(), "", "", path,
		fmt.Sprintf("%v", oldValue),
		fmt.Sprintf("%v", newValue),
	)
}
func (sl *StringLogger) Managed(
	path string,
	operation model.MergeOperation,
	oldValue, newValue model.ValueAndTimestamp,
	oldTimeStamp, newTimeStamp time.Time,
) {
	var ov, nv string

	if oldValue != nil {
		ov = fmt.Sprintf("%v", oldValue.GetValue())
	} else {
		ov = "nil"
	}

	if newValue != nil {
		nv = fmt.Sprintf("%v", newValue.GetValue())
	} else {
		nv = "nil"
	}

	sl.printHeader()

	fmt.Fprintf(&sl.log, lf,
		operation.String(),
		oldTimeStamp.Format(time.RFC3339),
		newTimeStamp.Format(time.RFC3339),
		path, ov, nv,
	)
}
