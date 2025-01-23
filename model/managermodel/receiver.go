package managermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type ReadOperation struct {
	// ID is the id of the model to read.
	ID persistencemodel.PersistenceID
	// Separation is the separation to use for this operation. If not set it will use the `Manager` default.
	Separation persistencemodel.ModelSeparation
	// ModelType is the type of the model. This is to explicitly direct the `Manager` to lookup the model by this name.
	// Otherwise, it will try to infer it via its `ID`.
	ModelType string
	// Version can be set to only return a certain version. If set to _zero_ it will return the latest version.
	Version int64
}

type ReadOperationResult struct {
	// ID is the id of the model
	ID persistencemodel.PersistenceID
	// Error is set when an error did occur during the operation.
	//
	// When error, only ID and this property may be valid
	Error error
	// Model is the resulting desired or reported model
	Model any
	// Version is version of the `Model`
	Version int64
	// TimeStamp is the timestamp of the model that was written. This is the main timestamp that gets updated
	// each time a model was created or updated. It is a Unix64 bit _UTC_ nanosecond timestamp.
	TimeStamp int64
}

// Receiver is a manager that allows clients to receive desired and reported states.
//
// When not combined is specified, it will return a result entry per _operation_.
type Receiver interface {
	Read(ctx context.Context, operations ...ReadOperation) []ReadOperationResult
}
