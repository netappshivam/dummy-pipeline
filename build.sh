#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Brings up binary
go mod tidy
GOOS=linux GOARCH=amd64 go build -o bin/build/linux/vsactl . # Build binary
echo "Binary built successfully at bin/build/linux/vsactl."

# Build Docker image using environment variables
docker buildx build --platform linux/amd64 \
  --build-arg RUNNER=github \
  --build-arg GO_VERSION=1.24.2 \
  --build-arg GO_FILENAME=go1.24.2.linux-amd64.tar.gz \
  --build-arg GO_FILENAME_SHA=68097bd680839cbc9d464a0edce4f7c333975e27a90246890e9f1078c7e702ad \
  -t vsactl:v1 .
echo "Docker image built successfully with tag vsactl:v1."

# Tag and push the Docker image
docker tag vsactl:v1 ghcr.io/vcp-vsa-control-plane/vsactl:v1
docker push ghcr.io/vcp-vsa-control-plane/vsactl:v1
echo "Docker image tagged and pushed successfully to ghcr.io/vcp-vsa-control-plane/vsactl:v1."

# Exit successfully
exit 0