package managermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type ReportOperation struct {
	// ClientID is a optional client ID.
	ClientID string
	// Version when set to zero -> report to latest version. Otherwise, it expects the specific version in persistence and if not,
	// it will fail with 409 (Conflict).
	Version int64
	// Model to report. If any desired values, those will be checked and acknowledged if matched.
	Model any
	// Separation is the separation to use for this operation. If not set it will use the `Manager` default.
	Separation persistencemodel.ModelSeparation
	// ID is the id of the model to report.
	ID persistencemodel.ID
	// ModelType is the type of the model. This is to explicitly direct the `Manager` to lookup the model by this name.
	// Otherwise, it will try to infer it via its `ID`.
	ModelType string
	// MergeLoggers will override the default merge loggers, for report function, in the `Manager`.
	MergeLoggers []model.CreatableMergeLogger
	// DesiredLoggers will override the default desired loggers, for report function, in the `Manager`.
	//
	// These loggers are invoked when the persisted desired model and the reported acknowledges the desired value and thus is removed,
	// from the desired model.
	DesiredLoggers []model.CreatableDesiredLogger
	// MergeMode is the merge mode to use when merging the reported models. If not set it will use the `merge.ServerIsMaster`.
	//
	// TIP: This is useful when removal of items in the reported model is wanted. When `merge.ServerIsMaster` is used, it will only
	// upsert the model. When `merge.ClientIsMaster` is used, it will add, remove and update items.
	MergeMode merge.MergeMode
}

type ReportOperationResult struct {
	// ID is the id of the model that was reported.
	ID persistencemodel.ID
	// MergeLoggers are those loggers that participated in the merge operation.
	MergeLoggers []model.MergeLogger
	// DesiredLoggers are those loggers that participated in the desired operation.
	DesiredLoggers []model.DesiredLogger
	// Error is set when an error did occur during the operation.
	//
	// When error, only ID and this property may be valid
	Error error
	// ReportedProcessed is set to `true` if there where changes and the reported model was persisted.
	//
	// If neither of those (reported, desired), nothing was changed.
	ReportedProcessed bool
	// DesiredProcessed is set to `true` if there where changes and the desired model was persisted.
	//
	// If neither of those (reported, desired), nothing was changed.
	DesiredProcessed bool
	// ReportModel is the resulting model after merge operation of the report model
	ReportModel any
	// DesiredModel is the resulting model after acknowledge operation of the desired model
	DesiredModel any
}

type Reportable interface {
	Report(ctx context.Context, operations ...ReportOperation) []ReportOperationResult
}
