package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"

	"github.com/stretchr/testify/assert"
)

const (
	updatedValue        = "updated-value"
	desiredValue        = "desired-value"
	desiredValue2       = "desired-value-2"
	newValue            = "new-value"
	field1key1          = "field1.key1"
	field1subField1     = "field2.sub_field1"
	newField            = "new-field"
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
	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
	}

	mockLogger := &MockLogger{} // Using the pre-existing MockLogger
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false,
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                   // Removed due to match
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()}, // since OnlyNewerTimeStamps is false
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

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "nested-match", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "nested-match", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: NestedStruct{
			SubField1: nil,                                                                                             // Removed due to match
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.SubField2.GetTimestamp()}, // Updated from reported
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1", field1subField1}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMapMerge(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now},
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now},
			"key2": MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: now},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1["key2"].GetTimestamp()}, // Updated from reported
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1["key3"].GetTimestamp()},      // Unchanged
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1}, mockLogger.AcknowledgedPaths)
}

func TestDesiredMixedDataTypes(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 string                  `json:"field2"`
		Field3 int                     `json:"field3"`
		Field4 model.ValueAndTimestamp `json:"field4"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: "constant",
		Field3: 42,
		Field4: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: "constant",
		Field3: 42,
		Field4: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                   // Removed due to match
		Field2: "constant",                                                                            // Unchanged
		Field3: 42,                                                                                    // Unchanged
		Field4: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field4.GetTimestamp()}, // Updated from reported
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredNestedMaps(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: now},
				"inner2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
			},
		},
	}
	desired := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: now},
				"inner2": MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: desiredValue2, Timestamp: now},
			},
			"outer3": {
				"inner4": MockValueAndTimestamp{Value: newValue, Timestamp: now},
			},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner2": MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1["outer1"]["inner2"].GetTimestamp()}, // Updated
			},
			"outer2": {
				"inner3": MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1["outer2"]["inner3"].GetTimestamp()}, // Updated
			},
			"outer3": {
				"inner4": MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1["outer3"]["inner4"].GetTimestamp()}, // Unchanged
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

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		Field3: nil, // Field3 is nil in reported
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
		Field3: nil, // Field3 is nil in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                   // Removed due to match
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()}, // Updated from reported
		Field3: nil,                                                                                   // Remains nil
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredUnexportedFields(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		field3 model.ValueAndTimestamp // Unexported field (not accessible)
	}

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: now},
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: now},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                         // Removed due to match
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()},       // Updated from reported
		field3: MockValueAndTimestamp{Value: shouldNotBeAccessed, Timestamp: desired.field3.GetTimestamp()}, // Unchanged
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

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match1", Timestamp: now},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match2", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
		Field3: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match3", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
	}

	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match1", Timestamp: now},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match2", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
		},
		Field3: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match3", Timestamp: now},
			SubField2: MockValueAndTimestamp{Value: desiredValue2, Timestamp: now},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil, // Removed due to match
		Field2: NestedStruct{
			SubField1: nil,                                                                                             // Removed due to match
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.SubField2.GetTimestamp()}, // Updated
		},
		Field3: NestedStruct{
			SubField1: nil,                                                                                             // Removed due to match
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field3.SubField2.GetTimestamp()}, // Updated
		},
	}, result)
	assert.ElementsMatch(t, []string{
		"field1",
		field1subField1,
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
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Remains unchanged
		Field2: NestedStruct{},                       // Remains unchanged
	}, result)
	assert.Empty(t, mockLogger.AcknowledgedPaths) // No paths should be logged
}

func TestDesiredMissingFieldsInDesired(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		Field3 model.ValueAndTimestamp `json:"field3"` // Only in reported
	}

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		Field3: MockValueAndTimestamp{Value: "reported-only", Timestamp: now}, // Extra field
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
		// Field3 is missing in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                      // Removed due to match
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()},    // Updated from reported
		Field3: MockValueAndTimestamp{Value: "reported-only", Timestamp: reported.Field3.GetTimestamp()}, // From reported
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

	now := time.Now()

	reported := Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: MockValueAndTimestamp{Value: "match", Timestamp: now},
				FieldB: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
			},
		},
	}
	desired := Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: MockValueAndTimestamp{Value: "match", Timestamp: now},
				FieldB: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
			},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, Level1{
		Level2: Level2{
			Level3: Level3{
				FieldA: nil,                                                                                                 // Removed due to match
				FieldB: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Level2.Level3.FieldB.GetTimestamp()}, // Updated
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

	now := time.Now()

	reported := TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: "match", Timestamp: now},
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
	}
	desired := TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: "match", Timestamp: now},
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
			{ID: "3", Group: "C"}: MockValueAndTimestamp{Value: newValue, Timestamp: now},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[ComplexKey]model.ValueAndTimestamp{
			{ID: "2", Group: "B"}: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1[ComplexKey{ID: "2", Group: "B"}].GetTimestamp()}, // Updated
			{ID: "3", Group: "C"}: MockValueAndTimestamp{Value: newValue, Timestamp: desired.Field1[ComplexKey{ID: "3", Group: "C"}].GetTimestamp()},      // Unchanged
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

	now := time.Now()

	reported := MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			"key1":                          MockValueAndTimestamp{Value: "match", Timestamp: now},
			42:                              MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: updatedValue, Timestamp: now},
		},
	}
	desired := MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			"key1":                          MockValueAndTimestamp{Value: "match", Timestamp: now},
			42:                              MockValueAndTimestamp{Value: desiredValue, Timestamp: now},
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: newValue, Timestamp: now},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, MixedKeyMap{
		Field1: map[interface{}]model.ValueAndTimestamp{
			42:                              MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1[42].GetTimestamp()},                              // Updated
			ComplexKey{ID: "1", Group: "A"}: MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field1[ComplexKey{ID: "1", Group: "A"}].GetTimestamp()}, // Updated
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1}, mockLogger.AcknowledgedPaths)
}

func TestDesiredStructWithPointerFields(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
	}

	now := time.Now()

	field1Reported := MockValueAndTimestamp{Value: "match", Timestamp: now}
	field2Reported := MockValueAndTimestamp{Value: updatedValue, Timestamp: now}

	field1Desired := MockValueAndTimestamp{Value: "match", Timestamp: now}
	field2Desired := MockValueAndTimestamp{Value: desiredValue, Timestamp: now}

	reported := TestStruct{
		Field1: &field1Reported,
		Field2: &field2Reported,
	}
	desired := TestStruct{
		Field1: &field1Desired,
		Field2: &field2Desired,
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                    // Removed due to match
		Field2: &MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()}, // Updated from reported
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
}

func TestDesiredNilPointersInStructs(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"`
		Field3 model.ValueAndTimestamp `json:"field3"`
	}

	now := time.Now()

	field2Reported := MockValueAndTimestamp{Value: updatedValue, Timestamp: now}
	field2Desired := MockValueAndTimestamp{Value: desiredValue, Timestamp: now}

	reported := TestStruct{
		Field1: nil,             // Field1 is nil in reported
		Field2: &field2Reported, // Field2 has a value in reported
		Field3: nil,             // Field3 is nil in reported
	}
	desired := TestStruct{
		Field1: nil,            // Field1 is nil in desired
		Field2: &field2Desired, // Field2 has a value in desired
		Field3: nil,            // Field3 is nil in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                                                    // Remains nil
		Field2: &MockValueAndTimestamp{Value: updatedValue, Timestamp: reported.Field2.GetTimestamp()}, // Updated from reported
		Field3: nil,                                                                                    // Remains nil
	}, result)
	assert.Empty(t, mockLogger.AcknowledgedPaths) // No paths should be logged
}

func TestDesiredAddNewField(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"` // Field2 exists only in reported
	}

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "existing", Timestamp: now},
		Field2: MockValueAndTimestamp{Value: newField, Timestamp: now}, // New field to add
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "existing", Timestamp: now}, // Unchanged
		// Field2 is missing
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                    // Acknowledged (unchanged)
		Field2: MockValueAndTimestamp{Value: newField, Timestamp: now}, // Added from reported
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field2"}, mockLogger.AddedPaths)
}

func TestDesiredAddNewKeyToMap(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now},  // Unchanged
			"key2": MockValueAndTimestamp{Value: newValue, Timestamp: now}, // New key to add
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now}, // Existing and unchanged
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: newValue, Timestamp: now}, // Added from reported
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field1.key2"}, mockLogger.AddedPaths)
}

func TestDesiredUpdateField(t *testing.T) {
	type TestStruct struct {
		Field1 model.ValueAndTimestamp `json:"field1"`
		Field2 model.ValueAndTimestamp `json:"field2"` // Field2 will be updated
	}

	now := time.Now()

	reported := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},      // Unchanged
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
	}
	desired := TestStruct{
		Field1: MockValueAndTimestamp{Value: "match", Timestamp: now},     // Existing and unchanged
		Field2: MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Existing but outdated
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: nil,                                                        // Acknowledged (unchanged)
		Field2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated from reported
	}, result)
	assert.ElementsMatch(t, []string{"field1"}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field2"}, mockLogger.UpdatePaths)
}

func TestDesiredUpdateValueInMap(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now},      // Unchanged
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "match", Timestamp: now},     // Existing and unchanged
			"key2": MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Existing but outdated
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated from reported
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field1.key2"}, mockLogger.UpdatePaths)
}

func TestDesiredAddAndUpdateInNestedStructs(t *testing.T) {
	type NestedStruct struct {
		SubField1 model.ValueAndTimestamp `json:"sub_field1"`
		SubField2 model.ValueAndTimestamp `json:"sub_field2"`
		SubField3 model.ValueAndTimestamp `json:"sub_field3"` // Only in reported
	}
	type TestStruct struct {
		Field1 NestedStruct `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match", Timestamp: now},      // Unchanged
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
			SubField3: MockValueAndTimestamp{Value: newField, Timestamp: now},     // New field
		},
	}
	desired := TestStruct{
		Field1: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match", Timestamp: now},     // Existing and unchanged
			SubField2: MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Existing but outdated
			// SubField3 is missing
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: NestedStruct{
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated
			SubField3: MockValueAndTimestamp{Value: newField, Timestamp: now},     // Added
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.sub_field1"}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field1.sub_field2"}, mockLogger.UpdatePaths)
	assert.ElementsMatch(t, []string{"field1.sub_field3"}, mockLogger.AddedPaths)
}

func TestDesiredAddAndUpdateInNestedMaps(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: now},      // Unchanged
				"inner2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
				"inner3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // New key
			},
		},
	}
	desired := TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner1": MockValueAndTimestamp{Value: "match", Timestamp: now},     // Existing and unchanged
				"inner2": MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Existing but outdated
				// "inner3" is missing
			},
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]map[string]model.ValueAndTimestamp{
			"outer1": {
				"inner2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated
				"inner3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // Added
			},
		},
	}, result)
	assert.ElementsMatch(t, []string{"field1.outer1.inner1"}, mockLogger.AcknowledgedPaths)
	assert.ElementsMatch(t, []string{"field1.outer1.inner2"}, mockLogger.UpdatePaths)
	assert.ElementsMatch(t, []string{"field1.outer1.inner3"}, mockLogger.AddedPaths)
}

func TestDesiredHandleEmptyMapInDesired(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now}, // New key
			"key2": MockValueAndTimestamp{Value: "value2", Timestamp: now}, // New key
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Empty map in desired
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now}, // Added
			"key2": MockValueAndTimestamp{Value: "value2", Timestamp: now}, // Added
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1, "field1.key2"}, mockLogger.AddedPaths)
}

func TestDesiredHandleEmptyMapInReported(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{}, // Empty map in reported
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now}, // Should remain unchanged
			"key2": MockValueAndTimestamp{Value: "value2", Timestamp: now}, // Should remain unchanged
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now}, // Unchanged
			"key2": MockValueAndTimestamp{Value: "value2", Timestamp: now}, // Unchanged
		},
	}, result)
	assert.Empty(t, mockLogger.AcknowledgedPaths) // No acknowledgments
	assert.Empty(t, mockLogger.AddedPaths)        // No additions
	assert.Empty(t, mockLogger.UpdatePaths)       // No updates
}

func TestDesiredPartialMatchesInMap(t *testing.T) {
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now},     // Matches desired
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // New key
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now},    // Matches reported
			"key2": MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Differs from reported
			// "key3" is missing
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // Added
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1}, mockLogger.AcknowledgedPaths) // Acknowledged
	assert.ElementsMatch(t, []string{"field1.key2"}, mockLogger.UpdatePaths)    // Updated
	assert.ElementsMatch(t, []string{"field1.key3"}, mockLogger.AddedPaths)     // Added
}

func TestDesiredMixedNestedMapsAndStructs(t *testing.T) {
	type NestedStruct struct {
		SubField1 model.ValueAndTimestamp `json:"sub_field1"`
		SubField2 model.ValueAndTimestamp `json:"sub_field2"`
	}
	type TestStruct struct {
		Field1 map[string]model.ValueAndTimestamp `json:"field1"`
		Field2 NestedStruct                       `json:"field2"`
	}

	now := time.Now()

	reported := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now},     // Matches desired
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // New key
		},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match", Timestamp: now},      // Matches desired
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated value
		},
	}
	desired := TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key1": MockValueAndTimestamp{Value: "value1", Timestamp: now},    // Matches reported
			"key2": MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Differs from reported
			// "key3" is missing
		},
		Field2: NestedStruct{
			SubField1: MockValueAndTimestamp{Value: "match", Timestamp: now},     // Matches reported
			SubField2: MockValueAndTimestamp{Value: "old-value", Timestamp: now}, // Differs from reported
		},
	}

	mockLogger := &MockLogger{}
	opts := merge.DesiredOptions{
		Loggers:             merge.DesiredLoggers{mockLogger},
		OnlyNewerTimeStamps: false, // Disable timestamp comparison
	}

	result, err := merge.Desired(reported, desired, opts)

	assert.NoError(t, err)
	assert.Equal(t, TestStruct{
		Field1: map[string]model.ValueAndTimestamp{
			"key2": MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated
			"key3": MockValueAndTimestamp{Value: newValue, Timestamp: now},     // Added
		},
		Field2: NestedStruct{
			SubField2: MockValueAndTimestamp{Value: updatedValue, Timestamp: now}, // Updated
		},
	}, result)
	assert.ElementsMatch(t, []string{field1key1, field1subField1}, mockLogger.AcknowledgedPaths)  // Acknowledged
	assert.ElementsMatch(t, []string{"field1.key2", "field2.sub_field2"}, mockLogger.UpdatePaths) // Updated
	assert.ElementsMatch(t, []string{"field1.key3"}, mockLogger.AddedPaths)                       // Added
}
