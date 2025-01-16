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
	"golang.org/x/exp/rand"

	"github.com/mariotoffia/godeviceshadow/utils"
)

// TestTableResource helps manage a DynamoDB table in tests.
type TestTableResource struct {
	Table  string
	Client *dynamodb.Client
}

// NewTestTableResource will create the table if needed and wait until it is ACTIVE.
func NewTestTableResource(ctx context.Context, table string) *TestTableResource {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	r := &TestTableResource{
		Table:  table,
		Client: dynamodb.NewFromConfig(cfg),
	}

	// Create table if it doesn't exist, then wait for ACTIVE.
	if err := r.createTableIfNotExists(ctx); err != nil {
		panic(fmt.Sprintf("failed to create or verify table %q: %v", table, err))
	}

	return r
}

// createTableIfNotExists checks if the table exists, creates it if needed,
// then waits until the table is ACTIVE.
func (r *TestTableResource) createTableIfNotExists(ctx context.Context) error {
	_, err := r.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &r.Table,
	})

	if err == nil {
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

	return r.waitForTableActive(ctx)
}

// waitForTableActive uses the TableExistsWaiter to wait until the table is ACTIVE.
func (r *TestTableResource) waitForTableActive(ctx context.Context) error {
	waiter := dynamodb.NewTableExistsWaiter(r.Client)
	// wait max 2 minutes
	return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: &r.Table,
	}, 2*time.Minute)
}

type DisposeOpts struct {
	// DeleteTable will delete the table if true, otherwise it will do nothing.
	//
	// NOTE: When delete table, it will wait until the table has been removed.
	DeleteTable bool
	// DeleteItems will remove all items from the table if true.
	DeleteItems bool
}

// Dispose deletes the table if 'delete' is true. Otherwise, it does nothing.
func (r *TestTableResource) Dispose(ctx context.Context, opts ...DisposeOpts) error {
	var opt DisposeOpts

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.DeleteItems && !opt.DeleteTable {
		if err := r.ClearTable(ctx); err != nil {
			return fmt.Errorf("failed to clear table %q: %w", r.Table, err)
		}
	}

	if !opt.DeleteTable {
		return nil // Keep the table
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

// ClearTable removes all items from the table.
func (r *TestTableResource) ClearTable(ctx context.Context) error {
	// 1) Scan the entire table
	var lastEvaluatedKey map[string]types.AttributeValue
	var allItems []map[string]types.AttributeValue

	for {
		scanOut, err := r.Client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         &r.Table,
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		allItems = append(allItems, scanOut.Items...)
		if len(scanOut.LastEvaluatedKey) == 0 {
			break
		}
		lastEvaluatedKey = scanOut.LastEvaluatedKey
	}

	if len(allItems) == 0 {
		return nil // no items to remove
	}

	var writeRequests []types.WriteRequest
	for _, item := range allItems {
		pk, pkOk := item["PK"].(*types.AttributeValueMemberS)
		sk, skOk := item["SK"].(*types.AttributeValueMemberS)
		if !pkOk || !skOk {
			continue
		}

		wr := types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: pk.Value},
					"SK": &types.AttributeValueMemberS{Value: sk.Value},
				},
			},
		}
		writeRequests = append(writeRequests, wr)
	}

	batches := utils.ToBatch(writeRequests, 25)
	for _, batchChunk := range batches {
		if err := r.deleteBatchChunk(ctx, batchChunk); err != nil {
			return err
		}
	}

	return nil
}

// deleteBatchChunk tries to delete up to 25 items in a single BatchWriteItem
// call, with basic retry on UnprocessedItems.
func (r *TestTableResource) deleteBatchChunk(ctx context.Context, chunk []types.WriteRequest) error {
	maxRetries := 5
	pending := chunk

	for attempt := 1; attempt <= maxRetries; attempt++ {
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.Table: pending,
			},
		}

		resp, err := r.Client.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("batchWriteItem chunk failed: %w", err)
		}

		unprocessed := resp.UnprocessedItems[r.Table]
		if len(unprocessed) == 0 {
			return nil // done
		}

		// unprocessed -> retry
		if attempt == maxRetries {
			return fmt.Errorf("unprocessed items remain after %d retries: %d items",
				maxRetries, len(unprocessed))
		}

		pending = unprocessed

		backoff := (time.Duration(1<<attempt) * 100 * time.Millisecond) +
			(time.Duration(rand.Int63n(int64(50 * time.Millisecond))))

		// Cap the maximum backoff at 30s
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}

		time.Sleep(backoff)
	}

	return nil
}
