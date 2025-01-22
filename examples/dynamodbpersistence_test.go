package examples

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"

	"github.com/stretchr/testify/require"
)

func TestReadUnversionedCombined(t *testing.T) {
	ctx := context.Background()

	res := dynamodbutils.NewTestTableResource(ctx, "MyTable")
	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	p, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:  res.Table,
		Client: res.Client,
	})
	require.NoError(t, err)

	clientID := persistutils.Id("test-")

	operations := p.Write(
		ctx,
		persistencemodel.WriteOptions{
			Config: persistencemodel.WriteConfig{
				Separation: persistencemodel.CombinedModels,
			},
		},
		persistencemodel.WriteOperation{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: "Europe/Stockholm",
				Sensors: map[string]Sensor{
					"temp": {Value: 21.5, TimeStamp: time.Now().UTC()},
				},
			},
		},
		persistencemodel.WriteOperation{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: TestModel{},
		},
	)

	require.Len(t, operations, 2)
	require.NoError(t, operations[0].Error)
	require.NoError(t, operations[1].Error)

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA"},
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 0, /*any -> non conditional read*/
	})

	require.Len(t, read, 2)
	require.Equal(t, "deviceA", read[0].ID.ID)
	require.Equal(t, "shadowA", read[0].ID.Name)
	require.NotNil(t, read[0].Model)

	require.Equal(t, "deviceA", read[1].ID.ID)
	require.Equal(t, "shadowA", read[1].ID.Name)
	require.NotNil(t, read[1].Model)

	model, ok := read[0].Model.(*TestModel)
	require.True(t, ok)
	require.NotNil(t, model)
}
