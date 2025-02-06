package desirelogger

import "github.com/mariotoffia/godeviceshadow/model"

// DesireLogger is used in a reported operation to log acknowledgement of a desire when
// a report is incoming. That is, this is done when `Report` is called.
type DesireLogger struct {
	acknowledged map[string]model.ValueAndTimestamp
}

func New() *DesireLogger {
	return &DesireLogger{acknowledged: map[string]model.ValueAndTimestamp{}}
}
