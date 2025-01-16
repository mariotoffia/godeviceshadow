package dynamodbpersistence

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Config struct {
	// Table is the name of the DynamoDB table to use. It is *required* to be a valid table.
	Table string `json:"table"`
	// Client to use when interacting with DynamoDB. If not set, it will create one from the default config.
	Client *dynamodb.Client `json:"-"`
	// ModelSeparation determines whether the desired and reported models should be stored separately or not.
	//
	// If nothing is set it will default `CombinedModels`. This may be overridden by the `separate` key in the
	// write request `Config` using _separation_ key.
	//
	// If `CombinedModels` and not both reported, desired models are passed in `Write`, it will return 400 (Bad Request).
	//
	// NOTE: This may be overridden in `WriteOptions` for all `WriteOperations` or per `WriteOperation` by setting the _Config_
	// with key `persistencemodel.ModelSeparationConfigKey` and value of `persistencemodel.ModelSeparation`.
	ModelSeparation persistencemodel.ModelSeparation `json:"separation,omitempty"`
	// MaxReadBatchSize is the maximum number of items to read in a single batch. If read exceeds this number, it
	// will be split into multiple requests. Default is 100.
	MaxReadBatchSize int `json:"read_batch,omitempty"`
	// MaxReadRetries is the maximum number of retries to make when reading from DynamoDB. Default is 3.
	//
	// This is when it return unprocessed keys and it will retry the request. All other errors are not retried.
	MaxReadRetries int `json:"read_retries,omitempty"`
	// MaxWriteParallelism is the maximum number of parallel requests to make to DynamoDB. If the number of items to write
	// exceeds this, it will be queued up and processed in parallel.
	//
	// It defaults to 1, i.e. no parallelism.
	MaxWriteParallelism int `json:"write_parallelism,omitempty"`
	// MaxWriteBatchSize is the maximum number of items to write in a single batch. If write exceeds this number, it
	// will be split into multiple requests. Default is 25. This is also used for delete operations.
	MaxWriteBatchSize int `json:"write_batch,omitempty"`
	// MaxWriteRetries is the maximum number of retries to make when writing to DynamoDB. Default is 3.
	//
	// This is when it return unprocessed keys and it will retry the request. All other errors are not retried.
	MaxWriteRetries int `json:"write_retries,omitempty"`
}
