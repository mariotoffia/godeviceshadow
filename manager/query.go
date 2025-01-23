package manager

import (
	"context"
	"fmt"

	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Query implements the read function in `managermodel.Receiver` interface.
func (mgr *Manager) Read(ctx context.Context, operations ...managermodel.ReadOperation) []managermodel.ReadOperationResult {
	if len(operations) == 0 {
		return nil
	}

	type readOperation struct {
		operation  *persistencemodel.ReadOperation
		separation persistencemodel.ModelSeparation
	}

	readOperations := make(map[string]*readOperation, len(operations))
	result := make([]managermodel.ReadOperationResult, 0, len(operations))

	for _, op := range operations {
		if op.ID.ModelType == 0 {
			result = append(result, managermodel.ReadOperationResult{
				ID:      op.ID,
				Version: op.Version,
				Error:   persistencemodel.Error400("combined model type is not supported, specify Separation instead"),
			})

			continue
		}

		te, ok := mgr.ResolveType(op.ModelType, op.ID.ToID())

		if !ok {
			result = append(result, managermodel.ReadOperationResult{
				ID:      op.ID,
				Version: op.Version,
				Error:   persistencemodel.Error400(fmt.Sprintf("could not resolve model for id: %s", op.ID)),
			})

			continue
		}

		sep := mgr.separation

		if op.Separation > 0 {
			sep = op.Separation
		}

		readOperations[op.ID.String()] = &readOperation{
			separation: sep,
			operation: &persistencemodel.ReadOperation{
				ID:      op.ID,
				Version: op.Version,
				Model:   te.Model,
			},
		}
	}

	readOps := make([]persistencemodel.ReadOperation, 0, len(readOperations))
	combined := make(map[string]struct{}, len(readOperations))

	// Make sure to combine if desired,reported and separation is combined as one readOp

	for _, ro := range readOperations {
		if ro.separation == persistencemodel.SeparateModels {
			readOps = append(readOps, *ro.operation)
		} else /*combined*/ {
			if _, ok := combined[ro.operation.ID.StringWithoutModelType()]; !ok {
				combined[ro.operation.ID.StringWithoutModelType()] = struct{}{}

				op := *ro.operation
				op.ID.ModelType = 0 // combined

				readOps = append(readOps, op)
			}
		}
	}

	readResult := mgr.persistence.Read(ctx, persistencemodel.ReadOptions{}, readOps...)

	toOperation := func(id persistencemodel.PersistenceID, mt persistencemodel.ModelType) *managermodel.ReadOperation {
		for _, op := range operations {
			if op.ID.ID == id.ID && op.ID.Name == id.Name {
				if mt == 0 {
					return &op
				}

				if op.ID.ModelType == mt {
					return &op
				}
			}
		}

		return nil
	}

	for _, rr := range readResult {
		// Special case when combined and error -> only one result is returned
		if rr.ID.ModelType == 0 {
			if op := toOperation(rr.ID, persistencemodel.ModelTypeReported); op != nil {
				result = append(result, managermodel.ReadOperationResult{
					ID:      rr.ID.ToPersistenceID(persistencemodel.ModelTypeReported),
					Version: rr.Version,
					Error:   rr.Error,
				})
			}

			if op := toOperation(rr.ID, persistencemodel.ModelTypeDesired); op != nil {
				result = append(result, managermodel.ReadOperationResult{
					ID:      rr.ID.ToPersistenceID(persistencemodel.ModelTypeDesired),
					Version: rr.Version,
					Error:   rr.Error,
				})
			}
		} else {
			result = append(result, managermodel.ReadOperationResult{
				ID:        rr.ID,
				Version:   rr.Version,
				Error:     rr.Error,
				Model:     rr.Model,
				TimeStamp: rr.TimeStamp,
			})
		}
	}

	return result
}
