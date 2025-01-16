package dynamodbpersistence

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils"
)

// Read uses BatchGetItem to fetch the items for each ReadOperation.
// If the number of operations exceeds `Config.MaxReadBatchSize`, it splits them into multiple calls.
func (p *Persistence) Read(
	ctx context.Context,
	opt persistencemodel.ReadOptions,
	operations ...persistencemodel.ReadOperation,
) []persistencemodel.ReadResult {
	if len(operations) == 0 {
		return nil
	}

	results := make([]persistencemodel.ReadResult, 0, len(operations))

	table := p.config.Table
	maxBatchSize := 100
	maxRetries := 3

	if p.config.MaxReadBatchSize > 0 {
		maxBatchSize = p.config.MaxReadBatchSize
	}

	if p.config.MaxReadRetries > 0 {
		maxRetries = p.config.MaxReadRetries
	}

	items, errors := prepareRead(operations, table, maxBatchSize)
	results = append(results, errors...)

	for _, req := range items {
		results = append(results, p.readBatch(ctx, req, table, maxRetries)...)
	}

	// Check for missing items
	results = append(results, missing(operations, results)...)

	return results
}

func (p *Persistence) readBatch(ctx context.Context, batch ReadBatch, table string, maxRetries int) []persistencemodel.ReadResult {
	retrieved, errors := p.batchGetItems(ctx, batch, table, maxRetries)

	// All errors -> return them
	if len(errors) > 0 {
		return errors
	}

	// Process retrieved items
	results := make([]persistencemodel.ReadResult, 0, len(retrieved))

	for _, item := range retrieved {
		// Extract PK & SK from the item
		id := primaryKeyToID(item["PK"].(*types.AttributeValueMemberS).Value)
		name := sortKeyToName(item["SK"].(*types.AttributeValueMemberS).Value)
		op := batch.OperationFromIDName(id, name)

		var stored PartialPersistenceObject

		if err := attributevalue.UnmarshalMap(item, &stored); err != nil {
			results = append(results, persistencemodel.ReadResult{
				ID: persistencemodel.PersistenceID{
					ID: op.ID.ID, Name: name, ModelType: op.ID.ModelType,
				},
				Error: fmt.Errorf("unmarshal persist object failed: %w", err),
			})

			continue
		}

		// Ensure version matches -> 409
		if op.Version > 0 && op.Version != stored.Version {
			results = append(results, persistencemodel.ReadResult{
				ID: persistencemodel.PersistenceID{ID: op.ID.ID, Name: name, ModelType: op.ID.ModelType},
				Error: persistencemodel.Error409(
					fmt.Sprintf("mismatching version, requested: %d, stored: %d", op.Version, stored.Version),
				),
			})

			continue
		}

		if isMapValue(item, "Desired") {
			if res, err := unmarshalFromMap(item["Desired"], op.Model); err != nil {
				results = append(results, persistencemodel.ReadResult{
					ID: op.ID.ToPersistenceID(persistencemodel.ModelTypeDesired), Error: fmt.Errorf("unmarshal desired failed: %w", err),
				})
			} else {
				results = append(results, persistencemodel.ReadResult{
					ID:          op.ID.ToPersistenceID(persistencemodel.ModelTypeDesired),
					Model:       res,
					Version:     stored.Version,
					TimeStamp:   stored.TimeStamp,
					ClientToken: stored.ClientToken,
				})
			}
		}

		if isMapValue(item, "Reported") {
			if res, err := unmarshalFromMap(item["Reported"], op.Model); err != nil {
				results = append(results, persistencemodel.ReadResult{
					ID: op.ID.ToPersistenceID(persistencemodel.ModelTypeReported), Error: fmt.Errorf("unmarshal reported failed: %w", err),
				})
			} else {
				results = append(results, persistencemodel.ReadResult{
					ID:          op.ID.ToPersistenceID(persistencemodel.ModelTypeReported),
					Model:       res,
					Version:     stored.Version,
					TimeStamp:   stored.TimeStamp,
					ClientToken: stored.ClientToken,
				})
			}
		}
	}

	return results
}

// batchGetItems will fetch the items from the database in batches.
func (p *Persistence) batchGetItems(
	ctx context.Context,
	batch ReadBatch,
	table string,
	maxRetries int,
) ([]map[string]types.AttributeValue, []persistencemodel.ReadResult) {
	// the variable current will always be the ones to fetch (initially all the keys) and updated with the unprocessed keys.
	current := &dynamodb.BatchGetItemInput{
		RequestItems: batch.Keys,
	}

	if keys, ok := batch.Keys[table]; !ok {
		return nil, p.buildFailedResults(batch, nil, fmt.Errorf("no keys to fetch for table %s", table))
	} else if len(keys.Keys) == 0 {
		return nil, nil // No keys to fetch
	}

	// items to return
	items := make([]map[string]types.AttributeValue, 0, len(batch.Keys[table].Keys))

	for attempts := 1; attempts <= maxRetries; attempts++ {
		// context canceled/deadline -> bail out
		if err := ctx.Err(); err != nil {
			return items, p.buildFailedResults(batch, current.RequestItems[table].Keys, err)
		}

		res, err := p.client.BatchGetItem(ctx, current)

		if err != nil {
			// Mark all "current" keys as failed
			return items, p.buildFailedResults(batch, current.RequestItems[table].Keys, err)
		}

		if list, ok := res.Responses[table]; ok && len(list) > 0 {
			items = append(items, list...)
		}

		// no unprocessed keys -> done
		unprocessed, hasUnprocessed := res.UnprocessedKeys[table]

		if !hasUnprocessed || len(unprocessed.Keys) == 0 {
			return items, nil
		}

		// max retries -> all unprocessed keys marked as failed
		if attempts == maxRetries {
			return items, p.buildFailedResults(batch, unprocessed.Keys,
				fmt.Errorf("unprocessed keys remain after %d retries", maxRetries))
		}

		// Retry with unprocessed keys
		current.RequestItems = res.UnprocessedKeys

		if err := exponentialBackoff(ctx, attempts); err != nil {
			return items, p.buildFailedResults(batch, unprocessed.Keys, err)
		}
	}

	return items, nil
}

// buildFailedResults is used by `batchGetItems` to build the failed results.
func (p *Persistence) buildFailedResults(
	batch ReadBatch,
	keys []map[string]types.AttributeValue,
	err error,
) []persistencemodel.ReadResult {
	failed := make([]persistencemodel.ReadResult, 0, len(keys))

	for _, key := range keys {
		pkAttr, pkOk := key["PK"].(*types.AttributeValueMemberS)
		skAttr, skOk := key["SK"].(*types.AttributeValueMemberS)

		if !pkOk || !skOk {
			continue // skip invalid keys
		}

		id := primaryKeyToID(pkAttr.Value)
		name := sortKeyToName(skAttr.Value)
		op := batch.OperationFromIDName(id, name)

		failed = append(failed, persistencemodel.ReadResult{
			ID:    op.ID,
			Error: readBatchErrorFixup(err),
		})
	}

	return failed
}

type ReadBatch struct {
	Operations []persistencemodel.ReadOperation
	Keys       map[string]types.KeysAndAttributes
}

func (rb *ReadBatch) OperationFromIDName(id, name string) persistencemodel.ReadOperation {
	for _, op := range rb.Operations {
		if op.ID.ID == id && op.ID.Name == name {
			return op
		}
	}

	return persistencemodel.ReadOperation{}
}

// prepareRead will prepare the items to read from the database. The _results_ map is used to store any
// errors that did occur during the preparation.
func prepareRead(operations []persistencemodel.ReadOperation, table string, maxBatchSize int) ([]ReadBatch, []persistencemodel.ReadResult) {
	batches := utils.ToBatch(operations, maxBatchSize)
	items := make([]ReadBatch, 0, len(batches))
	var errors []persistencemodel.ReadResult

	for _, batch := range batches {
		keys := make([]map[string]types.AttributeValue, 0, len(batch))

		// Build Keys
		for _, op := range batch {
			// Partition Key
			pk := toPartitionKey(op.ID)

			// Sort Key
			var sk string

			switch op.ID.ModelType {
			case 0:
				// If ModelType is zero, assume combined model type
				sk = "DSC#" + op.ID.Name
			case persistencemodel.ModelTypeDesired:
				sk = "DSD#" + op.ID.Name
			case persistencemodel.ModelTypeReported:
				sk = "DSR#" + op.ID.Name
			default:
				errors = append(errors, persistencemodel.ReadResult{
					ID:    op.ID,
					Error: persistencemodel.Error400("invalid model type"),
				})

				continue
			}

			keys = append(keys, map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pk},
				"SK": &types.AttributeValueMemberS{Value: sk},
			})
		}

		// Items to fetch
		items = append(items, ReadBatch{
			Operations: batch,
			Keys: map[string]types.KeysAndAttributes{
				table: {
					Keys: keys,
				},
			},
		})
	}

	return items, errors
}

// parseSKName is just a utility to strip off the prefix from SK.
func sortKeyToName(sk string) string {
	if len(sk) < 4 {
		return ""
	}

	return sk[4:]
}

func primaryKeyToID(pk string) string {
	if len(pk) < 3 {
		return ""
	}

	return pk[3:]
}

// missing will return a list of ReadResults that are missing from the results.
func missing(operations []persistencemodel.ReadOperation, results []persistencemodel.ReadResult) []persistencemodel.ReadResult {
	missing := make([]persistencemodel.ReadResult, 0, len(operations))

	for _, op := range operations {
		found := false

		for _, res := range results {
			if res.ID.ID == op.ID.ID && res.ID.Name == op.ID.Name {
				if op.ID.ModelType == 0 || res.ID.ModelType == op.ID.ModelType {
					found = true
					break
				}
			}
		}

		if !found {
			missing = append(missing, persistencemodel.ReadResult{
				ID:    op.ID,
				Error: persistencemodel.Error404("not found"),
			})
		}
	}

	return missing
}
