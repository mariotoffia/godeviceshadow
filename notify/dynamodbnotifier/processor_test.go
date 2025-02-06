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
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"
	"github.com/mariotoffia/godeviceshadow/types"
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

// PersistenceObject is taken from _dynamodbpersistence_ package
type PersistenceObject struct {
	Version     int64  `json:"version"`
	TimeStamp   int64  `json:"timestamp"`
	ClientToken string `json:"clientToken,omitempty"`
	Desired     any    `json:"desired,omitempty"`
	Reported    any    `json:"reported,omitempty"`
}

func TestReceiveOneEvent(t *testing.T) {
	// DynamoDB test utility
	tr := dynamodbutils.NewTestTableResource(context.Background(), TestTableName)
	defer tr.Dispose(context.Background(), dynamodbutils.DisposeOpts{DeleteItems: true})

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	var (
		oldImage, newImage any
		pcError            error
		wg                 sync.WaitGroup
	)

	processor, err := NewProcessorBuilder().
		WithTypeRegistry(
			types.NewRegistry().
				RegisterResolver(
					model.NewResolveFunc(
						func(id, name string) (model.TypeEntry, bool) {
							return model.TypeEntry{
								Model: reflect.TypeOf(&TestModel{}), Name: "homeHub"}, false
						},
					),
				),
		).
		WithStartDoneFunc(func(ctx context.Context, err error) {
			wg.Done() // signal that the processor is done
		}).
		WithEventModel(ProcessorEventModelManualAttach).
		WithStreamBuilder(
			stream.NewStreamBuilder().
				WithTable(TestTableName).
				UseClient(tr.Client).
				WithShardIteratorType(streamTypes.ShardIteratorTypeTrimHorizon),
		).
		WithCallbackProcessorCallback().
		WithProcessImageCallback(func(ctx context.Context, oldImg, newImg any, err error) error {
			oldImage = oldImg
			newImage = newImg
			pcError = err

			ctx.Done() // cancel the context to stop the processor

			return nil
		}).
		Build()

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

	_, err = tr.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TestTableName),
		Item:      item,
	})
	require.NoError(t, err)

	wg.Wait() // wait for the processor to finish

	require.NoError(t, pcError)
	require.NotNil(t, oldImage)
	require.NotNil(t, newImage)
}
