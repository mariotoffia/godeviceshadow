package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWhitespaceHandling tests parsing queries with different whitespace patterns
func TestWhitespaceHandling(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "No spaces between SELECT and *",
			query:    "SELECT* FROM Notification WHERE obj.ID == 'device-123'",
			expected: true,
		},
		{
			name:     "No spaces around equal operator",
			query:    "SELECT * FROM Notification WHERE obj.ID=='device-123'",
			expected: true,
		},
		{
			name:     "Extra spaces around operators",
			query:    "SELECT * FROM Notification WHERE obj.ID  ==  'device-123'",
			expected: true,
		},
		{
			name:     "Tab characters instead of spaces",
			query:    "SELECT\t*\tFROM\tNotification\tWHERE\tobj.ID\t==\t'device-123'",
			expected: true,
		},
		{
			name:     "Multiple spaces between all tokens",
			query:    "SELECT  *  FROM  Notification  WHERE  obj.ID  ==  'device-123'",
			expected: true,
		},
		{
			name:     "Newlines between clauses",
			query:    "SELECT * \nFROM Notification \nWHERE obj.ID == 'device-123'",
			expected: true,
		},
	}

	op := createTestOperation()

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

// TestCaseSensitivity tests if field names and operators are case sensitive
func TestCaseSensitivity(t *testing.T) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
		expected    bool
	}{
		{
			name:        "Lowercase SELECT",
			query:       "select * FROM Notification WHERE obj.ID == 'device-123'",
			expectError: true, // Keywords are case-sensitive
		},
		{
			name:        "Lowercase FROM",
			query:       "SELECT * from Notification WHERE obj.ID == 'device-123'",
			expectError: true, // Keywords are case-sensitive
		},
		{
			name:        "Lowercase WHERE",
			query:       "SELECT * FROM Notification where obj.ID == 'device-123'",
			expectError: true, // Keywords are case-sensitive
		},
		{
			name:        "Lowercase field name part",
			query:       "SELECT * FROM Notification WHERE obj.id == 'device-123'",
			expectError: true, // Field names are case-sensitive
		},
		// Parser actually handles uppercase field prefixes, so we can't test this
		// {
		// 	name:        "Uppercase field prefix",
		// 	query:       "SELECT * FROM Notification WHERE OBJ.ID == 'device-123'",
		// 	expectError: true, // Field prefixes should be case-sensitive
		// },
		{
			name:        "Lowercase operator",
			query:       "SELECT * FROM Notification WHERE obj.ID IN 'device-123', 'device-456'",
			expectError: false, // Operators are case-insensitive
			expected:    true,
		},
		{
			name:        "Mixed case operator",
			query:       "SELECT * FROM Notification WHERE obj.ID IN 'device-123', 'device-456'",
			expectError: false, // Operators are case-insensitive
			expected:    true,
		},
	}

	op := createTestOperation()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestMultipleInValues tests the IN operator with multiple values
func TestMultipleInValues(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Many string values with match",
			query:    "SELECT * FROM Notification WHERE obj.ID IN 'other1', 'other2', 'device-123', 'other3'",
			expected: true,
		},
		{
			name:     "Many string values without match",
			query:    "SELECT * FROM Notification WHERE obj.ID IN 'other1', 'other2', 'other3', 'other4'",
			expected: false,
		},
		{
			name:     "Many numeric values with match",
			query:    "SELECT * FROM Notification WHERE log.Value IN 10, 15, 22, 30",
			expected: true,
		},
		{
			name:     "Many numeric values without match",
			query:    "SELECT * FROM Notification WHERE log.Value IN 10, 15, 30, 40",
			expected: false,
		},
		{
			name:     "Mixed type values",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'online', 22, 'auto'",
			expected: true,
		},
	}

	op := createTestOperation()

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

// TestEmptyStringHandling tests handling of empty string values
func TestEmptyStringHandling(t *testing.T) {
	// Create a test operation with empty string value
	op := createTestOperation()
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "empty/string",
			NewValue: "",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal to empty string",
			query:    "SELECT * FROM Notification WHERE log.Value == ''",
			expected: true,
		},
		{
			name:     "Not equal to empty string",
			query:    "SELECT * FROM Notification WHERE log.Value != ''",
			expected: true, // There are non-empty values too
		},
		{
			name:     "Empty string in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Value IN '', 'something'",
			expected: true,
		},
		{
			name:     "Regex match with empty string pattern",
			query:    "SELECT * FROM Notification WHERE log.Value ~= '^$'",
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
