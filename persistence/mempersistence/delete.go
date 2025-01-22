package mempersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Delete deletes models from the in-memory persistence. Supports optional version constraints.
func (p *Persistence) Delete(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	results := make([]persistencemodel.WriteResult, len(operations))

	if opt.Tx != nil {
		for i, op := range operations {
			results[i] = persistencemodel.WriteResult{
				ID: persistencemodel.PersistenceID{
					ID:        op.ID.ID,
					Name:      op.ID.Name,
					ModelType: op.ID.ModelType,
				},
				Error: persistencemodel.Error400("Transactions are not supported"),
			}
		}
		return results
	}

	for i, op := range operations {
		results[i] = persistencemodel.WriteResult{
			ID: persistencemodel.PersistenceID{
				ID:        op.ID.ID,
				Name:      op.ID.Name,
				ModelType: op.ID.ModelType,
			},
			Error: p.store.DeleteEntry(op.ID.ID, op.ID.Name, op.Version),
		}
	}

	return results
}
