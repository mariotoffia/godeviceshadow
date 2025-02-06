package loggerutils

import "github.com/mariotoffia/godeviceshadow/model"

func FindCreatableMerge[T any](loggers []model.CreatableMergeLogger) T {
	for _, logger := range loggers {
		if sl, ok := logger.(any).(T); ok {
			return sl
		}
	}

	var zero T

	return zero
}

func FindCreatableDesire[T any](loggers []model.CreatableDesiredLogger) T {
	for _, logger := range loggers {
		if sl, ok := logger.(any).(T); ok {
			return sl
		}
	}

	var zero T

	return zero
}

func FindMerge[T any](loggers []model.MergeLogger) T {
	for _, logger := range loggers {
		if sl, ok := logger.(any).(T); ok {
			return sl
		}
	}

	var zero T

	return zero
}

func FindDesire[T any](loggers []model.DesiredLogger) T {
	for _, logger := range loggers {
		if sl, ok := logger.(any).(T); ok {
			return sl
		}
	}

	var zero T

	return zero
}
