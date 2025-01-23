package examples

import (
	"context"
	"reflect"
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

func TestDesireReportThatAcknowledgesAndReadAgain(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	const tz = "Europe/Stockholm"

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
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

	resDesire := mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now}, // We desire this to be set to 23.4
			},
		},
		ID: persistencemodel.ID{ID: "device1234", Name: "homeHub"},
	})

	require.Len(t, resDesire, 1)
	require.NoError(t, resDesire[0].Error)

	// Report the desired state -> Clears it in the desired
	resReport := mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now},
			},
		},
		ID: persistencemodel.ID{ID: "device1234", Name: "homeHub"},
	})

	require.Len(t, resReport, 1)
	require.NoError(t, resReport[0].Error)
	assert.True(t, resReport[0].ReportedProcessed)
	assert.True(t, resReport[0].DesiredProcessed)

	resRead := mgr.Read(ctx,
		managermodel.ReadOperation{
			ID: persistencemodel.PersistenceID{ID: "device1234", Name: "homeHub", ModelType: persistencemodel.ModelTypeReported},
		},
		managermodel.ReadOperation{
			ID: persistencemodel.PersistenceID{ID: "device1234", Name: "homeHub", ModelType: persistencemodel.ModelTypeDesired},
		},
	)
	require.Len(t, resRead, 2)
	require.NoError(t, resRead[0].Error)
	require.NoError(t, resRead[1].Error)

	var desired, reported TestModel

	if resRead[0].ID.ModelType == persistencemodel.ModelTypeReported {
		reported = resRead[0].Model.(TestModel)
		desired = resRead[1].Model.(TestModel)
	} else {
		reported = resRead[1].Model.(TestModel)
		desired = resRead[0].Model.(TestModel)
	}

	assert.Len(t, desired.Sensors, 0)
	require.NotNil(t, reported.Sensors)
	require.Len(t, reported.Sensors, 1)

	assert.Equal(t, 23.4, reported.Sensors["temp"].Value)
}
