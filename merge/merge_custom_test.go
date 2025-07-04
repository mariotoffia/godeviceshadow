package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCustomMerger tests the custom merger functionality
func TestCustomMerger(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldValue := &CustomMergeable{
		Name:      "old",
		Value:     10,
		Timestamp: oneHourAgo,
	}

	newValue := &CustomMergeable{
		Name:      "new",
		Value:     20,
		Timestamp: now,
	}

	// Test with ClientIsMaster
	merged, err := merge.Merge(oldValue, newValue, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)
	// No need for type assertion, Merge returns the same type as input
	assert.Equal(t, "old+new", merged.Name, "Name should be concatenated")
	assert.Equal(t, 20, merged.Value, "Value should be from client with ClientIsMaster")

	// Test with ServerIsMaster
	merged, err = merge.Merge(oldValue, newValue, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})
	require.NoError(t, err)

	// No need for type assertion, Merge returns the same type as input
	assert.Equal(t, "old+new", merged.Name, "Name should be concatenated")
	assert.Equal(t, 20, merged.Value, "Value should be from client because timestamp is newer")
}

// TestCustomMergerError tests error handling for custom mergers
func TestCustomMergerError(t *testing.T) {
	errorMergeable := &ErrorMergeable{ShouldError: true}
	otherMergeable := &ErrorMergeable{ShouldError: false}

	_, err := merge.Merge(errorMergeable, otherMergeable, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})

	require.Error(t, err, "Should return error from custom merger")
	assert.Contains(t, err.Error(), "simulated error from custom merger")
}

// TestCustomMergerWithDifferentTypes tests error handling for custom mergers with incompatible types
func TestCustomMergerWithDifferentTypes(t *testing.T) {
	customMerger := &CustomMergeable{
		Name:      "custom",
		Value:     10,
		Timestamp: time.Now(),
	}

	// Use MergeAny instead of Merge to allow different types
	_, err := merge.MergeAny(customMerger, "wrong type", merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be of the same type")
}

// TestNestedStructWithCustomMerger tests custom merger within a nested structure
func TestNestedStructWithCustomMerger(t *testing.T) {
	type Container struct {
		Name  string
		Merge *CustomMergeable
	}

	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldContainer := Container{
		Name: "old-container",
		Merge: &CustomMergeable{
			Name:      "old",
			Value:     10,
			Timestamp: oneHourAgo,
		},
	}

	newContainer := Container{
		Name: "new-container",
		Merge: &CustomMergeable{
			Name:      "new",
			Value:     20,
			Timestamp: now,
		},
	}

	merged, err := merge.Merge(oldContainer, newContainer, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)

	// No need for type assertion since Merge returns the same type
	result := merged
	assert.Equal(t, "new-container", result.Name)
	require.NotNil(t, result.Merge)
	assert.Equal(t, "new", result.Merge.Name, "The nested CustomMergeable doesn't receive special merge handling automatically")
	assert.Equal(t, 20, result.Merge.Value)
}
