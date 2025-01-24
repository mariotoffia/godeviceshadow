package notify

import "context"

// NotifierImpl is used to notify a external plugin such as _SQS_ or a in-memory
// queue.
//
// The notifier uses the `changelogger.ChangeMergeLogger` to get the changes in reported/desired
// state of the model.
type NotifierImpl struct {
}

type NotifierOperationType string

const (
	OperationTypeReport  NotifierOperationType = "report"
	OperationTypeDesired NotifierOperationType = "desired"
	OperationTypeDelete  NotifierOperationType = "delete"
)

type NotifierOperation struct {
}

type NotifierOperationResult struct {
}

func (n *NotifierImpl) Process(ctx context.Context, operations ...NotifierOperation) []NotifierOperationResult {
	panic("implement me")
}
