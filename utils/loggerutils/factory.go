package loggerutils

import "github.com/mariotoffia/godeviceshadow/model"

func CreateMergeLoggers(creators ...model.CreatableMergeLogger) []model.MergeLogger {
	loggers := make([]model.MergeLogger, 0, len(creators))

	for _, creator := range creators {
		loggers = append(loggers, creator.New())
	}

	return loggers
}

func CreateDesiredLoggers(creators ...model.CreatableDesiredLogger) []model.DesiredLogger {
	loggers := make([]model.DesiredLogger, 0, len(creators))

	for _, creator := range creators {
		loggers = append(loggers, creator.New())
	}

	return loggers
}
