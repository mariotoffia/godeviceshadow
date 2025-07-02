package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegexNotOperator tests the regex not operator (~!=)
func TestRegexNotOperator(t *testing.T) {
	op := createTestOperation()

	// Add test values to the operation
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/value1",
			NewValue: "abc-123",
		},
		changelogger.PlainValue{
			Path:     "test/value2",
			NewValue: "xyz-456",
		},
		changelogger.PlainValue{
			Path:     "test/value3",
			NewValue: "special*chars",
		},
		changelogger.PlainValue{
			Path:     "test/email",
			NewValue: "user@example.com",
		},
	)

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Basic regex not match",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= 'abc.*'",
			expected: true, // There are values that don't match the pattern
		},
		{
			name:     "Exact regex not match",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= 'abc-123'",
			expected: true, // There are values that don't match exactly abc-123
		},
		{
			name:     "Regex not match with path",
			query:    "SELECT * FROM Notification WHERE log.Path ~!= 'test/value[12]'",
			expected: true, // test/value3 and test/email don't match
		},
		{
			name:     "Regex not match with special characters",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= '.*\\*.*'",
			expected: true, // Values without * will match
		},
		{
			name:     "Email regex not match",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= '.*@example\\.com'",
			expected: true, // Non-email values will match
		},
		{
			name:     "No match with regex not",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= '.*'",
			expected: false, // .* matches everything, so ~!= '.*' should match nothing
		},
		{
			name:     "Complex regex not match",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= '[a-z]+-\\d+'",
			expected: true, // Special*chars and the email don't match the pattern
		},
		{
			name:     "Regex not with object ID",
			query:    "SELECT * FROM Notification WHERE obj.ID ~!= 'wrong-.*'",
			expected: true, // device-123 doesn't match wrong-.*
		},
		{
			name:     "Regex not with complex expression",
			query:    "SELECT * FROM Notification WHERE log.Value ~!= '.*\\d+.*' AND log.Path ~!= '.*value.*'",
			expected: true, // There are values that don't contain digits AND paths that don't contain 'value'
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

// TestRegexVsRegexNot tests the difference between regex match and regex not match
func TestRegexVsRegexNot(t *testing.T) {
	op := createTestOperation()

	// Add a specific test value
	op.MergeLogger.PlainLog[model.MergeOperationAdd] = append(
		op.MergeLogger.PlainLog[model.MergeOperationAdd],
		changelogger.PlainValue{
			Path:     "test/pattern",
			NewValue: "abc-123",
		},
	)

	testCases := []struct {
		name             string
		regexQuery       string
		notRegexQuery    string
		regexExpected    bool
		notRegexExpected bool
	}{
		{
			name:             "Simple pattern",
			regexQuery:       "SELECT * FROM Notification WHERE log.Value ~= 'abc-\\d+'",
			notRegexQuery:    "SELECT * FROM Notification WHERE log.Value ~!= 'abc-\\d+'",
			regexExpected:    true,
			notRegexExpected: true, // There are other values that don't match
		},
		{
			name:             "Exact match",
			regexQuery:       "SELECT * FROM Notification WHERE log.Value ~= '^abc-123$'",
			notRegexQuery:    "SELECT * FROM Notification WHERE log.Value ~!= '^abc-123$'",
			regexExpected:    true,
			notRegexExpected: true, // There are other values that don't match
		},
		{
			name:             "Match all",
			regexQuery:       "SELECT * FROM Notification WHERE log.Value ~= '.*'",
			notRegexQuery:    "SELECT * FROM Notification WHERE log.Value ~!= '.*'",
			regexExpected:    true,
			notRegexExpected: false, // Nothing should not match .*
		},
		{
			name:             "Complex pattern",
			regexQuery:       "SELECT * FROM Notification WHERE log.Value ~= '(abc|xyz)-\\d+'",
			notRegexQuery:    "SELECT * FROM Notification WHERE log.Value ~!= '(abc|xyz)-\\d+'",
			regexExpected:    true,
			notRegexExpected: true, // There are values that don't match
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" (regex)", func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.regexQuery)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.regexExpected, selected)
		})

		t.Run(tc.name+" (regex not)", func(t *testing.T) {
			sel, err := selectlang.ToSelection(tc.notRegexQuery)
			require.NoError(t, err)
			require.NotNil(t, sel)

			selected, _ := sel.Select(op, false)
			assert.Equal(t, tc.notRegexExpected, selected)
		})
	}
}
