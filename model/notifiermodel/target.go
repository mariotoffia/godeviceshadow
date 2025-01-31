package notifiermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/randutils"
)

// NotificationTarget is the target of one or
// more notification `Selection`.
type NotificationTarget interface {
	// Name is the name of the target.
	Name() string
	// Notify will notify the target with one or more operations that has passed `Selection`.
	// (or no selection). The optional _tx_ parameter is for plugins that support a transaction
	// that may be used to do mass notification in same transaction or even between the `Manager`
	// report/desired state changes and notification as a atomic unit.
	//
	// The _operations_ are passed as many as possible, to allow the implementation to do
	// batch processing of the operations (if possible).
	//
	// It must produces *exactly* the _same_ amount of results as the _operations_.
	Notify(ctx context.Context, tx *persistencemodel.TransactionImpl, operation ...NotifierOperation) []NotificationTargetResult
}

type NotificationTargetResult struct {
	Error     error
	Target    NotificationTarget
	Operation NotifierOperation
	Custom    map[string]any
}

type SelectionTargetImpl struct {
	// Selection is the selection that will be used to filter in/out
	// the `NotifierOperation` that will be passed to the `NotificationTarget`.
	//
	// If no selection, it will just call the target with all operations.
	Selection Selection
	// Target is the target that will be notified if any of the `Selections`
	// do match the `NotifierOperation`.
	Target NotificationTarget
}

// NotifyFunc is the same function as in `NotificationTarget`.
type NotifyFunc func(ctx context.Context, target NotificationTarget, tx *persistencemodel.TransactionImpl, operation ...NotifierOperation) []NotificationTargetResult

type FunctionTargetImpl struct {
	f    NotifyFunc
	name string
}

func (f *FunctionTargetImpl) Name() string {
	return f.name
}

func (f *FunctionTargetImpl) Notify(ctx context.Context, tx *persistencemodel.TransactionImpl, operation ...NotifierOperation) []NotificationTargetResult {
	return f.f(ctx, f, tx, operation...)
}

func FuncTarget(f NotifyFunc, name ...string) NotificationTarget {
	var n string

	if len(name) > 0 {
		n = name[0]
	} else {
		n, _ = randutils.GenerateId()
	}

	return &FunctionTargetImpl{
		f:    f,
		name: n,
	}
}
