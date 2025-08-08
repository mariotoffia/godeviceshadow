package pgxsql_test

import (
	"fmt"
	"net"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgproto3/v2"
)

// CapturedDbCommunication represents a slice of captured database queries with helper methods
type CapturedDbCommunication []string

// FindByRegex returns all statements that match the given regex pattern
func (c CapturedDbCommunication) FindByRegex(pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	var matches []string
	for _, stmt := range c {
		if re.MatchString(stmt) {
			matches = append(matches, stmt)
		}
	}
	return matches, nil
}

// HasRegex returns true if any statement matches the given regex pattern
func (c CapturedDbCommunication) HasRegex(pattern string) (bool, error) {
	matches, err := c.FindByRegex(pattern)
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

// MustHaveRegex panics if no statement matches the given regex pattern
func (c CapturedDbCommunication) MustHaveRegex(pattern string) CapturedDbCommunication {
	has, err := c.HasRegex(pattern)
	if err != nil {
		panic(fmt.Sprintf("regex error: %v", err))
	}
	if !has {
		panic(fmt.Sprintf("no statement found matching pattern: %s", pattern))
	}
	return c
}

// HasOrderedRegex returns true if statements matching the patterns appear in the specified order
func (c CapturedDbCommunication) HasOrderedRegex(patterns ...string) (bool, error) {
	if len(patterns) == 0 {
		return true, nil
	}

	// Compile all patterns
	regexes := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Errorf("invalid regex pattern %d: %w", i, err)
		}
		regexes[i] = re
	}

	// Track which pattern we're looking for
	currentPattern := 0

	for _, query := range c {
		if currentPattern < len(regexes) && regexes[currentPattern].MatchString(query) {
			currentPattern++
			if currentPattern == len(regexes) {
				return true, nil // Found all patterns in order
			}
		}
	}

	return currentPattern == len(regexes), nil
}

// MustHaveOrderedRegex panics if statement matching the patterns don't appear in the specified order
func (c CapturedDbCommunication) MustHaveOrderedRegex(patterns ...string) CapturedDbCommunication {
	has, err := c.HasOrderedRegex(patterns...)
	if err != nil {
		panic(fmt.Sprintf("regex error: %v", err))
	}
	if !has {
		panic(fmt.Sprintf("statements not found in expected order: %v", patterns))
	}
	return c
}

// Len returns the number of captured queries
func (c CapturedDbCommunication) Len() int {
	return len(c)
}

// BindCapture holds captured bind parameter values as strings
type BindCapture struct {
	Params []string
}

// MockServer represents a PostgreSQL protocol mock server for testing
type MockServer struct {
	listener           net.Listener
	capturedStatements []string
	capturedBinds      []BindCapture
	serverDone         chan bool
	t                  *testing.T
	mtx                *sync.Mutex
	closed             bool
	txActive           bool
}

// NewMockServer creates a new mock PostgreSQL server
func NewMockServer(t *testing.T) (*MockServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &MockServer{
		listener:           ln,
		capturedStatements: make([]string, 0),
		capturedBinds:      make([]BindCapture, 0),
		serverDone:         make(chan bool),
		t:                  t,
		mtx:                &sync.Mutex{},
		closed:             false,
		txActive:           false,
	}, nil
}

// safeLog logs only if the server hasn't been closed
func (m *MockServer) safeLog(format string, args ...interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if !m.closed {
		if len(args) > 0 {
			m.t.Logf(format, args...)
		} else {
			m.t.Log(format)
		}
	}
}

// Start begins the mock server in a goroutine
func (m *MockServer) Start() {
	go func() {
		defer close(m.serverDone)

		conn, err := m.listener.Accept()
		if err != nil {
			m.safeLog("Accept error: %v", err)
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
			m.safeLog("Error receiving startup: %v", err)
			return
		}
		m.safeLog("Received startup: %v", startupMessage)

		// Send authentication OK
		if err := backend.Send(&pgproto3.AuthenticationOk{}); err != nil {
			m.safeLog("Error sending auth: %v", err)
			return
		}

		if err := backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'}); err != nil {
			m.safeLog("Error sending ready: %v", err)
			return
		}

		// Main message loop
		for {
			msg, err := backend.Receive()
			if err != nil {
				m.safeLog("Error receiving message: %v", err)
				break
			}

			switch query := msg.(type) {
			case *pgproto3.Query:
				m.safeLog("Received query: %s", query.String)
				switch strings.ToLower(strings.TrimSpace(query.String)) {
				case "begin":
					m.txActive = true
				case "commit", "rollback":
					m.txActive = false
				}
				m.mtx.Lock()
				m.capturedStatements = append(m.capturedStatements, query.String)
				m.mtx.Unlock()

				// Send response
				if err := backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("OK")}); err != nil {
					m.safeLog("Error sending command complete: %v", err)
				}
				status := byte('I')
				if m.txActive {
					status = 'T'
				}
				if err := backend.Send(&pgproto3.ReadyForQuery{TxStatus: status}); err != nil {
					m.safeLog("Error sending ready: %v", err)
				}

			case *pgproto3.Parse:
				m.safeLog("Received parse: %s", query.Query)
				m.mtx.Lock()
				m.capturedStatements = append(m.capturedStatements, query.Query)
				m.mtx.Unlock()

				// Send ParseComplete
				if err := backend.Send(&pgproto3.ParseComplete{}); err != nil {
					m.safeLog("Error sending parse complete: %v", err)
				}

			case *pgproto3.Bind:
				m.safeLog("Received bind")
				// Capture parameters as strings (best-effort)
				params := make([]string, 0, len(query.Parameters))
				for _, p := range query.Parameters {
					if p == nil {
						params = append(params, "<NULL>")
						continue
					}
					params = append(params, string(p))
				}
				m.mtx.Lock()
				m.capturedBinds = append(m.capturedBinds, BindCapture{Params: params})
				m.mtx.Unlock()

				// Send BindComplete
				if err := backend.Send(&pgproto3.BindComplete{}); err != nil {
					m.safeLog("Error sending bind complete: %v", err)
				}

			case *pgproto3.Describe:
				m.safeLog("Received describe: %c %s", query.ObjectType, query.Name)
				// For statement descriptions, send ParameterDescription and RowDescription
				if query.ObjectType == 'S' { // Statement
					// Send ParameterDescription first
					if err := backend.Send(&pgproto3.ParameterDescription{
						ParameterOIDs: []uint32{25, 3802, 1184}, // text, jsonb, timestamptz
					}); err != nil {
						m.safeLog("Error sending parameter description: %v", err)
					}
					// Then send NoData since INSERT doesn't return rows
					if err := backend.Send(&pgproto3.NoData{}); err != nil {
						m.safeLog("Error sending no data: %v", err)
					}
				} else {
					// Send NoData for other describe types
					if err := backend.Send(&pgproto3.NoData{}); err != nil {
						m.safeLog("Error sending no data: %v", err)
					}
				}

			case *pgproto3.Execute:
				m.safeLog("Received execute")
				// Send CommandComplete for successful execution
				if err := backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("INSERT 0 1")}); err != nil {
					m.safeLog("Error sending command complete: %v", err)
				}

			case *pgproto3.Sync:
				m.safeLog("Received sync")
				// Send ReadyForQuery to indicate ready for next command
				status := byte('I')
				if m.txActive {
					status = 'T'
				}
				if err := backend.Send(&pgproto3.ReadyForQuery{TxStatus: status}); err != nil {
					m.safeLog("Error sending ready for query: %v", err)
				}

			case *pgproto3.Terminate:
				m.safeLog("Received terminate")
				return
			default:
				m.safeLog("Received other message: %T", query)
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

// GetCapturedStatements returns the SQL statements captured by the mock server
func (m *MockServer) GetCapturedStatements() CapturedDbCommunication {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.capturedStatements) == 0 {
		return CapturedDbCommunication(nil)
	}
	return CapturedDbCommunication(slices.Clone(m.capturedStatements))
}

// GetCapturedBinds returns copies of captured bind parameter values
func (m *MockServer) GetCapturedBinds() []BindCapture {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if len(m.capturedBinds) == 0 {
		return nil
	}
	cp := make([]BindCapture, len(m.capturedBinds))
	copy(cp, m.capturedBinds)
	return cp
}

// WaitForStatements waits until the specified number of statements are captured or timeout occurs
func (m *MockServer) WaitForStatements(n int, timeout time.Duration) CapturedDbCommunication {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		statements := m.GetCapturedStatements()
		if statements.Len() >= n {
			return statements
		}
		if time.Now().After(deadline) {
			return statements // Return what we have, even if incomplete
		}
		<-ticker.C
	}
}

// WaitForBinds waits until the specified number of binds are captured or timeout occurs
func (m *MockServer) WaitForBinds(n int, timeout time.Duration) []BindCapture {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		binds := m.GetCapturedBinds()
		if len(binds) >= n {
			return binds
		}
		if time.Now().After(deadline) {
			return binds
		}
		<-ticker.C
	}
}

// WaitForCompletion waits for the server to complete or timeout
func (m *MockServer) WaitForCompletion(timeout time.Duration) {
	select {
	case <-m.serverDone:
		m.safeLog("Server completed")
	case <-time.After(timeout):
		m.safeLog("Server timed out, but that's expected with our simple mock")
	}
}

// Close closes the mock server
func (m *MockServer) Close() error {
	m.mtx.Lock()
	m.closed = true
	m.mtx.Unlock()
	return m.listener.Close()
}
