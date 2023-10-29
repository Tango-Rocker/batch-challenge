#!/bin/bash
set -e
cd "$(dirname "$0")"
# Build the binary with all dependencies included
CGO_ENABLED=0 GOOS=linux go build -v -o batch-challenge-linux-amd64
echo "Build completed."
