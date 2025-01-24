package manager

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Delete implements the `managermodel.Remover` interface and will delete one or more models from the persistence specified in
// the _operations_. It will return the same amount of results as the operations provided.
//
// If model type is _zero_ in the operation it will delete both reported and desired model in one go since it signals a combined storage.
// If separate storage the model type *must* be provided.
func (mgr *Manager) Delete(ctx context.Context, operations ...managermodel.DeleteOperation) []managermodel.DeleteOperationOperationResult {
	if len(operations) == 0 {
		return nil
	}

	result := make([]managermodel.DeleteOperationOperationResult, 0, len(operations))
	deletes := make([]persistencemodel.WriteOperation, 0, len(operations))

	for _, op := range operations {
		deletes = append(deletes, persistencemodel.WriteOperation{
			ID:      op.ID,
			Version: op.Version,
		})
	}

	results := mgr.persistence.Delete(ctx, persistencemodel.WriteOptions{}, deletes...)

	for _, res := range results {
		result = append(result, managermodel.DeleteOperationOperationResult{
			ID:    res.ID,
			Error: res.Error,
		})
	}

	return result
}
