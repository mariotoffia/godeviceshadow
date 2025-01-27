package notifiermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/loggers/desirelogger"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type NotifierOperationType string

const (
	OperationTypeReport  NotifierOperationType = "report"
	OperationTypeDesired NotifierOperationType = "desired"
	OperationTypeDelete  NotifierOperationType = "delete"
)

type NotifierOperation struct {
	// ID is the id of the reported/desired model that should be notified. When
	// `OperationTypeDelete` the `persistencemodel.PersistenceID.ModelType` can be interpreted
	// as the type of the model to delete (0 if combined i.e. both).
	ID persistencemodel.PersistenceID
	// MergeLogger is the changelogger that contains the changes to be notified. When `OperationTypeReport`
	// it is the report changes and when `OperationTypeDesired` it is the desired changes.
	MergeLogger changelogger.ChangeMergeLogger
	// DesireLogger is the desirelogger that contains the desired state of the model.
	DesireLogger desirelogger.DesireLogger
	// Operation specifies which operation to notify about.
	Operation NotifierOperationType
	// Reported is set when in a `OperationTypeReport` operation. It is the reported state of the model.
	Reported any
	// Desired is set when in a `OperationTypeDesired` operation. It is the desired state of the model.
	//
	// NOTE: It may be part of a `OperationTypeReport` operation as well since desired is loaded and acknowledged.
	//       This is then the resulting desired state after acknowledging the desired state. Which can be discovered
	//       in `DesireLogger` instance.
	Desired any
}

type NotifierOperationResult struct {
	Error error
}

// Notifier will process the _operations_ and use registered `NotifyPlugin` to do the actual
// notification.
type Notifier interface {
	Process(ctx context.Context, operations ...NotifierOperation) []NotifierOperationResult
}
