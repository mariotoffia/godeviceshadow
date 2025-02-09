package dynamodbutils

import (
	"bytes"
	"context"
	"fmt"
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
	ContainerID  string
	Endpoint     string
	Client       *dynamodb.Client
	StreamClient *dynamodbstreams.Client
}

// StartLocalDynamoDB will start a local dynamodb table and make sure it has
// a table to persist in.
func StartLocalDynamoDB(ctx context.Context, table string) (*LocalDynamoDB, error) {
	containerID, endpoint, err := StartDynamoDBLocal()
	if err != nil {
		return nil, fmt.Errorf("Failed to start DynamoDB Local: %w", err)
	}

	cfg, err := DefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS config: %w", err)
	}

	client := DynamoDbClient(cfg, endpoint)
	streamClient := DynamoDbStreamClient(cfg, endpoint)

	time.Sleep(2 * time.Second)

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
func StartDynamoDBLocal() (string, string, error) {
	// Run Docker command to start container
	cmd := exec.Command("docker", "run", "-d", "-p", "8000:8000", "amazon/dynamodb-local")

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", "", fmt.Errorf("failed to start DynamoDB Local: %w", err)
	}

	// Get the container ID
	containerID := strings.TrimSpace(out.String())

	// Allow some time for DynamoDB to start
	time.Sleep(2 * time.Second)

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
