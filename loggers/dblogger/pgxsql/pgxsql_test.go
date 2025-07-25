package pgxsql

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPgxLoggerInitialize(t *testing.T) {
	// Create a mock PostgreSQL server following the official pgmock pattern
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}

	// Expect schema creation
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "CREATE SCHEMA IF NOT EXISTS public"}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("CREATE SCHEMA")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))

	// Expect table creation
	expectedTableSQL := `
		CREATE TABLE IF NOT EXISTS public.test_table (
			id SERIAL PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			value_json JSONB,
			timestamp TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: strings.ReplaceAll(expectedTableSQL, "\n", "")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))

	// Expect index creation
	expectedIndexSQL := `
		CREATE INDEX IF NOT EXISTS idx_test_table_path_timestamp 
		ON public.test_table (path, timestamp DESC)
	`
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: strings.ReplaceAll(expectedIndexSQL, "\n", "")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("CREATE INDEX")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))

	// Expect termination
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Terminate{}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer ln.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)

		conn, err := ln.Accept()
		if err != nil {
			serverErrChan <- err
			return
		}
		defer conn.Close()

		err = conn.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			serverErrChan <- err
			return
		}

		err = script.Run(pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn))
		if err != nil {
			serverErrChan <- err
			return
		}
	}()

	// Get the address and port
	parts := strings.Split(ln.Addr().String(), ":")
	host := parts[0]
	port := parts[1]
	connStr := fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pgConn, err := pgconn.Connect(ctx, connStr)
	require.NoError(t, err)

	// Execute schema creation
	results, err := pgConn.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS public").ReadAll()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Nil(t, results[0].Err)
	assert.Equal(t, "CREATE SCHEMA", string(results[0].CommandTag))

	// Execute table creation
	tableSQL := strings.ReplaceAll(expectedTableSQL, "\n", "")
	results, err = pgConn.Exec(ctx, tableSQL).ReadAll()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Nil(t, results[0].Err)
	assert.Equal(t, "CREATE TABLE", string(results[0].CommandTag))

	// Execute index creation
	indexSQL := strings.ReplaceAll(expectedIndexSQL, "\n", "")
	results, err = pgConn.Exec(ctx, indexSQL).ReadAll()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Nil(t, results[0].Err)
	assert.Equal(t, "CREATE INDEX", string(results[0].CommandTag))

	pgConn.Close(ctx)

	// Check server completed without error
	assert.NoError(t, <-serverErrChan)
}

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
