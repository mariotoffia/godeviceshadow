package dynamodbpersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

// Just for sonar-lint complaints
const (
	// conditionWriteExpression that is used for conditional writes on tag _version_. If change in `PersistenceObject.Version`
	// json tag name, this needs to be updated.
	conditionWriteExpression = "attribute_not_exists(version) OR version = :expected_version"
	// expectedVersionValueKey is the key used in the condition expression for the expected version value.
	expectedVersionValueKey = ":expected_version"
)

func (p *Persistence) Write(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	// Use default separation
	sep := p.config.ModelSeparation

	// Do we have a write operation override?
	if opt.Config.Separation != 0 {
		sep = opt.Config.Separation
	}

	groups := persistutils.Group(operations)

	for i := range groups {
		os := sep

		// Override on a per group basis
		if groups[i].ModelSeparation == 0 {
			groups[i].ModelSeparation = os
		}

		groups[i].Error = persistutils.Validate(groups[i])
	}

	maxParallelism := p.config.MaxParallelism
	if maxParallelism <= 0 {
		maxParallelism = 1
	}

	res := make([]persistencemodel.WriteResult, 0, len(operations))

	if maxParallelism == 1 {
		// Single thread
		for i := range groups {
			if groups[i].Error != nil {
				for j := range groups[i].Operations {
					res = append(res, persistencemodel.WriteResult{
						ID:      groups[i].Operations[j].ID,
						Version: groups[i].Operations[j].Version,
						Error:   groups[i].Error,
					})
				}

				continue
			}

			wr := p.WriteOperationGroup(ctx, opt, groups[i])
			res = append(res, wr...)
		}
	} else {
		// Parallel
		return p.writeParallel(ctx, opt, groups, maxParallelism)
	}

	return res
}

// WriteOperationGroup will write a single operation group to the persistence layer. It will return as many results as
// there are operations in the group.
//
// CAUTION: Verify the integrity of the _group_ before calling this function since it will not validate the group and
// the operations. Use `Write` instead!
func (p *Persistence) WriteOperationGroup(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	group persistutils.GroupedWriteOperation,
) []persistencemodel.WriteResult {
	if ctx.Err() != nil {
		return []persistencemodel.WriteResult{
			{Error: fmt.Errorf("context cancelled: %w", ctx.Err())},
		}
	}

	switch group.ModelSeparation {
	case persistencemodel.CombinedModels:
		return p.writeCombined(ctx, group)
	case persistencemodel.SeparateModels:
		return p.writeSeparate(ctx, group)
	}

	res := make([]persistencemodel.WriteResult, 0, len(group.Operations))

	for _, op := range group.Operations {
		res = append(res, persistencemodel.WriteResult{
			ID:        op.ID,
			Version:   op.Version,
			TimeStamp: time.Now().UTC().UnixNano(),
			Error:     persistencemodel.Error404(fmt.Sprintf("ModelSeparation '%s' is not supported", group.ModelSeparation.String())),
		})
	}

	return res
}

// dynamoDbPut performs a conditional write to DynamoDB.
func (p *Persistence) dynamoDbPut(
	ctx context.Context,
	pk string,
	sk string,
	obj PersistenceObject,
	expectedVersion int64,
) error {
	// Construct the DynamoDB item
	item, err := attributevalue.MarshalMap(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	// Add PK and SK attributes
	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}

	// Conditional expression for version matching
	expressionValues := map[string]types.AttributeValue{
		expectedVersionValueKey: &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expectedVersion)},
	}

	// Execute the conditional write
	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 aws.String(p.config.Table),
		Item:                      item,
		ConditionExpression:       aws.String(conditionWriteExpression),
		ExpressionAttributeValues: expressionValues,
	})

	if err != nil {
		return err
	}

	return nil
}
