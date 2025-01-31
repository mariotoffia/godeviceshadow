package notify_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/manager/stdmgr"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/notify"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestBuildSingleTarget(t *testing.T) {
	// Build Notification Manager
	notificationManager := notify.NewBuilder().
		TargetBuilder(
			notifiermodel.FuncTarget(
				func(
					ctx context.Context, target notifiermodel.NotificationTarget,
					tx *persistencemodel.TransactionImpl, operation ...notifiermodel.NotifierOperation,
				) []notifiermodel.NotificationTargetResult {
					// Target could e.g. be SQS, SNS, Email, SMS, etc.
					var res []notifiermodel.NotificationTargetResult

					for _, op := range operation {
						res = append(res, notifiermodel.NotificationTargetResult{
							Operation: op,
							Target:    target,
							Custom:    map[string]any{"pass": true},
						})
					}

					return res
				})).
		WithSelectionBuilder(
			notifiermodel.NewSelectionBuilder(
				notifiermodel.FuncSelection(
					func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
						return op.ID.Name == "homeHub" && op.Operation == notifiermodel.OperationTypeReport, nil
					})),
		).
		Build().
		Build()

	// Build a manager to do report on so we get a proper changelog et.al
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

	id := persistencemodel.ID{ID: "device123", Name: "homeHub"}
	res := mgr.Report(context.TODO(), managermodel.ReportOperation{
		ID: id, Model: TestModel{
			TimeZone: "Europe/Stockholm",
			Sensors: map[string]Sensor{
				"temp": {Value: 23.4, TimeStamp: time.Now().UTC()},
			},
		},
	})

	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)

	chl := changelogger.Find(res[0].MergeLoggers)

	nResult := notificationManager.Process(
		context.Background(), nil /*tx*/, notifiermodel.NotifierOperation{
			ID:          id.ToPersistenceID(persistencemodel.ModelTypeReported),
			MergeLogger: *chl,
			Operation:   notifiermodel.OperationTypeReport,
			Reported:    res[0].ReportModel,
			Desired:     res[0].DesiredModel,
		},
	)

	require.Len(t, nResult, 1)
	assert.NoError(t, nResult[0].Error)
	assert.Contains(t, nResult[0].Custom, "pass")
	assert.True(
		t, nResult[0].Operation.ID.Equal(id.ToPersistenceID(persistencemodel.ModelTypeReported)),
		"expected the in param operation")
}
