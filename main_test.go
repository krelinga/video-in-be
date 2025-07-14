// +build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

func TestEndToEndServer(t *testing.T) {
	// Set required environment variables before the test starts
	// This is needed because the tmdb package has an init() function that requires them
	tmpProjectDir := t.TempDir()
	tmpStateDir := t.TempDir()
	tmpUnclaimedDir := t.TempDir()
	tmpThumbsDir := t.TempDir()
	tmpLibraryDir := t.TempDir()

	os.Setenv("VIDEOIN_PROJECTDIR", tmpProjectDir)
	os.Setenv("VIDEOIN_STATEDIR", tmpStateDir)
	os.Setenv("VIDEOIN_UNCLAIMEDDIR", tmpUnclaimedDir)
	os.Setenv("VIDEOIN_THUMBSDIR", tmpThumbsDir)
	os.Setenv("VIDEOIN_TMDBKEY", "test-key")
	os.Setenv("VIDEOIN_LIBRARYDIR", tmpLibraryDir)

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("VIDEOIN_PROJECTDIR")
		os.Unsetenv("VIDEOIN_STATEDIR")
		os.Unsetenv("VIDEOIN_UNCLAIMEDDIR")
		os.Unsetenv("VIDEOIN_THUMBSDIR")
		os.Unsetenv("VIDEOIN_TMDBKEY")
		os.Unsetenv("VIDEOIN_LIBRARYDIR")
	}()
	ctx := context.Background()

	// Build the Docker image from the current directory
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    ".", // Current directory (repository root)
			Dockerfile: "Dockerfile",
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
		Cmd: []string{"-mode", "server"},
		WaitingFor: wait.ForHTTP("/").WithPort("25004/tcp").WithStartupTimeout(30 * time.Second),
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
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

	// Create Connect RPC client
	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
	client := pbconnect.NewServiceClient(http.DefaultClient, baseURL)

	// Test HelloWorld RPC call
	request := &pb.HelloWorldRequest{
		Name: "TestUser",
	}

	response, err := client.HelloWorld(ctx, connect.NewRequest(request))
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Hello, TestUser", response.Msg.Message)
}