package mempersistence

import (
	"context"
	"sync"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type InMemoryPersistence struct {
	// store is map[ID]map[Name]*modelEntry
	store map[string]map[string]*modelEntry
	mu    sync.RWMutex
}

type modelEntry struct {
	model       any
	modelType   persistencemodel.ModelType
	version     int64
	timestamp   int64
	clientToken string
}

// New creates a new instance of InMemoryReadonlyPersistence.
func New() *InMemoryPersistence {
	return &InMemoryPersistence{
		store: map[string]map[string]*modelEntry{},
	}
}

// List lists models in the in-memory persistence. SearchExpr is not supported.
func (p *InMemoryPersistence) List(
	ctx context.Context,
	opt persistencemodel.ListOptions,
) ([]persistencemodel.ListResult, error) {
	//
	if opt.SearchExpr != "" {
		return nil, persistencemodel.Error400("SearchExpr is not supported")
	}

	if opt.Token != "" {
		return nil, persistencemodel.Error400("Token is not supported")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if opt.ID != "" {
		if m, ok := p.store[opt.ID]; ok {

			res := make([]persistencemodel.ListResult, 0, len(m))
			i := 0

			for name, me := range m {

				res = append(res, persistencemodel.ListResult{
					ID:          persistencemodel.PersistenceID{ID: opt.ID, Name: name, ModelType: me.modelType},
					Version:     me.version,
					TimeStamp:   me.timestamp,
					ClientToken: me.clientToken,
				})

				i++
			}

			return res, nil
		}

		return nil, nil
	}

	var results []persistencemodel.ListResult

	for id, models := range p.store {
		if opt.ID != "" && opt.ID != id {
			continue
		}

		for name, entry := range models {
			results = append(results, persistencemodel.ListResult{
				ID:          persistencemodel.PersistenceID{ID: id, Name: name, ModelType: entry.modelType},
				Version:     entry.version,
				TimeStamp:   entry.timestamp,
				ClientToken: entry.clientToken,
			})
		}
	}

	return results, nil
}

// Read reads models from the in-memory persistence by ID and ModelType.
func (p *InMemoryPersistence) Read(
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

	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(operations) == 0 {
		return nil
	}

	for _, op := range operations {
		if op.ID.ModelType == 0 {
			results = append(results, persistencemodel.ReadResult{
				ID: persistencemodel.PersistenceID{
					ID:   op.ID.ID,
					Name: op.ID.Name,
				},
				Error: persistencemodel.Error400("ModelType is required"),
			})

			continue
		}

		result := persistencemodel.ReadResult{
			ID: persistencemodel.PersistenceID{
				ID:        op.ID.ID,
				Name:      op.ID.Name,
				ModelType: op.ID.ModelType,
			},
		}

		models, exists := p.store[op.ID.ID]
		if !exists {
			result.Error = persistencemodel.Error404("ID not found")
			results = append(results, result)

			continue
		}

		entry, exists := models[op.ID.Name]
		if !exists {
			result.Error = persistencemodel.Error404("ModelType not found")
			results = append(results, result)

			continue
		}

		if op.Version > 0 && op.Version != entry.version {
			result.Error = persistencemodel.Error404("Version not found")
			results = append(results, result)

			continue
		}

		result.Model = entry.model
		result.Version = entry.version
		result.TimeStamp = entry.timestamp
		result.ClientToken = entry.clientToken

		results = append(results, result)
	}

	return results
}
