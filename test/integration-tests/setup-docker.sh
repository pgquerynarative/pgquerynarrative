#!/bin/bash
set -e

echo "=== Test Setup - Docker Postgres ==="

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "ERROR: Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Start Postgres container
echo "Starting PostgreSQL container..."
docker run -d \
    --name pgquerynarrative-test \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=pgquerynarrative \
    -p 5433:5432 \
    postgres:18

# Wait for Postgres to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec pgquerynarrative-test pg_isready -U postgres >/dev/null 2>&1; then
        echo "PostgreSQL is ready!"
        break
    fi
    sleep 1
done

# Set database URL
export DATABASE_HOST=localhost
export DATABASE_PORT=5433
export DATABASE_NAME=pgquerynarrative
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_READONLY_USER=postgres
export DATABASE_READONLY_PASSWORD=postgres
export DATABASE_SSL_MODE=disable

DB_URL="postgres://postgres:postgres@localhost:5433/pgquerynarrative?sslmode=disable"

echo "Database URL: $DB_URL"
echo ""
echo "Test database is ready!"
echo "Run: export DB_URL=\"$DB_URL\""
