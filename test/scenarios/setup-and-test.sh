#!/bin/bash
set -e

# Setup and run test scenario
# This script ensures the environment is ready before testing

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

echo "=========================================="
echo "  PgQueryNarrative Test Setup & Run"
echo "=========================================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/ > /dev/null 2>&1; then
    echo "Server not running. Start with 'make start-docker' or 'make start-local'."
    echo ""
    echo "Please run in a separate terminal:"
    echo "  make start-docker   # or make start-local"
    echo ""
    echo "Or manually:"
    echo "  go run ./cmd/server"
    echo ""
    read -p "Press Enter once the server is running, or Ctrl+C to cancel..."
fi

# Run the test scenario
echo ""
echo "Running test scenario..."
echo ""
bash "$SCRIPT_DIR/test-scenario.sh"
