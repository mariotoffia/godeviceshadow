//go:build integration
// +build integration

package dynamodbnotifier

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"

	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/notify"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"
	"github.com/mariotoffia/godeviceshadow/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestTableName = "go-deviceshadow-test"
const tz = "Europe/Stockholm"

type Sensor struct {
	Value     any
	TimeStamp time.Time
}

type TestModel struct {
	TimeZone string
	Sensors  map[string]Sensor
}

func (sp *Sensor) GetTimestamp() time.Time {
	return sp.TimeStamp
}

func (sp *Sensor) GetValue() any {
	return sp.Value
}

func TestReceiveOneEvent(t *testing.T) {
	// DynamoDB test utility
	l, err := dynamodbutils.StartLocalDynamoDB(context.Background(), TestTableName, dynamodbutils.StartLocalDynamoDbOptions{
		Reuse: false,
	})
	require.NoError(t, err)

	defer l.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var (
		oldImage, newImage *PersistenceObject
		pcError            error
		wg                 sync.WaitGroup
		res                []notifiermodel.NotificationTargetResult
	)

	processor, err := NewProcessorBuilder().
		WithTypeRegistry(types.NewRegistry().
			RegisterResolver(
				model.NewResolveFunc(
					func(id, name string) (model.TypeEntry, bool) {
						return model.TypeEntry{
							Model: reflect.TypeOf(&TestModel{}), Name: "homeHub"}, true
					},
				),
			)).
		WithStartDoneFunc(func(ctx context.Context, err error) {
			wg.Done() // signal that the processor is done
		}).
		WithEventModel(ProcessorEventModelManualAttach).
		WithStreamBuilder(
			stream.NewStreamBuilder().
				WithTable(TestTableName).
				// IMPORTANT: use the local dynamodb client
				UseClient(l.Client).
				UseStreamClient(l.StreamClient).
				WithShardIteratorType(streamTypes.ShardIteratorTypeLatest),
		).
		WithCallbackProcessorCallback().
		WithProcessImageCallback(func(ctx context.Context, oldImg, newImg *PersistenceObject, err error) error {
			oldImage = oldImg
			newImage = newImg
			pcError = err

			return nil
		}).
		WithNotificationManager(notify.NewBuilder().
			TargetBuilder(
				notifiermodel.FuncTarget(
					func(
						ctx context.Context, target notifiermodel.NotificationTarget,
						tx *persistencemodel.TransactionImpl, operation ...notifiermodel.NotifierOperation,
					) []notifiermodel.NotificationTargetResult {
						// Target could e.g. be SQS, SNS, Email, SMS, etc.
						for _, op := range operation {
							res = append(res, notifiermodel.NotificationTargetResult{
								Operation: op,
								Target:    target,
								Custom:    map[string]any{"pass": true},
							})
						}

						return res
					})).
			WithSelectionBuilder(
				notifiermodel.NewSelectionBuilder(
					notifiermodel.FuncSelection(
						func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
							return true, nil // always include
						})),
			).
			Build().
			Build()).
		Build(context.Background())

	require.NoError(t, err)

	err = processor.Start(ctx, true /*async*/)
	require.NoError(t, err)

	time.Sleep(5 * time.Second) // wait for the processor to start

	// Write a PK: DS#{id}, SK: DSR#{name} item using the persistence object
	// to simulate a new reported state for a manager (would be done in
	// dynamodbpersistence package in real life)
	reported := PersistenceObject{
		Version:     1,
		TimeStamp:   time.Now().Unix(),
		ClientToken: "clientToken",
		Reported: &TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temperature": {Value: 23.4, TimeStamp: time.Now()},
			},
		},
	}

	item, err := attributevalue.MarshalMap(reported)
	require.NoError(t, err)

	item["PK"] = &ddbtypes.AttributeValueMemberS{Value: "DS#myDevice-123"}
	item["SK"] = &ddbtypes.AttributeValueMemberS{Value: "DSR#homeHub"}

	wg.Add(1)

	_, err = l.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TestTableName),
		Item:      item,
	})

	require.NoError(t, err)

	wg.Wait() // wait for the processor to finish

	require.NoError(t, pcError)
	require.Nil(t, oldImage)
	require.NotNil(t, newImage)
	assert.Equal(t, DynamoDbEventTypeInsert, newImage.EventType())
	assert.Len(t, res, 1)
}
