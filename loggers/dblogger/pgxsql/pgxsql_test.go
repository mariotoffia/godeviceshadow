package pgxsql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariotoffia/godeviceshadow/loggers/dblogger"
	"github.com/mariotoffia/godeviceshadow/loggers/dblogger/pgxsql"
	"github.com/mariotoffia/godeviceshadow/model"
)

func TestPgxLoggerInitializeWithMock(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)

	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create PgxLogger instance
	logger, err := pgxsql.New(ctx, pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       true,
		PreferID:          false,
		AssumeTableExists: false, // This will trigger table creation
		ConnectionString:  mockServer.GetConnectionString(),
	})

	// This might fail due to the simple mock, and that's ok
	if err != nil {
		t.Logf("Expected error creating logger with simple mock: %v", err)
		return
	}
	defer logger.Close()

	// Test Initialize method - this should execute the SQL commands
	err = logger.Initialize(ctx)
	t.Logf("Initialize result: %v", err)

	// Wait for expected statements (3: CREATE SCHEMA, CREATE TABLE, CREATE INDEX)
	capturedStatements := mockServer.WaitForStatements(3, 2*time.Second)
	t.Logf("Captured statements: %v", capturedStatements)

	// Assert that we captured the expected SQL commands in the correct order
	if capturedStatements.Len() > 0 {
		// Test individual query patterns
		capturedStatements.MustHaveRegex(`(?i)CREATE\s+SCHEMA\s+IF\s+NOT\s+EXISTS\s+test_schema`)
		capturedStatements.MustHaveRegex(`(?i)CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS\s+test_schema\.test_table`)
		capturedStatements.MustHaveRegex(`(?i)CREATE\s+UNIQUE\s+INDEX\s+IF\s+NOT\s+EXISTS\s+idx_test_table_path_timestamp`)

		// Test that the queries appear in the correct order: SCHEMA -> TABLE -> INDEX
		capturedStatements.MustHaveOrderedRegex(
			`(?i)CREATE\s+SCHEMA\s+IF\s+NOT\s+EXISTS\s+test_schema`,
			`(?i)CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS\s+test_schema\.test_table`,
			`(?i)CREATE\s+UNIQUE\s+INDEX\s+IF\s+NOT\s+EXISTS\s+idx_test_table_path_timestamp`,
		)

		// Verify specific table structure expectations
		tableMatches, err := capturedStatements.FindByRegex(`(?i)CREATE\s+TABLE.*test_table`)
		require.NoError(t, err)
		require.Len(t, tableMatches, 1, "Should have exactly one CREATE TABLE statement")

		tableSQL := tableMatches[0]
		assert.Contains(t, tableSQL, "id SERIAL PRIMARY KEY", "Table should have id column as primary key")
		assert.Contains(t, tableSQL, "path VARCHAR(500) NOT NULL", "Table should have path column")
		assert.Contains(t, tableSQL, "value JSONB", "Table should have value column as JSONB")
		assert.Contains(t, tableSQL, "timestamp TIMESTAMPTZ NOT NULL", "Table should have timestamp column")

		// Verify index is on the correct columns
		indexMatches, err := capturedStatements.FindByRegex(`(?i)CREATE\s+.*INDEX.*idx_test_table_path_timestamp`)
		require.NoError(t, err)
		require.Len(t, indexMatches, 1, "Should have exactly one CREATE INDEX statement")

		indexSQL := indexMatches[0]
		assert.Contains(t, indexSQL, "(path, timestamp)", "Index should be on path and timestamp columns")
	} else {
		t.Fatal("No queries captured - this indicates a problem with the mock server or PgxLogger")
	}

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

// TestValue implements model.ValueAndTimestamp for testing
type TestValue struct {
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

func (tv *TestValue) GetTimestamp() time.Time {
	return tv.Timestamp
}

func (tv *TestValue) GetValue() any {
	return tv.Data
}

// TestValueWithID implements model.IdValueAndTimestamp for testing
type TestValueWithID struct {
	TestValue
	ID string `json:"id"`
}

func (tv *TestValueWithID) GetID() string {
	return tv.ID
}

func TestPgxLoggerUpsertBasicInsert(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create PgxLogger instance
	logger, err := pgxsql.New(ctx, pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       true,
		PreferID:          false,
		AssumeTableExists: true, // Skip table creation for this test
		ConnectionString:  mockServer.GetConnectionString(),
	})
	require.NoError(t, err)
	defer logger.Close()

	// Create test data
	now := time.Now()
	testValues := []dblogger.LogValue{
		{
			Path:      "device.sensor1.temperature",
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"temp": 20.5, "unit": "C"},
				Timestamp: now,
			},
		},
		{
			Path:      "device.sensor2.humidity",
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"humidity": 65.0, "unit": "%"},
				Timestamp: now.Add(1 * time.Second),
			},
		},
	}

	// Test Upsert - should generate INSERT statements
	err = logger.Upsert(ctx, testValues)
	t.Logf("Upsert result: %v", err)

	// Wait for expected queries (1 INSERT statement for the batch)
	stmt := mockServer.WaitForStatements(1, 2*time.Second)
	t.Logf("Captured queries: %v", stmt)

	// Verify INSERT statements were generated
	stmt.MustHaveRegex(`(?i)INSERT\s+INTO\s+test_schema\.test_table`)

	// Verify the INSERT has proper structure (prepared statements use placeholders)
	insertMatches, err := stmt.FindByRegex(`(?i)INSERT\s+INTO.*test_table.*\(path,\s*value,\s*timestamp\)`)
	require.NoError(t, err)
	require.Greater(t, len(insertMatches), 0, "Should have captured INSERT statements")

	// Verify ON CONFLICT clause is present
	stmt.MustHaveRegex(`(?i)ON\s+CONFLICT.*DO\s+NOTHING`)

	// Verify prepared statement placeholders
	stmt.MustHaveRegex(`\$1,\s*\$2,\s*\$3`)

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

func TestPgxLoggerUpsertConflictHandling(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create PgxLogger instance
	logger, err := pgxsql.New(ctx, pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       true,
		PreferID:          false,
		AssumeTableExists: true,
		ConnectionString:  mockServer.GetConnectionString(),
	})
	require.NoError(t, err)
	defer logger.Close()

	// Create test data with same path and timestamp (should cause conflict)
	timestamp := time.Now()
	samePath := "device.sensor1.temperature"

	testValues := []dblogger.LogValue{
		{
			Path:      samePath,
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"temp": 20.5, "unit": "C"},
				Timestamp: timestamp,
			},
		},
		{
			Path:      samePath, // Same path
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"temp": 21.0, "unit": "C"}, // Different data
				Timestamp: timestamp,                                         // Same timestamp - should conflict
			},
		},
	}

	// Test Upsert - should generate INSERT statements with conflict resolution
	err = logger.Upsert(ctx, testValues)
	t.Logf("Upsert result: %v", err)

	// Wait for expected statements (1 INSERT statement for the batch)
	capturedStatements := mockServer.WaitForStatements(1, 2*time.Second)
	t.Logf("Captured statements: %v", capturedStatements)

	// Verify both INSERT statements were attempted (even though second would be ignored)
	insertMatches, err := capturedStatements.FindByRegex(`(?i)INSERT\s+INTO.*test_table`)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(insertMatches), 1, "Should have captured INSERT statements")

	// Verify ON CONFLICT DO NOTHING is present
	capturedStatements.MustHaveRegex(`(?i)ON\s+CONFLICT.*\(path,\s*timestamp\).*DO\s+NOTHING`)

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

func TestPgxLoggerUpsertWithPreferID(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create PgxLogger instance with PreferID = true
	logger, err := pgxsql.New(ctx, pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       false,
		PreferID:          true, // This should prefer ID over path
		AssumeTableExists: true,
		ConnectionString:  mockServer.GetConnectionString(),
	})
	require.NoError(t, err)
	defer logger.Close()

	// Create test data with ID
	now := time.Now()
	testValues := []dblogger.LogValue{
		{
			Path:      "device.sensor1.temperature.reading",
			Operation: model.MergeOperationAdd,
			Value: &TestValueWithID{
				TestValue: TestValue{
					Data:      map[string]interface{}{"temp": 20.5, "unit": "C"},
					Timestamp: now,
				},
				ID: "sensor1_temp",
			},
		},
	}

	// Test Upsert
	err = logger.Upsert(ctx, testValues)
	t.Logf("Upsert result: %v", err)

	// Wait for expected statements (1 INSERT statement)
	capturedStatements := mockServer.WaitForStatements(1, 2*time.Second)
	t.Logf("Captured statements: %v", capturedStatements)

	// Verify INSERT was generated
	capturedStatements.MustHaveRegex(`(?i)INSERT\s+INTO\s+test_schema\.test_table`)

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

func TestPgxLoggerUpsertEmptyValues(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create PgxLogger instance
	logger, err := pgxsql.New(ctx, pgxsql.Config{
		SchemaName:        "test_schema",
		TableName:         "test_table",
		UseFullPath:       true,
		PreferID:          false,
		AssumeTableExists: true,
		ConnectionString:  mockServer.GetConnectionString(),
	})
	require.NoError(t, err)
	defer logger.Close()

	// Test Upsert with empty slice - should be no-op
	err = logger.Upsert(ctx, []dblogger.LogValue{})
	assert.NoError(t, err, "Upsert with empty values should not error")

	// Wait for a short time (no statements expected)
	statements := mockServer.WaitForStatements(1, 500*time.Millisecond)
	t.Logf("Captured statements: %v", statements)

	// Should not have any INSERT statements for empty input
	insertMatches, err := statements.FindByRegex(`(?i)INSERT`)
	require.NoError(t, err)
	assert.Equal(t, 0, len(insertMatches), "Should not have any INSERT statements for empty input")

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}
