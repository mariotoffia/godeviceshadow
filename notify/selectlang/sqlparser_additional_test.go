package selectlang_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSpecialCharactersInStrings tests handling of special characters in string literals
func TestSpecialCharactersInStrings(t *testing.T) {
	// Create a test operation with values containing special characters
	specialCharsValues := make(map[string]any)
	specialCharsValues["quotes"] = "value with quotes"
	specialCharsValues["backslash"] = "value with backslash"
	specialCharsValues["unicode"] = "value with unicode: ñáéíóú"

	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     specialCharsValues,
	}

	// Create operation with special character values
	op := createTestOperation()
	op.MergeLogger.ManagedLog[model.MergeOperationAdd] = append(
		op.MergeLogger.ManagedLog[model.MergeOperationAdd],
		changelogger.ManagedValue{
			Path:     "test/specialchars",
			NewValue: mvs,
		},
	)

	// Add special character paths
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "path/with/special/characters/+*&^%$#@!",
			NewValue: "special path",
		},
		changelogger.PlainValue{
			Path:     "path/with/unicode/ñáéíóú",
			NewValue: "unicode path",
		},
		changelogger.PlainValue{
			Path:     "path/with/quotes/quotes",
			NewValue: "quotes in path",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Simple string with quotes",
			query:    "SELECT * FROM Notification WHERE log.Path == 'path/with/quotes/quotes'",
			expected: true,
		},
		{
			name:     "Backslash in regex pattern",
			query:    "SELECT * FROM Notification WHERE log.Path ~= 'path/with/special/characters/\\+'",
			expected: true,
		},
		{
			name:     "Unicode characters in path",
			query:    "SELECT * FROM Notification WHERE log.Path ~= 'path/with/unicode/.*'",
			expected: true,
		},
		{
			name:     "Special characters in IN operator",
			query:    "SELECT * FROM Notification WHERE log.Path IN 'path/with/special/characters/+*&^%$#@!', 'other/path'",
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

// TestFloatPrecision tests handling of floating point comparisons with different precision
func TestFloatPrecision(t *testing.T) {
	// Create test operation
	op := createTestOperation()

	// Add numeric values directly to the operation
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/float1",
			NewValue: 3.14159,
		},
		changelogger.PlainValue{
			Path:     "test/float2",
			NewValue: 10.5,
		},
		changelogger.PlainValue{
			Path:     "test/int1",
			NewValue: 42,
		},
		changelogger.PlainValue{
			Path:     "test/string_num",
			NewValue: "123.45",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Greater than integer",
			query:    "SELECT * FROM Notification WHERE log.Value > 40",
			expected: true,
		},
		{
			name:     "Less than float",
			query:    "SELECT * FROM Notification WHERE log.Value < 20",
			expected: true,
		},
		{
			name:     "Equal to float",
			query:    "SELECT * FROM Notification WHERE log.Value == 10.5",
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

// TestSyntaxErrorRecovery tests error handling and recovery from syntax errors
func TestSyntaxErrorRecovery(t *testing.T) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing closing quote",
			query:       "SELECT * FROM Notification WHERE obj.ID == 'missing-quote",
			expectError: true,
			errorMsg:    "syntax error",
		},
		{
			name:        "Invalid operator sequence",
			query:       "SELECT * FROM Notification WHERE obj.ID == == 'value'",
			expectError: true,
			errorMsg:    "syntax error",
		},
		{
			name:        "Missing operator",
			query:       "SELECT * FROM Notification WHERE obj.ID 'value'",
			expectError: true,
			errorMsg:    "syntax error",
		},
		{
			name:        "Invalid field access",
			query:       "SELECT * FROM Notification WHERE obj..ID == 'value'",
			expectError: true,
			errorMsg:    "syntax error",
		},
		{
			name:        "Missing FROM clause",
			query:       "SELECT * WHERE obj.ID == 'value'",
			expectError: true,
			errorMsg:    "syntax error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.query)
			if tc.expectError {
				assert.Error(t, err, "Expected an error for invalid syntax")
				if err != nil {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				assert.Nil(t, sel)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sel)
			}
		})
	}
}

// BenchmarkComplexQueries provides benchmark tests for complex query performance
func BenchmarkComplexQueries(b *testing.B) {
	op := createTestOperation()

	// Add test path for log.Path test
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/benchmark",
			NewValue: "benchmark-value",
		},
	)

	benchmarks := []struct {
		name  string
		query string
	}{
		{
			name: "MultipleANDConditions",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' AND " +
				"obj.Name == 'homeShadow' AND " +
				"obj.Operation == 'report' AND " +
				"log.Path ~= 'test/' AND " +
				"log.Value != 'offline'",
		},
		{
			name: "MultipleORConditions",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' OR " +
				"obj.ID == 'other-device' OR " +
				"obj.ID == 'another-device'",
		},
		{
			name: "NestedQueryWithParentheses",
			query: "SELECT * FROM Notification WHERE " +
				"(obj.ID == 'device-123' AND obj.Name == 'homeShadow') OR " +
				"(log.Path ~= 'test/' AND log.Value != 'offline')",
		},
	}

	for _, bm := range benchmarks {
		// Parse the query once outside the benchmark
		sel, err := selectlang.ToSelection(bm.query)
		if err != nil {
			b.Fatalf("Failed to parse query '%s': %v", bm.query, err)
		}

		b.Run(bm.name, func(b *testing.B) {
			// Reset the timer for the actual benchmark
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				_, _ = sel.Select(op, false)
			}
		})
	}
}

// TestComplexQueries tests the functionality of complex queries
func TestComplexQueries(t *testing.T) {
	op := createTestOperation()

	// Add test path for log.Path test
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/functionality",
			NewValue: "test-value",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name: "Multiple AND conditions",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' AND " +
				"obj.Name == 'homeShadow' AND " +
				"obj.Operation == 'report' AND " +
				"log.Path ~= 'test/'",
			expected: true,
		},
		{
			name: "Multiple OR conditions",
			query: "SELECT * FROM Notification WHERE " +
				"obj.ID == 'device-123' OR " +
				"obj.ID == 'other-device' OR " +
				"obj.ID == 'another-device'",
			expected: true,
		},
		{
			name: "Nested query with parentheses",
			query: "SELECT * FROM Notification WHERE " +
				"(obj.ID == 'device-123' AND obj.Name == 'homeShadow') OR " +
				"(log.Path ~= 'nonexistent/' AND log.Value == 'missing')",
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
