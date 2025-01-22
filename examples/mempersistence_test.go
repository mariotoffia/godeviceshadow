package examples

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	persistence := mempersistence.New(mempersistence.PersistenceOpts{
		Separation: persistencemodel.SeparateModels, // <1>
	})

	ctx := context.Background()

	writeRes := persistence.Write(ctx, // <2>
		persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID: "device123", Name: "HomeHub", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]IndoorTemperatureSensor{
				"temperature": {
					Floor:       1,
					Direction:   DirectionNorth,
					Temperature: 23.5,
					Humidity:    45.5,
					UpdatedAt:   time.Now(),
				},
			},
		})

	assert.Len(t, writeRes, 1)
	assert.NoError(t, writeRes[0].Error)

	res := persistence.Read(ctx, // <3>
		persistencemodel.ReadOptions{},
		persistencemodel.ReadOperation{
			ID: persistencemodel.PersistenceID{
				ID: "device123", Name: "HomeHub", ModelType: persistencemodel.ModelTypeReported,
			},
		})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)
	assert.NotNil(t, res[0].Model)

	temp := res[0].Model.(map[string]IndoorTemperatureSensor)["temperature"]
	assert.Equal(t, 1, temp.Floor)
	assert.Equal(t, DirectionNorth, temp.Direction)
	assert.Equal(t, 23.5, temp.Temperature)
	assert.Equal(t, 45.5, temp.Humidity)

	deleteRes := persistence.Delete(ctx, // <4>
		persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID: "device123", Name: "HomeHub", ModelType: persistencemodel.ModelTypeReported,
			},
		})

	assert.Len(t, deleteRes, 1)
	assert.NoError(t, deleteRes[0].Error)

	res = persistence.Read(ctx, // <5>
		persistencemodel.ReadOptions{},
		persistencemodel.ReadOperation{
			ID: persistencemodel.PersistenceID{
				ID: "device123", Name: "HomeHub", ModelType: persistencemodel.ModelTypeReported,
			},
		})

	assert.Len(t, res, 1)
	assert.Error(t, res[0].Error, "Read operation should return an error for a deleted model")
	assert.Equal(t, 404, res[0].Error.(persistencemodel.PersistenceError).Code)
}
