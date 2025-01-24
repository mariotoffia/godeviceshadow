package manager_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/manager"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllModels(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.CombinedModels).
		WithTypeRegistryResolver(
			types.NewRegistry().RegisterResolver(
				model.NewResolveFunc(func(id, name string) (model.TypeEntry, bool) {
					if name == "homeHub" {
						return model.TypeEntry{
							Name: "homeHub", Model: reflect.TypeOf(TestModel{}),
						}, true
					}

					return model.TypeEntry{}, false
				}),
			),
		).
		Build()

	for i := 0; i < 100; i++ {
		device := fmt.Sprintf("device%d", i+1)
		name := "homeHub"

		resReported := mgr.Report(ctx, managermodel.ReportOperation{
			Model: TestModel{
				TimeZone: "Europe/Stockholm",
				Sensors: map[string]Sensor{
					"temp": {Value: 23.4, TimeStamp: now},
					"sp":   {Value: 20.1, TimeStamp: now},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: name},
		})
		require.Len(t, resReported, 1)
		require.NoError(t, resReported[0].Error)

		resDesired := mgr.Desire(ctx, managermodel.DesireOperation{
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"sp": {Value: 23.4, TimeStamp: now},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: name},
		})
		require.Len(t, resDesired, 1)
		require.NoError(t, resDesired[0].Error)
	}

	results, err := mgr.List(ctx /*all*/)
	require.NoError(t, err)
	require.Len(t, results.Items, 200) // Both reported and desired
	assert.Equal(t, int64(2), results.Items[0].Version, "Since combined models")

	resDelete := mgr.Delete(ctx, managermodel.DeleteOperation{
		ID: persistencemodel.PersistenceID{ID: "device1", Name: "homeHub", ModelType: 0 /*since combined*/},
	})
	require.Len(t, resDelete, 1)
	require.NoError(t, resDelete[0].Error)

	results, err = mgr.List(ctx /*all*/)
	require.NoError(t, err)
	require.Len(t, results.Items, 198)
}

func TestListWithinID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
		WithTypeRegistryResolver(
			types.NewRegistry().RegisterResolver(
				model.NewResolveFunc(func(id, name string) (model.TypeEntry, bool) {
					if strings.HasPrefix(name, "homeHub") {
						return model.TypeEntry{
							Name: "homeHub", Model: reflect.TypeOf(TestModel{}),
						}, true
					}

					return model.TypeEntry{}, false
				}),
			),
		).
		Build()

	for i := 0; i < 100; i++ {
		device := "device123"
		name := fmt.Sprintf("homeHub-%d", i+1)

		resReported := mgr.Report(ctx, managermodel.ReportOperation{
			Model: TestModel{
				TimeZone: "Europe/Stockholm",
				Sensors: map[string]Sensor{
					"temp": {Value: 23.4, TimeStamp: now},
					"sp":   {Value: 20.1, TimeStamp: now},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: name},
		})
		require.Len(t, resReported, 1)
		require.NoError(t, resReported[0].Error)

		resDesired := mgr.Desire(ctx, managermodel.DesireOperation{
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"sp": {Value: 23.4, TimeStamp: now},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: name},
		})
		require.Len(t, resDesired, 1)
		require.NoError(t, resDesired[0].Error)
	}

	results, err := mgr.List(ctx /*all*/, managermodel.ListOptions{ID: "device123"})
	require.NoError(t, err)
	require.Len(t, results.Items, 200) // Both reported and desired
	assert.Equal(t, int64(1), results.Items[0].Version, "Since separated models")

	resDelete := mgr.Delete(ctx, managermodel.DeleteOperation{
		ID: persistencemodel.PersistenceID{ID: "device123", Name: "homeHub-99", ModelType: persistencemodel.ModelTypeReported},
	})

	require.Len(t, resDelete, 1)
	require.NoError(t, resDelete[0].Error)

	results, err = mgr.List(ctx /*all*/, managermodel.ListOptions{ID: "device123"})
	require.NoError(t, err)
	require.Len(t, results.Items, 199)

}
