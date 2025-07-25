package awsdsql

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dsql"
)

// DeleteCluster will delete a single region cluster identifier by the _identifier_ in the specified _region_.
// It returns an error if the deletion fails or if the cluster does not exist.
//
// If the _wait_ parameter is set to true, it will wait for the cluster to be deleted before returning. This
// may take considerable time depending on the cluster size and AWS region.
//
// TIP: If multi-region cluster, you may delete the cluster one by one and set _wait_ to `false` and then call,
// `WaitForDeletion` to wait for each _identifier_ to be deleted.
//
// If the cluster is delete protected, use the `RemoveDeleteProtection` method before calling this method.
func (cm *ClusterManagerImpl) DeleteCluster(ctx context.Context, region, identifier string, wait ...bool) error {
	if err := cm.ensureClient(ctx); err != nil {
		return err
	}

	if region == "" {
		return fmt.Errorf("region must be specified")
	}

	if identifier == "" {
		return fmt.Errorf("identifier must be specified")
	}

	client, err := cm.clientInRegion(region)

	if err != nil {
		return err
	}

	// Delete the cluster
	_, err = client.DeleteCluster(ctx, &dsql.DeleteClusterInput{
		Identifier: &identifier,
	})
	if err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	if len(wait) > 0 || wait[0] {
		return cm.WaitForDeletion(ctx, region, identifier)
	}

	return nil
}

// RemoveDeleteProtection will remove deletion protection from the cluster identified by _identifier_ in the specified _region_.
//
// It returns an error if the removal fails or if the cluster does not exist.
func (cm *ClusterManagerImpl) RemoveDeleteProtection(ctx context.Context, region, identifier string) error {
	if err := cm.ensureClient(ctx); err != nil {
		return err
	}

	if region == "" {
		return fmt.Errorf("region must be specified")
	}

	if identifier == "" {
		return fmt.Errorf("identifier must be specified")
	}

	client, err := cm.clientInRegion(region)

	if err != nil {
		return err
	}

	var deletionProtectionEnabled bool

	// Remove deletion protection
	update := dsql.UpdateClusterInput{
		Identifier:                &identifier,
		DeletionProtectionEnabled: &deletionProtectionEnabled,
	}

	_, err = client.UpdateCluster(ctx, &update)

	return err
}

func (cm *ClusterManagerImpl) WaitForDeletion(ctx context.Context, region, identifier string) error {
	if err := cm.ensureClient(ctx); err != nil {
		return err
	}

	if region == "" {
		return fmt.Errorf("region must be specified")
	}

	if identifier == "" {
		return fmt.Errorf("identifier must be specified")
	}

	client, err := cm.clientInRegion(region)

	if err != nil {
		return err
	}

	// Create waiter to check cluster deletion
	waiter := dsql.NewClusterNotExistsWaiter(client, func(options *dsql.ClusterNotExistsWaiterOptions) {
		options.MinDelay = 10 * time.Second
		options.MaxDelay = 30 * time.Second
		options.LogWaitAttempts = true
	})

	err = waiter.Wait(
		ctx, &dsql.GetClusterInput{Identifier: &identifier}, 5*time.Minute,
	)

	if err != nil {
		return fmt.Errorf("error waiting for cluster to be deleted: %w", err)
	}

	return nil
}
