package selectlang_test

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/loggers/desirelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// createTestOperation creates a standard test operation with common values for testing
func createTestOperation() notifiermodel.NotifierOperation {
	// Create a value with timestamp
	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": 22},
	}

	// Create a desire logger with some acknowledged values
	dl := desirelogger.New()
	dl.Acknowledge("device/settings/mode", &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     "auto",
	})

	// Create a merge logger with both managed and plain logs
	ml := changelogger.ChangeMergeLogger{
		ManagedLog: changelogger.ManagedLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "sensors/temperature/indoor",
					NewValue: mvs,
				},
			},
			model.MergeOperationUpdate: {
				{
					Path: "sensors/humidity/indoor",
					NewValue: &model.ValueAndTimestampImpl{
						Timestamp: time.Now().UTC(),
						Value:     map[string]any{"humidity": 45},
					},
				},
			},
		},
		PlainLog: changelogger.PlainLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "devices/status",
					NewValue: "online",
				},
			},
		},
	}

	return notifiermodel.NotifierOperation{
		ID:           persistencemodel.PersistenceID{ID: "device-123", Name: "homeShadow"},
		Operation:    notifiermodel.OperationTypeReport,
		MergeLogger:  ml,
		DesireLogger: *dl,
		Reported:     map[string]any{"status": "active"},
		Desired:      map[string]any{"mode": "auto"},
	}
}

// createComplexTestOperation creates a test operation with more complex data for nested queries
func createComplexTestOperation() notifiermodel.NotifierOperation {
	// Create a value with timestamp
	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": 22},
	}

	reValue := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": "re-123"},
	}

	// Create a merge logger with both managed and plain logs
	ml := changelogger.ChangeMergeLogger{
		ManagedLog: changelogger.ManagedLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "Sensors-123a-indoor",
					NewValue: mvs,
				},
				{
					Path:     "Sensors-456b-indoor",
					NewValue: reValue,
				},
			},
		},
		PlainLog: changelogger.PlainLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "Sensors-789c-indoor",
					NewValue: "temp", // Direct "temp" value for log.Value == 'temp' test
				},
			},
		},
	}

	// Create a test operation that should match the complex query
	return notifiermodel.NotifierOperation{
		ID:          persistencemodel.PersistenceID{ID: "myDevice-123", Name: "myShadow"},
		Operation:   notifiermodel.OperationTypeReport,
		MergeLogger: ml,
	}
}
