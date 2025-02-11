package dynamodbpersistence

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// List retrieves models according to the provided ListOptions.
//
// If opt.ID is provided, it will Query that specific partition key (DS#{ID}).
// Otherwise, it will Scan the entire table. If opt.Token is provided, it's
// used as the pagination token (ExclusiveStartKey). A new token will be
// placed in the last item’s Token field if more data is available.
func (p *Persistence) List(
	ctx context.Context,
	opt persistencemodel.ListOptions,
) (*persistencemodel.ListResults, error) {

	var (
		results      []persistencemodel.ListResult
		exclusiveKey map[string]types.AttributeValue
		pageSize     int32 = 100
		token        string
		err          error
	)

	if p.config.MaxReadBatchSize > 0 {
		pageSize = int32(p.config.MaxReadBatchSize)
	}

	// token -> decode it into ExclusiveStartKey
	if opt.Token != "" {
		exclusiveKey, err = p.decodeToken(opt.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to decode token: %w", err)
		}
	}

	// ID -> Query, no ID -> Scan
	if opt.ID != "" {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(p.config.Table),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "DS#" + opt.ID},
			},
			ExclusiveStartKey: exclusiveKey,
		}

		out, err := p.client.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}

		results, err = p.parseListResponse(out.Items)
		if err != nil {
			return nil, err
		}

		// more pages -> new token
		if len(out.LastEvaluatedKey) > 0 && len(results) > 0 {
			tok, err := p.encodeToken(out.LastEvaluatedKey)

			if err != nil {
				return nil, fmt.Errorf("encodeToken failed: %w", err)
			}

			token = tok
		}

	} else {
		// Scan all items
		input := &dynamodb.ScanInput{
			TableName:         aws.String(p.config.Table),
			ExclusiveStartKey: exclusiveKey,
			Limit:             &pageSize,
		}

		out, err := p.client.Scan(ctx, input)

		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		results, err = p.parseListResponse(out.Items)

		if err != nil {
			return nil, err
		}

		// If more pages exist -> attach token
		if len(out.LastEvaluatedKey) > 0 && len(results) > 0 {
			tok, err := p.encodeToken(out.LastEvaluatedKey)

			if err != nil {
				return nil, fmt.Errorf("encodeToken failed: %w", err)
			}

			token = tok
		}
	}

	return &persistencemodel.ListResults{
		Items: results,
		Token: token,
	}, nil
}

// parseListResponse converts a slice of DynamoDB items into ListResult entries.
//
// Each item should have:
//   - PK = "DS#{ID}"
//   - SK = "DSR#{Name}" or "DSD#{Name}" or "DSC#{Name}" etc.
//   - The rest of the fields (version, timestamp, clientToken, desired, reported)
//     are in the struct PersistenceObject (unmarshaled from the item).
//
// partialObject is a minimal struct to store only the fields you want to retrieve.
func (p *Persistence) parseListResponse(
	items []map[string]types.AttributeValue,
) ([]persistencemodel.ListResult, error) {
	//
	results := make([]persistencemodel.ListResult, 0, len(items))

	for _, it := range items {
		// Partition Key
		pkAttr, ok := it["PK"].(*types.AttributeValueMemberS)
		if !ok || !strings.HasPrefix(pkAttr.Value, "DS#") {
			// Skip if PK is missing or doesn't match DS#
			continue
		}

		id := pkAttr.Value[3:]

		// Sort Key
		skAttr, ok := it["SK"].(*types.AttributeValueMemberS)

		if !ok {
			continue
		}

		sk := skAttr.Value

		var (
			name      string
			modelType persistencemodel.ModelType
			combined  bool
		)
		switch {
		case strings.HasPrefix(sk, "DSR#"):
			name = sk[4:]
			modelType = persistencemodel.ModelTypeReported
		case strings.HasPrefix(sk, "DSD#"):
			name = sk[4:]
			modelType = persistencemodel.ModelTypeDesired
		case strings.HasPrefix(sk, "DSC#"):
			name = sk[4:]
			combined = true
		default:
			// Skip unknown SK patterns
			continue
		}

		var obj PartialPersistenceObject

		if err := attributevalue.UnmarshalMap(it, &obj); err != nil {
			return nil, fmt.Errorf("failed to unmarshal partialObject: %w", err)
		}

		// If "combined", we generate two results (desired + reported).
		// Otherwise we generate one result based on modelType.
		if combined {
			// Desired
			results = append(results, persistencemodel.ListResult{
				ID: persistencemodel.PersistenceID{
					ID:        id,
					Name:      name,
					ModelType: persistencemodel.ModelTypeDesired,
				},
				Version:     obj.Version,
				TimeStamp:   obj.TimeStamp,
				ClientToken: obj.ClientToken,
			})
			// Reported
			results = append(results, persistencemodel.ListResult{
				ID: persistencemodel.PersistenceID{
					ID:        id,
					Name:      name,
					ModelType: persistencemodel.ModelTypeReported,
				},
				Version:     obj.Version,
				TimeStamp:   obj.TimeStamp,
				ClientToken: obj.ClientToken,
			})
		} else {
			results = append(results, persistencemodel.ListResult{
				ID: persistencemodel.PersistenceID{
					ID:        id,
					Name:      name,
					ModelType: modelType,
				},
				Version:     obj.Version,
				TimeStamp:   obj.TimeStamp,
				ClientToken: obj.ClientToken,
			})
		}
	}

	return results, nil
}

// encodeToken converts a DynamoDB LastEvaluatedKey map into a string.
func (p *Persistence) encodeToken(keys map[string]types.AttributeValue) (string, error) {
	if len(keys) == 0 {
		return "", nil
	}

	raw, err := json.Marshal(keys)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(raw), nil
}

// decodeToken does the reverse: converts a string back into a map[string]types.AttributeValue.
func (p *Persistence) decodeToken(tok string) (map[string]types.AttributeValue, error) {
	if tok == "" {
		return nil, nil
	}

	b, err := base64.StdEncoding.DecodeString(tok)

	if err != nil {
		return nil, err
	}

	var first map[string]*types.AttributeValueMemberS

	if err := json.Unmarshal(b, &first); err != nil {
		return nil, err
	}

	out := make(map[string]types.AttributeValue, len(first))

	for k, v := range first {
		out[k] = v
	}

	return out, nil
}
