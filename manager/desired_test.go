package manager_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/manager"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDesiredCreateNew(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
		WithDesiredMergeLoggers(changelogger.New()).
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

	res := mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now}, // We desire this to be set to 23.4
			},
		},
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)

	chl := changelogger.FindLogger(res[0].MergeLoggers)
	require.NotNil(t, chl)
	require.Len(t, chl.PlainLog, 1)
	require.Len(t, chl.ManagedLog, 1)

	sns, err := chl.ManagedFromPath(`Sensors\..*`)
	require.NoError(t, err)

	sensors := sns.All()
	require.Len(t, sensors, 1)

	assert.Equal(t, "Sensors.temp", sensors[0].Path)
	assert.Nil(t, sensors[0].OldValue)
	assert.Equal(t, 23.4, sensors[0].NewValue.GetValue())
	assert.Equal(t, time.Time{}, sensors[0].OldTimeStamp)
	assert.Equal(t, now, sensors[0].NewValue.GetTimestamp())
}

func TestDesiredUpdateDesired(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
		WithDesiredMergeLoggers(changelogger.New()).
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

	res := mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now},
			},
		},
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)
	assert.True(t, res[0].Processed)

	res = mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.5, TimeStamp: now.Add(1 * time.Minute)},
			},
		},
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)
	assert.True(t, res[0].Processed)

	chl := changelogger.FindLogger(res[0].MergeLoggers)
	require.NotNil(t, chl)
	require.Len(t, chl.PlainLog, 1)
	require.Len(t, chl.ManagedLog, 1)

	assert.Len(t, chl.PlainLog[model.MergeOperationNotChanged], 1, "no change")
	assert.Len(t, chl.ManagedLog[model.MergeOperationUpdate], 1, "temp sensor shall have been updated")
}

func TestDesiredUpdateDesiredNotChanged(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := manager.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
		WithDesiredMergeLoggers(changelogger.New()).
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

	res := mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now},
			},
		},
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)

	res = mgr.Desire(ctx, managermodel.DesireOperation{
		ClientID: "myClient",
		Model: TestModel{
			TimeZone: tz,
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: now},
			},
		},
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)

	chl := changelogger.FindLogger(res[0].MergeLoggers)
	require.NotNil(t, chl)
	require.Len(t, chl.PlainLog, 1)
	require.Len(t, chl.ManagedLog, 1)

	// Nothing has changed
	assert.Len(t, chl.PlainLog[model.MergeOperationNotChanged], 1)
	assert.Len(t, chl.ManagedLog[model.MergeOperationNotChanged], 1)
}

func TestDesiredAcknowledge(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

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
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
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
		ID: persistencemodel.ID{ID: "device123", Name: "homeHub"},
	})

	require.Len(t, resReport, 1)
	require.NoError(t, resReport[0].Error)
}
