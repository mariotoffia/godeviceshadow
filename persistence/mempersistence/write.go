package mempersistence

import (
	"context"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Write writes a model into the in-memory persistence. It supports update and create operations.
func (p *InMemoryPersistence) Write(
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

	p.mu.Lock()
	defer p.mu.Unlock()

	for i, op := range operations {
		result := persistencemodel.WriteResult{
			ID: persistencemodel.PersistenceID{
				ID:        op.ID.ID,
				Name:      op.ID.Name,
				ModelType: op.ID.ModelType,
			},
		}

		// Check if the model already exists
		if _, exists := p.store[op.ID.ID]; !exists {
			p.store[op.ID.ID] = map[string]*modelEntry{}
		}

		entry, exists := p.store[op.ID.ID][op.ID.Name]

		// Handle version conflicts
		if exists && op.Version > 0 && entry.version != op.Version {
			result.Error = persistencemodel.Error409("Version conflict")
			results[i] = result
			continue
		}

		// Create or update the model
		var version int64

		if !exists {
			version = 1
		} else {
			version = entry.version + 1
		}

		timestamp := time.Now().UnixNano()

		p.store[op.ID.ID][op.ID.Name] = &modelEntry{
			model:       op.Model,
			modelType:   op.ID.ModelType,
			version:     version,
			timestamp:   timestamp,
			clientToken: op.ClientID,
		}

		result.Version = version
		result.TimeStamp = timestamp
		results[i] = result
	}

	return results
}

// Delete deletes models from the in-memory persistence. Supports optional version constraints.
func (p *InMemoryPersistence) Delete(
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

	p.mu.Lock()
	defer p.mu.Unlock()

	for i, op := range operations {
		result := persistencemodel.WriteResult{
			ID: persistencemodel.PersistenceID{
				ID:        op.ID.ID,
				Name:      op.ID.Name,
				ModelType: op.ID.ModelType,
			},
		}

		// Check if the model exists
		models, exists := p.store[op.ID.ID]
		if !exists {
			result.Error = persistencemodel.Error404("ID not found")
			results[i] = result
			continue
		}

		entry, exists := models[op.ID.Name]
		if !exists {
			result.Error = persistencemodel.Error404("ModelType not found")
			results[i] = result
			continue
		}

		// Handle version constraints
		if op.Version > 0 && entry.version != op.Version {
			result.Error = persistencemodel.Error409("Version conflict")
			results[i] = result
			continue
		}

		// Delete the model
		delete(models, op.ID.Name)
		if len(models) == 0 {
			delete(p.store, op.ID.ID)
		}

		result.Version = entry.version
		result.TimeStamp = entry.timestamp
		results[i] = result
	}

	return results
}
