package notify

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// NotifierImpl is used to notify a external plugin such as _SQS_ or a in-memory
// queue.
//
// The notifier uses the `changelogger.ChangeMergeLogger` to get the changes in reported/desired
// state of the model.
//
// The notifier is processing selections that is coupled with targets that in it's turn invokes a `NotifyPlugin`.
type NotifierImpl struct {
	// Targets is a list of `notifiermodel.Selection` and it's associated `notifiermodel.NotificationTarget`.
	//
	// When no selection, it will just call the target with all operations.
	Targets []notifiermodel.SelectionTargetImpl
}

// Process implements the `notifiermodel.Notifier` interface and will process the operations and notify
// attached `NotifyPlugin`.
func (n *NotifierImpl) Process(
	ctx context.Context,
	tx *persistencemodel.TransactionImpl,
	operations ...notifiermodel.NotifierOperation,
) []notifiermodel.NotifierOperationResult {
	// for fast lookup
	targets := make(map[string]notifiermodel.NotificationTarget, len(n.Targets))

	for _, target := range n.Targets {
		targets[target.Target.Name()] = target.Target
	}

	// record keeper of which operation to which target. Since some may be ignored but
	// we still want to allow for as much bach processing as possible (if target supports such).
	ops := make(map[string][]notifiermodel.NotifierOperation, len(targets))

	for _, operation := range operations {
		for _, target := range n.Targets {
			if target.Selection != nil {
				if selected, _ := target.Selection.Select(operation, false /*value*/); selected {
					ops[target.Target.Name()] = append(ops[target.Target.Name()], operation)
				}
			} else {
				ops[target.Target.Name()] = append(ops[target.Target.Name()], operation)
			}
		}
	}

	// "Max" allocate the result slice
	size := 0

	for _, operations := range ops {
		size += len(operations)
	}

	res := make([]notifiermodel.NotifierOperationResult, 0, size)

	for targetName, operations := range ops {
		target := targets[targetName]

		for _, tr := range target.Notify(ctx, tx, operations...) {
			res = append(res, notifiermodel.NotifierOperationResult{
				Target:    tr.Target,
				Operation: tr.Operation,
				Error:     tr.Error,
				Custom:    tr.Custom,
			})
		}
	}

	return res
}
