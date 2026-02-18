#!/bin/bash
# Safe command runner with 30-second timeout
# Usage: ./safe-run.sh "command to run"

set -e

TIMEOUT="${COMMAND_TIMEOUT:-30}"
COMMAND="$1"

if [ -z "$COMMAND" ]; then
    echo "Usage: $0 \"command to run\""
    exit 1
fi

echo "Running: $COMMAND"
echo "Timeout: ${TIMEOUT}s"
echo ""

# Run command with timeout
if timeout ${TIMEOUT}s bash -c "$COMMAND"; then
    echo ""
    echo "✅ Command completed successfully"
    exit 0
else
    EXIT_CODE=$?
    echo ""
    if [ $EXIT_CODE -eq 124 ]; then
        echo "❌ Command timed out after ${TIMEOUT} seconds"
    else
        echo "❌ Command failed with exit code: $EXIT_CODE"
    fi
    exit $EXIT_CODE
fi
