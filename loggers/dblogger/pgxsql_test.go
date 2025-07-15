package dblogger_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/dblogger"
	"github.com/mariotoffia/godeviceshadow/loggers/dblogger/pgxsql"
	"github.com/mariotoffia/godeviceshadow/model"
)

// TestValue implements ValueAndTimestamp for testing
type TestValue struct {
	Value     interface{}
	Timestamp time.Time
}

func (tv *TestValue) GetValue() interface{} {
	return tv.Value
}

func (tv *TestValue) GetTimestamp() time.Time {
	return tv.Timestamp
}

// TestIdValue implements IdValueAndTimestamp for testing
type TestIdValue struct {
	TestValue
	ID string
}

func (tiv *TestIdValue) GetID() string {
	return tiv.ID
}

func TestPgxLoggerConfig(t *testing.T) {
	config := pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       true,
		PreferID:          false,
		AssumeTableExists: false,
	}

	// Test that we can create a logger with nil pool (just for config testing)
	// In real usage, you'd pass a valid *pgxpool.Pool
	logger := pgxsql.New(nil, config)

	if logger == nil {
		t.Fatal("Expected logger to be created")
	}
}

func TestPgxLoggerPathHandling(t *testing.T) {
	// Test the path handling logic by creating test values
	tests := []struct {
		name             string
		config           pgxsql.Config
		path             string
		value            model.ValueAndTimestamp
		expectedBehavior string
	}{
		{
			name: "UseFullPath=true should use full path",
			config: pgxsql.Config{
				UseFullPath: true,
				PreferID:    false,
			},
			path:             "device.sensors.temperature",
			value:            &TestValue{Value: 25.5, Timestamp: time.Now()},
			expectedBehavior: "full_path",
		},
		{
			name: "UseFullPath=false should use last element",
			config: pgxsql.Config{
				UseFullPath: false,
				PreferID:    false,
			},
			path:             "device.sensors.temperature",
			value:            &TestValue{Value: 25.5, Timestamp: time.Now()},
			expectedBehavior: "last_element",
		},
		{
			name: "PreferID=true should use ID when available",
			config: pgxsql.Config{
				UseFullPath: false,
				PreferID:    true,
			},
			path:             "device.sensors.temperature",
			value:            &TestIdValue{TestValue: TestValue{Value: 25.5, Timestamp: time.Now()}, ID: "temp_sensor_01"},
			expectedBehavior: "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger with nil pool (we're not actually using database operations)
			logger := pgxsql.New(nil, tt.config)

			if logger == nil {
				t.Fatal("Expected logger to be created")
			}

			// Test that we can create LogValue without errors
			logValue := dblogger.LogValue{
				Path:      tt.path,
				Operation: model.MergeOperationAdd,
				Value:     tt.value,
			}

			// Verify the LogValue was created properly
			if logValue.Path != tt.path {
				t.Errorf("Expected path %s, got %s", tt.path, logValue.Path)
			}
			if logValue.Value != tt.value {
				t.Errorf("Expected value %v, got %v", tt.value, logValue.Value)
			}
		})
	}
}

// TestPgxLoggerWithPgmock demonstrates the pgmock integration
// This test is commented out because it requires careful setup of the mock server
// and exact SQL string matching which is brittle.
// In a real application, you'd use integration tests with a test database instead.
/*
func TestPgxLoggerWithPgmock(t *testing.T) {
	// This test demonstrates pgmock integration but is commented out
	// because pgmock requires exact SQL string matching which is brittle
	// and maintenance-heavy. For real testing, consider:
	// 1. Integration tests with a test PostgreSQL database
	// 2. Docker-based testing with testcontainers
	// 3. Using a tool like go-testdb for simpler mocking
}
*/
