package rdsql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"github.com/mariotoffia/godeviceshadow/loggers/dblogger"
	"github.com/mariotoffia/godeviceshadow/model"
)

// RdsqlLogger implements the DbLogger interface using AWS Aurora RDS Data API.
type RdsqlLogger struct {
	client        *rdsdata.Client
	resourceArn   string
	secretArn     string
	databaseName  string
	tableName     string
	schemaName    string
	useFullPath   bool
	preferID      bool
	transactionId *string
}

// Config contains the configuration for the RDS Data API logger.
type Config struct {
	// ResourceArn is the ARN of the RDS Data API resource.
	ResourceArn string
	// SecretArn is the ARN of the secret containing database credentials.
	SecretArn string
	// DatabaseName is the name of the database in RDS Data API.
	DatabaseName string
	// TableName is the name of the table to log into.
	TableName string
	// SchemaName is the name of the schema in RDS Data API.
	SchemaName string
	// UseFullPath, if true, store full path; if false, store only the last element (name) or ID
	UseFullPath bool
	// PreferID, if true, prefer ID over path when available
	PreferID bool
}

// New creates a new RdsqlLogger instance.
func New(client *rdsdata.Client, config Config) *RdsqlLogger {
	return &RdsqlLogger{
		client:       client,
		resourceArn:  config.ResourceArn,
		secretArn:    config.SecretArn,
		databaseName: config.DatabaseName,
		tableName:    config.TableName,
		schemaName:   config.SchemaName,
		useFullPath:  config.UseFullPath,
		preferID:     config.PreferID,
	}
}

// Initialize creates the required table structure if it doesn't exist.
func (r *RdsqlLogger) Initialize(ctx context.Context) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id SERIAL PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			value_json TEXT,
			timestamp TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, r.schemaName, r.tableName)

	_, err := r.client.ExecuteStatement(ctx, &rdsdata.ExecuteStatementInput{
		ResourceArn: aws.String(r.resourceArn),
		SecretArn:   aws.String(r.secretArn),
		Database:    aws.String(r.databaseName),
		Sql:         aws.String(createTableSQL),
	})

	return err
}

// Upsert performs batch insert of log values using RDS Data API.
func (r *RdsqlLogger) Upsert(ctx context.Context, table string, values []dblogger.LogValue) error {
	if len(values) == 0 {
		return nil
	}

	// Build batch insert SQL with parameters
	placeholders := make([]string, 0, len(values))
	parameters := make([]types.SqlParameter, 0, len(values)*4)

	for i, value := range values {
		placeholders = append(placeholders, fmt.Sprintf("(:path_%d, :value_json_%d, :timestamp_%d)", i, i, i))

		// Serialize value to JSON
		valueJson, err := json.Marshal(value.Value.GetValue())
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}

		// Determine path to store
		pathToStore := r.getPathToStore(value.Path, value.Value)

		parameters = append(parameters,
			types.SqlParameter{
				Name:  aws.String(fmt.Sprintf("path_%d", i)),
				Value: &types.FieldMemberStringValue{Value: pathToStore},
			},
			types.SqlParameter{
				Name:  aws.String(fmt.Sprintf("value_json_%d", i)),
				Value: &types.FieldMemberStringValue{Value: string(valueJson)},
			},
			types.SqlParameter{
				Name:  aws.String(fmt.Sprintf("timestamp_%d", i)),
				Value: &types.FieldMemberStringValue{Value: value.Value.GetTimestamp().Format(time.RFC3339)},
			},
		)
	}

	insertSQL := fmt.Sprintf(`
		INSERT INTO %s.%s (path, value_json, timestamp)
		VALUES %s
	`, r.schemaName, r.tableName, strings.Join(placeholders, ", "))

	input := &rdsdata.ExecuteStatementInput{
		ResourceArn: aws.String(r.resourceArn),
		SecretArn:   aws.String(r.secretArn),
		Database:    aws.String(r.databaseName),
		Sql:         aws.String(insertSQL),
		Parameters:  parameters,
	}

	if r.transactionId != nil {
		input.TransactionId = r.transactionId
	}

	_, err := r.client.ExecuteStatement(ctx, input)
	return err
}

// BeginTransaction starts a new transaction for batch operations.
func (r *RdsqlLogger) BeginTransaction(ctx context.Context) error {
	result, err := r.client.BeginTransaction(ctx, &rdsdata.BeginTransactionInput{
		ResourceArn: aws.String(r.resourceArn),
		SecretArn:   aws.String(r.secretArn),
		Database:    aws.String(r.databaseName),
	})
	if err != nil {
		return err
	}

	r.transactionId = result.TransactionId
	return nil
}

// CommitTransaction commits the current transaction.
func (r *RdsqlLogger) CommitTransaction(ctx context.Context) error {
	if r.transactionId == nil {
		return nil
	}

	_, err := r.client.CommitTransaction(ctx, &rdsdata.CommitTransactionInput{
		ResourceArn:   aws.String(r.resourceArn),
		SecretArn:     aws.String(r.secretArn),
		TransactionId: r.transactionId,
	})

	r.transactionId = nil
	return err
}

// RollbackTransaction rolls back the current transaction.
func (r *RdsqlLogger) RollbackTransaction(ctx context.Context) error {
	if r.transactionId == nil {
		return nil
	}

	_, err := r.client.RollbackTransaction(ctx, &rdsdata.RollbackTransactionInput{
		ResourceArn:   aws.String(r.resourceArn),
		SecretArn:     aws.String(r.secretArn),
		TransactionId: r.transactionId,
	})

	r.transactionId = nil
	return err
}

// getPathToStore returns either the full path, last element, or ID based on configuration and value type
func (r *RdsqlLogger) getPathToStore(path string, value model.ValueAndTimestamp) string {
	// Check if value implements IdValueAndTimestamp and we should prefer ID
	if idValue, ok := value.(model.IdValueAndTimestamp); ok {
		if r.preferID || !r.useFullPath {
			id := idValue.GetID()
			if id != "" {
				return id
			}
		}
	}

	// Use full path if configured
	if r.useFullPath || path == "" {
		return path
	}

	// Extract last element from dot-separated path
	parts := strings.Split(path, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}
