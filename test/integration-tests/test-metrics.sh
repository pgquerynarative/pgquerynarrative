#!/bin/bash
set -e

echo "=== Testing Metrics Calculator ==="
echo ""

cd "$(dirname "$0")/../.."

# Run metrics tests
echo "Running metrics unit tests..."
go test ./app/metrics/... -v

echo ""
echo "Metrics tests completed"
