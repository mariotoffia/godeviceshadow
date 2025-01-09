package dynamodbpersistence

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Config struct {
	// Table is the name of the DynamoDB table to use.
	Table string `json:"table"`
	// ModelSeparation determines whether the desired and reported models should be stored separately or not.
	//
	// If nothing is set it will default `CombinedModels`. This may be overridden by the `separate` key in the
	// write request `Config` using _separation_ key.
	//
	// If `CombinedModels` and not both reported, desired models are passed in `Write`, it will return 400 (Bad Request).
	//
	// NOTE: This may be overridden in `WriteOptions` for all `WriteOperations` or per `WriteOperation` by setting the _Config_
	// with key `persistencemodel.ModelSeparationConfigKey` and value of `persistencemodel.ModelSeparation`.
	ModelSeparation persistencemodel.ModelSeparation `json:"separate"`
	// MaxReadBatchSize is the maximum number of items to read in a single batch. If read exceeds this number, it
	// will be split into multiple requests.
	MaxReadBatchSize int `json:"batch"`
	// MaxParallelism is the maximum number of parallel requests to make to DynamoDB. If the number of items to write
	// exceeds this, it will be queued up and processed in parallel.
	//
	// It defaults to 1, i.e. no parallelism.
	MaxParallelism int `json:"parallel"`
	// AwsConfig to use when creating the client.
	AwsConfig aws.Config `json:"aws"`
}
