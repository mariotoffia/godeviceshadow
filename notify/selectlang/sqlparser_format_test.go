package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test edge cases and error handling
func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
		expected    bool
	}{
		{
			name:        "Empty WHERE clause",
			query:       "SELECT * FROM Notification",
			expectError: true,
		},
		{
			name:        "Unknown field",
			query:       "SELECT * FROM Notification WHERE unknown.Field == 'value'",
			expectError: true, // The parser should return an error for unknown fields
		},
		{
			name:        "Invalid operator",
			query:       "SELECT * FROM Notification WHERE obj.ID @ 'device-123'",
			expectError: true, // This should cause a syntax error
		},
		{
			name:        "Missing closing parenthesis",
			query:       "SELECT * FROM Notification WHERE (obj.ID == 'device-123'",
			expectError: true, // This should cause a syntax error
		},
	}

	op := createTestOperation()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				// Recover from any panics that might occur during parsing
				if r := recover(); r != nil {
					if !tc.expectError {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			sel, err := selectlang.ToSelection(tc.query)

			if tc.expectError {
				// Either err should be non-nil, or sel should be nil (or both)
				if err == nil && sel != nil {
					t.Error("Expected an error but got nil error and non-nil selection")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if sel == nil {
					t.Error("Expected value not to be nil.")
				} else {
					selected, _ := sel.Select(op, false)
					assert.Equal(t, tc.expected, selected)
				}
			}
		})
	}
}

// TestQueryFormatting tests various indentation and line break patterns
func TestQueryFormatting(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name: "Multi-line formatted query",
			query: `
				SELECT * 
				FROM Notification 
				WHERE obj.ID == 'device-123'
				  AND obj.Name == 'homeShadow'
			`,
			expected: true,
		},
		{
			name: "Indented complex conditions",
			query: `
				SELECT * FROM Notification WHERE
					(
						obj.ID == 'device-123' AND
						obj.Name == 'homeShadow'
					)
					OR
					(
						log.Operation == 'add' AND
						log.Path == 'sensors/temperature/indoor'
					)
			`,
			expected: true,
		},
		{
			name:     "Single line with no spaces",
			query:    "SELECT*FROM Notification WHERE(obj.ID=='device-123')",
			expected: true,
		},
		{
			name:     "Random line breaks",
			query:    "SELECT * FROM Notification WHERE obj.ID\n==\n'device-123'\nAND\n\nobj.Name=='homeShadow'",
			expected: true,
		},
	}

	op := createTestOperation()
	op.ID.ID = "device-123"
	op.ID.Name = "homeShadow"

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
