package dynamodbpersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

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
				Version:   reported.Version,
				TimeStamp: reportedObj.TimeStamp,
				Error:     err,
			},
			{
				ID:        desired.ID,
				Version:   desired.Version,
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

	reportedVersion := reportedObj.Version
	desiredVersion := desiredObj.Version

	if err != nil {
		reportedVersion = reported.Version
		desiredVersion = desired.Version
	}

	return []persistencemodel.WriteResult{
		{
			ID:        reported.ID,
			Version:   reportedVersion,
			TimeStamp: reportedObj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
		{
			ID:        desired.ID,
			Version:   desiredVersion,
			TimeStamp: desiredObj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
	}
}
