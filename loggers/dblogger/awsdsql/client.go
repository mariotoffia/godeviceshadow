package awsdsql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
)

// Will return a `dsql.Client` for the specified region. It will copy the existing configuration
// and set to the specified region. If the region is the same as the existing configuration.
//
// It will not cache the client, instead creates a new client each time this method is called.
func (cm *ClusterManagerImpl) clientInRegion(region string) (*dsql.Client, error) {
	if cm.Config == nil {
		return nil, fmt.Errorf("config is not initialized")
	}

	if cm.Config.Region == region {
		return cm.Client, nil
	}

	cfg := cm.Config.Copy()
	cfg.Region = region

	return dsql.NewFromConfig(cfg), nil
}

// ensureClient ensures that the `Client` is initialized. It will use the `Config` if set,
// otherwise it will load the default config.
func (cm *ClusterManagerImpl) ensureClient(ctx context.Context) error {
	if cm.Client != nil {
		return nil
	}

	if cm.Config == nil {
		var opts []func(*config.LoadOptions) error

		if cm.Region != "" {
			opts = append(opts, config.WithRegion(cm.Region))
		}

		if cfg, err := config.LoadDefaultConfig(ctx, opts...); err != nil {
			return err
		} else {
			cm.Config = &cfg
		}
	}

	cm.Client = dsql.NewFromConfig(*cm.Config)

	return nil
}
