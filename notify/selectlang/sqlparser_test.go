package selectlang_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/loggers/desirelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestOperation() notifiermodel.NotifierOperation {
	// Create a value with timestamp
	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": 22},
	}

	// Create a desire logger with some acknowledged values
	dl := desirelogger.New()
	dl.Acknowledge("device/settings/mode", &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     "auto",
	})

	// Create a merge logger with both managed and plain logs
	ml := changelogger.ChangeMergeLogger{
		ManagedLog: changelogger.ManagedLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "sensors/temperature/indoor",
					NewValue: mvs,
				},
			},
			model.MergeOperationUpdate: {
				{
					Path: "sensors/humidity/indoor",
					NewValue: &model.ValueAndTimestampImpl{
						Timestamp: time.Now().UTC(),
						Value:     map[string]any{"humidity": 45},
					},
				},
			},
		},
		PlainLog: changelogger.PlainLogMap{
			model.MergeOperationAdd: {
				{
					Path:     "devices/status",
					NewValue: "online",
				},
			},
		},
	}

	return notifiermodel.NotifierOperation{
		ID:           persistencemodel.PersistenceID{ID: "device-123", Name: "homeShadow"},
		Operation:    notifiermodel.OperationTypeReport,
		MergeLogger:  ml,
		DesireLogger: *dl,
		Reported:     map[string]any{"status": "active"},
		Desired:      map[string]any{"mode": "auto"},
	}
}

// Test obj.ID field operations
func TestObjIDOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'device-123'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE obj.ID == 'other-device'",
			expected: false,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE obj.ID != 'other-device'",
			expected: true,
		},
		{
			name:     "Not equal no match",
			query:    "SELECT * FROM Notification WHERE obj.ID != 'device-123'",
			expected: false,
		},
		{
			name:     "Regex match",
			query:    "SELECT * FROM Notification WHERE obj.ID ~= 'device-\\d+'",
			expected: true,
		},
		{
			name:     "Regex no match",
			query:    "SELECT * FROM Notification WHERE obj.ID ~= 'sensor-\\d+'",
			expected: false,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE obj.ID IN 'device-123', 'device-456'",
			expected: true,
		},
		{
			name:     "IN no match",
			query:    "SELECT * FROM Notification WHERE obj.ID IN 'device-456', 'device-789'",
			expected: false,
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

// Test obj.Name field operations
func TestObjNameOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match",
			query:    "SELECT * FROM Notification WHERE obj.Name == 'homeShadow'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE obj.Name == 'otherShadow'",
			expected: false,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE obj.Name != 'otherShadow'",
			expected: true,
		},
		{
			name:     "Not equal no match",
			query:    "SELECT * FROM Notification WHERE obj.Name != 'homeShadow'",
			expected: false,
		},
		{
			name:     "Regex match",
			query:    "SELECT * FROM Notification WHERE obj.Name ~= 'home.*'",
			expected: true,
		},
		{
			name:     "Regex no match",
			query:    "SELECT * FROM Notification WHERE obj.Name ~= 'office.*'",
			expected: false,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE obj.Name IN 'homeShadow', 'officeShadow'",
			expected: true,
		},
		{
			name:     "IN no match",
			query:    "SELECT * FROM Notification WHERE obj.Name IN 'officeShadow', 'kitchenShadow'",
			expected: false,
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

// Test obj.Operation field operations
func TestObjOperationOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match",
			query:    "SELECT * FROM Notification WHERE obj.Operation == 'report'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE obj.Operation == 'desired'",
			expected: false,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE obj.Operation != 'desired'",
			expected: true,
		},
		{
			name:     "Not equal no match",
			query:    "SELECT * FROM Notification WHERE obj.Operation != 'report'",
			expected: false,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE obj.Operation IN 'report', 'desired'",
			expected: true,
		},
		{
			name:     "IN no match",
			query:    "SELECT * FROM Notification WHERE obj.Operation IN 'desired', 'delete'",
			expected: false,
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

// Test log.Operation field operations
func TestLogOperationOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match (add)",
			query:    "SELECT * FROM Notification WHERE log.Operation == 'add'",
			expected: true,
		},
		{
			name:     "Equal match (update)",
			query:    "SELECT * FROM Notification WHERE log.Operation == 'update'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE log.Operation == 'delete'",
			expected: false,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE log.Operation != 'delete'",
			expected: true,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE log.Operation IN 'add', 'update'",
			expected: true,
		},
		{
			name:     "IN no match",
			query:    "SELECT * FROM Notification WHERE log.Operation IN 'delete', 'remove'",
			expected: false,
		},
		{
			name:     "Equal match acknowledge",
			query:    "SELECT * FROM Notification WHERE log.Operation == 'acknowledge'",
			expected: true,
		},
		{
			name:     "IN match acknowledge",
			query:    "SELECT * FROM Notification WHERE log.Operation IN 'acknowledge', 'update'",
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

// Test log.Path field operations
func TestLogPathOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match (managed log)",
			query:    "SELECT * FROM Notification WHERE log.Path == 'sensors/temperature/indoor'",
			expected: true,
		},
		{
			name:     "Equal match (plain log)",
			query:    "SELECT * FROM Notification WHERE log.Path == 'devices/status'",
			expected: true,
		},
		{
			name:     "Equal match (desire log)",
			query:    "SELECT * FROM Notification WHERE log.Path == 'device/settings/mode'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE log.Path == 'sensors/light/indoor'",
			expected: false,
		},
		{
			name:     "Not equal match",
			query:    "SELECT * FROM Notification WHERE log.Path != 'sensors/light/indoor'",
			expected: true,
		},
		{
			name:     "Regex match",
			query:    "SELECT * FROM Notification WHERE log.Path ~= 'sensors/.*/indoor'",
			expected: true,
		},
		{
			name:     "Regex no match",
			query:    "SELECT * FROM Notification WHERE log.Path ~= 'sensors/.*/outdoor'",
			expected: false,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE log.Path IN 'sensors/temperature/indoor', 'sensors/light/indoor'",
			expected: true,
		},
		{
			name:     "IN no match",
			query:    "SELECT * FROM Notification WHERE log.Path IN 'sensors/light/indoor', 'sensors/motion/indoor'",
			expected: false,
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

// Test log.Name field operations (using path as fallback)
func TestLogNameOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match (managed log)",
			query:    "SELECT * FROM Notification WHERE log.Name == 'sensors/temperature/indoor'",
			expected: true,
		},
		{
			name:     "Equal no match",
			query:    "SELECT * FROM Notification WHERE log.Name == 'unknown/path'",
			expected: false,
		},
		{
			name:     "Regex match",
			query:    "SELECT * FROM Notification WHERE log.Name ~= '.*temperature.*'",
			expected: true,
		},
		{
			name:     "IN match",
			query:    "SELECT * FROM Notification WHERE log.Name IN 'sensors/temperature/indoor', 'other/path'",
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
