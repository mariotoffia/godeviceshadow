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

	results := make(map[string]persistencemodel.ReadResult, len(operations))

	table := p.config.Table
	maxBatchSize := 100

	if p.config.MaxReadBatchSize > 0 {
		maxBatchSize = p.config.MaxReadBatchSize
	}

	items := prepareRead(operations, results, table, maxBatchSize)

	for _, req := range items {
		res, err := p.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: req.Keys,
		})

		if err != nil {
			// all have errors
			for _, op := range req.Operations {
				results[op.ID.String()] = persistencemodel.ReadResult{
					ID:    op.ID,
					Error: readBatchErrorFixup(err),
				}
			}

			continue
		}

		// Process retrieved items
		retrieved := res.Responses[table]

		for _, item := range retrieved {
			// Extract PK & SK from the item
			id := sortKeyToName(item["PK"].(*types.AttributeValueMemberS).Value)
			name := primaryKeyToID(item["SK"].(*types.AttributeValueMemberS).Value)
			op := req.FromID(id, name)

			var stored PartialPersistenceObject

			if err := attributevalue.UnmarshalMap(item, &stored); err != nil {
				results[op.ID.String()] = persistencemodel.ReadResult{
					ID: persistencemodel.PersistenceID{
						ID:   op.ID.ID,
						Name: name,
					},
					Error: fmt.Errorf("unmarshal persist object failed: %w", err),
				}

				continue
			}

			// Ensure version matches -> 404
			if op.Version > 0 && op.Version != stored.Version {
				results[op.ID.String()] = persistencemodel.ReadResult{
					ID: persistencemodel.PersistenceID{ID: op.ID.ID, Name: name},
					Error: persistencemodel.Error404(
						fmt.Sprintf("mismatching version, requested: %d, stored: %d", op.Version, stored.Version),
					),
				}

				continue
			}

			if item["Desired"] != nil {
				if res, err := unmarshalFromMap(item["Desired"], op.Model); err != nil {
					results[op.ID.String()] = persistencemodel.ReadResult{
						ID:    op.ID.ToPersistenceID(persistencemodel.ModelTypeDesired),
						Error: fmt.Errorf("unmarshal desired failed: %w", err),
					}
				} else {
					results[op.ID.String()] = persistencemodel.ReadResult{
						ID:          op.ID.ToPersistenceID(persistencemodel.ModelTypeDesired),
						Model:       res,
						Version:     stored.Version,
						TimeStamp:   stored.TimeStamp,
						ClientToken: stored.ClientToken,
					}
				}
			}

			if item["Reported"] != nil {
				if res, err := unmarshalFromMap(item["Reported"], op.Model); err != nil {
					results[op.ID.String()] = persistencemodel.ReadResult{
						ID:    op.ID.ToPersistenceID(persistencemodel.ModelTypeReported),
						Error: fmt.Errorf("unmarshal reported failed: %w", err),
					}
				} else {
					results[op.ID.String()] = persistencemodel.ReadResult{
						ID:          op.ID.ToPersistenceID(persistencemodel.ModelTypeReported),
						Model:       res,
						Version:     stored.Version,
						TimeStamp:   stored.TimeStamp,
						ClientToken: stored.ClientToken,
					}
				}
			}
		}
	}

	res := make([]persistencemodel.ReadResult, 0, len(operations))

	for _, r := range results {
		res = append(res, r)
	}

	return res
}

type ReadBatch struct {
	Operations []persistencemodel.ReadOperation
	Keys       map[string]types.KeysAndAttributes
}

func (rb *ReadBatch) FromID(pk, sk string) persistencemodel.ReadOperation {
	for _, op := range rb.Operations {
		if op.ID.ID == pk && op.ID.Name == sk {
			return op
		}
	}

	return persistencemodel.ReadOperation{}
}

// prepareRead will prepare the items to read from the database. The _results_ map is used to store any
// errors that did occur during the preparation.
func prepareRead(operations []persistencemodel.ReadOperation, results map[string]persistencemodel.ReadResult, table string, maxBatchSize int) []ReadBatch {
	batches := utils.ToBatch(operations, maxBatchSize)
	items := make([]ReadBatch, 0, len(batches))

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
				results[op.ID.String()] = persistencemodel.ReadResult{
					ID:    op.ID,
					Error: persistencemodel.Error400("invalid model type"),
				}

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

	return items
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
