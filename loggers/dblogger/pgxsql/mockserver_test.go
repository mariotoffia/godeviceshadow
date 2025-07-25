package pgxsql_test

import (
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgproto3/v2"
)

// MockServer represents a PostgreSQL protocol mock server for testing
type MockServer struct {
	listener        net.Listener
	capturedQueries []string
	serverDone      chan bool
	t               *testing.T
	mtx             *sync.Mutex
}

// NewMockServer creates a new mock PostgreSQL server
func NewMockServer(t *testing.T) (*MockServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &MockServer{
		listener:        ln,
		capturedQueries: make([]string, 0),
		serverDone:      make(chan bool),
		t:               t,
		mtx:             &sync.Mutex{},
	}, nil
}

// Start begins the mock server in a goroutine
func (m *MockServer) Start() {
	go func() {
		defer close(m.serverDone)

		conn, err := m.listener.Accept()
		if err != nil {
			m.t.Logf("Accept error: %v", err)
			return
		}
		defer conn.Close()

		// Set a deadline to avoid hanging
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		// Use PostgreSQL backend protocol
		backend := pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn)

		// Handle the startup message
		startupMessage, err := backend.ReceiveStartupMessage()
		if err != nil {
			m.t.Logf("Error receiving startup: %v", err)
			return
		}
		m.t.Logf("Received startup: %v", startupMessage)

		// Send authentication OK
		if err := backend.Send(&pgproto3.AuthenticationOk{}); err != nil {
			m.t.Logf("Error sending auth: %v", err)
			return
		}
		if err := backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'}); err != nil {
			m.t.Logf("Error sending ready: %v", err)
			return
		}

		// Main message loop
		for {
			msg, err := backend.Receive()
			if err != nil {
				m.t.Logf("Error receiving message: %v", err)
				break
			}

			switch query := msg.(type) {
			case *pgproto3.Query:
				m.t.Logf("Received query: %s", query.String)
				m.capturedQueries = append(m.capturedQueries, query.String)

				// Send response
				if err := backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("OK")}); err != nil {
					m.t.Logf("Error sending command complete: %v", err)
				}
				if err := backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'}); err != nil {
					m.t.Logf("Error sending ready: %v", err)
				}

			case *pgproto3.Terminate:
				m.t.Log("Received terminate")
				return
			default:
				m.t.Logf("Received other message: %T", query)
			}
		}
	}()
}

// GetConnectionString returns the connection string for connecting to the mock server
func (m *MockServer) GetConnectionString() string {
	parts := strings.Split(m.listener.Addr().String(), ":")
	host := parts[0]
	port := parts[1]
	return fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)
}

// GetCapturedQueries returns the SQL queries captured by the mock server
func (m *MockServer) GetCapturedQueries() []string {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if len(m.capturedQueries) == 0 {
		return nil
	}
	return slices.Clone(m.capturedQueries)
}

// WaitForCompletion waits for the server to complete or timeout
func (m *MockServer) WaitForCompletion(timeout time.Duration) {
	select {
	case <-m.serverDone:
		m.t.Log("Server completed")
	case <-time.After(timeout):
		m.t.Log("Server timed out, but that's expected with our simple mock")
	}
}

// Close closes the mock server
func (m *MockServer) Close() error {
	return m.listener.Close()
}
