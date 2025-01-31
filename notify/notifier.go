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
type NotifierImpl struct {
	Filters []notifiermodel.Selection
}

// Process implements the `notifiermodel.Notifier` interface and will process the operations and notify
// attached `NotifyPlugin`.
func (n *NotifierImpl) Process(ctx context.Context, operations ...notifiermodel.NotifierOperation) []notifiermodel.NotifierOperationResult {
	panic("implement me")
}
