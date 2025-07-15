package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"

	"github.com/stretchr/testify/assert"
)

const (
	noMatch             = "no-match"
	desiredValue        = "desired-value"
	desiredValue2       = "desired-value-2"
	newValue            = "new-value"
	shouldNotBeAccessed = "should-not-be-accessed"
)

type TestStruct struct {
	Field1 model.ValueAndTimestamp `json:"field1"`
	Field2 model.ValueAndTimestamp `json:"field2"`
}

type MockValueAndTimestamp struct {
	Value     any
	Timestamp time.Time
}

func (m MockValueAndTimestamp) GetTimestamp() time.Time {
	return m.Timestamp
}

func (m MockValueAndTimestamp) GetValue() any {
	return m.Value
}

func TestDesiredSimpleStruct(t *testing.T) {
	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
	}

	mockLogger := &MockLogger{} // Using the pre-existing MockLogger
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field2.GetTimestamp()},
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredNestedStruct(t *testing.T) {
	type NestedStruct struct {
		SubField1 model.ValueAndTimestamp `json:"sub_field1"`
		SubField2 model.ValueAndTimestamp `json:"sub_field2"`
	}

	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 NestedStruct            `json:"field2"`
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "nested-match", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: "nested-no-match", Timestamp: time.Now()},
		},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "nested-match", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: "desired-nested-value", Timestamp: time.Now()},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: NestedStruct{
			SubField1: nil, // Removed due to match
			SubField2: MockValueAndTimestamp{Value: "desired-nested-value", Timestamp: desired.Field2.SubField2.GetTimestamp()},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1", "field2.sub_field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMapMerge(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			"key2": MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			"key2": MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: time.Now()},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field1["key2"].GetTimestamp()},
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1["key3"].GetTimestamp()},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.key1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMixedDataTypes(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 string                  `json:"field2"`
		Field3 int                     `json:"field3"`
		Field4 model.ValueAndTimestamp `json:"field4"`
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: "constant",
		Field3: 42,
		Field4: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: "constant",
		Field3: 42,
		Field4: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,        // Removed due to match
		Field2: "constant", // Unchanged
		Field3: 42,         // Unchanged
		Field4: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field4.GetTimestamp()},
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredNestedMaps(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]map[string]model.ValueAndTimestamp `json:"field1"`
	}

	reported := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
				"inner2": MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
			},
		},
	}
	desired := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
				"inner2": MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: desiredValue2, Timestamp: time.Now()},
			},
			"outer3": {
				"inner4": MockValueAndTimestamp{Value: newValue, Timestamp: time.Now()},
			},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner2": MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field1["outer1"]["inner2"].GetTimestamp()},
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: desiredValue2, Timestamp: desired.Field1["outer2"]["inner3"].GetTimestamp()},
			},
			"outer3": {
				"inner4": MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1["outer3"]["inner4"].GetTimestamp()},
			},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.outer1.inner1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredEmptyAndNullFields(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		Field3 model.ValueAndTimestamp `json:"field3"`
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		Field3: nil, // Field3 is reported as nil
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
		Field3: nil, // Field3 is desired as nil
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field2.GetTimestamp()},
		Field3: nil, // Field3 remains nil
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredUnexportedFields(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		field3 model.ValueAndTimestamp // Unexported field (not accessible)
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: time.Now()},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: time.Now()},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field2.GetTimestamp()},
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: desired.field3.GetTimestamp()}, // Should remain unchanged
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMultipleMatchesInNestedStructs(t *testing.T) {

	type NestedStruct struct {
		SubField1 model.ValueAndTimestamp `json:"sub_field1"`
		SubField2 model.ValueAndTimestamp `json:"sub_field2"`
	}

	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 NestedStruct            `json:"field2"`
		Field3 NestedStruct            `json:"field3"`
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match1", Timestamp: time.Now()},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match2", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
		Field3: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match3", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
	}

	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match1", Timestamp: time.Now()},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match2", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
		},
		Field3: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match3", Timestamp: time.Now()},
			SubField2: MockValueAndTimestamp{Value: desiredValue2, Timestamp: time.Now()},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: NestedStruct{
			SubField1: nil, // Removed due to match
			SubField2: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field2.SubField2.GetTimestamp()},
		},
		Field3: NestedStruct{
			SubField1: nil, // Removed due to match
			SubField2: MockValueAndTimestamp{Value: desiredValue2, Timestamp: desired.Field3.SubField2.GetTimestamp()},
		},
	}, result)
	assert.ElementsMatch(t, []string{
		"field1",
		"field2.sub_field1",
		"field3.sub_field1",
	}, mockLogger.AcknowledgedPaths)
}

func TestDesiredEmptyMapsAndStructs(t *testing.T) {
	type NestedStruct struct {
		SubField1 model.ValueAndTimestamp `json:"sub_field1"`
	}

	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
		Field2 NestedStruct                       `json:"field2"`
	}

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Empty map
		Field2: NestedStruct{},                       // Empty struct
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Empty map (same)
		Field2: NestedStruct{},                       // Empty struct (same)
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Unchanged
		Field2: NestedStruct{},                       // Unchanged
	}, result)
	assert.Empty(t, mockLogger.AcknowledgedPaths) // No paths should be logged
}

func TestDesiredMissingFieldsInDesired(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		Field3 model.ValueAndTimestamp `json:"field3"` // Only in reported
	}

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		Field3: MockValueAndTimestamp{Value: "reported-only", Timestamp: time.Now()},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
		// Field3 is missing in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field2.GetTimestamp()},
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMultiLevelNestedStructs(t *testing.T) {
	type Level3 struct {
		FieldA model.ValueAndTimestamp `json:"field_a"`
		FieldB model.ValueAndTimestamp `json:"field_b"`
	}
	type Level2 struct {
		Level3 Level3 `json:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `json:"level2"`
	}

	reported := Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
				FieldB: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
			},
		},
	}
	desired := Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
				FieldB: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: nil, // Removed due to match
				FieldB: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Level2.Level3.FieldB.GetTimestamp()},
			},
		},
	}, result)
	assert.ElementsMatch(t, []string{"level2.level3.field_a"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMapWithComplexKeys(t *testing.T) {
	type ComplexKey struct {
		ID    string
		Group string
	}
	type TestStruct struct {
		Field1 map[ComplexKey]model.ValueAndTimestamp `json:"field1"`
	}

	reported := TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
	}
	desired := TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			{ID: "3", Group: "C"}: MockValueAndTimestamp{Value: newValue, Timestamp: time.Now()},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field1[ComplexKey{ID: "2", Group: "B"}].GetTimestamp()},
			{ID: "3", Group: "C"}: MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1[ComplexKey{ID: "3", Group: "C"}].GetTimestamp()},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.{ID:1,Group:A}"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMixedKeyTypesInMap(t *testing.T) {
	type ComplexKey struct {
		ID    string
		Group string
	}
	type MixedKeyMap struct {
		Field1 map[interface{}]model.ValueAndTimestamp `json:"field1"`
	}

	reported := MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			"key1":                          MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			42:                              MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()},
		},
	}
	desired := MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			"key1":                          MockValueAndTimestamp{Value: "match", Timestamp: time.Now()},
			42:                              MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()},
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: newValue, Timestamp: time.Now()},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			42:                              MockValueAndTimestamp{Value: desiredValue, Timestamp: desired.Field1[42].GetTimestamp()},
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1[ComplexKey{ID: "1", Group: "A"}].GetTimestamp()},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.key1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredStructWithPointerFields(t *testing.T) {
	// Arrange
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
	}

	field1Reported := MockValueAndTimestamp{Value: "match", Timestamp: time.Now()}
	field2Reported := MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()}

	field1Desired := MockValueAndTimestamp{Value: "match", Timestamp: time.Now()}
	field2Desired := MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()}

	reported := TestStruct{
		Field1: &field1Reported,
		Field2: &field2Reported,
	}
	desired := TestStruct{
		Field1: &field1Desired,
		Field2: &field2Desired,
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: &field2Desired,
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredNilPointersInStructs(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		Field3 model.ValueAndTimestamp `json:"field3"`
	}

	field2Reported := MockValueAndTimestamp{Value: noMatch, Timestamp: time.Now()}

	field2Desired := MockValueAndTimestamp{Value: desiredValue, Timestamp: time.Now()}

	reported := TestStruct{
		Field1: nil, // Field1 is nil in reported
		Field2: &field2Reported,
		Field3: nil, // Field3 is nil in reported
	}
	desired := TestStruct{
		Field1: nil, // Field1 is nil in desired
		Field2: &field2Desired,
		Field3: nil, // Field3 is nil in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{Loggers: merge.DesiredLoggers{mockLogger}}

	result, err := merge.Desired(context.Background(), reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Remains nil
		Field2: &field2Desired,
		Field3: nil, // Remains nil
	}, result)
	assert.Empty(t, mockLogger.AcknowledgedPaths) // No paths should be logged
}
