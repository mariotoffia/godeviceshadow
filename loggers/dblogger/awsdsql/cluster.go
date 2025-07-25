package awsdsql

import (
	"context"
	"maps"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
	"github.com/aws/aws-sdk-go-v2/service/dsql/types"
)

// https://github.com/aws-samples/aurora-dsql-samples/blob/main/go/cluster_management/internal/util/find_cluster.go

type ClusterManagerImpl struct {
	// Region to place the administrative cluster instance. This may be one
	// of many regions the cluster may participate in - however, this is where
	// the cluster is primarily managed.
	//
	// If not set, the default region will be used.
	Region string
	// Config is the AWS SDK configuration used to create clients.
	//
	// If not, set, the default AWS SDK configuration will be used.
	Config *aws.Config
	// Client is the client created and used in order to not create a new client for each request.
	Client *dsql.Client
	// cache is a map of clusters keyed by their identifier to avoid unnecessary API roundtrips.
	cache map[string]dsql.GetClusterOutput
}

// NonIntractableStatuses are statuses that the `NonIntractable` method will regard
// as non-intractable.
var NonIntractableStatuses = []types.ClusterStatus{
	types.ClusterStatusCreating,
	types.ClusterStatusUpdating,
	types.ClusterStatusDeleting,
	types.ClusterStatusDeleted,
	types.ClusterStatusFailed,
	types.ClusterStatusPendingSetup,
	types.ClusterStatusPendingDelete,
}

// IntractableStatuses are statuses that the `Intractable` method will regard
// as intractable.
var IntractableStatuses = []types.ClusterStatus{
	types.ClusterStatusActive,
	types.ClusterStatusIdle,
	types.ClusterStatusInactive,
}

type NewClusterManagerOptions struct {
	// Region is the AWS region where the DSQL clusters are managed. If not
	// set it will use the default region (if set).
	Region string
	// Config, if set, will be used and the _Region_ will be ignored.
	Config *aws.Config
}

// NewClusterManager creates a new ClusterManagerImpl instance with the provided options.
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

// Clusters returns a map of clusters that match the specified tags and the status.
//
// NOTE: It will use the internal cached clusters from `FetchClusters` method. If you want to
// ensure the clusters are up-to-date, you should call `FetchClusters` before calling this method.
//
// If no _status_ is provided, it will return all clusters matching the _tags_. If no _tags_ it will
// set it to auto include all. Thus to get all clusters, just leave the _tags_ and _status_ empty.
//
// TIP: The _status_ are ORed together and the _tags_ are ANDed together.
func (cm *ClusterManagerImpl) Clusters(tags map[string]string, status ...types.ClusterStatus) map[string]dsql.GetClusterOutput {
	result := map[string]dsql.GetClusterOutput{}

	// Checks that the _cluster_ has all the tags and tag value specified.
	hasAllTagAndValue := func(cluster dsql.GetClusterOutput) bool {
		if len(tags) == 0 {
			return true // No tags specified, include all clusters
		}
		for key, value := range tags {
			if tagValue, ok := cluster.Tags[key]; !ok || tagValue != value {
				return false
			}
		}

		return true
	}

	hasStatus := func(cluster dsql.GetClusterOutput) bool {
		if len(status) == 0 {
			return true // No status specified, include all clusters
		}

		return slices.ContainsFunc(status, func(s types.ClusterStatus) bool {
			return cluster.Status == s
		})
	}

	for _, cluster := range cm.cache {
		if hasStatus(cluster) && hasAllTagAndValue(cluster) {
			result[*cluster.Identifier] = cluster
		}
	}

	return result
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
		return cm.Clusters(tags[0], IntractableStatuses...)
	} else {
		return cm.Clusters(nil, IntractableStatuses...)
	}
}

// NonIntractable returns all _tags_ filtered clusters that are non-intractable (creating, updating,
// deleting, deleted, failed, pending setup/delete).
func (cm *ClusterManagerImpl) NonIntractable(tags ...map[string]string) map[string]dsql.GetClusterOutput {

	if len(tags) > 0 {
		return cm.Clusters(tags[0], NonIntractableStatuses...)
	} else {
		return cm.Clusters(nil, NonIntractableStatuses...)
	}
}

// Cache returns a shallow copy of the internal cache of clusters.
func (cm *ClusterManagerImpl) Cache() map[string]dsql.GetClusterOutput {
	return maps.Clone(cm.cache)
}
