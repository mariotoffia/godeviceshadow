package dynamodbpersistence

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/rand"

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

func readBatchErrorFixup(err error) error {
	if err == nil {
		return nil
	}

	var resourceNotFound *types.ResourceNotFoundException

	if errors.As(err, &resourceNotFound) {
		return persistencemodel.Error404("item not found")
	}

	return err
}

func isMapValue(m map[string]types.AttributeValue, key string) bool {
	if m, ok := m[key]; ok && m != nil {
		if des, ok := m.(*types.AttributeValueMemberM); ok && des != nil {
			return true
		}
	}

	return false
}

func unmarshalFromMap(m types.AttributeValue, t reflect.Type) (any, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	value := reflect.New(t)

	if !value.IsValid() {
		return nil, fmt.Errorf("invalid type: %s", t.String())
	}

	inst := value.Interface()

	if des, ok := m.(*types.AttributeValueMemberM); ok && des != nil {
		if err := attributevalue.UnmarshalMap(des.Value, inst); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		} else {
			return inst, nil
		}
	}

	return nil, fmt.Errorf("expected AttributeValueMemberM but got: %T", m)
}

// diffWriteRequest returns the items in `write` that are NOT in `unprocessed`.
func diffWriteRequest(write, unprocessed []types.WriteRequest) []types.WriteRequest {
	if len(unprocessed) == 0 {
		return write
	}

	subsetMap := make(map[string]bool, len(unprocessed))

	for _, wr := range unprocessed {
		key := writeRequestToID(wr)
		subsetMap[key] = true
	}

	// For each item in all, if it's not in subsetMap, it was processed
	var diff []types.WriteRequest

	for _, wr := range write {
		key := writeRequestToID(wr)
		if !subsetMap[key] {
			diff = append(diff, wr)
		}
	}

	return diff
}

// writeRequestToID is a small helper that builds e.g. "PK=<val>,SK=<val>"
// to uniquely identify the request in a map
func writeRequestToID(wr types.WriteRequest) string {
	if wr.DeleteRequest == nil {
		return ""
	}
	pkVal := ""
	skVal := ""

	if pk, ok := wr.DeleteRequest.Key["PK"].(*types.AttributeValueMemberS); ok {
		pkVal = pk.Value
	}
	if sk, ok := wr.DeleteRequest.Key["SK"].(*types.AttributeValueMemberS); ok {
		skVal = sk.Value
	}
	return fmt.Sprintf("PK=%s,SK=%s", pkVal, skVal)
}

// exponentialBackoff sleeps for an exponentially growing duration based on `attempt`.
// It also includes a small random jitter (to avoid thundering herd), and respects context cancellation.
func exponentialBackoff(ctx context.Context, attempt int) error {
	backoff := (time.Duration(1<<attempt) * 100 * time.Millisecond) +
		(time.Duration(rand.Int63n(int64(50 * time.Millisecond))))

	// Cap the maximum backoff at 30s
	if backoff > 30*time.Second {
		backoff = 30 * time.Second
	}

	// Use a select to wait or cancel
	select {
	case <-ctx.Done():
		// Context was cancelled or deadline exceeded
		return ctx.Err()
	case <-time.After(backoff):
		// Successfully slept the entire backoff duration
		return nil
	}
}
