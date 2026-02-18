#!/bin/bash
set -e

echo "=== Teardown Test Database ==="

# Stop and remove container
if docker ps -a | grep -q pgquerynarrative-test; then
    echo "Stopping and removing test container..."
    docker stop pgquerynarrative-test >/dev/null 2>&1 || true
    docker rm pgquerynarrative-test >/dev/null 2>&1 || true
    echo "Test container removed"
else
    echo "No test container found"
fi
