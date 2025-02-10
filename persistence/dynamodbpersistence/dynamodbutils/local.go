package dynamodbutils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
)

type LocalDynamoDB struct {
	// ContainerID is the ID of the local DynamoDB container.
	ContainerID string
	// Endpoint is the endpoint of the local DynamoDB instance.
	Endpoint string
	// Client is a DynamoDB client configured to use the local DynamoDB instance.
	Client *dynamodb.Client
	// StreamClient is a DynamoDB Streams client configured to use the local DynamoDB instance.
	StreamClient *dynamodbstreams.Client
}

type StartLocalDynamoDbOptions struct {
	// Reuse specifies whether to reuse an existing local DynamoDB instance.
	Reuse bool
	// MaxWarmupWaitTime is the maximum time to wait for the local DynamoDB instance to become ready.
	// Default is _10_ seconds.
	MaxWarmupWaitTime time.Duration
}

// StartLocalDynamoDB will start a local dynamodb and make sure it has  a_table_ to persist in.
//
// This uses docker to start the local dynamodb instance via `docker run -d -p 8000:8000", "amazon/dynamodb-local`.
//
// Use `LocalDynamoDB.Close()` to stop the local dynamodb instance.
//
// The `LocalDynamoDB.Client` and `LocalDynamoDB.StreamClient` are configured to use the local dynamodb instance.
//
// If there are other _amazon/dynamodb-local_ instances running, all are killed and then a new one is started.
//
// If _reuse_ is `true` then it will keep one instance running (if any).
func StartLocalDynamoDB(ctx context.Context, table string, opts ...StartLocalDynamoDbOptions) (*LocalDynamoDB, error) {
	var opt StartLocalDynamoDbOptions

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.MaxWarmupWaitTime == 0 {
		opt.MaxWarmupWaitTime = 10 * time.Second
	}

	containerID, endpoint, err := StartDynamoDBLocal(opt.Reuse)
	if err != nil {
		return nil, fmt.Errorf("Failed to start DynamoDB Local: %w", err)
	}

	cfg, err := DefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS config: %w", err)
	}

	client := DynamoDbClient(cfg, endpoint)
	streamClient := DynamoDbStreamClient(cfg, endpoint)

	if err := WaitForDynamoDB(ctx, client, opt.MaxWarmupWaitTime); err != nil {
		return nil, fmt.Errorf("dynamodb did not become ready: %w", err)
	}

	if _, err := NewTestTableResourceWithClient(ctx, table, client); err != nil {
		return nil, err
	}

	return &LocalDynamoDB{
		ContainerID:  containerID,
		Endpoint:     endpoint,
		Client:       client,
		StreamClient: streamClient,
	}, nil
}

// Close implements the `io.Closer` interface to stop the local DynamoDB container.
func (l *LocalDynamoDB) Close() error {
	if err := StopDynamoDBLocal(l.ContainerID); err != nil {
		return fmt.Errorf("Failed to stop DynamoDB Local: %w", err)
	}

	return nil
}

func DefaultConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		// Required but ignored for local DynamoDB
		config.WithRegion("eu-west-1"),
		// Dummy credentials
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
	)
}

func DynamoDbClient(cfg aws.Config, endpoint string) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}

func DynamoDbStreamClient(cfg aws.Config, endpoint string) *dynamodbstreams.Client {
	return dynamodbstreams.NewFromConfig(cfg, func(o *dynamodbstreams.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}

// StartDynamoDBLocal starts the DynamoDB local container and returns its container ID and endpoint.
//
// It will remove any existing DynamoDB Local containers before starting a new one. When _reuse_ is `true`
// it will keep one (if running).
//
// Use the `WaitForDynamoDB` function to wait for the local DynamoDB instance to become ready.
func StartDynamoDBLocal(reuse bool) (string, string, error) {
	existingContainers, err := ListDynamoDBContainers()

	if err != nil {
		return "", "", fmt.Errorf("failed to list existing DynamoDB containers: %w", err)
	}

	// If more than one instance exists, stop all but one (if reuse is `true`)
	var existingContainerID string

	for _, id := range existingContainers {
		if !reuse {
			StopAndCleanupDynamoDBContainer(id)
			continue
		}

		if running, err := IsContainerRunning(id); err != nil {
			return "", "", fmt.Errorf("failed to check if container is running: %w", err)
		} else if running {
			if existingContainerID != "" {
				StopAndCleanupDynamoDBContainer(id)
			} else {
				existingContainerID = id
			}
		}
	}

	if existingContainerID != "" {
		return existingContainerID, "http://localhost:8000", nil
	}

	// Run Docker command to start container
	cmd := exec.Command("docker", "run", "-d", "-p", "8000:8000", "amazon/dynamodb-local")

	var out bytes.Buffer

	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		return "", "", fmt.Errorf("failed to start DynamoDB Local: %w", err)
	}

	// Get the container ID
	containerID := strings.TrimSpace(out.String())

	// Return container ID and endpoint
	return containerID, "http://localhost:8000", nil
}

// StopDynamoDBLocal stops and removes the given container ID
func StopDynamoDBLocal(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", containerID)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to stop DynamoDB Local: %w", err)
	}

	return nil
}

// ListDynamoDBContainers returns the IDs of all running DynamoDB Local containers.
func ListDynamoDBContainers() ([]string, error) {
	cmd := exec.Command("docker", "ps", "-q", "--filter", "ancestor=amazon/dynamodb-local")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list DynamoDB Local containers: %w", err)
	}

	ids := strings.Fields(out.String()) // Split by whitespace to get container IDs
	return ids, nil
}

// StopAndCleanupDynamoDBContainer stops and removes a specific DynamoDB container and its volumes.
func StopAndCleanupDynamoDBContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", "-v", containerID)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stop and remove DynamoDB Local: %w", err)
	}
	log.Printf("Stopped and removed DynamoDB container: %s", containerID)
	return nil
}

func IsContainerRunning(containerID string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", containerID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to check if container is running: %w", err)
	}

	if strings.TrimSpace(out.String()) == "true" {
		return true, nil
	}

	return false, nil
}

func StartExistingContainer(containerID string) (string, error) {
	cmd := exec.Command("docker", "start", containerID)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to start existing container: %w", err)
	}

	return "http://localhost:8000", nil
}

// WaitForDynamoDB polls the provided DynamoDB client until it becomes available or a timeout is reached.
func WaitForDynamoDB(ctx context.Context, client *dynamodb.Client, timeout ...time.Duration) error {
	var to <-chan time.Time

	if len(timeout) == 0 {
		to = time.After(10 * time.Second)
	} else {
		to = time.After(timeout[0])
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-to:
			return fmt.Errorf("timeout waiting for DynamoDB endpoint to become available")
		case <-ticker.C:
			// A lightweight API call to test if DynamoDB is ready.
			_, err := client.ListTables(ctx, &dynamodb.ListTablesInput{
				Limit: aws.Int32(1),
			})
			if err == nil {
				return nil
			}
		}
	}
}
