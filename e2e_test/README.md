# End-to-End Testing

This directory contains end-to-end tests for the video-in-be server.

## Running Tests

### Docker Testing
To run the full Docker-based test (may fail in some network environments):

```bash
cd e2e_test
go test -v -run TestEndToEndServer
```

### All Tests
To run all tests:

```bash
cd e2e_test
go test -v
```

## What the Tests Do

1. **TestEndToEndServer**: 
   - Builds a Docker container from the repository
   - Starts the container in server mode
   - Makes a HelloWorld RPC call to the containerized server
   - Verifies the response

2. **TestDockerBuild**:
   - Tests that the Dockerfile can build successfully

## Test Configuration

The tests use a special `test-key` value for the TMDB API key, which the server detects and uses to skip actual API calls while still allowing the server to start normally.