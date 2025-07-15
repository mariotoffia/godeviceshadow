package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMergeSliceByID tests merging of slices using ID-based merging
func TestMergeSliceByID(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

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
			{ID: "temp", TimeStamp: now, Value: 23.5},            // newer timestamp
			{ID: "pressure", TimeStamp: now, Value: 1013.25},     // new sensor
			{ID: "humidity", TimeStamp: oneHourAgo, Value: 46.0}, // same timestamp
		},
	}

	// Test with ID-based slice merging enabled
	merged, err := merge.Merge(context.Background(), oldContainer, newContainer, merge.MergeOptions{
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
	assert.Equal(t, oneHourAgo, humiditySensor.TimeStamp)
	assert.Equal(t, 45.0, humiditySensor.Value, "Original value kept when timestamps are equal")
	assert.Equal(t, now, pressureSensor.TimeStamp)
	assert.Equal(t, 1013.25, pressureSensor.Value)
}
