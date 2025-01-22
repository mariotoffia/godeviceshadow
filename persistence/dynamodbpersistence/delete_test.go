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

func TestDeleteAnyVersion(t *testing.T) {
	ctx := context.Background()

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

	delete := p.Delete(ctx, persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Version: 0, /*any version*/
	})

	require.Len(t, delete, 1)
	assert.NoError(t, delete[0].Error)
}

func TestDeleteIncorrectVersion(t *testing.T) {
	ctx := context.Background()

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

	delete := p.Delete(ctx, persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Version: 14, /*not found -> error */
	})

	require.Len(t, delete, 1)
	assert.Error(t, delete[0].Error)
	assert.ErrorAs(t, delete[0].Error, &persistencemodel.PersistenceError{})
	assert.Equal(t, 409, delete[0].Error.(persistencemodel.PersistenceError).Code)
}

func TestDeleteCorrectVersion(t *testing.T) {
	ctx := context.Background()

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

	delete := p.Delete(ctx, persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID:      persistencemodel.PersistenceID{ID: "deviceA", Name: "shadowA", ModelType: 0 /*combined*/},
		Version: 1, /* version condition ok -> success */
	})

	require.Len(t, delete, 1)
	assert.NoError(t, delete[0].Error)
}
