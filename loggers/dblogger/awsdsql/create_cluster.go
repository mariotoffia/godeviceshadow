package awsdsql

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
	"github.com/aws/aws-sdk-go-v2/service/dsql/types"
)

type CreateClusterOptions struct {
	// Identifier is the unique identifier for the cluster. It must be unique
	// across all clusters in the AWS account and region.
	Identifier string
	// The Region that serves as the witness region for a multi-Region cluster.
	// The witness Region helps maintain cluster consistency and quorum. It stores a limited
	// window of encrypted transaction logs, which is used to improve multi-Region durability
	// and availability. Multi-Region witness Regions do not have endpoints.
	//
	// NOTE: This has to be separate from the _Regions_ property but within a region set (see _Regions_).
	Witness string
	// Regions are which regions that serves active connections. Where the witness (if any)
	// is the one that stores encrypted transaction logs but doesn’t provide client endpoints.
	//
	// CAUTION: The regions (and witness) must be specified in a region set
	// see https://docs.aws.amazon.com/aurora-dsql/latest/userguide/what-is-aurora-dsql.html for
	// details.
	//
	// When the number of regions is one it is a single region cluster, if it is two it is a valid
	// multi-region cluster. Three and more is not supported.
	Regions []string
	// Tags are the tags to apply to the cluster. These are used for filtering
	// clusters in the `ClusterManagerImpl` instance.
	Tags map[string]string
	// NonWitnessTags are any optional tags on the non-witness cluster participant, matching _Regions_,
	// property name. If not found here, it will revert to the _Tags_ property for those as well.
	NonWitnessTags map[string]map[string]string
	// DeleteProtect, if set to true, will enable delete protection for the cluster.
	DeleteProtect bool
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

// CreateCluster creates a new cluster with the specified options.
//
// CAUTION: This can take up to 10 minutes to complete for multi-region clusters and up to 5 minutes for single region clusters.
//
// It may return an error along with a map of region keys and the cluster identifiers created. This is to ensure
// cleanup if necessary. This is when the wait of the cluster to become active fails but the actual creation worked. Sometime
// cluster one may been created but cluster two may fail, only the first region cluster identifier will be present
// in the map.
func (cm *ClusterManagerImpl) CreateCluster(ctx context.Context, opts CreateClusterOptions) (map[string]string, error) {
	switch len(opts.Regions) {
	case 1:
		if id, err := cm.createClusterSingle(ctx, opts); id != "" {
			return map[string]string{opts.Regions[0]: id}, err
		} else {
			return nil, err
		}
	case 2:
		return cm.createClusterMulti(ctx, opts)
	}

	return nil, fmt.Errorf("invalid number of regions specified: %d, must be 1 or 2", len(opts.Regions))
}

// CreateCluster creates a new cluster with the specified options.
//
// It will return the identifier of the cluster created or an error if it fails.
//
// If it succeeded creating but not waiting, it will return the identifier of the cluster and a error.
//
// CAUTION: This can take up to 5 minutes to complete.
func (cm *ClusterManagerImpl) createClusterSingle(ctx context.Context, opts CreateClusterOptions) (string, error) {
	if err := cm.ensureClient(ctx); err != nil {
		return "", err
	}

	if len(opts.Regions) != 1 {
		return "", fmt.Errorf("invalid number of regions specified: %d, must be 1 for single region cluster", len(opts.Regions))
	}

	client, err := cm.clientInRegion(opts.Regions[0])

	if err != nil {
		return "", fmt.Errorf("failed to create client for region %s: %v", opts.Regions[0], err)
	}

	clusterProperties, err := client.CreateCluster(ctx, &dsql.CreateClusterInput{
		DeletionProtectionEnabled: &opts.DeleteProtect,
		Tags:                      opts.Tags,
	})

	if err != nil {
		return "", fmt.Errorf("error creating cluster: %w", err)
	}

	// Create the waiter with our custom options
	waiter := dsql.NewClusterActiveWaiter(client, func(o *dsql.ClusterActiveWaiterOptions) {
		o.MaxDelay = 30 * time.Second
		o.MinDelay = 10 * time.Second
		o.LogWaitAttempts = true
	})

	id := clusterProperties.Identifier

	// Create the input for the clusterProperties
	getInput := &dsql.GetClusterInput{
		Identifier: id,
	}

	// Wait for the cluster to become active
	err = waiter.Wait(ctx, getInput, 5*time.Minute)
	if err != nil {
		return *id, fmt.Errorf("error waiting for cluster to become active: %w", err)
	}

	return *id, nil
}

// createClusterMulti creates a multi-region cluster with two regions specified in the options.
//
// If any errors occurs, it will return an error and possibly a map with the instances created so
// it is possible to clean up the resources created.
//
// If it succeeds, it will return a map with the identifiers of the clusters created in each region and `nil` error.
//
// CAUTION: This can take up to 10 minutes to complete.
func (cm *ClusterManagerImpl) createClusterMulti(ctx context.Context, opts CreateClusterOptions) (map[string]string, error) {
	if len(opts.Regions) != 2 {
		return nil, fmt.Errorf(
			"invalid number of regions specified: %d, must be 2 for multi-region cluster",
			len(opts.Regions),
		)
	}

	getTags := func(region string) map[string]string {
		if opts.NonWitnessTags != nil {
			if tags, ok := opts.NonWitnessTags[region]; ok {
				return tags
			}
		}

		return opts.Tags
	}

	if err := cm.ensureClient(ctx); err != nil {
		return nil, err
	}

	// Create first cluster with no peer (will update it later)
	input := &dsql.CreateClusterInput{
		DeletionProtectionEnabled: &opts.DeleteProtect,
		Tags:                      getTags(opts.Regions[0]),
	}

	// witness is specified -> use it
	if len(opts.Witness) > 0 {
		input.MultiRegionProperties = &types.MultiRegionProperties{
			WitnessRegion: aws.String(opts.Witness),
		}
	}

	client1, err := cm.clientInRegion(opts.Regions[0])

	if err != nil {
		return nil, fmt.Errorf("failed to create client for region %s: %v", opts.Regions[0], err)
	}

	node1, err := client1.CreateCluster(ctx, input)

	if err != nil {
		return nil, fmt.Errorf("failed to create first cluster: %v", err)
	}

	// Second cluster, set first cluster as peer
	input2 := &dsql.CreateClusterInput{
		DeletionProtectionEnabled: &opts.DeleteProtect,
		MultiRegionProperties: &types.MultiRegionProperties{
			Clusters: []string{*node1.Arn},
		},
		Tags: getTags(opts.Regions[1]),
	}

	// if witness -> set it
	if opts.Witness != "" {
		input2.MultiRegionProperties.WitnessRegion = aws.String(opts.Witness)
	}

	client2, err := cm.clientInRegion(opts.Regions[1])
	node2, err := client2.CreateCluster(ctx, input2)

	if err != nil {
		return nil, fmt.Errorf("failed to create second cluster: %v", err)
	}

	// update first cluster with the second cluster as a peer
	update := dsql.UpdateClusterInput{
		Identifier: node1.Identifier,
		MultiRegionProperties: &types.MultiRegionProperties{
			Clusters: []string{*node2.Arn},
		}}

	// if witness -> set it
	if opts.Witness != "" {
		update.MultiRegionProperties.WitnessRegion = aws.String(opts.Witness)
	}

	_, err = client1.UpdateCluster(ctx, &update)

	if err != nil {
		return nil, fmt.Errorf("failed to update cluster to associate with first cluster. %v", err)
	}

	// Create the waiter with our custom options for first cluster
	waiter := dsql.NewClusterActiveWaiter(client1, func(o *dsql.ClusterActiveWaiterOptions) {
		o.MaxDelay = 30 * time.Second // Creating a multi-region cluster can take a few minutes
		o.MinDelay = 10 * time.Second
		o.LogWaitAttempts = true
	})

	// Now that multiRegionProperties is fully defined for both clusters
	// they'll begin the transition to ACTIVE

	// Create the input for the clusterProperties to monitor for first cluster
	getInput := &dsql.GetClusterInput{
		Identifier: node1.Identifier,
	}

	// Wait for the first cluster to become active - max 5 minutes
	err = waiter.Wait(ctx, getInput, 5*time.Minute)
	if err != nil {
		return map[string]string{
			opts.Regions[0]: *node1.Identifier,
			opts.Regions[1]: *node2.Identifier,
		}, fmt.Errorf("error waiting for first cluster to become active: %w", err)
	}

	// Create the waiter with our custom options
	waiter2 := dsql.NewClusterActiveWaiter(client2, func(o *dsql.ClusterActiveWaiterOptions) {
		o.MaxDelay = 30 * time.Second // Creating a multi-region cluster can take a few minutes
		o.MinDelay = 10 * time.Second
		o.LogWaitAttempts = true
	})

	// Create the input for the clusterProperties to monitor for second
	getInput2 := &dsql.GetClusterInput{
		Identifier: node2.Identifier,
	}

	// Wait for the second cluster to become active - max 5 minutes
	err = waiter2.Wait(ctx, getInput2, 5*time.Minute)
	if err != nil {
		return map[string]string{
			opts.Regions[0]: *node1.Identifier,
			opts.Regions[1]: *node2.Identifier,
		}, fmt.Errorf("error waiting for second cluster to become active: %w", err)
	}

	return map[string]string{
		opts.Regions[0]: *node1.Identifier,
		opts.Regions[1]: *node2.Identifier,
	}, nil
}
