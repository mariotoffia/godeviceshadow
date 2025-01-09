package dynamodbpersistence

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func toPartitionKey(id persistencemodel.PersistenceID) string {
	return fmt.Sprintf("DS#%s", id.ID)
}

func toSortKey(id persistencemodel.PersistenceID, mt persistencemodel.ModelType) string {
	if mt == persistencemodel.ModelTypeReported {
		return fmt.Sprintf("DSR#%s", id.Name)
	} else if mt == persistencemodel.ModelTypeDesired {
		return fmt.Sprintf("DSD#%s", id.Name)
	}

	return fmt.Sprintf("unknown model type: %s", mt.String())
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
