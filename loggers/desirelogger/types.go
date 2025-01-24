package desirelogger

import "github.com/mariotoffia/godeviceshadow/model"

type DesireLogger struct {
	acknowledged map[string]model.ValueAndTimestamp
}

func New() *DesireLogger {
	return &DesireLogger{acknowledged: map[string]model.ValueAndTimestamp{}}
}
