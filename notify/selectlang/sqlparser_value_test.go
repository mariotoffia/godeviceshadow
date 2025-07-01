package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test log.Value field operations
func TestLogValueOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match string (plain log)",
			query:    "SELECT * FROM Notification WHERE log.Value == 'online'",
			expected: true,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE log.Value != 'offline'",
			expected: true,
		},
		{
			name:     "Equal match map value",
			query:    "SELECT * FROM Notification WHERE log.Value == '22'",
			expected: true,
		},
		{
			name:     "Equal match desire log",
			query:    "SELECT * FROM Notification WHERE log.Value == 'auto'",
			expected: true,
		},
		{
			name:     "Numeric greater than match",
			query:    "SELECT * FROM Notification WHERE log.Value > 20",
			expected: true,
		},
		{
			name:     "Numeric greater than no match",
			query:    "SELECT * FROM Notification WHERE log.Value > 30",
			expected: false,
		},
		{
			name:     "Numeric less than match",
			query:    "SELECT * FROM Notification WHERE log.Value < 30",
			expected: true,
		},
		{
			name:     "Numeric less than no match",
			query:    "SELECT * FROM Notification WHERE log.Value < 20",
			expected: false,
		},
		{
			name:     "Numeric greater than or equal match",
			query:    "SELECT * FROM Notification WHERE log.Value >= 22",
			expected: true,
		},
		{
			name:     "Numeric less than or equal match",
			query:    "SELECT * FROM Notification WHERE log.Value <= 22",
			expected: true,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE log.Value IN 'online', 'offline'",
			expected: true,
		},
		{
			name:     "IN match numeric",
			query:    "SELECT * FROM Notification WHERE log.Value IN 22, 23",
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
			assert.Equal(t, tc.expected, selected, "Query: %s", tc.query)
		})
	}
}

// Test log.Value HAS operator operations
func TestLogValueHasOperator(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "HAS operator with existing key",
			query:    "SELECT * FROM Notification WHERE log.Value HAS 'temp'",
			expected: true,
		},
		{
			name:     "HAS operator with non-existing key",
			query:    "SELECT * FROM Notification WHERE log.Value HAS 'nonexistent'",
			expected: false,
		},
		{
			name:     "Complex query with HAS operator",
			query:    "SELECT * FROM Notification WHERE (obj.ID ~= 'device-\\d+' AND log.Value HAS 'temp') OR log.Path == 'devices/status'",
			expected: true,
		},
		{
			name:     "Combined HAS and other conditions",
			query:    "SELECT * FROM Notification WHERE log.Value HAS 'temp' AND log.Operation == 'add'",
			expected: true,
		},
		{
			name:     "Multiple HAS conditions - one passing",
			query:    "SELECT * FROM Notification WHERE log.Value HAS 'temp' OR log.Value HAS 'nonexistent'",
			expected: true,
		},
		{
			name:     "Multiple HAS conditions - none passing",
			query:    "SELECT * FROM Notification WHERE log.Value HAS 'nonexistent1' OR log.Value HAS 'nonexistent2'",
			expected: false,
		},
	}

	op := createTestOperation()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			selection, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)

			result, _ := selection.Select(op, true)
			assert.Equal(t, tc.expected, result)
		})
	}
}
