package mempersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Read reads models from the in-memory persistence by ID and ModelType.
func (p *Persistence) Read(
	ctx context.Context,
	opt persistencemodel.ReadOptions,
	operations ...persistencemodel.ReadOperation,
) []persistencemodel.ReadResult {
	//
	results := make([]persistencemodel.ReadResult, 0, len(operations))

	if opt.Tx != nil {
		for _, op := range operations {
			results = append(results, persistencemodel.ReadResult{
				ID: persistencemodel.PersistenceID{
					ID:        op.ID.ID,
					Name:      op.ID.Name,
					ModelType: op.ID.ModelType,
				},
				Error: persistencemodel.Error400("Transactions are not supported"),
			})
		}

		return results
	}

	if len(operations) == 0 {
		return nil
	}

	toResult := func(entry *modelEntry, id persistencemodel.PersistenceID, mt persistencemodel.ModelType, model any) persistencemodel.ReadResult {
		return persistencemodel.ReadResult{
			ID: persistencemodel.PersistenceID{
				ID:        id.ID,
				Name:      id.Name,
				ModelType: mt,
			},
			Model:       model,
			Version:     entry.version,
			TimeStamp:   entry.timestamp,
			ClientToken: entry.clientToken,
		}
	}

	for _, op := range operations {
		entry, err := p.store.GetEntry(op.ID.ID, op.ID.Name, op.Version)

		if err != nil {
			results = append(results, persistencemodel.ReadResult{
				ID: persistencemodel.PersistenceID{
					ID:        op.ID.ID,
					Name:      op.ID.Name,
					ModelType: op.ID.ModelType,
				},
				Error: err,
			})

			continue
		}

		if op.ID.ModelType == 0 /*combined*/ {
			if entry.desired != nil {
				results = append(results, toResult(entry, op.ID, persistencemodel.ModelTypeDesired, entry.desired))
			}

			if entry.reported != nil {
				results = append(results, toResult(entry, op.ID, persistencemodel.ModelTypeReported, entry.reported))
			}

			continue
		}

		if op.ID.ModelType == persistencemodel.ModelTypeDesired {
			results = append(results, toResult(entry, op.ID, op.ID.ModelType, entry.desired))
		} else {
			results = append(results, toResult(entry, op.ID, op.ID.ModelType, entry.reported))
		}
	}

	return results
}
