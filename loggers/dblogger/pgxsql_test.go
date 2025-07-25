package dblogger

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

func TestPgxLoggerUpsert(t *testing.T) {
	// Create a mock PostgreSQL server following the official pgmock pattern
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}

	// Expect table creation query
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, data JSONB NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))

	// Expect upsert query
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "INSERT INTO test_table (data) VALUES ('{\"key\":\"value\"}') ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data"}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("INSERT 0 1")}))
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

	// Execute table creation
	results, err := pgConn.Exec(ctx, "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, data JSONB NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)").ReadAll()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Nil(t, results[0].Err)
	assert.Equal(t, "CREATE TABLE", string(results[0].CommandTag))

	// Execute upsert
	results, err = pgConn.Exec(ctx, "INSERT INTO test_table (data) VALUES ('{\"key\":\"value\"}') ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data").ReadAll()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Nil(t, results[0].Err)
	assert.Equal(t, "INSERT 0 1", string(results[0].CommandTag))

	pgConn.Close(ctx)

	// Check server completed without error
	assert.NoError(t, <-serverErrChan)
}
