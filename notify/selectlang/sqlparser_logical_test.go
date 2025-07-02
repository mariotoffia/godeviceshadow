package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComplexExpressions tests complex logical expressions with parentheses
func TestComplexExpressions(t *testing.T) {
	op := createTestOperation()

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "AND expression (both true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name:     "AND expression (one false)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND obj.Name == 'wrongName'",
			expected: false,
		},
		{
			name:     "OR expression (both true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' OR obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name:     "OR expression (one true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' OR obj.Name == 'wrongName'",
			expected: true,
		},
		{
			name:     "OR expression (both false)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'wrong-id' OR obj.Name == 'wrongName'",
			expected: false,
		},
		{
			name:     "Complex expression with parentheses",
			query:    "SELECT * FROM Notification WHERE (obj.ID == 'device-123' AND obj.Name == 'homeShadow') OR obj.Operation == 'report'",
			expected: true,
		},
		{
			name:     "Complex expression with multiple ANDs and ORs",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND (obj.Name == 'homeShadow' OR obj.Operation == 'wrong')",
			expected: true,
		},
		{
			name:     "Nested parentheses",
			query:    "SELECT * FROM Notification WHERE (obj.ID == 'device-123' AND (obj.Name == 'homeShadow' OR obj.Operation == 'wrong'))",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestMultipleConditionsOnSameField tests multiple conditions on the same field
func TestMultipleConditionsOnSameField(t *testing.T) {
	// Create a test operation with a range of numeric values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/range1",
			NewValue: 25,
		},
		changelogger.PlainValue{
			Path:     "test/range2",
			NewValue: 75,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Multiple conditions with AND",
			query:    "SELECT * FROM Notification WHERE log.Value > 20 AND log.Value < 30",
			expected: true,
		},
		{
			name:     "Multiple conditions with OR",
			query:    "SELECT * FROM Notification WHERE log.Value < 10 OR log.Value > 70",
			expected: true,
		},
		{
			name:     "Conflicting conditions (no match)",
			query:    "SELECT * FROM Notification WHERE log.Value > 100 AND log.Value < 50",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestNonStringRegex tests regex operations on non-string values
func TestNonStringRegex(t *testing.T) {
	// Create a test operation with non-string values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/numeric",
			NewValue: 12345,
		},
		changelogger.PlainValue{
			Path:     "test/boolean",
			NewValue: true,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Regex on numeric value",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '123.*'",
			expected: true,
		},
		{
			name:     "Regex on boolean value",
			query:    "SELECT * FROM Notification WHERE log.Value ~= 'true'",
			expected: true,
		},
		{
			name:     "Regex with numeric pattern",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '\\d+'",
			expected: true,
		},
		{
			name:     "Complex regex on non-string values",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '^(true|\\d+)$'",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestNullValueHandling tests handling of null values in various expressions
func TestNullValueHandling(t *testing.T) {
	// Create a test operation with null value
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/null",
			NewValue: nil,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal to null string",
			query:    "SELECT * FROM Notification WHERE log.Value == 'null'",
			expected: true,
		},
		{
			name:     "Not equal to null",
			query:    "SELECT * FROM Notification WHERE log.Value != 'null'",
			expected: true, // Other values in the log are not null
		},
		{
			name:     "Null in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'null', 'value'",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestBooleanValueHandling tests handling of boolean values
func TestBooleanValueHandling(t *testing.T) {
	// Create a test operation with boolean values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/bool/true",
			NewValue: true,
		},
		changelogger.PlainValue{
			Path:     "test/bool/false",
			NewValue: false,
		},
		changelogger.PlainValue{
			Path:     "test/bool/string/true",
			NewValue: "true",
		},
		changelogger.PlainValue{
			Path:     "test/bool/string/false",
			NewValue: "false",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal to true as string",
			query:    "SELECT * FROM Notification WHERE log.Value == 'true'",
			expected: true,
		},
		{
			name:     "Equal to false as string",
			query:    "SELECT * FROM Notification WHERE log.Value == 'false'",
			expected: true,
		},
		{
			name:     "Boolean in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'true', 'false'",
			expected: true,
		},
		{
			name:     "Equal with mixed case",
			query:    "SELECT * FROM Notification WHERE log.Value == 'True'",
			expected: false, // Case-sensitive comparison
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}
