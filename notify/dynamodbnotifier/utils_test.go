//go:build integration
// +build integration

package dynamodbnotifier

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startDynamoDBContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForListeningPort("8000").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	// Get the mapped port
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := container.MappedPort(ctx, "8000")
	if err != nil {
		return nil, "", err
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())
	return container, endpoint, nil
}

func TestDynamoDBLocal(t *testing.T) {
	ctx := context.Background()

	// Start DynamoDB Local
	container, endpoint, err := startDynamoDBContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to start DynamoDB Local: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("Failed to stop DynamoDB Local: %v", err)
		}
	}()

	// Configure AWS SDK v2 client for local DynamoDB
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"), // Region is required, but irrelevant for local DynamoDB
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")), // Use dummy credentials
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{URL: endpoint, SigningRegion: "us-west-2"}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create a DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Test by listing tables
	out, err := svc.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}

	t.Logf("Tables: %v", out.TableNames)
}
