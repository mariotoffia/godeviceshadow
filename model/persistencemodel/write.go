package persistencemodel

import "context"

// Persistence allows a client to read and write a model to a persistent storage.
type Persistence interface {
	ReadonlyPersistence
	// Write will write a _model_ onto the persistent storage. If the model already exists, it will
	// update the model. If the model does not exist, it will create a new model.
	//
	// If a conflict occurs, it will return a `PersistenceError` with code 409 (Conflict).
	Write(ctx context.Context, opt WriteOptions, operation ...WriteOperation) []WriteResult
	// Delete will delete one or more models from the persistent storage. If the model does not exist, it will
	// return a `PersistenceError` with code 404 (Not Found).
	//
	// When `WriteOperation.Version` is set to zero or less, it will delete the model regardless of the version.
	Delete(ctx context.Context, opt WriteOptions, operation ...WriteOperation) []WriteResult
}

type WriteOptions struct {
	// Tx is a optional transaction that the write operation shall be performed in.
	Tx *Transaction
	// Config are custom setting/config specific for the `Persistence` operation.
	Config map[string]any
}

type WriteOperation struct {
	// ClientID is a optional clientID (if `Persistence` supports it). When client sets,
	// this it is often needed to be a unique identifier for each transaction or write operation.
	ClientID string
	// ID is a unique identifier e.g. MyCar 22 or a UUID.
	ID PersistenceID
	// Model is the model that will be written.
	Model any
	// ModelType is the model type that this `PersistenceID` refers to.
	ModelType ModelType
	// Version is the version of the model that will be written. If this version is not matching, it
	// will fail with a `PersistenceError` with code 409 (Conflict). This is the version read from the
	// `Persistence.Read` operation.
	//
	// If a "real" transaction is used by the `Persistence` implementation, this is not the main distiguisher
	// if the model will be created/updated but the semantics is still upheld.
	//
	// This is primarily for implementations that uses optimistic locking.
	Version int64
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
