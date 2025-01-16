package dynamodbpersistence

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// Persistence is the DynamoDB persistence plugin.
type Persistence struct {
	config Config
	client *dynamodb.Client
}

// New creates a new DynamoDB persistence plugin.
func New(ctx context.Context, config Config) (*Persistence, error) {
	cfg := config

	if cfg.MaxWriteParallelism <= 0 {
		cfg.MaxWriteParallelism = 1
	}

	if cfg.MaxReadBatchSize <= 0 {
		cfg.MaxReadBatchSize = 25
	}

	if cfg.ModelSeparation == 0 {
		cfg.ModelSeparation = persistencemodel.CombinedModels
	}

	if cfg.AwsConfig.Region == "" {
		return nil, persistencemodel.Error400("not a valid AWS configuration")
	}

	return &Persistence{
		config: cfg,
		client: dynamodb.NewFromConfig(config.AwsConfig),
	}, nil
}

// PersistenceObject is the object that is stored in the persistence layer.
//
// Depending on the configuration of the persistence layer, the desired and reported models may be
// stored separately.
type PersistenceObject struct {
	Version int64 `json:"version"`
	// TimeStamp is a unix64 nano timestamp with UTC time
	TimeStamp int64 `json:"timestamp"`
	// ClientToken is a unique token for the client that initiated the request
	ClientToken string `json:"clientToken,omitempty"`
	// Desired is the desired model (if such is present). Depending on how the persistence is
	// configured, this may be stored separately from the reported model.
	Desired any `json:"desired,omitempty"`
	// Reported is the reported model (if such is present). Depending on how the persistence is
	// configured, this may be stored separately from the desired model.
	Reported any `json:"reported,omitempty"`
}

// PartialPersistenceObject is a object that can be used when desired and reported is not wanted (just metadata).
type PartialPersistenceObject struct {
	// Version is the version of the object
	Version int64 `json:"version"`
	// TimeStamp is a unix64 nano timestamp with UTC time
	TimeStamp int64 `json:"timestamp"`
	// ClientToken is a unique token for the client that initiated the request
	ClientToken string `json:"clientToken,omitempty"`
}
