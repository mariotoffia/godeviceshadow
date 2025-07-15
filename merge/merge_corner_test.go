package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMixedImplementationSlice tests behavior when slices have mixed implementation of interfaces
func TestMixedImplementationSlice(t *testing.T) {
	now := time.Now().UTC()

	// Create a slice with mixed implementation - some implement IdValueAndTimestamp, some don't
	mixedSlice1 := []interface{}{
		&IdSensor{ID: "temp1", TimeStamp: now, Value: 22.5},
		42,       // doesn't implement IdValueAndTimestamp
		"string", // doesn't implement IdValueAndTimestamp
	}

	mixedSlice2 := []interface{}{
		&IdSensor{ID: "temp2", TimeStamp: now, Value: 23.5},
		true, // doesn't implement IdValueAndTimestamp
	}

	// Should fall back to position-based merging
	merged, err := merge.Merge(context.Background(), mixedSlice1, mixedSlice2, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true, // Even though this is true, it should use position-based due to mixed types
	})

	require.NoError(t, err)
	require.Len(t, merged, 2, "With ClientIsMaster and position-based merging, should have client length")

	// First element should be the IdSensor from slice2
	sensor, ok := merged[0].(*IdSensor)
	require.True(t, ok, "First element should be *IdSensor")
	assert.Equal(t, "temp1", sensor.ID, "With position-based merging, the ID doesn't affect the merge")

	// Second element should be the boolean from slice2
	boolean, ok := merged[1].(bool)
	require.True(t, ok, "Second element should be bool")
	assert.True(t, boolean)
}

// TestNestedIDBasedSliceMerging tests merging of ID-based slices in nested structures
func TestNestedIDBasedSliceMerging(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Create nested structures with slices that should be merged by ID
	oldContainer := SensorContainer{
		Name: "OldContainer",
		Sensors: []IdSensor{
			{ID: "temp1", TimeStamp: oneHourAgo, Value: 22.5},
			{ID: "temp2", TimeStamp: oneHourAgo, Value: 18.0},
		},
	}

	newContainer := SensorContainer{
		Name: "NewContainer",
		Sensors: []IdSensor{
			{ID: "temp2", TimeStamp: now, Value: 19.0}, // Updated
			{ID: "temp3", TimeStamp: now, Value: 25.0}, // New
			// temp1 is missing
		},
	}

	// Test with ClientIsMaster mode and ID-based merging
	merged, err := merge.Merge(context.Background(), oldContainer, newContainer, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "NewContainer", merged.Name, "Container name should be updated")
	require.Len(t, merged.Sensors, 2, "Should have 2 sensors (temp1 removed, temp3 added)")

	// Map by ID for easier testing
	sensorMap := make(map[string]IdSensor)
	for _, s := range merged.Sensors {
		sensorMap[s.ID] = s
	}

	require.Contains(t, sensorMap, "temp2", "temp2 should be present")
	require.Contains(t, sensorMap, "temp3", "temp3 should be present")
	require.NotContains(t, sensorMap, "temp1", "temp1 should be removed with ClientIsMaster")

	assert.Equal(t, 19.0, sensorMap["temp2"].Value, "temp2 value should be updated")
	assert.Equal(t, now, sensorMap["temp2"].TimeStamp, "temp2 timestamp should be updated")

	// Test with ServerIsMaster mode
	mergedServer, err := merge.Merge(context.Background(), oldContainer, newContainer, merge.MergeOptions{
		Mode:            merge.ServerIsMaster,
		MergeSlicesByID: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "NewContainer", mergedServer.Name, "Container name should be updated")
	require.Len(t, mergedServer.Sensors, 3, "Should have 3 sensors (temp1 kept with ServerIsMaster)")

	// Map by ID for easier testing
	sensorMapServer := make(map[string]IdSensor)
	for _, s := range mergedServer.Sensors {
		sensorMapServer[s.ID] = s
	}

	require.Contains(t, sensorMapServer, "temp1", "temp1 should be present with ServerIsMaster")
	require.Contains(t, sensorMapServer, "temp2", "temp2 should be present")
	require.Contains(t, sensorMapServer, "temp3", "temp3 should be present")

	assert.Equal(t, 22.5, sensorMapServer["temp1"].Value, "temp1 value should be kept")
}
