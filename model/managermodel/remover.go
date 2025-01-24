package managermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type DeleteOperation struct {
	// ID is the id of the model to delete. If _zero_ it is a combined persistence and hence it will delete both
	// reported and desired model in one go. If separate use specific model type.
	ID persistencemodel.PersistenceID
	// Version can be set to only delete a certain version. If set to _zero_ it will delete any version.
	Version int64
}

type DeleteOperationOperationResult struct {
	// ID is the id of the model that was deleted. If _zero_ it was a combined persistence and hence it deleted both
	// reported and desired model in one go. If separate it will have either desired or reported model type.
	ID persistencemodel.PersistenceID
	// Error is set when an error did occur during the operation.
	//
	// When error, only ID and this property may be valid
	Error error
}

// Remover is the interface that a model manager that can delete models must implement.
type Remover interface {
	// Delete will delete one or more models from the persistence. If combined persistence it will delete both
	// reported and desired model in one go. If separate it will delete either desired or reported model type.
	//
	// It will produce _exactly_ the same amount of results as the operations provided.
	Delete(ctx context.Context, operations ...DeleteOperation) []DeleteOperationOperationResult
}
