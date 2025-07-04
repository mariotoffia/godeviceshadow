package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/stretchr/testify/assert"
)

func TestDesiredWithDifferentTypes(t *testing.T) {
	// Different types for reportedModel and desiredModel
	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "value", Timestamp: time.Now()},
	}
	desired := "not a struct"

	_, err := merge.DesiredAny(reported, desired, merge.DesiredOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reported and desired model must be of the same kind")

	// Using the generic version which should also fail
	var s string = "not a struct"
	_, err = merge.DesiredAny(reported, s, merge.DesiredOptions{})
	assert.Error(t, err)
}

func TestDesiredWithNilValues(t *testing.T) {
	// Using nil pointers with DesiredAny should be handled gracefully
	var nilPtr *TestStruct
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "value", Timestamp: time.Now()},
	}

	// Test with nil pointer as reported
	result, err := merge.DesiredAny(nilPtr, &desired, merge.DesiredOptions{})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDesiredSliceHandling(t *testing.T) {
	type TestStruct struct {
		Slice []model.ValueAndTimestamp `json:"slice"`
	}

	reported := TestStruct{
		Slice: []model.ValueAndTimestamp{
			MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
	}
	desired := TestStruct{
		Slice: []model.ValueAndTimestamp{
			MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			MockValueAndTimestamp{Value: newValue, Timestamp: time.Now()}, // Extra element
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result.Slice))
	assert.Nil(t, result.Slice[0]) // Matched and removed
	assert.Equal(t, desiredValue, result.Slice[1].(MockValueAndTimestamp).Value)
	assert.Equal(t, newValue, result.Slice[2].(MockValueAndTimestamp).Value)
	assert.ElementsMatch(t, []string{"slice.0"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMakeAddressable(t *testing.T) {
	// This test indirectly tests the makeAddressable function by passing a non-addressable value
	type TestStruct struct {
		Field model.ValueAndTimestamp `json:"field"`
	}

	// Create a value that's not directly addressable by using a map value
	reportedMap := map[string]TestStruct{
		"key": {Field: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()}},
	}
	desiredMap := map[string]TestStruct{
		"key": {Field: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()}},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	// Use the map values directly, which are not addressable
	result, err := merge.DesiredAny(reportedMap["key"], desiredMap["key"], opts)

	assert.NoError(t, err)
	resultStruct, ok := result.(TestStruct)
	assert.True(t, ok)
	assert.Nil(t, resultStruct.Field) // Should be nil due to match
	assert.ElementsMatch(t, []string{"field"}, mockLogger.AcknowledgedPaths)
}
