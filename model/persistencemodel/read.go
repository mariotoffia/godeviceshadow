package persistencemodel

import (
	"context"
	"reflect"
)

// ReadonlyPersistence is when a persistence only allows for listing and reading of models but no writes.
type ReadonlyPersistence interface {
	// List will list models from the persistent storage. It will only return an error
	// if there was a problem in the persistence or if the `Persistence` does not support
	// `ListOptions.SearchExpr` or incorrect `ListOptions.Token` is set.
	List(ctx context.Context, opt ListOptions) ([]ListResult, error)
	// Read will read a model from the persistent storage. If not found, it will return
	// a `PersistenceError` with code 404 (Not Found).
	//
	// NOTE: It will return a separate error for each _operation_ even if e.g. the storage is completely down and
	// thus each _id_ will return the same error.
	Read(ctx context.Context, opt ReadOptions, operation ...ReadOperation) []ReadResult
}

// ReadOperation is a read operation that will be performed.
type ReadOperation struct {
	// ID is a unique identifier e.g. MyCar 22 or a UUID.
	ID PersistenceID
	// Model is the model `type` that will be read.
	Model reflect.Type
	// Version is the version of the model that will be read. If 0 or less it will be ignored,
	// otherwise it will only return the model with the version that matches the `Version`
	// or a `PersistenceError` with code 404 (Not Found) is returned.
	Version int64
}

// ReadConfig is the configuration for the `Persistence.Read` operation.
type ReadConfig struct {
	// AdditionalProperties are custom setting/config specific for the `Persistence` operation.
	AdditionalProperties map[string]any
}

type ReadOptions struct {
	// Tx is a optional transaction that the read operation shall be performed in.
	Tx *Transaction
	// Config is where any common or `Persistence` specific configuration is set.
	Config ReadConfig
}

// ListConfig is the configuration for the `Persistence.List` operation.
type ListConfig struct {
	// AdditionalProperties are custom setting/config specific for the `Persistence` operation.
	AdditionalProperties map[string]any
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
	// Config is where any common or `Persistence` specific configuration is set.
	Config ListConfig
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
	// ClientToken is a optional last client token write operation (not all `Persistence` will return this and thus
	// set it to "").
	//
	// NOTE: if no client token has been used in the write operation, it will not be persisted and hence not visible here either.
	ClientToken string
	// Token is set to "something" when there's a additional page to be fetched of results.
	Token string
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
	// ClientToken is the last client token write operation (if any)
	ClientToken string
	// Error is set when the operation failed.
	Error error
}
