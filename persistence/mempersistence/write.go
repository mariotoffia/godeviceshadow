package mempersistence

import (
	"context"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

// Write writes a model into the in-memory persistence. It supports update and create operations.
func (p *Persistence) Write(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	//
	results := make([]persistencemodel.WriteResult, 0, len(operations))
	sep := p.opt.Separation

	if opt.Config.Separation != 0 {
		sep = opt.Config.Separation
	}

	groups := persistutils.Group(operations, sep)

	if opt.Tx != nil {
		for _, op := range operations {
			results = append(results, persistencemodel.WriteResult{
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

	for _, op := range groups {
		if err := persistutils.Validate(op); err != nil {
			for _, o := range op.Operations {
				results = append(results, persistencemodel.WriteResult{
					ID:      o.ID,
					Version: o.Version,
					Error:   err,
				})
			}

			continue
		}

		results = append(results, p.writeGroup(op)...)
	}

	return results
}

func (p *Persistence) writeGroup(group persistutils.GroupedWriteOperation) []persistencemodel.WriteResult {
	if group.ModelSeparation == persistencemodel.CombinedModels {
		return p.writeCombined(group)
	}

	res := make([]persistencemodel.WriteResult, 0, len(group.Operations))

	for _, op := range group.Operations {
		res = append(res, p.writeSingle(op))
	}

	return res
}

func (p *Persistence) writeCombined(group persistutils.GroupedWriteOperation) []persistencemodel.WriteResult {
	// Get the desired and reported models
	desired := group.GetByModelType(persistencemodel.ModelTypeDesired)
	reported := group.GetByModelType(persistencemodel.ModelTypeReported)

	if desired == nil && reported == nil {
		return []persistencemodel.WriteResult{
			{
				ID:    persistencemodel.PersistenceID{ID: group.ID, Name: group.Name, ModelType: persistencemodel.ModelTypeDesired},
				Error: persistencemodel.Error400("Neither desired or reported is set, use delete to delete models"),
			},
			{
				ID:    persistencemodel.PersistenceID{ID: group.ID, Name: group.Name, ModelType: persistencemodel.ModelTypeReported},
				Error: persistencemodel.Error400("Neither desired or reported is set, use delete to delete models"),
			},
		}
	}

	var (
		version  int64
		des, rep any
	)

	if desired != nil {
		version = desired.Version
		des = desired.Model
	}

	if reported != nil {
		version = reported.Version
		rep = reported.Model
	}

	now := time.Now().UTC().UnixNano()
	entry, err := p.store.StoreEntry(0 /*combined*/, group.ID, group.Name, &modelEntry{
		version:   version,
		timestamp: now,
		modelType: 0, // Combined
		desired:   des,
		reported:  rep,
	})

	if entry == nil {
		entry = &modelEntry{
			version:   version,
			timestamp: now,
		}
	}

	return []persistencemodel.WriteResult{
		{
			ID:        persistencemodel.PersistenceID{ID: group.ID, Name: group.Name, ModelType: persistencemodel.ModelTypeDesired},
			Version:   entry.version,
			TimeStamp: entry.timestamp,
			Error:     err,
		},
		{
			ID:        persistencemodel.PersistenceID{ID: group.ID, Name: group.Name, ModelType: persistencemodel.ModelTypeReported},
			Version:   entry.version,
			TimeStamp: entry.timestamp,
			Error:     err,
		},
	}
}

func (p *Persistence) writeSingle(op persistencemodel.WriteOperation) persistencemodel.WriteResult {
	now := time.Now().UTC().UnixNano()

	entry := modelEntry{
		version:   op.Version,
		timestamp: now,
		modelType: op.ID.ModelType,
	}

	if op.ID.ModelType == persistencemodel.ModelTypeDesired {
		entry.desired = op.Model
	} else if op.ID.ModelType == persistencemodel.ModelTypeReported {
		entry.reported = op.Model
	}

	res, err := p.store.StoreEntry(op.ID.ModelType, op.ID.ID, op.ID.Name, &entry)

	if res == nil {
		res = &modelEntry{
			version:   op.Version,
			timestamp: now,
		}
	}

	return persistencemodel.WriteResult{
		ID:        op.ID,
		Version:   res.version,
		TimeStamp: res.timestamp,
		Error:     err,
	}
}
