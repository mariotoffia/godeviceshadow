package dblogger

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model"
)

// DsqlLogger uses Aurora DSQL to log changes into a database.
//
// It's exposes begin, commit, and rollback methods to manage transaction.
//
// It will only log managed values when add, update.
type DsqlLogger struct {
	// client is the `DbLogger` that implements the database interactions.
	client DbLogger
	// batchSize is the number of values to log in a single batch operation.
	batchSize int
	// currentBatch is the current batch of values to log.
	currentBatch []LogValue
}

// LogValue is a value to be logged into the database.
type LogValue struct {
	Path      string
	Operation model.MergeOperation
	Value     model.ValueAndTimestamp
}

// DbLogger is the worker of the `DsqlLogger` that implements the
// database interactions.
//
// For example, it may implement the _AWS Aurora_ `rdsdata.Client` API interactions
// or the. `package github.com/jackc/pgx/v5` native _PostgreSQL_ client interactions.
type DbLogger interface {
	// Initialize make sure that the database has correct schema and is ready to use.
	Initialize(ctx context.Context) error
	// Upsert will insert or update the _values_ in a single batch operation.
	Upsert(ctx context.Context, table string, values []LogValue) error
}
