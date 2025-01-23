package manager

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func (mgr *Manager) Desire(ctx context.Context, operations ...managermodel.DesireOperation) []managermodel.DesireOperationResult {
	if len(operations) == 0 {
		return nil
	}

	res := make(map[string]*managermodel.DesireOperationResult)
	readOps := make([]persistencemodel.ReadOperation, 0, len(operations))

	readOperation := func(id persistencemodel.PersistenceID) *persistencemodel.ReadOperation {
		for _, op := range readOps {
			if op.ID.Equal(id) {
				return &op
			}
		}

		return nil
	}

	desiredOperation := func(id persistencemodel.PersistenceID) *managermodel.DesireOperation {
		for _, op := range operations {
			if op.ID.ID == id.ID && op.ID.Name == id.Name {
				return &op
			}
		}

		return nil
	}

	// Create read operations
	for _, op := range operations {
		te, ok := mgr.ResolveType(op.ModelType, op.ID)

		if !ok {
			res[op.ID.String()] = &managermodel.DesireOperationResult{
				ID:    op.ID,
				Error: fmt.Errorf("could not resolve model for id: %s", op.ID),
			}

			continue
		}

		sep := mgr.separation

		if op.Separation > 0 {
			sep = op.Separation
		}

		// We always read the last version
		if sep == persistencemodel.CombinedModels {
			readOps = append(readOps, persistencemodel.ReadOperation{
				ID:    persistencemodel.PersistenceID{ID: op.ID.ID, Name: op.ID.Name, ModelType: 0 /*combined*/},
				Model: te.Model,
			})
		} else /*separate*/ {
			readOps = append(readOps, persistencemodel.ReadOperation{
				ID:    persistencemodel.PersistenceID{ID: op.ID.ID, Name: op.ID.Name, ModelType: persistencemodel.ModelTypeDesired},
				Model: te.Model,
			})
		}
	}

	// Read all models
	readResults := mgr.persistence.Read(ctx, persistencemodel.ReadOptions{}, readOps...)

	// Group all read results (to handle combined storage as well)
	grouped := make(map[string]*groupedPersistenceResult)

	for _, rr := range readResults {
		if rr.Error != nil {
			var pe persistencemodel.PersistenceError

			if errors.As(rr.Error, &pe); pe.Code == 404 {
				// Special case Not found -> create a new model
				op := readOperation(rr.ID)

				rr.Model = reflect.New(op.Model).Elem().Interface()
				rr.Error = nil
			}

			if rr.Error != nil {
				res[rr.ID.String()] = &managermodel.DesireOperationResult{
					ID: rr.ID.ToID(), Error: rr.Error,
				}

				continue
			}
		}

		var reported, desired *persistencemodel.ReadResult

		if rr.ID.ModelType == persistencemodel.ModelTypeReported {
			reported = &rr
		} else {
			desired = &rr
		}

		if grp, ok := grouped[rr.ID.StringWithoutModelType()]; !ok {
			grouped[rr.ID.StringWithoutModelType()] = &groupedPersistenceResult{
				id:       rr.ID.ToID(),
				reported: reported,
				desired:  desired,
				dop:      desiredOperation(rr.ID),
			}
		} else {
			if reported != nil {
				grp.reported = reported
			}

			if desired != nil {
				grp.desired = desired
			}
		}
	}

	// Remove all that has no desired operation
	keys := maps.Keys(grouped)
	for k := range keys {
		grp := grouped[k]

		if grp.dop == nil {
			delete(grouped, k)
		}
	}

	for _, rr := range grouped {
		// Create merge loggers
		ml := mgr.createMergeLoggers(false /*report*/, rr.dop.MergeLoggers)

		// Merge the model
		newDesired, err := merge.MergeAny(rr.desired.Model, rr.dop.Model, merge.MergeOptions{
			Mode:    merge.ServerIsMaster,
			Loggers: ml,
		})

		if err != nil {
			res[rr.dop.ID.String()] = &managermodel.DesireOperationResult{
				ID:    rr.dop.ID,
				Error: err,
			}

			continue
		}

		dl, _ := FindMergeDirtyLogger(ml)

		res[rr.dop.ID.String()] = &managermodel.DesireOperationResult{
			ID:           rr.dop.ID,
			Model:        newDesired,
			MergeLoggers: ml,
		}

		if dl.Dirty {
			rr.queueDesired = newDesired

			// If combined persistence reported has to be written
			if rr.reported != nil {
				rr.queueReported = rr.reported.Model
			}
		}
	}

	writes := make([]persistencemodel.WriteOperation, 0, len(grouped)*2)

	for _, grp := range grouped {
		if grp.queueDesired == nil {
			continue
		}

		sep := mgr.separation

		if grp.dop.Separation > 0 {
			sep = grp.dop.Separation
		}

		writes = append(writes, persistencemodel.WriteOperation{
			ID:       persistencemodel.PersistenceID{ID: grp.id.ID, Name: grp.id.Name, ModelType: persistencemodel.ModelTypeDesired},
			Model:    grp.queueDesired,
			ClientID: grp.dop.ClientID,
			Version:  grp.desired.Version,
			Config:   persistencemodel.WriteOperationConfig{Separation: sep},
		})

		if sep == persistencemodel.CombinedModels && grp.queueReported != nil {
			writes = append(writes, persistencemodel.WriteOperation{
				ID:       persistencemodel.PersistenceID{ID: grp.id.ID, Name: grp.id.Name, ModelType: persistencemodel.ModelTypeReported},
				Model:    grp.queueReported,
				ClientID: grp.dop.ClientID,
				Version:  grp.reported.Version,
				Config:   persistencemodel.WriteOperationConfig{Separation: sep},
			})
		}
	}

	writeResults := mgr.persistence.Write(ctx, persistencemodel.WriteOptions{}, writes...)

	for _, wr := range writeResults {
		if wr.Error != nil {
			res[wr.ID.String()] = &managermodel.DesireOperationResult{
				ID:        wr.ID.ToID(),
				Error:     wr.Error,
				Version:   wr.Version,
				TimeStamp: wr.TimeStamp,
			}

			continue
		}

		if dop, ok := res[wr.ID.StringWithoutModelType()]; ok {
			dop.Processed = true
		}
	}

	all := make([]managermodel.DesireOperationResult, 0, len(res))

	for _, v := range res {
		all = append(all, *v)
	}

	return all
}
