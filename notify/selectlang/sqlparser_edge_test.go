package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParserEdgeCases tests various edge cases for the SQL parser
func TestParserEdgeCases(t *testing.T) {
	op := createTestOperation()

	testCases := []struct {
		name     string
		query    string
		isValid  bool
		expected bool
	}{
		{
			name:     "Unnecessary spaces in query",
			query:    "  SELECT   *   FROM   Notification   WHERE   obj.ID   ==   'device-123'  ",
			isValid:  true,
			expected: true,
		},
		{
			name:     "Multiple spaces between tokens",
			query:    "SELECT * FROM Notification WHERE obj.ID    ==    'device-123'",
			isValid:  true,
			expected: true,
		},
		{
			name:    "Empty query",
			query:   "",
			isValid: false,
		},
		{
			name:    "Incomplete query - only SELECT",
			query:   "SELECT",
			isValid: false,
		},
		{
			name:    "Missing table name",
			query:   "SELECT * FROM",
			isValid: false,
		},
		{
			name:    "Invalid table name",
			query:   "SELECT * FROM Invalid WHERE obj.ID == 'device-123'",
			isValid: false,
		},
		{
			name:    "Invalid field name format",
			query:   "SELECT * FROM Notification WHERE obj-ID == 'device-123'",
			isValid: false,
		},
		{
			name:    "Non-existent field",
			query:   "SELECT * FROM Notification WHERE obj.NonExistent == 'device-123'",
			isValid: false,
		},
		{
			name:    "Invalid operator",
			query:   "SELECT * FROM Notification WHERE obj.ID === 'device-123'",
			isValid: false,
		},
		{
			name:     "Valid query with comments-like characters",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123 /* comment */'",
			isValid:  true,
			expected: false, // No device with that ID including comment text
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)

			if !tc.isValid {
				assert.Error(t, err, "Expected an error for invalid query")
				assert.Nil(t, sel)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, selected)
		})
	}
}

// TestMalformedQueries tests handling of malformed queries
func TestMalformedQueries(t *testing.T) {
	testCases := []struct {
		name  string
		query string
	}{
		{
			name:  "SQL injection attempt 1",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123'; DROP TABLE Notification; --",
		},
		{
			name:  "SQL injection attempt 2",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' OR 1=1",
		},
		{
			name:  "Malformed parentheses 1",
			query: "SELECT * FROM Notification WHERE (obj.ID == 'device-123'",
		},
		{
			name:  "Malformed parentheses 2",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123')",
		},
		{
			name:  "Invalid WHERE clause combination",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' UNION SELECT * FROM Other",
		},
		{
			name:  "Non-SQL query",
			query: "GET /api/devices/device-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)

			// All these queries should be rejected with an error
			assert.Error(t, err, "Expected an error for malformed query")
			assert.Nil(t, sel)
		})
	}
}

// TestUnsupportedSQLFeatures tests that unsupported SQL features are rejected
func TestUnsupportedSQLFeatures(t *testing.T) {
	testCases := []struct {
		name  string
		query string
	}{
		{
			name:  "JOIN clause",
			query: "SELECT * FROM Notification JOIN Other ON Notification.ID = Other.ID",
		},
		{
			name:  "GROUP BY clause",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' GROUP BY obj.ID",
		},
		{
			name:  "ORDER BY clause",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' ORDER BY obj.ID",
		},
		{
			name:  "LIMIT clause",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' LIMIT 10",
		},
		{
			name:  "HAVING clause",
			query: "SELECT * FROM Notification WHERE obj.ID == 'device-123' HAVING COUNT(*) > 1",
		},
		{
			name:  "Subquery",
			query: "SELECT * FROM Notification WHERE obj.ID IN (SELECT ID FROM Other)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)

			// All these features should be rejected with an error
			assert.Error(t, err, "Expected an error for unsupported SQL feature")
			assert.Nil(t, sel)
		})
	}
}
