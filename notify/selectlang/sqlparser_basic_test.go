package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNullHandling tests how the parser handles NULL values
func TestNullHandling(t *testing.T) {
	// Create a test operation with null value
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "null/value",
			NewValue: nil,
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal to null string representation",
			query:    "SELECT * FROM Notification WHERE log.Value == 'null'",
			expected: true,
		},
		{
			name:     "Not equal to empty string",
			query:    "SELECT * FROM Notification WHERE log.Value != ''",
			expected: true,
		},
		{
			name:     "NULL value in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'null', 'something'",
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

// TestBooleanValues tests handling of boolean values in SQL queries
func TestBooleanValues(t *testing.T) {
	// Create a test operation with boolean values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "boolean/true",
			NewValue: true,
		},
		changelogger.PlainValue{
			Path:     "boolean/false",
			NewValue: false,
		},
		changelogger.PlainValue{
			Path:     "boolean/string/true",
			NewValue: "true",
		},
		changelogger.PlainValue{
			Path:     "boolean/string/false",
			NewValue: "false",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal to boolean true as string",
			query:    "SELECT * FROM Notification WHERE log.Value == 'true'",
			expected: true,
		},
		{
			name:     "Equal to boolean false as string",
			query:    "SELECT * FROM Notification WHERE log.Value == 'false'",
			expected: true,
		},
		{
			name:     "Boolean value in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'true', 'false'",
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

// TestComparisonOperators tests all comparison operators with different data types
func TestComparisonOperators(t *testing.T) {
	// Create a test operation with various values
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "number/integer",
			NewValue: 100,
		},
		changelogger.PlainValue{
			Path:     "number/float",
			NewValue: 10.5,
		},
		changelogger.PlainValue{
			Path:     "string/numeric",
			NewValue: "50",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		// Equal operator
		{
			name:     "Equal operator with number",
			query:    "SELECT * FROM Notification WHERE log.Value == 100",
			expected: true,
		},
		{
			name:     "Equal operator with float",
			query:    "SELECT * FROM Notification WHERE log.Value == 10.5",
			expected: true,
		},
		{
			name:     "Equal operator with string",
			query:    "SELECT * FROM Notification WHERE log.Value == '50'",
			expected: true,
		},

		// Not equal operator
		{
			name:     "Not equal operator with number",
			query:    "SELECT * FROM Notification WHERE log.Value != 999",
			expected: true,
		},
		{
			name:     "Not equal operator with float",
			query:    "SELECT * FROM Notification WHERE log.Value != 99.9",
			expected: true,
		},
		{
			name:     "Not equal operator with string",
			query:    "SELECT * FROM Notification WHERE log.Value != 'wrong'",
			expected: true,
		},

		// Greater than operator
		{
			name:     "Greater than operator with smaller number",
			query:    "SELECT * FROM Notification WHERE log.Value > 50",
			expected: true,
		},
		{
			name:     "Greater than operator with larger number",
			query:    "SELECT * FROM Notification WHERE log.Value > 200",
			expected: false,
		},

		// Less than operator
		{
			name:     "Less than operator with larger number",
			query:    "SELECT * FROM Notification WHERE log.Value < 150",
			expected: true,
		},
		{
			name:     "Less than operator with smaller number",
			query:    "SELECT * FROM Notification WHERE log.Value < 5",
			expected: false, // Expecting false since our test value 100 is not < 5
		},

		// Greater than or equal operator
		{
			name:     "Greater than or equal with exact match",
			query:    "SELECT * FROM Notification WHERE log.Value >= 100",
			expected: true,
		},
		{
			name:     "Greater than or equal with smaller value",
			query:    "SELECT * FROM Notification WHERE log.Value >= 50",
			expected: true,
		},

		// Less than or equal operator
		{
			name:     "Less than or equal with exact match",
			query:    "SELECT * FROM Notification WHERE log.Value <= 10.5",
			expected: true,
		},
		{
			name:     "Less than or equal with larger value",
			query:    "SELECT * FROM Notification WHERE log.Value <= 150",
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
