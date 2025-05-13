# Use the official Golang image to build the Go binary
FROM golang:1.23.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Use a minimal base image to run the binary
FROM debian:bookworm-slim

# Install ffmpeg
RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 25004

# Set these environment variables before running the container.
ENV VIDEOIN_PROJECTDIR= \
    VIDEOIN_STATEDIR= \
    VIDEOIN_UNCLAIMEDDIR=

# Command to run the binary
CMD ["./main"]