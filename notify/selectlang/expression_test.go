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
	// Very simple statement that should parse correctly
	stmt := `
        SELECT * FROM Notification WHERE
        obj.ID ~= 'myDevice-\\d+'
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

	// Test with matching ID regex
	selected, value := sel.Select(oper, false /*value*/)
	assert.True(t, selected, "Basic ID regex should match")
	assert.Len(t, value, 0)

	// Now test with a numeric comparison
	stmt2 := `
        SELECT * FROM Notification WHERE
        log.Value > 20
    `
	sel2, err := selectlang.ToSelection(stmt2)
	require.NoError(t, err)
	require.NotNil(t, sel2)

	// With temp = 20, should not match
	selected, value = sel2.Select(oper, false /*value*/)
	assert.False(t, selected, "With temp = 20, value > 20 should be false")
	assert.Len(t, value, 0)

	// Change to temp = 21, should match
	mvs.Value = map[string]any{"temp": 21}
	selected, value = sel2.Select(oper, false /*value*/)
	assert.True(t, selected, "With temp = 21, value > 20 should be true")
	assert.Len(t, value, 0)

	// Now test with regex comparison
	stmt3 := `
        SELECT * FROM Notification WHERE
        log.Value ~= 're-\\d+'
    `
	sel3, err := selectlang.ToSelection(stmt3)
	require.NoError(t, err)
	require.NotNil(t, sel3)

	// Change to re-123, should match
	mvs.Value = map[string]any{"temp": "re-123"}
	selected, value = sel3.Select(oper, false /*value*/)
	assert.True(t, selected, "With temp = 're-123', value ~= 're-\\d+' should be true")
	assert.Len(t, value, 0)

	// Change to rea-123, should not match
	mvs.Value = map[string]any{"temp": "rea-123"}
	selected, value = sel3.Select(oper, false /*value*/)
	assert.False(t, selected, "With temp = 'rea-123', value ~= 're-\\d+' should be false")
	assert.Len(t, value, 0)
}
