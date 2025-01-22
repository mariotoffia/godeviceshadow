package dynamodbpersistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils"
)

func (p *Persistence) Delete(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	//
	if len(operations) == 0 {
		return nil
	}

	maxBatchSize := 25
	maxRetries := 3
	table := p.config.Table

	if p.config.MaxWriteBatchSize > 0 {
		maxBatchSize = p.config.MaxWriteBatchSize
	}

	if p.config.MaxWriteRetries > 0 {
		maxRetries = p.config.MaxWriteRetries
	}

	var ucOperations, condOperations []persistencemodel.WriteOperation

	for _, op := range operations {
		if op.Version > 0 {
			condOperations = append(condOperations, op)
		} else {
			ucOperations = append(ucOperations, op)
		}
	}

	results := make([]persistencemodel.WriteResult, 0, len(operations))
	batches, errors := prepareDelete(ucOperations, table, maxBatchSize)

	results = append(results, errors...)

	// Batch Delete (version = 0)
	for _, b := range batches {
		rs := p.deleteBatch(ctx, b, table, maxRetries)
		results = append(results, rs...)
	}

	// Do conditional deletes (version > 0) - if any
	return append(results, p.deleteConditionalItems(ctx, condOperations, table)...)
}

func prepareDelete(
	operations []persistencemodel.WriteOperation,
	table string,
	maxBatchSize int,
) ([]DeleteBatch, []persistencemodel.WriteResult) {

	var (
		results []DeleteBatch
		errors  []persistencemodel.WriteResult
	)

	batches := utils.ToBatch(operations, maxBatchSize)

	for _, batch := range batches {
		var writes []types.WriteRequest

		for _, op := range batch {
			pk := toPartitionKey(op.ID)

			var sk string

			switch op.ID.ModelType {
			case 0: // combined
				sk = "DSC#" + op.ID.Name
			case persistencemodel.ModelTypeDesired:
				sk = "DSD#" + op.ID.Name
			case persistencemodel.ModelTypeReported:
				sk = "DSR#" + op.ID.Name
			default:
				errors = append(errors, persistencemodel.WriteResult{
					ID:    op.ID,
					Error: persistencemodel.Error400("invalid model type"),
				})

				continue
			}

			writes = append(writes, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: pk},
						"SK": &types.AttributeValueMemberS{Value: sk},
					},
				},
			})
		}

		results = append(results, DeleteBatch{
			Operations: batch,
			Items: map[string][]types.WriteRequest{
				table: writes,
			},
		})
	}

	return results, errors
}

type DeleteBatch struct {
	Operations []persistencemodel.WriteOperation
	Items      map[string][]types.WriteRequest
}

func (rb *DeleteBatch) OperationFromIDName(id, name string) persistencemodel.WriteOperation {
	for _, op := range rb.Operations {
		if op.ID.ID == id && op.ID.Name == name {
			return op
		}
	}

	return persistencemodel.WriteOperation{}
}

func (p *Persistence) deleteBatch(
	ctx context.Context,
	batch DeleteBatch,
	table string,
	maxRetries int,
) []persistencemodel.WriteResult {

	current := &dynamodb.BatchWriteItemInput{
		RequestItems: batch.Items,
	}

	// items will hold the results of the batch deletes
	var items []persistencemodel.WriteResult

	if len(batch.Items[table]) == 0 {
		return nil
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// context cancelled/errored -> mark all "current" items as failed
		if err := ctx.Err(); err != nil {
			return append(
				items, p.buildDeleteFailedResults(batch, current.RequestItems[table], err)...,
			)
		}

		res, err := p.client.BatchWriteItem(ctx, current)

		if err != nil {
			// mark all "current" items failed
			return append(
				items, p.buildDeleteFailedResults(batch, current.RequestItems[table], err)...,
			)
		}

		// no unprocessed items -> success
		unprocessed := res.UnprocessedItems

		processed := diffWriteRequest(current.RequestItems[table], unprocessed[table])

		items = append(items, p.buildDeleteSuccessResults(batch, processed)...)

		if len(unprocessed) == 0 || len(unprocessed[table]) == 0 {
			return items
		}

		if attempt == maxRetries {
			return append(
				items, p.buildDeleteFailedResults(
					batch, unprocessed[table],
					fmt.Errorf("unprocessed items remain after %d retries", maxRetries))...,
			)
		}

		// retry unprocessed
		current.RequestItems = unprocessed

		if err := exponentialBackoff(ctx, attempt); err != nil {
			return append(
				items, p.buildDeleteFailedResults(batch, current.RequestItems[table], err)...,
			)
		}
	}

	// Should never reach here
	return items
}

// buildDeleteFailedResults maps a slice of WriteRequests back to their
// corresponding WriteOperation and returns error results for each.
func (p *Persistence) buildDeleteFailedResults(
	batch DeleteBatch,
	requests []types.WriteRequest,
	err error,
) []persistencemodel.WriteResult {

	failed := make([]persistencemodel.WriteResult, 0, len(requests))

	for _, wr := range requests {
		del := wr.DeleteRequest

		if del == nil {
			continue
		}

		pkAttr, pkOk := del.Key["PK"].(*types.AttributeValueMemberS)
		skAttr, skOk := del.Key["SK"].(*types.AttributeValueMemberS)

		if !pkOk || !skOk {
			continue // skip invalid
		}

		id := primaryKeyToID(pkAttr.Value)
		name := sortKeyToName(skAttr.Value)
		op := batch.OperationFromIDName(id, name)

		failed = append(failed, persistencemodel.WriteResult{
			ID:    op.ID,
			Error: readBatchErrorFixup(err),
		})
	}

	return failed
}

// buildDeleteSuccessResults marks items as successfully deleted.
func (p *Persistence) buildDeleteSuccessResults(
	batch DeleteBatch,
	requests []types.WriteRequest,
) []persistencemodel.WriteResult {

	success := make([]persistencemodel.WriteResult, 0, len(requests))

	for _, wr := range requests {
		del := wr.DeleteRequest
		if del == nil {
			continue
		}
		pkAttr, pkOk := del.Key["PK"].(*types.AttributeValueMemberS)
		skAttr, skOk := del.Key["SK"].(*types.AttributeValueMemberS)
		if !pkOk || !skOk {
			continue
		}
		id := primaryKeyToID(pkAttr.Value)
		name := sortKeyToName(skAttr.Value)
		op := batch.OperationFromIDName(id, name)

		success = append(success, persistencemodel.WriteResult{ID: op.ID, Version: 0})
	}

	return success
}

func (p *Persistence) deleteConditionalItems(
	ctx context.Context,
	operations []persistencemodel.WriteOperation,
	table string,
) []persistencemodel.WriteResult {
	//
	if len(operations) == 0 {
		return nil
	}

	results := make([]persistencemodel.WriteResult, 0, len(operations))

	for _, op := range operations {
		pk := toPartitionKey(op.ID)

		var sk string

		switch op.ID.ModelType {
		case 0:
			sk = "DSC#" + op.ID.Name
		case persistencemodel.ModelTypeDesired:
			sk = "DSD#" + op.ID.Name
		case persistencemodel.ModelTypeReported:
			sk = "DSR#" + op.ID.Name
		default:
			results = append(results, persistencemodel.WriteResult{
				ID:    op.ID,
				Error: persistencemodel.Error400("invalid model type"),
			})

			continue
		}

		input := &dynamodb.DeleteItemInput{
			TableName: &table,
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pk},
				"SK": &types.AttributeValueMemberS{Value: sk},
			},
			ConditionExpression: aws.String("Version = :ver"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":ver": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", op.Version)},
			},
		}

		// Execute the delete
		_, err := p.client.DeleteItem(ctx, input)

		if err != nil {
			var cfe *types.ConditionalCheckFailedException

			if ok := errors.As(err, &cfe); ok {
				results = append(results, persistencemodel.WriteResult{
					ID: op.ID,
					Error: persistencemodel.Error409(
						fmt.Sprintf("conditional delete failed, expected version = %d", op.Version),
					),
				})
			} else {
				// Other error
				results = append(results, persistencemodel.WriteResult{
					ID:      op.ID,
					Version: op.Version,
					Error:   err,
				})
			}
		} else {
			// Success
			results = append(results, persistencemodel.WriteResult{
				ID:      op.ID,
				Version: op.Version,
			})
		}
	}

	return results
}
