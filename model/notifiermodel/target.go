package notifiermodel

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

// NotificationTarget is the target of one or
// more notification `Selection`.
type NotificationTarget interface {
	// Notify will notify the target with one or more operations that has passed `Selection`.
	// (or no selection). The optional _tx_ parameter is for plugins that support a transaction
	// that may be used to do mass notification in same transaction or even between the `Manager`
	// report/desired state changes and notification as a atomic unit.
	//
	// The _operations_ are passed as many as possible, to allow the implementation to do
	// batch processing of the operations (if possible).
	Notify(tx *persistencemodel.TransactionImpl, operation ...NotifierOperation) error
}
