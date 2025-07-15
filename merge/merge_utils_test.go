package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMapKeyFormatting tests the internal formatKey function through maps with various key types
func TestMapKeyFormatting(t *testing.T) {
	// Test with string keys
	strMap := map[string]int{"a": 1, "b": 2}
	merged, err := merge.Merge(context.Background(), strMap, strMap, merge.MergeOptions{})
	require.NoError(t, err)
	require.Len(t, merged, 2)

	// Test with int keys
	intMap := map[int]int{1: 10, 2: 20}
	merged2, err := merge.Merge(context.Background(), intMap, intMap, merge.MergeOptions{})
	require.NoError(t, err)
	require.Len(t, merged2, 2)

	// Test with struct keys
	type StructKey struct {
		ID   int
		Name string
	}
	structMap := map[StructKey]int{
		{ID: 1, Name: "first"}:  100,
		{ID: 2, Name: "second"}: 200,
	}
	merged3, err := merge.Merge(context.Background(), structMap, structMap, merge.MergeOptions{})
	require.NoError(t, err)
	require.Len(t, merged3, 2)

	// Test with pointer keys
	key1 := &StructKey{ID: 1, Name: "first"}
	key2 := &StructKey{ID: 2, Name: "second"}
	ptrMap := map[*StructKey]int{
		key1: 1000,
		key2: 2000,
	}
	merged4, err := merge.Merge(context.Background(), ptrMap, ptrMap, merge.MergeOptions{})
	require.NoError(t, err)
	require.Len(t, merged4, 2)

	// Test with interface keys containing various types
	interfaceMap := map[interface{}]int{
		"string": 1,
		42:       2,
		key1:     3,
	}
	merged5, err := merge.Merge(context.Background(), interfaceMap, interfaceMap, merge.MergeOptions{})
	require.NoError(t, err)
	require.Len(t, merged5, 3)
}

// TestIdInterfaceUnwrapping tests the ID interface unwrapping functionality
func TestIdInterfaceUnwrapping(t *testing.T) {
	now := time.Now().UTC()

	// Create slices with different types implementing IdValueAndTimestamp
	// Some directly, some through pointers
	slice1 := []interface{}{
		&IdSensor{ID: "direct-ptr", TimeStamp: now, Value: 22.5},
		IdSensor{ID: "value-type", TimeStamp: now, Value: 18.0}, // This won't implement the interface
		&MockIdValueType{ID: "ptr-required", TimeStamp: now, Value: "test"},
	}

	slice2 := []interface{}{
		&IdSensor{ID: "direct-ptr2", TimeStamp: now, Value: 23.5},
		&MockIdValueType{ID: "ptr-required2", TimeStamp: now, Value: "test2"},
	}

	// Merge should still work even though some types require pointers
	_, err := merge.Merge(context.Background(), slice1, slice2, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err, "Should handle different implementation patterns")

	// Test unwrapping through the mergeSliceByID functionality
	// Create slices with all elements properly implementing IdValueAndTimestamp
	idSensorSlice1 := []*IdSensor{
		{ID: "sensor1", TimeStamp: now, Value: 22.5},
		{ID: "sensor2", TimeStamp: now, Value: 18.0},
	}

	idSensorSlice2 := []*IdSensor{
		{ID: "sensor2", TimeStamp: now, Value: 19.0}, // Updated value
		{ID: "sensor3", TimeStamp: now, Value: 25.0}, // New sensor
	}

	merged, err := merge.Merge(context.Background(), idSensorSlice1, idSensorSlice2, merge.MergeOptions{
		Mode:            merge.ClientIsMaster,
		MergeSlicesByID: true,
	})
	require.NoError(t, err)

	// Should have merged correctly by ID
	require.Len(t, merged, 2)

	// Map results by ID
	sensorMap := make(map[string]*IdSensor)
	for _, s := range merged {
		sensorMap[s.ID] = s
	}

	require.Contains(t, sensorMap, "sensor2")
	require.Contains(t, sensorMap, "sensor3")
	require.NotContains(t, sensorMap, "sensor1") // Removed due to ClientIsMaster

	assert.Equal(t, 18.0, sensorMap["sensor2"].Value, "Original value kept when timestamps are equal")
}

// TestCustomMergerErrorPropagation tests error handling for custom merger implementations
func TestCustomMergerErrorPropagation(t *testing.T) {
	// Test a custom merger that returns an error
	errorMergeable := &ErrorMergeable{ShouldError: true}
	_, err := merge.Merge(context.Background(), errorMergeable, errorMergeable, merge.MergeOptions{})
	require.Error(t, err, "Should propagate error from custom merger")
	assert.Contains(t, err.Error(), "simulated error", "Error message should be preserved")

	// Test with mismatched types for custom merger
	customMergeable := &CustomMergeable{}
	// Use MergeAny instead of Merge to allow different types
	_, err = merge.MergeAny(context.Background(), customMergeable, 42, merge.MergeOptions{})
	require.Error(t, err, "Should error with mismatched types")
}
