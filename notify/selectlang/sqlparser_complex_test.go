package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test complex expressions with AND, OR, and parentheses
func TestComplexExpressions(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name: "AND expression (both true)",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' AND obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name: "AND expression (one false)",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' AND obj.Name == 'wrongName'",
			expected: false,
		},
		{
			name: "OR expression (both true)",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' OR obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name: "OR expression (one true)",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'wrong-id' OR obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name: "OR expression (both false)",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'wrong-id' OR obj.Name == 'wrongName'",
			expected: false,
		},
		{
			name: "Complex expression with parentheses",
			query: "SELECT * FROM Notification WHERE " +
				"(obj.ID == 'device-123' AND obj.Operation == 'report') OR " +
				"(log.Path == 'sensors/temperature/indoor' AND log.Value > 20)",
			expected: true,
		},
		{
			name: "Complex expression with multiple ANDs and ORs",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' AND " +
				"(log.Operation == 'add' OR log.Operation == 'update') AND " +
				"(log.Path ~= 'sensors/.*' OR log.Path ~= 'devices/.*')",
			expected: true,
		},
		{
			name: "Nested parentheses",
			query: "SELECT * FROM Notification WHERE " +
				"(obj.ID == 'device-123' AND (log.Value > 20 OR log.Value == 'online')) OR " +
				"(obj.Name == 'wrongName' AND log.Operation == 'delete')",
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

// TestComplexNestedQuery tests a complex query with nested expressions and multiple conditions
func TestComplexNestedQuery(t *testing.T) {
	// Create a test operation with matching data using the helper from test_utils.go
	op := createComplexTestOperation()

	// The complex query to test
	complexQuery := `
        SELECT * FROM Notification WHERE
        (
            obj.ID ~= 'myDevice-\\d+' AND
            obj.Name == 'myShadow' AND
            obj.Operation IN 'report', 'desired'
        )
        AND
        (
            log.Operation IN 'add', 'update' AND
            log.Path ~= '^Sensors-.*-indoor$' AND
            log.Value == 'temp' AND
            (
                log.Value > 20 OR (log.Value ~= 're-\\d+' AND log.Value != 'apa' OR (log.Value > 99 AND log.Value != 'bubben-\\d+'))
            )
        )
        OR
        (log.Operation IN 'add', 'update')
    `

	// Parse the query
	sel, err := selectlang.ToSelection(complexQuery)
	require.NoError(t, err, "Failed to parse complex query")
	require.NotNil(t, sel, "Selection should not be nil")

	// Test the query against the operation
	matched, _ := sel.Select(op, false)
	assert.True(t, matched, "Complex query should match the operation")

	// Test individual components to help diagnose any issues
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "ID regex match",
			query:    "SELECT * FROM Notification WHERE obj.ID ~= 'myDevice-\\d+'",
			expected: true,
		},
		{
			name:     "Name match",
			query:    "SELECT * FROM Notification WHERE obj.Name == 'myShadow'",
			expected: true,
		},
		{
			name:     "Operation IN match",
			query:    "SELECT * FROM Notification WHERE obj.Operation IN 'report', 'desired'",
			expected: true,
		},
		{
			name:     "Log operation match",
			query:    "SELECT * FROM Notification WHERE log.Operation IN 'add', 'update'",
			expected: true,
		},
		{
			name:     "Path regex match",
			query:    "SELECT * FROM Notification WHERE log.Path ~= '^Sensors-.*-indoor$'",
			expected: true,
		},
		{
			name:     "Value temp match",
			query:    "SELECT * FROM Notification WHERE log.Value == 'temp'",
			expected: true,
		},
		{
			name:     "Value > 20 match",
			query:    "SELECT * FROM Notification WHERE log.Value > 20",
			expected: true,
		},
		{
			name:     "Value regex match",
			query:    "SELECT * FROM Notification WHERE log.Value ~= 're-\\d+'",
			expected: true,
		},
		{
			name:     "Value != 'apa' match",
			query:    "SELECT * FROM Notification WHERE log.Value != 'apa'",
			expected: true,
		},
		{
			name:     "First part of complex AND",
			query:    "SELECT * FROM Notification WHERE obj.ID ~= 'myDevice-\\d+' AND obj.Name == 'myShadow'",
			expected: true,
		},
		{
			name:     "Second part of complex AND",
			query:    "SELECT * FROM Notification WHERE log.Operation == 'add' AND log.Path ~= '^Sensors-.*-indoor$'",
			expected: true,
		},
		{
			name:     "OR with Value conditions",
			query:    "SELECT * FROM Notification WHERE log.Value > 20 OR log.Value ~= 're-\\d+'",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			matched, _ := sel.Select(op, false)
			assert.Equal(t, tc.expected, matched, "Query: %s", tc.query)
		})
	}
}
