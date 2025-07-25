package pgxsql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mariotoffia/godeviceshadow/loggers/dblogger"
	"github.com/mariotoffia/godeviceshadow/model"
)

// PgxLogger implements the DbLogger interface using native PostgreSQL via pgx/v5.
type PgxLogger struct {
	pool              *pgxpool.Pool
	internalPool      bool
	schemaName        string
	tableName         string
	useFullPath       bool
	preferID          bool
	assumeTableExists bool
	tx                pgx.Tx
}

// Config contains the configuration for the PostgreSQL logger.
type Config struct {
	// SchemaName is the name of the schema in _PostgreSQL_.
	SchemaName string
	// TableName is the name of the table to log into.
	TableName string
	// UseFullPath is if true, store full path; if false, store only the last element (name) or ID.
	UseFullPath bool
	// PreferID is if true, prefer ID over path when available.
	PreferID bool
	// AssumeTableExists is if true, assume the table already exists and skip creation.
	AssumeTableExists bool
	// Pool is optional connection pool to use, if `nil`, a new pool will be created with the
	// supplied `ConnectionString`. This pool is managed by this instance.
	//
	// NOTE: If you provide a pool, you must manage its lifecycle (close it when done).
	Pool *pgxpool.Pool
	// ConnectionString is the connection string to use if no pool is provided.
	ConnectionString string
}

// New creates a new PgxLogger instance.
func New(ctx context.Context, config Config) (*PgxLogger, error) {
	var pool *pgxpool.Pool

	var internalPool bool

	if config.Pool != nil {
		pool = config.Pool
	} else {
		// Create a new connection pool if none provided
		if config.ConnectionString == "" {
			return nil, fmt.Errorf("pgxsql: ConnectionString must be provided if Pool is not set")
		}

		var err error

		pool, err = pgxpool.New(ctx, config.ConnectionString)
		if err != nil {
			return nil, fmt.Errorf("pgxsql: failed to create connection pool: %w", err)
		}

		internalPool = true
	}

	return &PgxLogger{
		pool:              pool,
		internalPool:      internalPool,
		schemaName:        config.SchemaName,
		tableName:         config.TableName,
		useFullPath:       config.UseFullPath,
		preferID:          config.PreferID,
		assumeTableExists: config.AssumeTableExists,
	}, nil
}

// Initialize creates the required table structure if it doesn't exist.
func (p *PgxLogger) Initialize(ctx context.Context) error {
	if p.assumeTableExists {
		return nil
	}

	// Create schema if it doesn't exist
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", p.schemaName)
	_, err := p.pool.Exec(ctx, createSchemaSQL)

	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Create table if it doesn't exist
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id SERIAL PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			value JSONB,
			timestamp TIMESTAMPTZ NOT NULL
		)
	`, p.schemaName, p.tableName)

	_, err = p.pool.Exec(ctx, createTableSQL)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Create index for better query performance and to support conflict resolution
	createIndexSQL := fmt.Sprintf(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_%s_path_timestamp 
		ON %s.%s (path, timestamp)
	`, p.tableName, p.schemaName, p.tableName)

	_, err = p.pool.Exec(ctx, createIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// Upsert performs batch insert of log values using native PostgreSQL.
//
// If the table already contains a record with the same path and timestamp,
// it will silently skip inserting the value, otherwise it will insert
// the new value.
func (p *PgxLogger) Upsert(ctx context.Context, values []dblogger.LogValue) error {
	if len(values) == 0 {
		return nil
	}

	// Build batch insert SQL with conflict resolution
	insertSQL := fmt.Sprintf(`
		INSERT INTO %s.%s (path, value, timestamp)
		VALUES ($1, $2, $3)
		ON CONFLICT (path, timestamp) DO NOTHING
	`, p.schemaName, p.tableName)

	// Use batch for better performance
	batch := &pgx.Batch{}

	for _, value := range values {
		// Serialize value to JSON
		valueJson, err := json.Marshal(value.Value.GetValue())
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}

		// Determine path to store
		pathToStore := p.valuePath(value.Path, value.Value)

		batch.Queue(
			insertSQL, pathToStore, string(valueJson), value.Value.GetTimestamp(),
		)
	}

	// Execute batch
	var sender interface {
		SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	}

	if p.tx != nil {
		sender = p.tx
	} else {
		sender = p.pool
	}

	// Do all operation here
	results := sender.SendBatch(ctx, batch)
	defer results.Close()

	// Process all results to ensure no errors
	for i := 0; i < len(values); i++ {
		_, err := results.Exec() // will increment the result cursor
		if err != nil {
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	return nil
}

// Begin starts a new transaction for batch operations.
func (p *PgxLogger) Begin(ctx context.Context) error {
	if p.tx != nil {
		return fmt.Errorf("transaction already active")
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	p.tx = tx
	return nil
}

// Commit commits the current transaction.
func (p *PgxLogger) Commit(ctx context.Context) error {
	if p.tx == nil {
		return nil
	}

	err := p.tx.Commit(ctx)
	p.tx = nil
	return err
}

// Rollback rolls back the current transaction.
func (p *PgxLogger) Rollback(ctx context.Context) error {
	if p.tx == nil {
		return nil
	}

	err := p.tx.Rollback(ctx)
	p.tx = nil
	return err
}

// valuePath returns either the full path, last element, or ID
// based on configuration and value type.
func (p *PgxLogger) valuePath(path string, value model.ValueAndTimestamp) string {
	// Check if value implements IdValueAndTimestamp and we should prefer ID
	if idValue, ok := value.(model.IdValueAndTimestamp); ok {
		if p.preferID || !p.useFullPath {
			id := idValue.GetID()
			if id != "" {
				return id
			}
		}
	}

	// Use full path if configured
	if p.useFullPath || path == "" {
		return path
	}

	// Extract last element from dot-separated path
	parts := strings.Split(path, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}

// Close will close the connection pool if created in `New`. If supplied,
// the caller is responsible for managing the lifecycle of the pool.
//
// If a ongoing transaction, it will rollback such.
//
// This implements the `io.Closer` interface.
func (p *PgxLogger) Close() {
	if p.tx != nil {
		p.tx.Rollback(context.Background())
		p.tx = nil
	}

	if p.pool != nil && p.internalPool {
		p.pool.Close()
	}
}
