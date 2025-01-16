package dynamodbpersistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestTableName = "go-deviceshadow-test"

type Sensor struct {
	Value     any
	TimeStamp time.Time
}

type TestModel struct {
	TimeZone string
	Sensors  map[string]Sensor
}

func TestListEmpty(t *testing.T) {
	ctx := context.TODO()
	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)

	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	persistence, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:  res.Table,
		Client: res.Client,
	})

	require.NoError(t, err)

	results, err := persistence.List(ctx, persistencemodel.ListOptions{
		ID: "test",
	})

	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestListSingleCombinedDeviceShadow(t *testing.T) {
	ctx := context.TODO()
	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)

	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	persistence, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table: res.Table, Client: res.Client,
	})

	require.NoError(t, err)

	writes := persistence.Write(ctx, persistencemodel.WriteOptions{
		Config: persistencemodel.WriteConfig{
			Separation: persistencemodel.CombinedModels,
		}},
		persistencemodel.WriteOperation{ // Reported
			ClientID: persistutils.Id("test"),
			ID:       persistencemodel.PersistenceID{ID: "test", Name: "test-ds", ModelType: persistencemodel.ModelTypeReported},
			Model:    TestModel{TimeZone: "Europe/Stockholm", Sensors: map[string]Sensor{"temp": {Value: 23.4, TimeStamp: time.Now().UTC()}}},
		},
		persistencemodel.WriteOperation{ // Desired
			ClientID: persistutils.Id("test"),
			ID:       persistencemodel.PersistenceID{ID: "test", Name: "test-ds", ModelType: persistencemodel.ModelTypeDesired},
			Model:    TestModel{},
		})

	require.Len(t, writes, 2)
	require.NoError(t, writes[0].Error)
	require.NoError(t, writes[1].Error)

	results, err := persistence.List(ctx, persistencemodel.ListOptions{
		ID: "test",
	})

	require.NoError(t, err)
	assert.Len(t, results, 2)
}
