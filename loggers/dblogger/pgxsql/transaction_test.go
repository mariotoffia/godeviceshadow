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

func TestPgxLoggerTransactionCommit(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	// Begin transaction
	err = logger.Begin(ctx)
	assert.NoError(t, err, "Begin transaction should not error")

	// Create test data
	now := time.Now()
	testValues := []dblogger.LogValue{
		{
			Path:      "device.sensor1.temperature",
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"temp": 22.5, "unit": "C"},
				Timestamp: now,
			},
		},
	}

	// Test Upsert within transaction
	err = logger.Upsert(ctx, testValues)
	assert.NoError(t, err, "Upsert within transaction should not error")

	// Commit transaction
	err = logger.Commit(ctx)
	assert.NoError(t, err, "Commit transaction should not error")

	// Wait for expected statements (BEGIN + INSERT + COMMIT)
	capturedStatements := mockServer.WaitForStatements(3, 3*time.Second)
	t.Logf("Captured statements: %v", capturedStatements)

	// Verify BEGIN was executed
	capturedStatements.MustHaveRegex(`(?i)BEGIN`)

	// Verify INSERT was executed
	capturedStatements.MustHaveRegex(`(?i)INSERT\s+INTO\s+test_schema\.test_table`)

	// Verify COMMIT was executed
	capturedStatements.MustHaveRegex(`(?i)COMMIT`)

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

func TestPgxLoggerTransactionRollback(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	// Begin transaction
	err = logger.Begin(ctx)
	assert.NoError(t, err, "Begin transaction should not error")

	// Create test data
	now := time.Now()
	testValues := []dblogger.LogValue{
		{
			Path:      "device.sensor1.temperature",
			Operation: model.MergeOperationAdd,
			Value: &TestValue{
				Data:      map[string]interface{}{"temp": 25.0, "unit": "C"},
				Timestamp: now,
			},
		},
	}

	// Test Upsert within transaction
	err = logger.Upsert(ctx, testValues)
	assert.NoError(t, err, "Upsert within transaction should not error")

	// Rollback transaction
	err = logger.Rollback(ctx)
	assert.NoError(t, err, "Rollback transaction should not error")

	// Wait for expected statements (BEGIN + INSERT + ROLLBACK)
	capturedStatements := mockServer.WaitForStatements(3, 3*time.Second)
	t.Logf("Captured statements: %v", capturedStatements)

	// Verify BEGIN was executed
	capturedStatements.MustHaveRegex(`(?i)BEGIN`)

	// Verify INSERT was executed
	capturedStatements.MustHaveRegex(`(?i)INSERT\s+INTO\s+test_schema\.test_table`)

	// Verify ROLLBACK was executed
	capturedStatements.MustHaveRegex(`(?i)ROLLBACK`)

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}

func TestPgxLoggerTransactionDoubleBeginError(t *testing.T) {
	// Create and start the mock server
	mockServer, err := NewMockServer(t)
	require.NoError(t, err)
	defer mockServer.Close()

	mockServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	// Begin first transaction
	err = logger.Begin(ctx)
	assert.NoError(t, err, "First Begin should not error")

	// Try to begin second transaction - should error
	err = logger.Begin(ctx)
	assert.Error(t, err, "Second Begin should error")
	assert.Contains(t, err.Error(), "transaction already active", "Error should mention active transaction")

	// Rollback to clean up
	err = logger.Rollback(ctx)
	assert.NoError(t, err, "Rollback should not error")

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}
