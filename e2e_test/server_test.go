package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	pbconnect "buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDockerBuild tests that the Dockerfile can build successfully
func TestDockerBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker build test in short mode")
	}

	ctx := context.Background()

	// Test building Docker image from the parent directory
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "..", // Parent directory (repository root)
			},
		},
		Started: false, // Only build, don't start
	})

	assert.NoError(t, err, "Docker build failed")
}

func TestEndToEndServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker-based test in short mode")
	}

	ctx := context.Background()

	// Use the vendor-based Dockerfile to avoid network issues in tests
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "..", // Parent directory (repository root)
		},
		ExposedPorts: []string{"25004/tcp"},
		Env: map[string]string{
			"VIDEOIN_PROJECTDIR":   "/tmp/project",
			"VIDEOIN_STATEDIR":     "/tmp/state",
			"VIDEOIN_UNCLAIMEDDIR": "/tmp/unclaimed",
			"VIDEOIN_THUMBSDIR":    "/tmp/thumbs",
			"VIDEOIN_TMDBKEY":      "test-key",
			"VIDEOIN_LIBRARYDIR":   "/tmp/library",
		},
		Cmd:        []string{"-mode", "server"},
		WaitingFor: wait.ForLog("Hello, World!").WithStartupTimeout(30 * time.Second),
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if !assert.NoError(t, err, "Docker build failed") {
		return
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	// Get the container's host and port
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "25004")
	require.NoError(t, err)

	// Wait a bit for server to be ready to accept connections
	time.Sleep(2 * time.Second)

	// Create Connect RPC client
	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
	client := pbconnect.NewServiceClient(http.DefaultClient, baseURL)

	// Test HelloWorld RPC call
	request := &pb.HelloWorldRequest{
		Name: "DockerTestUser",
	}

	response, err := client.HelloWorld(ctx, connect.NewRequest(request))
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Hello, DockerTestUser", response.Msg.Message)
}
