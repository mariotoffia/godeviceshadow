package stdmgr

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Report will report one or more models to the `Manager`. It will merge the reported model with the existing model
// (from persistence) and if any desired (from persistence) values are matched, it will acknowledge them. Then it
// will update both reported and desired (if any changes).
//
// If another process / go routine is updating the same id, it may fail, and return an error. If e.g. 409 (Conflict)
// the caller may safely re-try the operation.
//
// TIP: It will *always* return a slice of `  managermodel.ReportOperationResult` with the same length as the input `operations`.
//
// This implements the `managermodel.Reportable` interface.
func (mgr *ManagerImpl) Report(ctx context.Context, operations ...managermodel.ReportOperation) []managermodel.ReportOperationResult {
	if len(operations) == 0 {
		return nil
	}

	toResults := func(results map[string]*managermodel.ReportOperationResult) []managermodel.ReportOperationResult {
		res := make([]managermodel.ReportOperationResult, 0, len(operations))

		for op := range results {
			res = append(res, *results[op])
		}

		return res
	}

	results := make(map[string]*managermodel.ReportOperationResult, len(operations))

	// Prepare for read
	readOps := mgr.reportPrepareForRead(operations, results)

	if len(readOps) == 0 {
		return toResults(results)
	}

	// Read the models
	readResults := mgr.reportReadFromPersistence(ctx, readOps, results, true /*create*/)

	if len(readResults) == 0 {
		return toResults(results)
	}

	// Merge the models
	readResults = mgr.reportMergeModels(readResults, operations, results)

	// Now we may have queueDesired|Reported models to persist.
	writes := reportCreateWrites(readResults)

	if len(writes) == 0 {
		return toResults(results)
	}

	// Write
	mgr.reportWriteBack(ctx, writes, results)

	return toResults(results)
}

func (mgr *ManagerImpl) reportMergeModels(
	readResults []groupedPersistenceResult,
	operations []managermodel.ReportOperation,
	results map[string]*managermodel.ReportOperationResult,
) []groupedPersistenceResult {
	//
	getOperation := func(id, name string) *managermodel.ReportOperation {
		for _, op := range operations {
			if op.ID.ID == id && op.ID.Name == name {
				return &op
			}
		}

		return nil
	}

	for i, rdr := range readResults {
		op := getOperation(rdr.id.ID, rdr.id.Name)
		sep := mgr.separation

		if op.Separation > 0 {
			sep = op.Separation
		}

		readResults[i].op = op

		// Merge the Reported models
		var (
			reported any // Merge into this to use below: desired
			err      error
		)

		if rdr.reported != nil {
			ml := mgr.createMergeLoggers(true /*report*/, op.MergeLoggers)

			mergeMode := merge.ServerIsMaster

			if op.MergeMode > 0 {
				mergeMode = op.MergeMode
			}

			reported, err = merge.MergeAny(rdr.reported.Model, op.Model, merge.MergeOptions{
				Mode:    mergeMode,
				Loggers: ml,
			})

			if err != nil {
				results[rdr.id.String()] = &managermodel.ReportOperationResult{
					ID:    rdr.id,
					Error: err,
				}

				continue
			}

			dl, _ := FindMergeDirtyLogger(ml)

			if !dl.Dirty {
				// Nothing to do (no changes)
				results[rdr.id.String()] = &managermodel.ReportOperationResult{
					ID:           rdr.id,
					MergeLoggers: ml,
				}

				continue
			}

			// need persist -> queue
			readResults[i].queueReported = reported

			results[rdr.id.String()] = &managermodel.ReportOperationResult{
				ID:           rdr.id,
				MergeLoggers: ml,
				Model:        reported,
			}
		}

		if rdr.desired != nil && reported != nil {
			dl := mgr.createDesiredLoggers(op.DesiredLoggers)
			model, err := merge.DesiredAny(reported, rdr.desired.Model, merge.DesiredOptions{
				Loggers: dl,
			})

			if err != nil {
				// Skip this
				results[rdr.id.String()] = &managermodel.ReportOperationResult{
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
				results[rdr.id.String()] = &managermodel.ReportOperationResult{
					ID:             rdr.id,
					DesiredLoggers: dl,
				}
			}

			// Always needed when combined independent on dirty
			if dla.Dirty || sep == persistencemodel.CombinedModels {
				// need persist -> queue
				readResults[i].queueDesired = model

				// Make sure reported is persisted as well when combined models
				if sep == persistencemodel.CombinedModels && readResults[i].queueReported == nil {
					readResults[i].queueReported = reported
				}
			}
		}
	}

	return readResults
}

func (mgr *ManagerImpl) reportWriteBack(ctx context.Context, writes []persistencemodel.WriteOperation, results map[string]*managermodel.ReportOperationResult) {
	result := mgr.persistence.Write(ctx, persistencemodel.WriteOptions{
		Config: persistencemodel.WriteConfig{
			Separation: mgr.separation,
		},
	}, writes...)

	for _, wr := range result {
		if wr.Error != nil {
			if r, ok := results[wr.ID.StringWithoutModelType()]; ok {
				r.Error = wr.Error
			} else {
				results[wr.ID.StringWithoutModelType()] = &managermodel.ReportOperationResult{
					ID:    wr.ID.ToID(),
					Error: wr.Error,
				}
			}

			continue
		}

		if wr.ID.ModelType == persistencemodel.ModelTypeReported {
			if r, ok := results[wr.ID.StringWithoutModelType()]; ok {
				r.ReportedProcessed = true
			} else {
				results[wr.ID.StringWithoutModelType()] = &managermodel.ReportOperationResult{
					ID:                wr.ID.ToID(),
					ReportedProcessed: true,
				}
			}
		} else if wr.ID.ModelType == persistencemodel.ModelTypeDesired {
			if r, ok := results[wr.ID.StringWithoutModelType()]; ok {
				r.DesiredProcessed = true
			} else {
				results[wr.ID.StringWithoutModelType()] = &managermodel.ReportOperationResult{
					ID:               wr.ID.ToID(),
					DesiredProcessed: true,
				}
			}
		}
	}
}

// createDesiredLoggers will create logger instance from _loggers_ (if any), if none where submitted, it will use the `Manager.desiredLoggers`.
//
// If the `DesiredAckLogger` is not present in the _loggers_ it will be automatically added.
func (mgr *ManagerImpl) createDesiredLoggers(loggers []model.CreatableDesiredLogger) []model.DesiredLogger {
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

func (mgr *ManagerImpl) reportReadFromPersistence(
	ctx context.Context,
	readOps []persistencemodel.ReadOperation,
	results map[string]*managermodel.ReportOperationResult,
	create bool,
) []groupedPersistenceResult {
	readResults := mgr.persistence.Read(ctx, persistencemodel.ReadOptions{}, readOps...)
	res := make(map[string]*groupedPersistenceResult, len(readResults))

	findOp := func(id persistencemodel.PersistenceID) *persistencemodel.ReadOperation {
		for _, op := range readOps {
			if op.ID.Equal(id) {
				return &op
			}
		}

		return nil
	}

	for _, rdr := range readResults {
		if rdr.Error != nil {
			var pe persistencemodel.PersistenceError

			if errors.As(rdr.Error, &pe); create && pe.Code == 404 /*not found*/ {
				// Special case -> create new empty model
				if op := findOp(rdr.ID); op != nil && op.Version == 0 /*only when version is 0*/ {
					rdr.Model = reflect.New(op.Model).Elem().Interface()
					rdr.Error = nil

					if rdr.ID.ModelType == 0 /*combined*/ {
						// Special case, we need to init desired as well
						desired := rdr
						reported := rdr

						desired.ID = desired.ID.ToPersistenceID(persistencemodel.ModelTypeDesired)
						reported.ID = reported.ID.ToPersistenceID(persistencemodel.ModelTypeReported)

						res[rdr.ID.StringWithoutModelType()] = &groupedPersistenceResult{
							id:       rdr.ID.ToID(),
							reported: &reported,
							desired:  &desired,
						}

						continue
					}
				}
			}

			if rdr.Model == nil {
				results[rdr.ID.StringWithoutModelType()] = &managermodel.ReportOperationResult{
					ID:    rdr.ID.ToID(),
					Error: rdr.Error,
				}

				continue
			}
		}

		// In case where there where no model in persistence -> create one
		if rdr.Model == nil {
			if op := findOp(rdr.ID); op != nil {
				rdr.Model = reflect.New(op.Model).Elem().Interface()
			}
		}

		if r, ok := res[rdr.ID.StringWithoutModelType()]; !ok {
			if rdr.ID.ModelType == persistencemodel.ModelTypeReported {
				res[rdr.ID.StringWithoutModelType()] = &groupedPersistenceResult{id: rdr.ID.ToID(), reported: &rdr}
			} else {
				res[rdr.ID.StringWithoutModelType()] = &groupedPersistenceResult{id: rdr.ID.ToID(), desired: &rdr}
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

func (mgr *ManagerImpl) reportPrepareForRead(operations []managermodel.ReportOperation, results map[string]*managermodel.ReportOperationResult) []persistencemodel.ReadOperation {
	// Prepare for read
	readOps := make([]persistencemodel.ReadOperation, 0, len(operations))

	for _, op := range operations {
		te, ok := mgr.ResolveType(op.ModelType, op.ID)

		if !ok {
			results[op.ID.String()] = &managermodel.ReportOperationResult{
				ID:    op.ID,
				Error: fmt.Errorf("unable to resolve type: %s", op.ModelType),
			}

			continue
		}

		if op.Model == nil {
			results[op.ID.String()] = &managermodel.ReportOperationResult{
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

func reportCreateWrites(readResults []groupedPersistenceResult) []persistencemodel.WriteOperation {
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

	return writes
}
