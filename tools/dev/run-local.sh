#!/bin/bash
set -e

echo "==> Starting PgQueryNarrative (Local Postgres)"

# Check if Postgres is running
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "ERROR: PostgreSQL is not running on localhost:5432"
    echo "Please start PostgreSQL first:"
    echo "  - macOS: brew services start postgresql@18"
    echo "  - Linux: sudo systemctl start postgresql"
    exit 1
fi

# Database details
DB_NAME="pgquerynarrative"
DB_USER="$(whoami)"
DB_URL="postgres://${DB_USER}@localhost:5432/${DB_NAME}?sslmode=disable"

echo "==> Creating database if not exists..."
createdb ${DB_NAME} 2>/dev/null || echo "Database already exists"

echo "==> Running init SQL..."
psql -d ${DB_NAME} -f infra/postgres-init/00-init.sql

echo "==> Running migrations..."
DATABASE_URL="${DB_URL}" sh ./tools/db/migrate.sh up "${DB_URL}"

echo "==> Seeding data (optional)..."
if [ "${PGQUERYNARRATIVE_SEED}" = "true" ]; then
    psql "${DB_URL}" -f ./tools/db/seed.sql
    echo "Seed data loaded"
fi

echo "==> Generating Goa code..."
go install goa.design/goa/v3/cmd/goa@latest
goa gen github.com/pgquerynarrative/pgquerynarrative/api/design

echo "==> Starting server..."
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_NAME=${DB_NAME}
export DATABASE_USER=${DB_USER}
export DATABASE_PASSWORD=""
export DATABASE_READONLY_USER=${DB_USER}
export DATABASE_READONLY_PASSWORD=""
export PGQUERYNARRATIVE_HOST=${PGQUERYNARRATIVE_HOST:-0.0.0.0}
export PGQUERYNARRATIVE_PORT=${PGQUERYNARRATIVE_PORT:-8080}

go run ./cmd/server
