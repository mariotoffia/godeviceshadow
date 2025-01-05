package str

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// StringLogger is a `model.MergeLogger` that logs the operations to a bytes buffer
// where it is possible to retrieve the log as a string.
type StringLogger struct {
	// log is the buffer that contains the log.
	log bytes.Buffer
}

// NewStringLogger creates a new StringLogger.
func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

// String returns the log as a string.
func (sl *StringLogger) String() string {
	return sl.log.String()
}

func (sl *StringLogger) Processed(
	path string,
	operation model.MergeOperation,
	oldValue, newValue any,
	oldTimeStamp, newTimeStamp time.Time,
) {
	var ov, nv string

	if oldValue != nil {
		ov = fmt.Sprintf("%v", oldValue)
	} else {
		ov = "nil"
	}

	if newValue != nil {
		nv = fmt.Sprintf("%v", newValue)
	} else {
		nv = "nil"
	}

	sl.log.WriteString(path)
	sl.log.WriteString(" ")
	sl.log.WriteString(operation.String())
	sl.log.WriteString(" ")
	sl.log.WriteString(ov)
	sl.log.WriteString(" ")
	sl.log.WriteString(nv)
	sl.log.WriteString(" ")
	sl.log.WriteString(oldTimeStamp.Format(time.RFC3339))
	sl.log.WriteString(" ")
	sl.log.WriteString(newTimeStamp.Format(time.RFC3339))
}
