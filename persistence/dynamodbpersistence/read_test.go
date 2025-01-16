package dynamodbpersistence_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tz = "Europe/Stockholm"

func TestReadUnversionedCombined(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
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
				Separation: persistencemodel.CombinedModels, // store them combined
			},
		},
		persistencemodel.WriteOperation{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: tz,
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

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

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

	if model.TimeZone != "" {
		assert.Equal(t, tz, model.TimeZone)
		assert.Contains(t, model.Sensors, "temp")

		sensor, ok := model.Sensors["temp"]
		require.True(t, ok)

		assert.Equal(t, 21.5, sensor.Value)
	} else {
		model, ok := read[1].Model.(*TestModel)
		require.True(t, ok)

		assert.Equal(t, tz, model.TimeZone)
		assert.Contains(t, model.Sensors, "temp")

		sensor, ok := model.Sensors["temp"]
		require.True(t, ok)

		assert.Equal(t, 21.5, sensor.Value)
	}
}

func TestReadIncorrectVersionIsConflictSeparateModels(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
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
				Separation: persistencemodel.SeparateModels,
			},
		},
		persistencemodel.WriteOperation{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"temp": {Value: 21.5, TimeStamp: time.Now().UTC()},
				},
			},
		},
	)

	require.Len(t, operations, 1)

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported}, // <- reported!!
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 2, /*conditional read -> incorrect version -> fail */
	})

	require.Len(t, read, 1)
	assert.Error(t, read[0].Error)
	assert.ErrorAs(t, read[0].Error, &persistencemodel.PersistenceError{})
	assert.Equal(t, read[0].Error.(persistencemodel.PersistenceError).Code, 409)
}

func TestReadIncorrectVersionIsConflictCombinedModels(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
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
				TimeZone: tz,
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

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 2, /*conditional read -> incorrect version -> fail */
	})

	require.Len(t, read, 1)
	assert.Error(t, read[0].Error)
	assert.ErrorAs(t, read[0].Error, &persistencemodel.PersistenceError{})
	assert.Equal(t, read[0].Error.(persistencemodel.PersistenceError).Code, 409)
}

func TestReadCorrectVersionSuccessSeparateModels(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
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
				Separation: persistencemodel.SeparateModels,
			},
		},
		persistencemodel.WriteOperation{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"temp": {Value: 21.5, TimeStamp: time.Now().UTC()},
				},
			},
		},
	)

	require.Len(t, operations, 1)

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported}, // <- reported!!
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 1, /*conditional read -> correct version -> success */
	})

	require.Len(t, read, 1)
	assert.NoError(t, read[0].Error)
}

func TestReadCorrectVersionSuccessCombinedModels(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
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
				TimeZone: tz,
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

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 1, /*conditional read -> correct version -> success */
	})

	require.Len(t, read, 2)
	assert.NoError(t, read[0].Error)
	assert.NoError(t, read[1].Error)
}

func TestReadCorrectVersionSuccessCombinedModelsUpdate(t *testing.T) {
	ctx := context.TODO()

	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	p, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:  res.Table,
		Client: res.Client,
	})
	require.NoError(t, err)

	clientID := persistutils.Id("test-")
	writeOperations := []persistencemodel.WriteOperation{
		{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"temp": {Value: 21.5, TimeStamp: time.Now().UTC()},
				},
			},
		},
		{
			ClientID: clientID,
			ID: persistencemodel.PersistenceID{
				ID: "deviceA", Name: "shadowA", ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: TestModel{},
		},
	}

	operations := p.Write(
		ctx,
		persistencemodel.WriteOptions{
			Config: persistencemodel.WriteConfig{
				Separation: persistencemodel.CombinedModels,
			},
		},
		writeOperations...,
	)

	require.Len(t, operations, 2)

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read := p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 1, /*conditional read -> correct version -> success */
	})

	require.Len(t, read, 2)
	assert.NoError(t, read[0].Error)
	assert.NoError(t, read[1].Error)

	writeOperations[0].Version = 1
	writeOperations[1].Version = 1

	operations = p.Write(
		ctx,
		persistencemodel.WriteOptions{
			Config: persistencemodel.WriteConfig{
				Separation: persistencemodel.CombinedModels,
			},
		},
		writeOperations...,
	)

	require.Len(t, operations, 2)

	for _, w := range operations {
		require.NoError(t, w.Error)
	}

	read = p.Read(ctx, persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Model:   reflect.TypeOf(&TestModel{}),
		Version: 2, /*conditional read -> correct version -> success */
	})

	require.Len(t, read, 2)
	assert.NoError(t, read[0].Error)
	assert.NoError(t, read[1].Error)

}
