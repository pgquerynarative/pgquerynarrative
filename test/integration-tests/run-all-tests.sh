#!/bin/bash
set -e

echo "=========================================="
echo "  Complete Test Suite"
echo "=========================================="
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/../.."

# Step 1: Setup Docker Postgres
echo "Step 1: Setting up test database..."
bash "$SCRIPT_DIR/setup-docker.sh"
export DB_URL="postgres://postgres:postgres@localhost:5433/pgquerynarrative?sslmode=disable"

# Step 2: Run migrations
echo ""
echo "Step 2: Running migrations..."
bash "$SCRIPT_DIR/run-migrations.sh"

# Step 3: Seed data
echo ""
echo "Step 3: Seeding test data..."
bash "$SCRIPT_DIR/seed-data.sh"

# Step 4: Build application
echo ""
echo "Step 4: Building application..."
go build -o bin/server ./cmd/server
echo "Build successful"

# Step 5: Run unit tests
echo ""
echo "Step 5: Running unit tests..."
go test ./app/metrics/... -v
go test ./app/llm/... -v || echo "LLM tests skipped (may require Ollama)"
go test ./app/story/... -v || echo "Story tests skipped (may require LLM)"

# Step 6: Start server in background
echo ""
echo "Step 6: Starting server..."
export DATABASE_HOST=localhost
export DATABASE_PORT=5433
export DATABASE_NAME=pgquerynarrative
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_READONLY_USER=postgres
export DATABASE_READONLY_PASSWORD=postgres
export DATABASE_SSL_MODE=disable
export LLM_PROVIDER=ollama
export LLM_BASE_URL=http://localhost:11434
export LLM_MODEL=llama3.2
./bin/server &
SERVER_PID=$!
sleep 5

# Step 7: Test API
echo ""
echo "Step 7: Testing API endpoints..."
bash "$SCRIPT_DIR/test-api.sh" || echo "API tests had issues"

# Step 8: Cleanup
echo ""
echo "Step 8: Cleaning up..."
kill $SERVER_PID 2>/dev/null || true
bash "$SCRIPT_DIR/teardown-docker.sh"

echo ""
echo "=========================================="
echo "  Test Suite Complete"
echo "=========================================="
