package stdmgr_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/manager/stdmgr"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tz = "Europe/Stockholm"

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

func TestReportCreateNew(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := stdmgr.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.SeparateModels).
		WithReportLoggers(changelogger.New()).
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

	res := mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Version:  0,
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

func TestReportUpdateReport(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := stdmgr.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.CombinedModels).
		WithReportLoggers(changelogger.New()).
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

	res := mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Version:  0,
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
	assert.True(t, res[0].ReportedProcessed)

	res = mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Version:  0, // update latest
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
	assert.True(t, res[0].ReportedProcessed)

	chl := changelogger.FindLogger(res[0].MergeLoggers)
	require.NotNil(t, chl)
	require.Len(t, chl.PlainLog, 1)
	require.Len(t, chl.ManagedLog, 1)

	assert.Len(t, chl.PlainLog[model.MergeOperationNotChanged], 1, "no change")
	assert.Len(t, chl.ManagedLog[model.MergeOperationUpdate], 1, "temp sensor shall have been updated")
}

func TestReportUpdateReportNotChanged(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mgr := stdmgr.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.CombinedModels).
		WithReportLoggers(changelogger.New()).
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

	res := mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Version:  0,
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

	res = mgr.Report(ctx, managermodel.ReportOperation{
		ClientID: "myClient",
		Version:  0, // update latest
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

// BenchmarkNewReportAndUpdateReport benchmarks the creation of a new report and then
// update the report with a new value.
//
// On my machine it takes about 14μs to perform this benchmark.
func BenchmarkNewReportAndUpdateReport(t *testing.B) {
	ctx := context.Background()
	now := time.Now()

	mgr := stdmgr.New().
		WithPersistence(mempersistence.New()).
		WithSeparation(persistencemodel.CombinedModels).
		WithReportLoggers(changelogger.New()).
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

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		device := fmt.Sprintf("device%d", i)
		res := mgr.Report(ctx, managermodel.ReportOperation{
			ClientID: "myClient",
			Version:  0,
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"temp": {Value: 23.4, TimeStamp: now},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: "homeHub"},
		})

		require.Len(t, res, 1)
		require.NoError(t, res[0].Error)
		assert.True(t, res[0].ReportedProcessed)

		res = mgr.Report(ctx, managermodel.ReportOperation{
			ClientID: "myClient",
			Version:  0, // update latest
			Model: TestModel{
				TimeZone: tz,
				Sensors: map[string]Sensor{
					"temp": {Value: 23.5, TimeStamp: now.Add(1 * time.Minute)},
				},
			},
			ID: persistencemodel.ID{ID: device, Name: "homeHub"},
		})

		require.Len(t, res, 1)
		require.NoError(t, res[0].Error)
		assert.True(t, res[0].ReportedProcessed)

		chl := changelogger.FindLogger(res[0].MergeLoggers)
		require.NotNil(t, chl)
		require.Len(t, chl.PlainLog, 1)
		require.Len(t, chl.ManagedLog, 1)

		assert.Len(t, chl.PlainLog[model.MergeOperationNotChanged], 1, "no change")
		assert.Len(t, chl.ManagedLog[model.MergeOperationUpdate], 1, "temp sensor shall have been updated")
	}
}
