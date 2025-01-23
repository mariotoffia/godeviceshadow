package manager

import "github.com/mariotoffia/godeviceshadow/model"

// createMergeLoggers will create logger instance from _loggers_ (if any), if none where submitted, it will use the `Manager.reportedLoggers`.
//
// If the `MergeDirtyLogger` is not present in the _loggers_ it will be automatically added.
//
// If _report_ is `true` it will use the `Manager.reportedLoggers` when _loggers_ is empty. Otherwise it will use `Manager.desiredMergeLoggers`.
func (mgr *Manager) createMergeLoggers(report bool, loggers []model.CreatableMergeLogger) []model.MergeLogger {
	if len(loggers) == 0 {
		if report {
			loggers = mgr.reportedMergeLoggers
		} else {
			loggers = mgr.desiredMergeLoggers
		}
	}

	// Add dirty detection
	if !HasMergeDirtyLoggerCreator(loggers) {
		loggers = append(loggers, &MergeDirtyLogger{})
	}

	res := make([]model.MergeLogger, 0, len(loggers))

	for _, lg := range loggers {
		res = append(res, lg.New())
	}

	return res
}
