package dynamodbpersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils"

	"golang.org/x/exp/rand"
)

func (p *Persistence) Delete(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {

	// If no operations, return empty
	if len(operations) == 0 {
		return nil
	}

	// We’ll define a maxBatchSize and maxRetries (similar to Read).
	// Adjust or load from p.config as needed.
	maxBatchSize := 25 // DynamoDB BatchWriteItem limit is 25
	maxRetries := 3    // or pull from p.config.MaxWriteRetries if that exists
	table := p.config.Table

	// Prepare the batch requests
	batches, prepErrors := prepareDelete(operations, table, maxBatchSize)

	// Our final list of results
	results := make([]persistencemodel.WriteResult, 0, len(operations))
	results = append(results, prepErrors...) // any model-type or input errors

	// For each batch, do the actual deletes
	for _, b := range batches {
		wrs := p.deleteBatch(ctx, b, table, maxRetries)
		results = append(results, wrs...)
	}

	// Return aggregated results
	return results
}

func prepareDelete(
	operations []persistencemodel.WriteOperation,
	table string,
	maxBatchSize int,
) ([]DeleteBatch, []persistencemodel.WriteResult) {

	// We’ll build a “batch” of up to maxBatchSize "WriteRequests" at a time
	// to feed into a DynamoDB.BatchWriteItem call.
	var (
		batches []DeleteBatch
		errors  []persistencemodel.WriteResult
	)

	// Use some function like your utils.ToBatch to chunk operations
	opBatches := utils.ToBatch(operations, maxBatchSize)

	for _, chunk := range opBatches {
		var writes []types.WriteRequest

		for _, op := range chunk {
			// Build one or more DeleteRequests for each operation
			pk := toPartitionKey(op.ID)

			// Figure out which SK(s) to delete
			var sortKeys []string
			switch op.ID.ModelType {
			case 0: // combined => delete all three
				sortKeys = []string{
					"DSC#" + op.ID.Name,
					"DSD#" + op.ID.Name,
					"DSR#" + op.ID.Name,
				}
			case persistencemodel.ModelTypeDesired:
				sortKeys = []string{"DSD#" + op.ID.Name}
			case persistencemodel.ModelTypeReported:
				sortKeys = []string{"DSR#" + op.ID.Name}
			default:
				// Invalid model type => skip & record error
				errors = append(errors, persistencemodel.WriteResult{
					ID:    op.ID,
					Error: persistencemodel.Error400("invalid model type"),
				})
				continue
			}

			// For each SK needed, add a DeleteRequest
			for _, sk := range sortKeys {
				delReq := types.WriteRequest{
					DeleteRequest: &types.DeleteRequest{
						Key: map[string]types.AttributeValue{
							"PK": &types.AttributeValueMemberS{Value: pk},
							"SK": &types.AttributeValueMemberS{Value: sk},
						},
					},
				}
				writes = append(writes, delReq)
			}
		}

		// Put them in a map for BatchWriteItem
		batchMap := map[string][]types.WriteRequest{
			table: writes,
		}

		// Build our DeleteBatch
		db := DeleteBatch{
			Operations: chunk,
			Items:      batchMap,
		}
		batches = append(batches, db)
	}

	return batches, errors
}

type DeleteBatch struct {
	// The original subset of WriteOperations that produced these items
	Operations []persistencemodel.WriteOperation
	// Items is the map[tableName] => []WriteRequest for DynamoDB.BatchWriteItem
	Items map[string][]types.WriteRequest
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

	// Build a final slice of results
	var finalResults []persistencemodel.WriteResult

	// Possibly short-circuit if no items
	if len(batch.Items[table]) == 0 {
		return nil
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Check if context is cancelled
		if err := ctx.Err(); err != nil {
			// Mark all current items as failed
			fr := p.buildDeleteFailedResults(batch, current.RequestItems[table], err)
			finalResults = append(finalResults, fr...)
			return finalResults
		}

		// Do the batch write
		res, err := p.client.BatchWriteItem(ctx, current)
		if err != nil {
			// All current items failed
			fr := p.buildDeleteFailedResults(batch, current.RequestItems[table], err)
			finalResults = append(finalResults, fr...)
			return finalResults
		}

		// Check unprocessed
		unprocessed := res.UnprocessedItems

		// If no unprocessed => success
		if len(unprocessed) == 0 || len(unprocessed[table]) == 0 {
			// We assume everything is deleted.
			// Return success results for each item that was in "current"
			// or do so item-by-item if you want them individually.
			fr := p.buildDeleteSuccessResults(batch, current.RequestItems[table])
			finalResults = append(finalResults, fr...)
			return finalResults
		}

		// If we have unprocessed items but reached maxRetries => fail them
		if attempt == maxRetries {
			fr := p.buildDeleteFailedResults(batch, unprocessed[table],
				fmt.Errorf("unprocessed items remain after %d retries", maxRetries))
			finalResults = append(finalResults, fr...)
			return finalResults
		}

		// Otherwise, continue retry with unprocessed
		current.RequestItems = unprocessed

		// Exponential backoff with jitter
		backoff := (time.Duration(1<<attempt) * 100 * time.Millisecond) +
			(time.Duration(rand.Int63n(int64(50 * time.Millisecond))))

		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
		time.Sleep(backoff)
	}

	// theoretically never reached if we return inside loop
	return finalResults
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
			continue // skip if somehow not a delete
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
			Error: readBatchErrorFixup(err), // or a specialized "deleteBatchErrorFixup"
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

		// It's typical to not return a "version" if there's no concurrency check,
		// but you can do so if you track version changes. Here we keep it simple.
		success = append(success, persistencemodel.WriteResult{
			ID: op.ID,
			// No Error => success
		})
	}

	return success
}
