package manager

import (
	"context"
	"fmt"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Report will report one or more models to the `Manager`. It will merge the reported model with the existing model
// (from persistence) and if any desired (from persistence) values are matched, it will acknowledge them. Then it
// will update both reported and desired (if any changes).
//
// If another process / go routine is updating the same id, it may fail, and return an error. If e.g. 409 (Conflict)
// the caller may safely re-try the operation.
//
// TIP: It will *always* return a slice of `ReportOperationResult` with the same length as the input `operations`.
func (mgr *Manager) Report(ctx context.Context, operations ...ReportOperation) []ReportOperationResult {
	if len(operations) == 0 {
		return nil
	}

	toResults := func(results map[string]*ReportOperationResult) []ReportOperationResult {
		res := make([]ReportOperationResult, 0, len(operations))

		for op := range results {
			res = append(res, *results[op])
		}

		return res
	}

	getOperation := func(id, name string) *ReportOperation {
		for _, op := range operations {
			if op.ID.ID == id && op.ID.Name == name {
				return &op
			}
		}

		return nil
	}

	results := make(map[string]*ReportOperationResult, len(operations))

	// Prepare for read
	readOps := mgr.prepareForRead(operations, results)

	if len(readOps) == 0 {
		return toResults(results)
	}

	// Read the models
	readResults := mgr.readFromPersistence(ctx, readOps, results)

	if len(readResults) == 0 {
		return toResults(results)
	}

	for i, rdr := range readResults {
		op := getOperation(rdr.id.ID, rdr.id.Name)

		readResults[i].op = op

		// Merge the Reported models
		if rdr.reported != nil {
			ml := mgr.createMergeLoggers(op.MergeLoggers)

			model, err := merge.MergeAny(rdr.reported.Model, op.Model, merge.MergeOptions{
				Mode:    merge.ServerIsMaster,
				Loggers: ml,
			})

			if err != nil {
				results[rdr.id.String()] = &ReportOperationResult{
					ID:    rdr.id,
					Error: err,
				}

				continue
			}

			dl, _ := FindMergeDirtyLogger(ml)

			if !dl.Dirty {
				// Nothing to do (no changes)
				results[rdr.id.String()] = &ReportOperationResult{
					ID:           rdr.id,
					MergeLoggers: ml,
				}

				continue
			}

			// need persist -> queue
			readResults[i].queueReported = model

			results[rdr.id.String()] = &ReportOperationResult{
				ID:           rdr.id,
				MergeLoggers: ml,
				Model:        model,
			}
		}

		if rdr.desired != nil {
			dl := mgr.createDesiredLoggers(op.DesiredLoggers)
			model, err := merge.DesiredAny(rdr.reported.Model, rdr.desired.Model, merge.DesiredOptions{
				Loggers: dl,
			})

			if err != nil {
				// Skip this
				results[rdr.id.String()] = &ReportOperationResult{
					ID:    rdr.id,
					Error: err,
				}

				// Do not persist
				readResults[i].queueReported = nil

				continue
			}

			dla, _ := FindDesiredAckLogger(dl)

			if r, ok := results[rdr.id.String()]; ok {
				r.DesiredLoggers = dl
			} else {
				results[rdr.id.String()] = &ReportOperationResult{
					ID:             rdr.id,
					DesiredLoggers: dl,
				}
			}

			if dla.Dirty {
				// need persist -> queue
				readResults[i].queueDesired = model
			}
		}
	}

	// Now we may have queueDesired|Reported models to persist.
	writes := make([]persistencemodel.WriteOperation, 0, len(readResults))

	for _, rdr := range readResults {
		if rdr.queueReported != nil {
			writes = append(writes, persistencemodel.WriteOperation{
				ClientID: rdr.op.ClientID,
				ID:       rdr.reported.ID,
				Model:    rdr.queueReported,
				Version:  rdr.reported.Version,
				Config: persistencemodel.WriteOperationConfig{
					Separation: rdr.op.Separation,
				},
			})
		}

		if rdr.queueDesired != nil {
			writes = append(writes, persistencemodel.WriteOperation{
				ClientID: rdr.op.ClientID,
				ID:       rdr.desired.ID,
				Model:    rdr.queueDesired,
				Version:  rdr.desired.Version,
				Config: persistencemodel.WriteOperationConfig{
					Separation: rdr.op.Separation,
				},
			})
		}
	}

	if len(writes) == 0 {
		return toResults(results)
	}

	// Write
	result := mgr.persistence.Write(ctx, persistencemodel.WriteOptions{
		Config: persistencemodel.WriteConfig{
			Separation: mgr.separation,
		},
	}, writes...)

	for _, wr := range result {
		if wr.Error != nil {
			if r, ok := results[wr.ID.String()]; ok {
				r.Error = wr.Error
			} else {
				results[wr.ID.String()] = &ReportOperationResult{
					ID:    wr.ID.ToModelIndependentPersistenceID(),
					Error: wr.Error,
				}
			}

			continue
		}

		if wr.ID.ModelType == persistencemodel.ModelTypeReported {
			if r, ok := results[wr.ID.String()]; ok {
				r.ReportedProcessed = true
			} else {
				results[wr.ID.String()] = &ReportOperationResult{
					ID:                wr.ID.ToModelIndependentPersistenceID(),
					ReportedProcessed: true,
				}
			}
		} else if wr.ID.ModelType == persistencemodel.ModelTypeDesired {
			if r, ok := results[wr.ID.String()]; ok {
				r.DesiredProcessed = true
			} else {
				results[wr.ID.String()] = &ReportOperationResult{
					ID:               wr.ID.ToModelIndependentPersistenceID(),
					DesiredProcessed: true,
				}
			}
		}
	}

	return toResults(results)
}

// createDesiredLoggers will create logger instance from _loggers_ (if any), if none where submitted, it will use the `Manager.desiredLoggers`.
//
// If the `DesiredAckLogger` is not present in the _loggers_ it will be automatically added.
func (mgr *Manager) createDesiredLoggers(loggers []model.CreatableDesiredLogger) []model.DesiredLogger {
	if len(loggers) == 0 {
		loggers = mgr.reportedDesiredLoggers
	}

	// Add dirty detection
	if !HasDesiredAckLoggerCreator(loggers) {
		loggers = append(loggers, &DesiredAckLogger{})
	}

	res := make([]model.DesiredLogger, 0, len(loggers))

	for _, lg := range loggers {
		res = append(res, lg.New())
	}

	return res
}

// createMergeLoggers will create logger instance from _loggers_ (if any), if none where submitted, it will use the `Manager.reportedLoggers`.
//
// If the `MergeDirtyLogger` is not present in the _loggers_ it will be automatically added.
func (mgr *Manager) createMergeLoggers(loggers []model.CreatableMergeLogger) []model.MergeLogger {
	if len(loggers) == 0 {
		loggers = mgr.reportedMergeLoggers
	}

	// Add dirty detection
	if !HasMergeDirtyLoggerCreator(loggers) {
		loggers = append(loggers, &MergeDirtyLogger{})
	}

	res := make([]model.MergeLogger, 0, len(loggers))

	for _, lg := range loggers {
		res = append(res, lg.New())
	}

	return res
}

func (mgr *Manager) readFromPersistence(ctx context.Context, readOps []persistencemodel.ReadOperation, results map[string]*ReportOperationResult) []groupedPersistenceResult {
	readResults := mgr.persistence.Read(ctx, persistencemodel.ReadOptions{}, readOps...)
	res := make(map[string]*groupedPersistenceResult, len(readResults))

	for _, rdr := range readResults {
		if rdr.Error != nil {
			results[rdr.ID.String()] = &ReportOperationResult{
				ID:    rdr.ID.ToModelIndependentPersistenceID(),
				Error: rdr.Error,
			}

			continue
		}

		if r, ok := res[rdr.ID.String()]; !ok {
			if rdr.ID.ModelType == persistencemodel.ModelTypeReported {
				res[rdr.ID.String()] = &groupedPersistenceResult{id: rdr.ID.ToModelIndependentPersistenceID(), reported: &rdr}
			} else {
				res[rdr.ID.String()] = &groupedPersistenceResult{id: rdr.ID.ToModelIndependentPersistenceID(), desired: &rdr}
			}
		} else {
			if rdr.ID.ModelType == persistencemodel.ModelTypeReported {
				r.reported = &rdr
			} else {
				r.desired = &rdr
			}
		}
	}

	r := make([]groupedPersistenceResult, 0, len(res))

	for _, v := range res {
		r = append(r, *v)
	}

	return r
}

func (mgr *Manager) prepareForRead(operations []ReportOperation, results map[string]*ReportOperationResult) []persistencemodel.ReadOperation {
	// Prepare for read
	readOps := make([]persistencemodel.ReadOperation, 0, len(operations))

	for _, op := range operations {
		te, ok := mgr.ResolveType(op.ModelType, op.ID)

		if !ok {
			results[op.ID.String()] = &ReportOperationResult{
				ID:    op.ID,
				Error: fmt.Errorf("unable to resolve type: %s", op.ModelType),
			}

			continue
		}

		if op.Model == nil {
			results[op.ID.String()] = &ReportOperationResult{
				ID:    op.ID,
				Error: fmt.Errorf("model is nil (use delete to remove model)"),
			}

			continue
		}

		sep := op.Separation

		if op.Separation == 0 {
			sep = mgr.separation
		}

		if sep == persistencemodel.SeparateModels {
			readOps = append(readOps,
				persistencemodel.ReadOperation{
					ID:      persistencemodel.PersistenceID{ID: op.ID.ID, Name: op.ID.Name, ModelType: persistencemodel.ModelTypeReported},
					Version: op.Version,
					Model:   te.Model,
				},
				persistencemodel.ReadOperation{
					ID:      persistencemodel.PersistenceID{ID: op.ID.ID, Name: op.ID.Name, ModelType: persistencemodel.ModelTypeDesired},
					Version: op.Version,
					Model:   te.Model,
				},
			)
		} else {
			readOps = append(readOps, persistencemodel.ReadOperation{
				ID:      persistencemodel.PersistenceID{ID: op.ID.ID, Name: op.ID.Name, ModelType: 0 /*combined*/},
				Version: op.Version,
				Model:   te.Model,
			})
		}
	}

	return readOps
}
