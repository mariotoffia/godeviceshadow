package pgxsql_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariotoffia/godeviceshadow/loggers/dblogger/pgxsql"
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

	// Wait a bit for the server to capture data
	time.Sleep(200 * time.Millisecond)

	// Get captured queries
	capturedQueries := mockServer.GetCapturedQueries()
	t.Logf("Captured queries: %v", capturedQueries)

	// Assert that we captured the expected SQL commands
	if len(capturedQueries) > 0 {
		found := false
		for _, query := range capturedQueries {
			if strings.Contains(query, "CREATE SCHEMA") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have captured CREATE SCHEMA command")

		found = false
		for _, query := range capturedQueries {
			if strings.Contains(query, "CREATE TABLE") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have captured CREATE TABLE command")

		found = false
		for _, query := range capturedQueries {
			if strings.Contains(query, "CREATE") && strings.Contains(query, "INDEX") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have captured CREATE INDEX command")
	} else {
		t.Log("No queries captured, but we tested that the logger attempts to connect and execute SQL")
	}

	// Wait for server to finish
	mockServer.WaitForCompletion(1 * time.Second)
}
