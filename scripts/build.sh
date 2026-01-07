#!/bin/bash
set -e

echo "Building Lambda function for Connectra API..."

# Build for Linux/amd64 (Lambda runtime)
GOOS=linux GOARCH=amd64 go build \
  -o bootstrap \
  -ldflags="-s -w" \
  ./cmd/lambda/main.go

echo "Build complete! Binary: bootstrap"
echo "Binary size: $(du -h bootstrap | cut -f1)"
