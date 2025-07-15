package awsdsql

import (
	"context"
	"fmt"
	"maps"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
	"github.com/aws/aws-sdk-go-v2/service/dsql/types"
)

// https://github.com/aws-samples/aurora-dsql-samples/blob/main/go/cluster_management/internal/util/find_cluster.go

type ClusterManagerImpl struct {
	// Region to place the administrative cluster instance. This may be one
	// of many regions the cluster may participate in - however, this is where
	// the cluster is primarily managed.
	Region string
	// Config is the AWS SDK configuration used to create clients.
	Config *aws.Config
	// Client is the cached client used
	Client *dsql.Client
	// cache is a map of clusters keyed by their identifier.
	cache map[string]dsql.GetClusterOutput
}

type NewClusterManagerOptions struct {
	// Region is the AWS region where the DSQL clusters are managed. If not
	// set it will use the default region (if set).
	Region string
	// Config, if set, will be used and the _Region_ will be ignored.
	Config *aws.Config
}

func NewClusterManager(opts ...NewClusterManagerOptions) *ClusterManagerImpl {
	var opt NewClusterManagerOptions

	if len(opts) > 0 {
		opt = opts[0]
	}

	return &ClusterManagerImpl{Region: opt.Region, Config: opt.Config}
}

// FetchClusters will populate clusters from the AWS DSQL service and they
// are cached in the `ClusterManagerImpl` instance for queries. If an error occurs, the
// internal cache is not modified and `nil` is returned along the `error`.
func (cm *ClusterManagerImpl) FetchClusters(ctx context.Context) (map[string]dsql.GetClusterOutput, error) {
	if err := cm.ensureClient(ctx); err != nil {
		return nil, err
	}

	input := &dsql.ListClustersInput{}
	result := map[string]dsql.GetClusterOutput{}

	for {
		clusters, err := cm.Client.ListClusters(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, val := range clusters.Clusters {
			if cluster, err := cm.Client.GetCluster(ctx, &dsql.GetClusterInput{Identifier: val.Identifier}); err != nil {
				return nil, err
			} else {
				// Store the cluster in the result map
				result[*val.Identifier] = *cluster
			}
		}

		if clusters.NextToken == nil || len(*clusters.NextToken) == 0 {
			break
		}
	}

	cm.cache = result

	return result, nil
}

// ByIdentifier returns the cluster identified by _identifier_ if it exists in the cache.
//
// It returns the cluster and a boolean indicating whether the cluster was found.
//
// TIP: Use the `FetchClusters` method to populate the cache before calling this method.
func (cm *ClusterManagerImpl) ByIdentifier(identifier string) (dsql.GetClusterOutput, bool) {
	if cm.cache == nil {
		return dsql.GetClusterOutput{}, false
	}

	cluster, ok := cm.cache[identifier]
	return cluster, ok
}

// Intractable returns all _tags_ filtered clusters that are intractable (active, idle, inactive). That is,
// clusters that are or will automatically be active. If no _tags_ all clusters are considered for active, idle, and inactive.
func (cm *ClusterManagerImpl) Intractable(tags ...map[string]string) map[string]dsql.GetClusterOutput {
	if len(tags) > 0 {
		return cm.Clusters(tags[0], types.ClusterStatusActive, types.ClusterStatusIdle, types.ClusterStatusInactive)
	} else {
		return cm.Clusters(nil, types.ClusterStatusActive, types.ClusterStatusIdle, types.ClusterStatusInactive)
	}
}

// NonIntractable returns all _tags_ filtered clusters that are non-intractable (creating, updating,
// deleting, deleted, failed, pending setup/delete).
func (cm *ClusterManagerImpl) NonIntractable(tags ...map[string]string) map[string]dsql.GetClusterOutput {
	statuses := []types.ClusterStatus{
		types.ClusterStatusCreating,
		types.ClusterStatusUpdating,
		types.ClusterStatusDeleting,
		types.ClusterStatusDeleted,
		types.ClusterStatusFailed,
		types.ClusterStatusPendingSetup,
		types.ClusterStatusPendingDelete,
	}

	if len(tags) > 0 {
		return cm.Clusters(tags[0], statuses...)
	} else {
		return cm.Clusters(nil, statuses...)
	}
}

// Cache returns a shallow copy of the internal cache of clusters.
func (cm *ClusterManagerImpl) Cache() map[string]dsql.GetClusterOutput {
	return maps.Clone(cm.cache)
}

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
