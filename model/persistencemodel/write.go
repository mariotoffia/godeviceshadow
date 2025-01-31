package persistencemodel

import "context"

// Persistence allows a client to read and write a model to a persistent storage.
type Persistence interface {
	ReadonlyPersistence
	// Write will write a _model_ onto the persistent storage. If the model already exists, it will
	// update the model. If the model does not exist, it will create a new model.
	//
	// If a conflict occurs, it will return a `PersistenceError` with code 409 (Conflict).
	//
	// If the `WriteConfig.Separation` is set to `CombinedModels`, it will write the model as a combined model and
	// thus it *REQUIRES* two _operation_, one for reported and one for desired. If it is set to `SeparateModels`,
	// it will write the model as separate models and thus the reported and can be written independently.
	Write(ctx context.Context, opt WriteOptions, operation ...WriteOperation) []WriteResult
	// Delete will delete one or more models from the persistent storage. If the model does not exist, it will
	// return a `PersistenceError` with code 404 (Not Found).
	//
	// When `WriteOperation.Version` is set to zero or less, it will delete the model regardless of the version.
	//
	// Use `PersistenceID.ModelType` as _zero_ to delete a combined model. Do not supply with both desired and reported
	// as in write operation. Otherwise it will interpret the persistence as separate and hence each model type needs
	// to be deleted separately.
	Delete(ctx context.Context, opt WriteOptions, operation ...WriteOperation) []WriteResult
}

// WriteConfig is the configuration for the `Persistence.Write` operation.
type WriteConfig struct {
	// Separation is the model separation that will be used for the write operation. The `CombinedModels` is the
	// default separation if not set.
	Separation ModelSeparation `json:"separation,omitempty"`

	// AdditionalProperties are custom setting/config specific for the `Persistence` operation.
	AdditionalProperties map[string]any
}

type WriteOptions struct {
	// Tx is a optional transaction that the write operation shall be performed in.
	Tx *TransactionImpl
	// Config is where any common or `Persistence` specific configuration is set.
	Config WriteConfig
}

type WriteOperationConfig struct {
	// Separation is the model separation that will be used for the write operation. If not set, it will use
	// the `WriteOptions.Config.Separation` or the default in the `Persistence`. If neither, it will use
	// `CombinedModels`.
	Separation ModelSeparation `json:"separation,omitempty"`
	// AdditionalProperties are custom setting/config specific for the current `WriteOperation` and `Persistence`.
	AdditionalProperties map[string]any
}

type WriteOperation struct {
	// ClientID is a optional clientID (if `Persistence` supports it). When client sets,
	// this it is often needed to be a unique identifier for each transaction or write operation.
	ClientID string
	// ID is a unique identifier e.g. MyCar 22 or a UUID.
	ID PersistenceID
	// Model is the model that will be written.
	Model any
	// Version is the version of the model that will be written. If this version is not matching, it
	// will fail with a `PersistenceError` with code 409 (Conflict). This is the version read from the
	// `Persistence.Read` operation.
	//
	// The version will always be updated with 1 when the model was successfully written.
	Version int64
	// Config is where any common or `Persistence` specific configuration is set.
	Config WriteOperationConfig
}

type WriteResult struct {
	// ID is the id of the model that was written.
	ID PersistenceID
	// Version is the version of the model that was written.
	Version int64
	// TimeStamp is the timestamp of the model that was written. This is the main timestamp that gets updated
	// each time a model was created or updated. It is a Unix64 bit nanosecond timestamp.
	TimeStamp int64
	// Error is set when the operation failed.
	Error error
}
