package selectlang_test

import (
	"testing"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// Test log.Name field operations for map keys
func TestLogNameMapKeyOperations(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Equal match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name == 'temp'",
			expected: true,
		},
		{
			name:     "Equal no match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name == 'nonexistent'",
			expected: false,
		},
		{
			name:     "Regex match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name ~= '^te.*'",
			expected: true,
		},
		{
			name:     "Regex no match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name ~= '^xyz.*'",
			expected: false,
		},
		{
			name:     "IN match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name IN 'temp', 'other'",
			expected: true,
		},
		{
			name:     "IN no match for map key",
			query:    "SELECT * FROM Notification WHERE log.Name IN 'nonexistent1', 'nonexistent2'",
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
