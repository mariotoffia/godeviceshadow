package selectlang_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogNameWithoutPathFallback tests the log.Name field operations without using path as fallback
// This tests the new behavior where log.Name only checks map keys
func TestLogNameWithoutPathFallback(t *testing.T) {

	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match with map key",
			query:    "SELECT * FROM Notification WHERE log.Name == 'temp'",
			expected: true,
		},
		{
			name:     "Equal no match with map key",
			query:    "SELECT * FROM Notification WHERE log.Name == 'unknown_key'",
			expected: false,
		},
		{
			name:     "Equal no match with path (should not check path)",
			query:    "SELECT * FROM Notification WHERE log.Name == 'sensors/temperature/indoor'",
			expected: false,
		},
		{
			name:     "Regex match with map key",
			query:    "SELECT * FROM Notification WHERE log.Name ~= '^te.*'",
			expected: true,
		},
		{
			name:     "IN match with map key",
			query:    "SELECT * FROM Notification WHERE log.Name IN 'temp', 'other'",
			expected: true,
		},
	}

	op := createTestOperationForName()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip per-log evaluation by forcing global evaluation
			sel, err := selectlang.ToSelection(tc.query)
			require.NoError(t, err)
			require.NotNil(t, sel)

			// Using global evaluation for testing backward compatibility
			result, _ := sel.Select(op, true)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Create an updated test operation specifically for log.Name tests
func createTestOperationForName() notifiermodel.NotifierOperation {
	op := createTestOperation()

	// Add a map with 'temp' key for log.Name checks
	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": 22},
	}

	// Ensure ManagedLog has entries with map keys
	op.MergeLogger.ManagedLog[model.MergeOperationAdd] = append(
		op.MergeLogger.ManagedLog[model.MergeOperationAdd],
		changelogger.ManagedValue{
			Path:     "some/random/path",
			NewValue: mvs,
		},
	)

	return op
}
