package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComplexExpressionCombinations tests additional complex logical expressions with parentheses
func TestComplexExpressionCombinations(t *testing.T) {
	op := createTestOperation()

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Triple AND condition (all true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND obj.Name == 'homeShadow' AND obj.Operation == 'report'",
			expected: true,
		},
		{
			name:     "Triple AND condition (one false)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND obj.Name == 'homeShadow' AND obj.Operation == 'wrong'",
			expected: false,
		},
		{
			name:     "Triple OR condition (all true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' OR obj.Name == 'homeShadow' OR obj.Operation == 'report'",
			expected: true,
		},
		{
			name:     "Triple OR condition (one true)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'wrong-id' OR obj.Name == 'homeShadow' OR obj.Operation == 'wrong'",
			expected: true,
		},
		{
			name:     "Triple OR condition (all false)",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'wrong-id' OR obj.Name == 'wrongName' OR obj.Operation == 'wrong'",
			expected: false,
		},
		{
			name:     "Complex nested expression",
			query:    "SELECT * FROM Notification WHERE (obj.ID == 'device-123' AND (obj.Name == 'homeShadow' OR obj.Operation == 'wrong')) OR (log.Value == 'online')",
			expected: true,
		},
		{
			name:     "Mixed operators without parentheses",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123' AND obj.Name == 'homeShadow' OR obj.Operation == 'wrong'",
			expected: true, // AND has precedence over OR
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

// TestFieldValueRangeConditions tests conditions with value ranges
func TestFieldValueRangeConditions(t *testing.T) {
	// Create a test operation with a range of numeric values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/temp1",
			NewValue: 18.5,
		},
		changelogger.PlainValue{
			Path:     "test/temp2",
			NewValue: 22.0,
		},
		changelogger.PlainValue{
			Path:     "test/temp3",
			NewValue: 25.5,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Value in comfortable range",
			query:    "SELECT * FROM Notification WHERE log.Value >= 20.0 AND log.Value <= 24.0",
			expected: true, // 22.0 is in range
		},
		{
			name:     "Value in cold range",
			query:    "SELECT * FROM Notification WHERE log.Value < 20.0",
			expected: true, // 18.5 is below 20.0
		},
		{
			name:     "Value in hot range",
			query:    "SELECT * FROM Notification WHERE log.Value > 24.0",
			expected: true, // 25.5 is above 24.0
		},
		{
			name:     "Complex range condition",
			query:    "SELECT * FROM Notification WHERE (log.Value < 19.0 OR log.Value > 24.0)",
			expected: true, // 18.5 and 25.5 match
		},
		{
			name:     "Range with OR and AND mixed",
			query:    "SELECT * FROM Notification WHERE log.Value > 18.0 AND log.Value < 19.0 OR log.Value > 25.0",
			expected: true, // 18.5 and 25.5 match
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

// TestRegexPatternVariations tests different regex pattern variations
func TestRegexPatternVariations(t *testing.T) {
	// Create a test operation with values for regex testing
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/email",
			NewValue: "user@example.com",
		},
		changelogger.PlainValue{
			Path:     "test/phone",
			NewValue: "+1-555-123-4567",
		},
		changelogger.PlainValue{
			Path:     "test/date",
			NewValue: "2023-07-01",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Email regex pattern",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '.*@example\\.com'",
			expected: true,
		},
		{
			name:     "Phone number regex pattern",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '\\+\\d+-\\d+-\\d+-\\d+'",
			expected: true,
		},
		{
			name:     "Date format regex pattern",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '\\d{4}-\\d{2}-\\d{2}'",
			expected: true,
		},
		{
			name:     "Complex regex with multiple patterns",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '(\\d{4}-\\d{2}-\\d{2}|\\+\\d+.*)'",
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

// TestMixedTypeComparisons tests comparisons between different types
func TestMixedTypeComparisons(t *testing.T) {
	// Create a test operation with mixed types
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/number_as_string",
			NewValue: "42",
		},
		changelogger.PlainValue{
			Path:     "test/string_that_looks_like_bool",
			NewValue: "true",
		},
		changelogger.PlainValue{
			Path:     "test/actual_number",
			NewValue: 42,
		},
		changelogger.PlainValue{
			Path:     "test/actual_bool",
			NewValue: true,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Numeric string equal to number",
			query:    "SELECT * FROM Notification WHERE log.Value == 42",
			expected: true, // Both "42" string and 42 number should match
		},
		{
			name:     "Numeric string compared as number",
			query:    "SELECT * FROM Notification WHERE log.Value > 40",
			expected: true, // "42" string should be compared as number and match
		},
		{
			name:     "Boolean string and boolean value",
			query:    "SELECT * FROM Notification WHERE log.Value == 'true'",
			expected: true, // Both "true" string and true boolean should match
		},
		{
			name:     "Mixed types in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN 42, 'true', 'online'",
			expected: true, // Multiple types should match
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
