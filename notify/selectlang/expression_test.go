package selectlang_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpressionPrimaryAndLogger(t *testing.T) {
	stmt := `
		(
			id: /myDevice-\d+/ AND 
			name: 'myShadow' AND 
			operation: report,desired
		)
		AND
		(add,update:/^Sensors-.*-indoor$/ == 'temp'  
		WHERE (
			value > 20 OR (value == /re-\d+/ AND value != 'apa' OR (value > 99 AND value != /bubben-\d+/)))
		)
		OR 
		(acknowledge)
	`

	sel, err := selectlang.ToSelection(stmt)
	require.NoError(t, err)
	require.NotNil(t, sel)

	mvs := &model.ValueAndTimestampImpl{
		Timestamp: time.Now().UTC(),
		Value:     map[string]any{"temp": 20},
	}

	oper := notifiermodel.NotifierOperation{
		ID:        persistencemodel.PersistenceID{ID: "myDevice-123", Name: "myShadow"},
		Operation: notifiermodel.OperationTypeReport,
		MergeLogger: changelogger.ChangeMergeLogger{
			ManagedLog: changelogger.ManagedLogMap{
				model.MergeOperationAdd: {
					{
						Path:     "Sensors-123a-indoor",
						NewValue: mvs,
					},
				},
			},
		},
	}

	selected, value := sel.Select(oper, false /*value*/)

	assert.False(t, selected, "Since v == 20 -> Fail")
	assert.Len(t, value, 0)

	mvs.Value = map[string]any{"temp": 21}

	selected, value = sel.Select(oper, false /*value*/)

	assert.True(t, selected, "Since v > 20 -> Success")
	assert.Len(t, value, 0)

	mvs.Value = map[string]any{"temp": "re-123"}

	selected, value = sel.Select(oper, false /*value*/)
	assert.True(t, selected, `Since v == /re-\d+/ -> Success`)
	assert.Len(t, value, 0)

	mvs.Value = map[string]any{"temp": "rea-123"}

	selected, value = sel.Select(oper, false /*value*/)
	assert.False(t, selected, `Since v == /re-\d+/ -> Fail`)
	assert.Len(t, value, 0)
}
