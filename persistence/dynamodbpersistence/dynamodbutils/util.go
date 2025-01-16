package dynamodbutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TestTableResource helps manage a DynamoDB table in tests.
type TestTableResource struct {
	Table  string
	Client *dynamodb.Client
}

// NewTestTableResource will creating the table if needed and waiting until it is ACTIVE
func NewTestTableResource(ctx context.Context, table string) *TestTableResource {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	r := &TestTableResource{
		Table:  table,
		Client: dynamodb.NewFromConfig(cfg),
	}

	// Create table if it doesn't exist, then wait for ACTIVE
	if err := r.createTableIfNotExists(ctx); err != nil {
		panic(fmt.Sprintf("failed to create or verify table %q: %v", table, err))
	}

	return r
}

// createTableIfNotExists checks if the table exists, creates it if needed,
// then waits until the table is ACTIVE.
func (r *TestTableResource) createTableIfNotExists(ctx context.Context) error {
	// check if the table exists
	_, err := r.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &r.Table,
	})
	if err == nil {
		// table exists, ensure it's ACTIVE
		return r.waitForTableActive(ctx)
	}

	var nf *types.ResourceNotFoundException

	if ok := errors.As(err, &nf); !ok {
		return err
	}

	// table does not exist -> create
	_, err = r.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: &r.Table,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		return fmt.Errorf("createTable failed: %w", err)
	}

	// wait for ACTIVE
	return r.waitForTableActive(ctx)
}

// waitForTableActive uses the TableExistsWaiter to wait until table is ACTIVE.
func (r *TestTableResource) waitForTableActive(ctx context.Context) error {
	waiter := dynamodb.NewTableExistsWaiter(r.Client)
	// wait max 2 minutes
	return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: &r.Table,
	}, 2*time.Minute)
}

// dispose deletes the table if 'delete' is true.
func (r *TestTableResource) Dispose(ctx context.Context, delete bool) error {
	if !delete {
		// Keep the table
		return nil
	}

	_, err := r.Client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: &r.Table,
	})
	if err != nil {
		return fmt.Errorf("failed to delete table %q: %w", r.Table, err)
	}

	waiter := dynamodb.NewTableNotExistsWaiter(r.Client)
	// wait for the table to be removed
	return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: &r.Table,
	}, 2*time.Minute)
}
