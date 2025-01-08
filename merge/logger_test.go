package merge_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	mock.Mock
	AcknowledgedPaths []string
	AddedPaths        []string
	UpdatePaths       []string
}

func (m *MockLogger) Desired(path string, operation model.MergeOperation, value model.ValueAndTimestamp) {
	switch operation {
	case model.MergeOperationRemove:
		m.AcknowledgedPaths = append(m.AcknowledgedPaths, path)
	case model.MergeOperationAdd:
		m.AddedPaths = append(m.AddedPaths, path)
	case model.MergeOperationUpdate:
		m.UpdatePaths = append(m.UpdatePaths, path)
	default:
		panic(fmt.Sprintf("unexpected operation: %s", operation.String()))
	}
}

func (m *MockLogger) Managed(path string, operation model.MergeOperation, oldValue, newValue model.ValueAndTimestamp, oldTimeStamp, newTimeStamp time.Time) {
	m.Called(path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
}

func (m *MockLogger) Plain(path string, operation model.MergeOperation, oldValue, newValue any) {
	m.Called(path, operation, oldValue, newValue)
}

func TestLoggerProcessedCalled(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo},
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Updated timestamp
			}},
		},
	}

	// Expect the Processed method to be called with specific arguments
	mockLogger := &MockLogger{}

	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once().
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0],
			&newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerPlainCalledForUnchangedValue(t *testing.T) {
	now := time.Now().UTC()

	oldDevice := Device{
		Name: "SameDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},
			}},
		},
	}

	newDevice := Device{
		Name: "SameDevice", // Unchanged plain value
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for unchanged values
	mockLogger.
		On("Plain", "Name", model.MergeOperationNotChanged, "SameDevice", "SameDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect Processed to be called for Circuits.Sensors
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationNotChanged,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			now, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedCalledForAddedValue(t *testing.T) {
	now := time.Now().UTC()

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // No sensors in old device
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // New sensor added
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for updated plain fields
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method to be called for the added sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[0],
			time.Time{}, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedCalledForUpdatedValue(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo}, // Older timestamp
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Updated timestamp
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method to be called for the updated sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedCalledForRemovedValue(t *testing.T) {
	now := time.Now().UTC()

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Sensor present in old device
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Sensor removed in new device
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method to be called for the removed sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationRemove,
			&oldDevice.Circuits[0].Sensors[0], nil,
			now, time.Time{}).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedCalledForNestedValues(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo}, // Nested managed value
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: oneHourAgo}, // Another nested managed value
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Updated timestamp for first sensor
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: oneHourAgo}, // No change for second sensor
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once().
		On("Plain", "Circuits.1.ID", model.MergeOperationNotChanged, 2, 2).Once()

	// Expect the Processed method for nested managed values
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once().
		On("Managed", "Circuits.1.Sensors.0", model.MergeOperationNotChanged,
			&oldDevice.Circuits[1].Sensors[0], &newDevice.Circuits[1].Sensors[0],
			oneHourAgo, oneHourAgo).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerPlainCalledForAddedPlainValue(t *testing.T) {
	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Empty sensors in the old device
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}},
			{ID: 2, Sensors: []Sensor{}}, // New circuit added
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field and first circuit entry
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Plain method for the added circuit
	mockLogger.
		On("Plain", "Circuits.1.ID", model.MergeOperationAdd, nil, 2).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerPlainCalledForRemovedPlainValue(t *testing.T) {
	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Circuit present in old device
			{ID: 2, Sensors: []Sensor{}}, // Circuit to be removed
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Circuit remains
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field and first circuit entry
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Plain method for the removed circuit
	mockLogger.
		On("Plain", "Circuits.1.ID", model.MergeOperationRemove, 2, nil).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedAndPlainForMixedUpdates(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo}, // To be updated
			}},
			{ID: 2, Sensors: []Sensor{}}, // To be removed
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},        // Updated timestamp
				{ID: 11, TimeStamp: oneHourAgo}, // New sensor added
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for the updated Name field and first circuit entry
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Plain method for the removed circuit
	mockLogger.
		On("Plain", "Circuits.1.ID", model.MergeOperationRemove, 2, nil).Once()

	// Expect the Processed method for the updated sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once()

	// Expect the Processed method for the added sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[1],
			time.Time{}, oneHourAgo).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForEqualTimestamps(t *testing.T) {
	now := time.Now().UTC()

	oldDevice := Device{
		Name: "Device",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Sensor with a specific timestamp
			}},
		},
	}

	newDevice := Device{
		Name: "Device", // Unchanged name
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Same timestamp as old device
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method to be called for unchanged Name
	mockLogger.
		On("Plain", "Name", model.MergeOperationNotChanged, "Device", "Device").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method to indicate no change for the sensor
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationNotChanged,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			now, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMultipleNestedChanges(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo},
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: twoHoursAgo},
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},        // Updated timestamp
				{ID: 12, TimeStamp: oneHourAgo}, // New sensor added
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: now}, // Updated timestamp
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once().
		On("Plain", "Circuits.1.ID", model.MergeOperationNotChanged, 2, 2).Once()

	// Expect the Processed method for sensor updates and additions
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once().
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[1],
			time.Time{}, oneHourAgo).Once().
		On("Managed", "Circuits.1.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[1].Sensors[0], &newDevice.Circuits[1].Sensors[0],
			twoHoursAgo, now).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMultipleAdditions(t *testing.T) {
	now := time.Now().UTC()

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Initially empty
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},                     // First added sensor
				{ID: 11, TimeStamp: now},                     // Second added sensor
				{ID: 12, TimeStamp: now.Add(-1 * time.Hour)}, // Third added sensor
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for updated Name
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method for added sensors
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[0],
			time.Time{}, now).Once().
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[1],
			time.Time{}, now).Once().
		On("Managed", "Circuits.0.Sensors.2", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[2],
			time.Time{}, now.Add(-1*time.Hour)).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMultipleRemovals(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo}, // First sensor to be removed
				{ID: 11, TimeStamp: oneHourAgo}, // Second sensor to be removed
				{ID: 12, TimeStamp: oneHourAgo}, // Third sensor to be removed
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // All sensors removed
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for the updated Name
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method for removed sensors
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationRemove,
			&oldDevice.Circuits[0].Sensors[0], nil,
			oneHourAgo, time.Time{}).Once().
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationRemove,
			&oldDevice.Circuits[0].Sensors[1], nil,
			oneHourAgo, time.Time{}).Once().
		On("Managed", "Circuits.0.Sensors.2", model.MergeOperationRemove,
			&oldDevice.Circuits[0].Sensors[2], nil,
			oneHourAgo, time.Time{}).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMixedOperationsAcrossCircuits(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo},
				{ID: 11, TimeStamp: oneHourAgo},
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: twoHoursAgo},
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now}, // Updated timestamp
				// Removed
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: now},        // Updated timestamp
				{ID: 21, TimeStamp: oneHourAgo}, // New sensor added
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once().
		On("Plain", "Circuits.1.ID", model.MergeOperationNotChanged, 2, 2).Once()

	// Circuit 1: Sensors
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[0].Sensors[0], &newDevice.Circuits[0].Sensors[0],
			oneHourAgo, now).Once().
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationRemove,
			&oldDevice.Circuits[0].Sensors[1], nil,
			oneHourAgo, time.Time{}).Once()

	// Circuit 2: Sensors
	mockLogger.
		On("Managed", "Circuits.1.Sensors.0", model.MergeOperationUpdate,
			&oldDevice.Circuits[1].Sensors[0], &newDevice.Circuits[1].Sensors[0],
			twoHoursAgo, now).Once().
		On("Managed", "Circuits.1.Sensors.1", model.MergeOperationAdd,
			nil, &newDevice.Circuits[1].Sensors[1],
			time.Time{}, oneHourAgo).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerPlainForMultiplePlainValueUpdates(t *testing.T) {
	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1},
			{ID: 2},
		},
	}

	newDevice := Device{
		Name: "NewDevice", // Updated
		Circuits: []Circuit{
			{ID: 10}, // Updated
			{ID: 20}, // Updated
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for updates in Name and Circuits IDs
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationUpdate, 1, 10).Once().
		On("Plain", "Circuits.1.ID", model.MergeOperationUpdate, 2, 20).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerForEmptyCircuits(t *testing.T) {
	oldDevice := Device{
		Name:     "OldDevice",
		Circuits: []Circuit{}, // Empty circuits
	}

	newDevice := Device{
		Name:     "NewDevice",
		Circuits: []Circuit{}, // Also empty circuits
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for the updated Name field
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once()

	// Expect no calls to Processed for circuits, since both are empty
	mockLogger.AssertNotCalled(t, "Processed")

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerForAddingSensorsToEmptyList(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{}}, // Initially empty
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: now},        // First added sensor
				{ID: 11, TimeStamp: oneHourAgo}, // Second added sensor
			}},
		},
	}

	mockLogger := &MockLogger{}

	// Expect the Plain method for the updated Name
	mockLogger.
		On("Plain", "Name", model.MergeOperationUpdate, "OldDevice", "NewDevice").Once().
		On("Plain", "Circuits.0.ID", model.MergeOperationNotChanged, 1, 1).Once()

	// Expect the Processed method for added sensors
	mockLogger.
		On("Managed", "Circuits.0.Sensors.0", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[0],
			time.Time{}, now).Once().
		On("Managed", "Circuits.0.Sensors.1", model.MergeOperationAdd,
			nil, &newDevice.Circuits[0].Sensors[1],
			time.Time{}, oneHourAgo).Once()

	_, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForUpdatedMapValue(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue", UpdatedAt: oneHourAgo}, // Updated
	}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue", UpdatedAt: now}, // Newer timestamp
	}

	mockLogger := &MockLogger{}

	mockLogger.
		On("Managed", "M.key1", model.MergeOperationUpdate,
			oldMap["key1"], newMap["key1"],
			oneHourAgo, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForAddedMapKey(t *testing.T) {
	now := time.Now().UTC()

	oldMap := map[string]*TimestampedMapVal{}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue", UpdatedAt: now}, // Added key
	}

	mockLogger := &MockLogger{}

	mockLogger.
		On("Managed", "M.key1", model.MergeOperationAdd,
			nil, newMap["key1"],
			time.Time{}, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForRemovedMapKey(t *testing.T) {
	now := time.Now().UTC()

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue", UpdatedAt: now}, // Key to be removed
	}

	newMap := map[string]*TimestampedMapVal{}

	mockLogger := &MockLogger{}

	mockLogger.
		On("Managed", "M.key1", model.MergeOperationRemove,
			oldMap["key1"], nil,
			now, time.Time{}).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForUnchangedMapValue(t *testing.T) {
	now := time.Now().UTC()

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "sameValue", UpdatedAt: now}, // Unchanged
	}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "sameValue", UpdatedAt: now}, // Same value and timestamp
	}

	mockLogger := &MockLogger{}

	mockLogger.
		On("Managed", "M.key1", model.MergeOperationNotChanged,
			oldMap["key1"], newMap["key1"],
			now, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMixedMapOperations(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue1", UpdatedAt: twoHoursAgo}, // To be updated
		"key2": {Value: "oldValue2", UpdatedAt: oneHourAgo},  // To be removed
	}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue1", UpdatedAt: now},        // Updated
		"key3": {Value: "newValue3", UpdatedAt: oneHourAgo}, // Added
	}

	mockLogger := &MockLogger{}

	// Expect the Processed method for each operation
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationUpdate,
			oldMap["key1"], newMap["key1"],
			twoHoursAgo, now).Once().
		On("Managed", "M.key2", model.MergeOperationRemove,
			oldMap["key2"], nil,
			oneHourAgo, time.Time{}).Once().
		On("Managed", "M.key3", model.MergeOperationAdd,
			nil, newMap["key3"],
			time.Time{}, oneHourAgo).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMultipleMapUpdates(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue1", UpdatedAt: twoHoursAgo}, // To be updated
		"key2": {Value: "oldValue2", UpdatedAt: oneHourAgo},  // To be updated
	}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue1", UpdatedAt: now}, // Updated
		"key2": {Value: "newValue2", UpdatedAt: now}, // Updated
	}

	mockLogger := &MockLogger{}

	// Expect the Processed method for each updated key
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationUpdate,
			oldMap["key1"], newMap["key1"],
			twoHoursAgo, now).Once().
		On("Managed", "M.key2", model.MergeOperationUpdate,
			oldMap["key2"], newMap["key2"],
			oneHourAgo, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForAdditionsAndRemovalsInMap(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue1", UpdatedAt: oneHourAgo}, // To be removed
		"key2": {Value: "oldValue2", UpdatedAt: oneHourAgo}, // To be removed
	}

	newMap := map[string]*TimestampedMapVal{
		"key3": {Value: "newValue3", UpdatedAt: now}, // Added
		"key4": {Value: "newValue4", UpdatedAt: now}, // Added
	}

	mockLogger := &MockLogger{}

	// Expect the Processed method for removed keys
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationRemove,
			oldMap["key1"], nil,
			oneHourAgo, time.Time{}).Once().
		On("Managed", "M.key2", model.MergeOperationRemove,
			oldMap["key2"], nil,
			oneHourAgo, time.Time{}).Once()

	// Expect the Processed method for added keys
	mockLogger.
		On("Managed", "M.key3", model.MergeOperationAdd,
			nil, newMap["key3"],
			time.Time{}, now).Once().
		On("Managed", "M.key4", model.MergeOperationAdd,
			nil, newMap["key4"],
			time.Time{}, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForMixedOperationsAndUnchangedValuesInMap(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue1", UpdatedAt: twoHoursAgo},     // To be updated
		"key2": {Value: "unchangedValue", UpdatedAt: oneHourAgo}, // Unchanged
		"key3": {Value: "oldValue3", UpdatedAt: oneHourAgo},      // To be removed
	}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue1", UpdatedAt: now},             // Updated
		"key2": {Value: "unchangedValue", UpdatedAt: oneHourAgo}, // Unchanged
		"key4": {Value: "newValue4", UpdatedAt: now},             // Added
	}

	mockLogger := &MockLogger{}

	// Expect the Processed method for updated key
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationUpdate,
			oldMap["key1"], newMap["key1"],
			twoHoursAgo, now).Once()

	// Expect the Processed method for unchanged key
	mockLogger.
		On("Managed", "M.key2", model.MergeOperationNotChanged,
			oldMap["key2"], newMap["key2"],
			oneHourAgo, oneHourAgo).Once()

	// Expect the Processed method for removed key
	mockLogger.
		On("Managed", "M.key3", model.MergeOperationRemove,
			oldMap["key3"], nil,
			oneHourAgo, time.Time{}).Once()

	// Expect the Processed method for added key
	mockLogger.
		On("Managed", "M.key4", model.MergeOperationAdd,
			nil, newMap["key4"],
			time.Time{}, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForOnlyAdditionsInMap(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{}

	newMap := map[string]*TimestampedMapVal{
		"key1": {Value: "newValue1", UpdatedAt: now},        // Added
		"key2": {Value: "newValue2", UpdatedAt: oneHourAgo}, // Added
		"key3": {Value: "newValue3", UpdatedAt: now},        // Added
	}

	mockLogger := &MockLogger{}

	// Expect the Processed method for each added key
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationAdd,
			nil, newMap["key1"],
			time.Time{}, now).Once().
		On("Managed", "M.key2", model.MergeOperationAdd,
			nil, newMap["key2"],
			time.Time{}, oneHourAgo).Once().
		On("Managed", "M.key3", model.MergeOperationAdd,
			nil, newMap["key3"],
			time.Time{}, now).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestLoggerProcessedForOnlyRemovalsInMap(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldMap := map[string]*TimestampedMapVal{
		"key1": {Value: "oldValue1", UpdatedAt: twoHoursAgo}, // To be removed
		"key2": {Value: "oldValue2", UpdatedAt: oneHourAgo},  // To be removed
		"key3": {Value: "oldValue3", UpdatedAt: now},         // To be removed
	}

	newMap := map[string]*TimestampedMapVal{}

	mockLogger := &MockLogger{}

	// Expect the Processed method for each removed key
	mockLogger.
		On("Managed", "M.key1", model.MergeOperationRemove,
			oldMap["key1"], nil,
			twoHoursAgo, time.Time{}).Once().
		On("Managed", "M.key2", model.MergeOperationRemove,
			oldMap["key2"], nil,
			oneHourAgo, time.Time{}).Once().
		On("Managed", "M.key3", model.MergeOperationRemove,
			oldMap["key3"], nil,
			now, time.Time{}).Once()

	type MapHolder struct {
		M map[string]*TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	_, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Loggers: merge.MergeLoggers{mockLogger},
		Mode:    merge.ClientIsMaster,
	})

	require.NoError(t, err)
	mockLogger.AssertExpectations(t)
}
