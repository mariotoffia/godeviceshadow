package managermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type DesireOperation struct {
	// ClientID is a optional client ID.
	ClientID string
	// Model is the new desired model to merge into existing one.
	Model any
	// Separation is the separation to use for this operation. If not set it will use the `Manager` default.
	Separation persistencemodel.ModelSeparation
	// ID is the id of the model to desire.
	ID persistencemodel.ID
	// ModelType is the type of the model. This is to explicitly direct the `Manager` to lookup the model by this name.
	// Otherwise, it will try to infer it via its `ID`.
	ModelType string
	// MergeLoggers will override the default merge loggers, for desired function.
	MergeLoggers []model.CreatableMergeLogger
}

type DesireOperationResult struct {
	// ID is the id of the model that was reported.
	ID persistencemodel.ID
	// MergeLoggers are those loggers that participated in the merge operation.
	MergeLoggers []model.MergeLogger
	// Error is set when an error did occur during the operation.
	//
	// When error, only ID and this property may be valid
	Error error
	// Model is the resulting desired model after merge operation
	Model any
	// Processed is set to `true` if it was changed and persisted.
	Processed bool
	// Version is the possibly new version of the model.
	Version int64
	// TimeStamp is the timestamp of the model that was written. This is the main timestamp that gets updated
	// each time a model was created or updated. It is a Unix64 bit _UTC_ nanosecond timestamp.
	TimeStamp int64
}

// Desireable is when a manager supports upserting a desired model.
type Desireable interface {
	Desire(ctx context.Context, operations ...DesireOperation) []DesireOperationResult
}
