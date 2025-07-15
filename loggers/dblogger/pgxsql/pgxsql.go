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
	// UseFullPath is if true, store full path; if false, store only the last element (name) or ID
	UseFullPath bool
	// PreferID is if true, prefer ID over path when available
	PreferID bool
	// AssumeTableExists is if true, assume the table already exists and skip creation
	AssumeTableExists bool
}

// New creates a new PgxLogger instance.
func New(pool *pgxpool.Pool, config Config) *PgxLogger {
	return &PgxLogger{
		pool:              pool,
		schemaName:        config.SchemaName,
		tableName:         config.TableName,
		useFullPath:       config.UseFullPath,
		preferID:          config.PreferID,
		assumeTableExists: config.AssumeTableExists,
	}
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
			value_json JSONB,
			timestamp TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`, p.schemaName, p.tableName)

	_, err = p.pool.Exec(ctx, createTableSQL)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Create index for better query performance
	createIndexSQL := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_%s_path_timestamp 
		ON %s.%s (path, timestamp DESC)
	`, p.tableName, p.schemaName, p.tableName)

	_, err = p.pool.Exec(ctx, createIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// Upsert performs batch insert of log values using native PostgreSQL.
func (p *PgxLogger) Upsert(ctx context.Context, values []dblogger.LogValue) error {
	if len(values) == 0 {
		return nil
	}

	// Build batch insert SQL
	insertSQL := fmt.Sprintf(`
		INSERT INTO %s.%s (path, value_json, timestamp)
		VALUES ($1, $2, $3)
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
		pathToStore := p.getPathToStore(value.Path, value.Value)

		batch.Queue(insertSQL,
			pathToStore,
			string(valueJson),
			value.Value.GetTimestamp(),
		)
	}

	// Execute batch
	var err error

	if p.tx != nil {
		// Use transaction if available
		results := p.tx.SendBatch(ctx, batch)
		defer results.Close()

		// Process all results to ensure no errors
		for range values {
			_, err = results.Exec()
			if err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}
		}
	} else {
		// Use connection from pool
		results := p.pool.SendBatch(ctx, batch)
		defer results.Close()

		// Process all results to ensure no errors
		for range values {
			_, err = results.Exec()
			if err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}
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

// getPathToStore returns either the full path, last element, or ID based on configuration and value type
func (p *PgxLogger) getPathToStore(path string, value model.ValueAndTimestamp) string {
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

// Close closes the connection pool.
func (p *PgxLogger) Close() {
	if p.tx != nil {
		p.tx.Rollback(context.Background())
		p.tx = nil
	}
	if p.pool != nil {
		p.pool.Close()
	}
}
