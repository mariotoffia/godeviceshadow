package model

import "context"

// ReadonlyPersistence is when a persistence only allows for listing and reading of models but no writes.
type ReadonlyPersistence interface {
	// List will list models from the persistent storage. It will only return an error
	// if there was a problem in the persistence or if the `Persistence` does not support
	// `ListOptions.SearchExpr` or incorrect `ListOptions.Token` is set.
	List(ctx context.Context, opt ListOptions) ([]ListResult, error)
	// Read will read a model from the persistent storage. If not found, it will return
	// a `PersistenceError` with code 404 (Not Found).
	//
	// NOTE: It will return a separate error for each _id_ even if e.g. the storage is completely down and
	// thus each _id_ will return the same error.
	Read(ctx context.Context, opt ReadOptions, id ...PersistenceID) []ReadResult
}

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

type PersistenceID struct {
	// Is a unique identifier e.g. MyCar 22 or a UUID.
	ID string
	// Name is the persistence name so it is possible to have multiple device shadows.
	Name string
	// ModelType is the model type that this `PersistenceID` refers to.
	ModelType ModelType
}

type WriteOptions struct {
	// Custom are custom setting specific for the `Persistence` operation.
	Custom map[string]any
}

type ReadOptions struct {
	// Custom are custom setting specific for the `Persistence` operation.
	Custom map[string]any
}

type ListOptions struct {
	// ID is the id of the models to list. Under a ID there may be many named models.
	// If omitted and the SearchExpr is omitted, all IDs and their named models will be returned.
	ID string
	// SearchExpr is a search expression to use to filter the IDs. This may not be supported by
	// the `Persistence`.
	SearchExpr string
	// Token is a ID that the `Persistence` did return when there's a additional page to be fetched of results.
	Token string
	// Custom are custom setting specific for the `Persistence` operation.
	Custom map[string]any
}

// ListResult is returned when  list operation is performed.
type ListResult struct {
	// ID is the id of the model
	ID PersistenceID
	// Version is a optional version of the model that was listed (not all `Persistence` will return this and thus
	// set it to -1).
	Version int64
	// TimeStamp is a optional timestamp of the model that was listed (not all `Persistence` will return this and thus
	// set it to -1).
	TimeStamp int64
	// Token is set to "something" when there's a additional page to be fetched of results.
	Token string
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

// ReadResult is returned when a model is read from the persistent storage.
type ReadResult struct {
	// ID is the id of the model that was read.
	ID PersistenceID
	// Model is the model that was read.
	Model any
	// Version is the version of the model that was read.
	Version int64
	// TimeStamp is the timestamp of the model that was read. This is the main timestamp that gets updated
	// each time a model was created or updated. It is a Unix64 bit nanosecond timestamp.
	TimeStamp int64
	// Error is set when the operation failed.
	Error error
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

// ModelType stipulates the model type in e.g. a `PersistenceID`.
type ModelType int

const (
	// ModelTypeReported is a the reported portion.
	ModelTypeReported ModelType = 1
	// ModelTypeDesired is the desired portion.
	ModelTypeDesired ModelType = 2
)

// PersistenceError is returned when error in persistence has occurred.
// It uses http error codes but the `Message` is a human readable message.
type PersistenceError struct {
	// Code is a _HTTP_ error code.
	Code int
	// Custom is a custom code to be more precise.
	Custom int
	// Message is a custom message that is human readable.
	Message string
}

func (e PersistenceError) Error() string {
	return e.Message
}

func Error400(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 400, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 400, Message: message}
}

func Error404(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 404, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 404, Message: message}
}

func Error409(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 409, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 409, Message: message}
}

func Error500(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 500, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 500, Message: message}
}
