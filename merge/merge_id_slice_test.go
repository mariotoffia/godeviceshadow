package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSliceByIDMerging tests the ID-based slice merging functionality
func TestSliceByIDMerging(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	// Create container with ID sensors
	oldContainer := SensorContainer{
		Name: "old-container",
		Sensors: []IdSensor{
			{ID: "temp", TimeStamp: oneHourAgo, Value: 22.5},
			{ID: "humidity", TimeStamp: oneHourAgo, Value: 45.0},
		},
	}

	newContainer := SensorContainer{
		Name: "new-container",
		Sensors: []IdSensor{
			{ID: "temp", TimeStamp: now, Value: 23.5},             // newer timestamp
			{ID: "pressure", TimeStamp: now, Value: 1013.25},      // new sensor
			{ID: "humidity", TimeStamp: twoHoursAgo, Value: 46.0}, // older timestamp
		},
	}

	// Test with ID-based slice merging enabled
	merged, err := merge.Merge(oldContainer, newContainer, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err)

	// Check results
	assert.Equal(t, "new-container", merged.Name)
	require.Len(t, merged.Sensors, 3, "Should have 3 sensors after merge")

	// The order is not guaranteed, so find by ID
	var tempSensor, humiditySensor, pressureSensor *IdSensor
	for i := range merged.Sensors {
		switch merged.Sensors[i].ID {
		case "temp":
			tempSensor = &merged.Sensors[i]
		case "humidity":
			humiditySensor = &merged.Sensors[i]
		case "pressure":
			pressureSensor = &merged.Sensors[i]
		}
	}

	require.NotNil(t, tempSensor, "Temperature sensor should be present")
	require.NotNil(t, humiditySensor, "Humidity sensor should be present")
	require.NotNil(t, pressureSensor, "Pressure sensor should be present")

	assert.Equal(t, now, tempSensor.TimeStamp, "Temp sensor should have newer timestamp")
	assert.Equal(t, 23.5, tempSensor.Value, "Temp sensor should have new value")
	assert.Equal(t, oneHourAgo, humiditySensor.TimeStamp, "Humidity should keep old timestamp which is newer")
	assert.Equal(t, 45.0, humiditySensor.Value, "Humidity should keep old value because its timestamp is newer")
	assert.Equal(t, now, pressureSensor.TimeStamp)
	assert.Equal(t, 1013.25, pressureSensor.Value)
}

// TestNestedSliceByIDMerging tests ID-based merging in nested structures
func TestNestedSliceByIDMerging(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	type DeviceWithIDSensors struct {
		Name    string
		Sensors []IdSensor
		Nested  struct {
			InnerSensors []IdSensor
		}
	}

	oldDevice := DeviceWithIDSensors{
		Name: "device-old",
		Sensors: []IdSensor{
			{ID: "temp", TimeStamp: oneHourAgo, Value: 22.5},
		},
		Nested: struct {
			InnerSensors []IdSensor
		}{
			InnerSensors: []IdSensor{
				{ID: "inner-humid", TimeStamp: oneHourAgo, Value: 45.0},
			},
		},
	}

	newDevice := DeviceWithIDSensors{
		Name: "device-new",
		Sensors: []IdSensor{
			{ID: "temp", TimeStamp: now, Value: 23.5},        // newer
			{ID: "pressure", TimeStamp: now, Value: 1013.25}, // new
		},
		Nested: struct {
			InnerSensors []IdSensor
		}{
			InnerSensors: []IdSensor{
				{ID: "inner-temp", TimeStamp: now, Value: 23.0}, // new
			},
		},
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err)

	// Check top level
	assert.Equal(t, "device-new", merged.Name)
	require.Len(t, merged.Sensors, 2)

	// Check nested
	require.Len(t, merged.Nested.InnerSensors, 1, "With ClientIsMaster, we should only have elements from client")

	// Find by ID in nested
	var innerTemp *IdSensor
	for i := range merged.Nested.InnerSensors {
		if merged.Nested.InnerSensors[i].ID == "inner-temp" {
			innerTemp = &merged.Nested.InnerSensors[i]
			break
		}
	}
	require.NotNil(t, innerTemp, "inner-temp sensor should be present")
	assert.Equal(t, now, innerTemp.TimeStamp)
	assert.Equal(t, 23.0, innerTemp.Value)
}

// TestEmptySliceByIDMerging tests merging with empty slices
func TestEmptySliceByIDMerging(t *testing.T) {
	now := time.Now().UTC()

	// Old has values, new is empty
	oldContainer := SensorContainer{
		Name: "old-container",
		Sensors: []IdSensor{
			{ID: "temp", TimeStamp: now, Value: 22.5},
		},
	}

	emptyContainer := SensorContainer{
		Name:    "empty-container",
		Sensors: []IdSensor{}, // Explicitly empty
	}

	// With ClientIsMaster, the result should be empty
	merged, err := merge.Merge(oldContainer, emptyContainer, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "empty-container", merged.Name)
	assert.Empty(t, merged.Sensors, "With ClientIsMaster, should use empty slice from client")

	// With ServerIsMaster, should keep the server values
	merged, err = merge.Merge(oldContainer, emptyContainer, merge.MergeOptions{
		Mode:            merge.ServerIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "empty-container", merged.Name) // Name still gets overridden
	require.Len(t, merged.Sensors, 1, "With ServerIsMaster, should keep server values")
	assert.Equal(t, "temp", merged.Sensors[0].ID)
}
