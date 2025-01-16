package dynamodbpersistence_test

import (
	"context"
	"fmt"
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
	require.NotNil(t, results)
	assert.Len(t, results.Items, 0)
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
	require.NotNil(t, results)
	assert.Len(t, results.Items, 2)
}

func TestListMultipleDifferentIDs(t *testing.T) {
	ctx := context.TODO()

	// Create or verify the test table
	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	// Initialize our DynamoDB persistence
	p, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:  res.Table,
		Client: res.Client,
	})
	require.NoError(t, err)

	// We'll create two device IDs, each with two states (Reported, Desired)
	writeOps := p.Write(
		ctx,
		persistencemodel.WriteOptions{
			Config: persistencemodel.WriteConfig{
				Separation: persistencemodel.CombinedModels, // store them combined
			},
		},

		// Device #1, reported
		persistencemodel.WriteOperation{
			ClientID: persistutils.Id("test-"),
			ID: persistencemodel.PersistenceID{
				ID:        "deviceA",
				Name:      "shadowA",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: "UTC",
				Sensors: map[string]Sensor{
					"temp": {Value: 21.5, TimeStamp: time.Now().UTC()},
				},
			},
		},
		// Device #1, desired
		persistencemodel.WriteOperation{
			ClientID: persistutils.Id("test-"),
			ID: persistencemodel.PersistenceID{
				ID:        "deviceA",
				Name:      "shadowA",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: TestModel{
				TimeZone: "UTC",
			},
		},

		// Device #2, reported
		persistencemodel.WriteOperation{
			ClientID: persistutils.Id("test-"),
			ID: persistencemodel.PersistenceID{
				ID:        "deviceB",
				Name:      "shadowB",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: TestModel{
				TimeZone: "CET",
			},
		},
		// Device #2, desired
		persistencemodel.WriteOperation{
			ClientID: persistutils.Id("test-"),
			ID: persistencemodel.PersistenceID{
				ID:        "deviceB",
				Name:      "shadowB",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: TestModel{
				TimeZone: "CET",
			},
		},
	)

	require.Len(t, writeOps, 4)
	for _, w := range writeOps {
		assert.NoError(t, w.Error, "Expected each write operation to succeed")
	}

	// 1) List with no ID => should return all items (4 items).
	allResults, err := p.List(ctx, persistencemodel.ListOptions{})
	require.NoError(t, err)
	require.NotNil(t, allResults)
	assert.Len(t, allResults.Items, 4, "Expected 4 items total (2 for deviceA, 2 for deviceB)")

	// 2) List by ID="deviceA" => only deviceA items (2 items).
	aResults, err := p.List(ctx, persistencemodel.ListOptions{ID: "deviceA"})
	require.NoError(t, err)
	require.NotNil(t, aResults)
	assert.Len(t, aResults.Items, 2, "Expected 2 items for deviceA")

	// 3) List by ID="deviceB" => only deviceB items (2 items).
	bResults, err := p.List(ctx, persistencemodel.ListOptions{ID: "deviceB"})
	require.NoError(t, err)
	require.NotNil(t, bResults)
	assert.Len(t, bResults.Items, 2, "Expected 2 items for deviceB")
}

func TestListPagination(t *testing.T) {
	ctx := context.TODO()

	// Create or verify the test table
	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	pageSize := 11

	// Initialize our DynamoDB persistence
	p, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:            res.Table,
		Client:           res.Client,
		MaxReadBatchSize: pageSize,
	})
	require.NoError(t, err)

	var writeOps []persistencemodel.WriteOperation

	clientId := persistutils.Id("test-")
	for i := 1; i <= 10; i++ {
		deviceID := persistutils.Id("device-")
		writeOps = append(writeOps,
			persistencemodel.WriteOperation{
				ClientID: clientId,
				ID: persistencemodel.PersistenceID{
					ID:        deviceID,
					Name:      "shadow",
					ModelType: persistencemodel.ModelTypeReported,
				},
				Model: TestModel{
					TimeZone: "Europe/Stockholm",
					Sensors: map[string]Sensor{
						"temp": {Value: 23.4, TimeStamp: time.Now().UTC()},
					},
				},
			},
			persistencemodel.WriteOperation{
				ClientID: clientId,
				ID: persistencemodel.PersistenceID{
					ID:        deviceID,
					Name:      "shadow",
					ModelType: persistencemodel.ModelTypeDesired,
				},
				Model: TestModel{},
			},
		)
	}

	// Perform writes
	writeResults := p.Write(ctx, persistencemodel.WriteOptions{
		Config: persistencemodel.WriteConfig{Separation: persistencemodel.SeparateModels},
	}, writeOps...)

	require.Len(t, writeResults, 20)

	for _, w := range writeResults {
		require.NoError(t, w.Error, "Expected each write operation to succeed")
	}

	firstPage, err := p.List(ctx, persistencemodel.ListOptions{Token: ""}) // Start at first page
	require.NoError(t, err)
	require.NotNil(t, firstPage)

	// Ensure the first page has pageSize items and a continuation token
	assert.Len(t, firstPage.Items, pageSize)
	assert.NotEmpty(t, firstPage.Token)

	// Fetch the second page using the returned token
	secondPage, err := p.List(ctx, persistencemodel.ListOptions{Token: firstPage.Token})
	require.NoError(t, err)
	require.NotNil(t, secondPage)

	// Ensure second page has the remaining items
	assert.Len(t, secondPage.Items, 9)
	assert.Empty(t, secondPage.Token)
}

func TestListMultipleShadowsOnSameID(t *testing.T) {
	ctx := context.TODO()

	// Create or verify the test table
	res := dynamodbutils.NewTestTableResource(ctx, TestTableName)
	defer res.Dispose(ctx, dynamodbutils.DisposeOpts{DeleteItems: true})

	pageSize := 100

	// Initialize our DynamoDB persistence
	p, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:            res.Table,
		Client:           res.Client,
		MaxReadBatchSize: pageSize,
	})

	require.NoError(t, err)

	var operations []persistencemodel.WriteOperation

	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("limhamn-%d", i+1)

		operations = append(operations,
			persistencemodel.WriteOperation{
				ClientID: persistutils.Id("test-"),
				ID: persistencemodel.PersistenceID{
					ID:        id,
					Name:      "shadow",
					ModelType: persistencemodel.ModelTypeReported,
				},
				Model: TestModel{
					TimeZone: "UTC",
				},
			},
		)
	}

	results := p.Write(ctx, persistencemodel.WriteOptions{
		Config: persistencemodel.WriteConfig{Separation: persistencemodel.SeparateModels},
	}, operations...)

	require.Len(t, results, 3)

	for _, w := range results {
		require.NoError(t, w.Error)
	}

	page, err := p.List(ctx, persistencemodel.ListOptions{
		ID:         "limhamn-2",
		SearchExpr: "shadow-1",
	})

	require.NoError(t, err)
	require.NotNil(t, page)

	assert.Len(t, page.Items, 1)
	assert.Empty(t, page.Token)
}
