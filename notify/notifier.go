package notify

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
)

// NotifierImpl is used to notify a external plugin such as _SQS_ or a in-memory
// queue.
//
// The notifier uses the `changelogger.ChangeMergeLogger` to get the changes in reported/desired
// state of the model.
//
// The notifier is processing selections that is coupled with targets that in it's turn invokes a `NotifyPlugin`.
//
// The selection is done in the following manner:
// * All below can be combined using: AND, OR, NOT, (, )
// * ID: id, name <- regexp match
// * Operation: report, desired, delete | all
// * Logger[n]:
// ** merge/desire op: add,remove,update,no-change,acknowledge | all
// ** merge/desire path <- regexp match
// ** merge/desire [name], value <- operators: ==, !=, >, <, >=, <=, regexp match, before, after.
// ** Supported value types: (float,int,uint, string, bool, time.Time)
//
// The logger value name is optional and can only be applied to `model.ValueAndTimestamp.GetValue()` that returns a `map[string]any`.
// If no value name it is assumed that the value is a scalar value.
//
// Since there may be many logger entries in the expression it is possible to combine them with AND, OR and NOT. It is also possible
// use parenthesis to group expressions. Logger expressions always need a grouping since it is always a qualifier to operation + optional path
// and then the value expression. That means that all before the first _AND_ is a logger entry selection expression.
//
// Operator: regexp '/{regex expression}/' can only be applied to string value. The _before_, _after_ can only be used with `time.Time` values and expects
// a `time.RFC3339` formatted string.
//
// Example ID expressions:
// ```(id: /myDevice-\d+/ AND name: 'myShadow')```
// The above example checks if the id and name matches the regexp `/myDevice-\d+/` and the name is equal to `myShadow`.
//
// Example ID and Operation expressions:
// ```(id: /myDevice-\d+/ AND operation: report,delete)```
// The above example checks if the id matches the regexp `/myDevice-\d+/` and the operation is either `report` or `delete`.
//
// Example Logger expressions:
//
// ```
// (add,update:/^Sensors-.*-indoor$/ == 'temp' AND (value > 20 OR value < 10))
// ```
// The above example checks the add, update operations if any path matches the regexp `^Sensors-.*-indoor$` and the value is a map that has
// "temp" as key and the value is either greater than 20 or less than 10.
//
// ```
// (add,update:/^Sensors-.*-indoor$/ AND (value > 20 OR value < 10))
// ```
// All add, update operations that has a path matching the regexp `^Sensors-.*-indoor$` and the value is either greater than 20 or less than 10.
// This requires the value to be a scalar value.
//
// ```
// (add,update AND (value > 20 OR value < 10))
// ```
// This captures all add, update operation and checks all scalars if they are greater than 20 or less than 10.
//
// ```
// (add,update,acknowledge)
// ```
// This captures all add, update, acknowledge operations.
//
// Combine all Examples:
// ```
// (id: /myDevice-\d+/ AND operation: report,delete) AND (add,update:/^Sensors-.*-indoor$/ AND (value > 20 OR value < 10))
// ```
//
// ```
// (id: /myDevice-\d+/ AND name: 'myShadow') AND (add,update:/^Sensors-.*-indoor$/ == 'temp' AND (value > 20 OR value < 10)) OR (add,update)
// ```
type NotifierImpl struct {
}

// Process implements the `notifiermodel.Notifier` interface and will process the operations and notify
// attached `NotifyPlugin`.
func (n *NotifierImpl) Process(ctx context.Context, operations ...notifiermodel.NotifierOperation) []notifiermodel.NotifierOperationResult {
	panic("implement me")
}
