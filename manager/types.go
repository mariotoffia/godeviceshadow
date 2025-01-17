package manager

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Manager struct {
	// persistence to use for _CRUD_ operations
	persistence persistencemodel.Persistence
	// typeRegistry is the type registry that is used when no `TypeRegistryResolver` is
	// registered. It will use it to resolve the types based on the name only or if passed
	// as argument on read operation.
	typeRegistry model.TypeRegistry
	// typeRegistryResolver is primarily used to resolve a type based on it's id and name.
	typeRegistryResolver model.TypeRegistryResolver
	// reportedMergeLoggers is a slice of instantiator to produce MergeLogger(s) of which is all applied when merge operations when
	// user invokes the `Report` function.
	reportedMergeLoggers []model.CreatableMergeLogger
	// reportedDesiredLoggers is a slice of instantiator to produce DesiredLogger(s) of which is all applied when merge operations when
	// user invokes the `Report` function.
	reportedDesiredLoggers []model.CreatableDesiredLogger
	// desiredMergeLoggers is a slice of instantiator to produce MergeLogger(s) of which is all applied when merge operations when
	// user invokes the `Desired` function.
	desiredMergeLoggers []model.CreatableMergeLogger
	// separation is the default separation.
	separation persistencemodel.ModelSeparation
}

type ReportOperation struct {
	// ClientID is a optional client ID.
	ClientID string
	// Separation is the separation to use for this operation. If not set it will use the `Manager` default.
	Separation persistencemodel.ModelSeparation
	// ID is the id of the model to report.
	ID persistencemodel.ModelIndependentPersistenceID
	// ModelType is the type of the model. This is to explicitly direct the `Manager` to lookup the model by this name.
	// Otherwise, it will try to infer it via its `ID`.
	ModelType string
	// MergeLoggers will override the default merge loggers, for report function, in the `Manager`.
	MergeLoggers []model.CreatableMergeLogger
	// DesiredLoggers will override the default desired loggers, for report function, in the `Manager`.
	DesiredLoggers []model.CreatableDesiredLogger
	// Model to report. If any desired values, those will be checked and acknowledged if matched.
	Model any
	// Version when set to zero -> report to latest version. Otherwise, it expects the specific version in persistence and if not,
	// it will fail with 409 (Conflict).
	Version int64
}

type ReportOperationResult struct {
	// ID is the id of the model that was reported.
	ID persistencemodel.ModelIndependentPersistenceID
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
	// Model is the resulting model after merge operation
	Model any
}

type groupedPersistenceResult struct {
	id            persistencemodel.ModelIndependentPersistenceID
	reported      *persistencemodel.ReadResult
	desired       *persistencemodel.ReadResult
	queueReported any
	queueDesired  any
	op            *ReportOperation
}
