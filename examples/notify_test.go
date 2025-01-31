package examples

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
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationWithDSL(t *testing.T) {
	stmt := `
		(
			id: /myDevice-\d+/ AND 
			name: 'homeHub' AND 
			operation: report,desired
		)
		AND
		(add,update:/^Sensors.indoor-\d+$/ == 'temp'  
		WHERE (
			value > 20 OR (value == /^re-\d+/ AND value != 'apa' OR (value > 99 AND value != /^bubben-\d+$/)))
		)
		OR 
		(acknowledge)
	`

	sel, err := selectlang.ToSelection(stmt)
	require.NoError(t, err)

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
		WithSelection(sel).
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

	id := persistencemodel.ID{ID: "myDevice-992", Name: "homeHub"}
	res := mgr.Report(context.TODO(), managermodel.ReportOperation{
		ID: id, Model: TestModel{
			TimeZone: "Europe/Stockholm",
			Sensors: map[string]Sensor{
				"indoor-991": {Value: map[string]any{"temp": 23.4, "rh": 45.6}, TimeStamp: time.Now().UTC()},
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
