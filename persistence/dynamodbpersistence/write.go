package dynamodbpersistence

import (
	"context"
	"errors"
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
	conditionWriteExpression = "attribute_not_exists(version) OR version = :expected_version"
	expectedVersionValueKey  = ":expected_version"
)

func (p *Persistence) Write(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	// Use default separation
	sep := p.config.ModelSeparation

	// Do we have a write operation override?
	if os, ok := persistutils.FromConfig[persistencemodel.ModelSeparation](opt.Config, persistencemodel.ModelSeparationConfigKey); ok {
		sep = os
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
		var gc int

		for i := range groups {
			if groups[i].Error != nil {
				for j := range groups[i].Operations {
					res = append(res, persistencemodel.WriteResult{
						ID:      groups[i].Operations[j].ID,
						Version: groups[i].Operations[j].Version,
						Error:   groups[i].Error,
					})

					gc++
				}

				continue
			}

			// Write the operation group
			wr := p.WriteOperationGroup(ctx, opt, groups[i])

			res = append(res, make([]persistencemodel.WriteResult, len(groups[i].Operations))...)
			gc += len(wr)
		}
	} else {
		// Parallel execution
		return p.parallelWrite(ctx, opt, groups, maxParallelism)
	}

	return res
}

// WriteOperationGroup will write a single operation group to the persistence layer. It will return as many results as
// there are operations in the group.
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

// writeCombined will combine both reported and desired models into a single persistence object and write it to the
// persistence layer. It will only do one conditional check on version and write the increased version in same operation.
func (p *Persistence) writeCombined(
	ctx context.Context,
	group persistutils.GroupedWriteOperation,
) []persistencemodel.WriteResult {
	// Determine PK and SK for combined models
	pk := fmt.Sprintf("DS#%s", group.ID)
	sk := fmt.Sprintf("DSC#%s", group.Name)

	reported := group.GetByModelType(persistencemodel.ModelTypeReported)
	desired := group.GetByModelType(persistencemodel.ModelTypeDesired)

	// Create PersistenceObject with both Desired and Reported
	now := time.Now().UTC().UnixNano()

	obj := PersistenceObject{
		Version:     reported.Version + 1,
		TimeStamp:   now,
		ClientToken: reported.ClientID,
		Desired:     desired.Model,
		Reported:    desired.Model,
	}

	// Perform conditional write
	err := p.dynamoDbPut(ctx, pk, sk, obj, reported.Version)

	return []persistencemodel.WriteResult{
		{
			ID:        reported.ID,
			Version:   obj.Version,
			TimeStamp: obj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
		{
			ID:        desired.ID,
			Version:   obj.Version,
			TimeStamp: obj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
	}
}

// writeSeparate will write reported and desired models separately to the persistence layer. It will use a transaction
// if both models are present, otherwise it will use a single conditional write.
func (p *Persistence) writeSeparate(
	ctx context.Context,
	group persistutils.GroupedWriteOperation,
) []persistencemodel.WriteResult {
	reported := group.GetByModelType(persistencemodel.ModelTypeReported)
	desired := group.GetByModelType(persistencemodel.ModelTypeDesired)

	if reported != nil && desired != nil {
		// Both reported and desired models present, use a transaction
		return p.writeTransactional(ctx, group, reported, desired)
	}

	// Only one model is present -> use "plain" conditional writes
	results := make([]persistencemodel.WriteResult, 0, len(group.Operations))
	if reported != nil {
		results = append(results, p.writeSingle(ctx, reported, persistencemodel.ModelTypeReported))
	}
	if desired != nil {
		results = append(results, p.writeSingle(ctx, desired, persistencemodel.ModelTypeDesired))
	}

	return results
}

// writeTransactional performs a transactional write for both reported and desired models.
func (p *Persistence) writeTransactional(
	ctx context.Context,
	group persistutils.GroupedWriteOperation,
	reported, desired *persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	pk := fmt.Sprintf("DS#%s", group.ID)

	reportedKey := fmt.Sprintf("DSR#%s", group.Name)
	desiredKey := fmt.Sprintf("DSD#%s", group.Name)

	now := time.Now().UTC().UnixNano()

	reportedObj := PersistenceObject{
		Version:     reported.Version + 1,
		TimeStamp:   now,
		ClientToken: reported.ClientID,
		Reported:    reported.Model,
	}
	desiredObj := PersistenceObject{
		Version:     desired.Version + 1,
		TimeStamp:   now,
		ClientToken: desired.ClientID,
		Desired:     desired.Model,
	}

	reportItem, err := marshalDynamoDBItem(reportedKey, pk, reportedObj)
	desiredItem, err2 := marshalDynamoDBItem(desiredKey, pk, desiredObj)

	if err == nil && err2 != nil {
		err = err2
	}

	if err != nil {
		return []persistencemodel.WriteResult{
			{
				ID:        reported.ID,
				Version:   reportedObj.Version,
				TimeStamp: reportedObj.TimeStamp,
				Error:     err,
			},
			{
				ID:        desired.ID,
				Version:   desiredObj.Version,
				TimeStamp: desiredObj.TimeStamp,
				Error:     err,
			},
		}
	}

	// Build transaction input
	transactions := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName:           aws.String(p.config.Table),
				Item:                reportItem,
				ConditionExpression: aws.String(conditionWriteExpression),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					expectedVersionValueKey: &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", reported.Version)},
				},
			},
		},
		{
			Put: &types.Put{
				TableName:           aws.String(p.config.Table),
				Item:                desiredItem,
				ConditionExpression: aws.String(conditionWriteExpression),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					expectedVersionValueKey: &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", desired.Version)},
				},
			},
		},
	}

	_, err = p.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactions,
	})

	return []persistencemodel.WriteResult{
		{
			ID:        reported.ID,
			Version:   reportedObj.Version,
			TimeStamp: reportedObj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
		{
			ID:        desired.ID,
			Version:   desiredObj.Version,
			TimeStamp: desiredObj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
	}
}

// writeSingle performs a single conditional put for either reported or desired models.
func (p *Persistence) writeSingle(
	ctx context.Context,
	op *persistencemodel.WriteOperation,
	modelType persistencemodel.ModelType,
) persistencemodel.WriteResult {
	pk := fmt.Sprintf("DS#%s", op.ID.ID)
	var sk string

	if modelType == persistencemodel.ModelTypeReported {
		sk = fmt.Sprintf("DSR#%s", op.ID.Name)
	} else {
		sk = fmt.Sprintf("DSD#%s", op.ID.Name)
	}

	now := time.Now().UTC().UnixNano()

	obj := PersistenceObject{
		Version:     op.Version + 1,
		TimeStamp:   now,
		ClientToken: op.ClientID,
	}

	if modelType == persistencemodel.ModelTypeReported {
		obj.Reported = op.Model
	} else {
		obj.Desired = op.Model
	}

	err := p.dynamoDbPut(ctx, pk, sk, obj, op.Version)

	return persistencemodel.WriteResult{
		ID:        op.ID,
		Version:   obj.Version,
		TimeStamp: obj.TimeStamp,
		Error:     conditionalWriteErrorFixup(err),
	}
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
	condition := conditionWriteExpression
	expressionValues := map[string]types.AttributeValue{
		expectedVersionValueKey: &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expectedVersion)},
	}

	// Execute the conditional write
	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 aws.String(p.config.Table),
		Item:                      item,
		ConditionExpression:       aws.String(condition),
		ExpressionAttributeValues: expressionValues,
	})
	if err != nil {
		return err
	}

	return nil
}

func marshalDynamoDBItem(sk, pk string, obj PersistenceObject) (map[string]types.AttributeValue, error) {
	item, err := attributevalue.MarshalMap(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item: %w", err)
	}
	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}
	return item, nil
}

// conditionalWriteErrorFixup will check the error and correct it to http error if appropriate. This is only
// for conditional writes.
func conditionalWriteErrorFixup(err error) error {
	if err == nil {
		return nil
	}

	// Handle ConditionalCheckFailedException
	var conditionalErr *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalErr) {
		return persistencemodel.Error409("conditional check failed")
	}

	// Handle TransactionCanceledException
	var transactionErr *types.TransactionCanceledException
	if errors.As(err, &transactionErr) {
		msg := "transaction canceled: conditional check failed"
		if transactionErr.Message != nil {
			msg = *transactionErr.Message
		}
		return persistencemodel.Error409(msg)
	}

	// Return the original error if not handled
	return err
}
