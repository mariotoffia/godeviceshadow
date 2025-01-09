package dynamodbpersistence

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

type Config struct {
	// Table is the name of the DynamoDB table to use.
	Table string `json:"table"`
	// ModelSeparation determines whether the desired and reported models should be stored separately or not.
	//
	// If nothing is set it will default `CombinedModels`. This may be overridden by the `separate` key in the
	// write request `Config` using _separation_ key.
	//
	// If `CombinedModels` and not both reported, desired models are passed in `Write`, it will return 400 (Bad Request).
	ModelSeparation persistencemodel.ModelSeparation `json:"separate"`
	// MaxReadBatchSize is the maximum number of items to read in a single batch. If read exceeds this number, it
	// will be split into multiple requests.
	MaxReadBatchSize int `json:"batch"`
}
