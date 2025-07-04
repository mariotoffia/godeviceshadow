package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// This test focuses on the notifyRecursive function indirectly
func TestNotifyRecursive(t *testing.T) {
	now := time.Now().UTC()

	// Test Add operation
	t.Run("NotifyAdd", func(t *testing.T) {
		device := Device{
			Name: "TestDevice",
			Circuits: []Circuit{
				{ID: 1, Sensors: []Sensor{
					{ID: 10, TimeStamp: now},
				}},
			},
		}

		mockLogger := &MockLogger{}
		// For merging an empty device with a populated device:
		// 1. Plain fields (Name) use MergeOperationUpdate
		// 2. The array/slice is first notified as Add, then its elements
		// 3. ValueAndTimestamp fields use MergeOperationAdd
		mockLogger.
			On("Plain", "Name", model.MergeOperationUpdate, "", "TestDevice").Once().
			On("Plain", "Circuits.0.ID", model.MergeOperationAdd, nil, 1).Once().
			On("Managed", circuits0Sensors0, model.MergeOperationAdd,
				nil, &device.Circuits[0].Sensors[0],
				time.Time{}, now).Once()

		// Pass empty device as base and the device as client to simulate an Add operation
		emptyDevice := Device{}
		_, err := merge.Merge(emptyDevice, device, merge.MergeOptions{
			Loggers: merge.MergeLoggers{mockLogger},
			Mode:    merge.ClientIsMaster,
		})

		assert.NoError(t, err)
		mockLogger.AssertExpectations(t)
	})

	// Test Remove operation
	t.Run("NotifyRemove", func(t *testing.T) {
		device := Device{
			Name: "TestDevice",
			Circuits: []Circuit{
				{ID: 1, Sensors: []Sensor{
					{ID: 10, TimeStamp: now},
				}},
			},
		}

		mockLogger := &MockLogger{}
		// For merging a populated device with an empty device (in ClientIsMaster mode):
		// The implementation treats empty string as an update rather than a remove
		mockLogger.
			On("Plain", "Name", model.MergeOperationUpdate, "TestDevice", "").Once().
			On("Plain", circuits0ID, model.MergeOperationRemove, 1, nil).Once().
			On("Managed", circuits0Sensors0, model.MergeOperationRemove,
				&device.Circuits[0].Sensors[0], nil,
				now, time.Time{}).Once()

		// Pass the device as base and empty device as client to simulate a Remove operation
		emptyDevice := Device{}
		_, err := merge.Merge(device, emptyDevice, merge.MergeOptions{
			Loggers: merge.MergeLoggers{mockLogger},
			Mode:    merge.ClientIsMaster,
		})

		assert.NoError(t, err)
		mockLogger.AssertExpectations(t)
	})

	// Test NotChanged operation
	t.Run("NotifyNotChanged", func(t *testing.T) {
		device1 := Device{
			Name: "TestDevice",
			Circuits: []Circuit{
				{ID: 1, Sensors: []Sensor{
					{ID: 10, TimeStamp: now},
				}},
			},
		}

		device2 := Device{
			Name: "TestDevice",
			Circuits: []Circuit{
				{ID: 1, Sensors: []Sensor{
					{ID: 10, TimeStamp: now},
				}},
			},
		}

		mockLogger := &MockLogger{}
		// When merging identical devices, we expect NotChanged operations for all fields
		mockLogger.
			On("Plain", "Name", model.MergeOperationNotChanged, "TestDevice", "TestDevice").Once().
			On("Plain", circuits0ID, model.MergeOperationNotChanged, 1, 1).Once().
			On("Managed", circuits0Sensors0, model.MergeOperationNotChanged,
				&device1.Circuits[0].Sensors[0], &device2.Circuits[0].Sensors[0],
				now, now).Once()

		// Pass identical devices to trigger NotChanged operations
		_, err := merge.Merge(device1, device2, merge.MergeOptions{
			Loggers: merge.MergeLoggers{mockLogger},
			Mode:    merge.ClientIsMaster,
		})

		assert.NoError(t, err)
		mockLogger.AssertExpectations(t)
	})

	// Test with nested maps
	t.Run("NotifyNestedMap", func(t *testing.T) {
		type NestedMapStruct struct {
			NestedMap map[string]map[string]int `json:"nested_map"`
		}

		nestedMapStruct1 := NestedMapStruct{
			NestedMap: map[string]map[string]int{
				"outer1": {
					"inner1": 10,
				},
			},
		}

		nestedMapStruct2 := NestedMapStruct{}

		mockLogger := &MockLogger{}
		// The implementation only notifies for the innermost value
		// (not for outer map or top-level map as might be expected)
		mockLogger.
			On("Plain", "nested_map.outer1.inner1", model.MergeOperationRemove, 10, mock.Anything).Once()

		// Merge from nested map to empty to simulate removal
		_, err := merge.MergeAny(nestedMapStruct1, nestedMapStruct2, merge.MergeOptions{
			Loggers: merge.MergeLoggers{mockLogger},
			Mode:    merge.ClientIsMaster,
		})

		assert.NoError(t, err)
		mockLogger.AssertExpectations(t)
	})

	// Test with zero loggers
	t.Run("NoLoggers", func(t *testing.T) {
		device := Device{
			Name: "TestDevice",
		}
		emptyDevice := Device{}

		// No loggers should not cause any errors
		_, err := merge.Merge(emptyDevice, device, merge.MergeOptions{
			Loggers: merge.MergeLoggers{},
			Mode:    merge.ClientIsMaster,
		})

		assert.NoError(t, err)
	})
}
