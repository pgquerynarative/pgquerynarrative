#!/bin/bash
set -e

if [ -z "$DB_URL" ]; then
    echo "ERROR: DB_URL environment variable not set"
    echo "Run: export DB_URL=\"postgres://user:pass@host:port/db?sslmode=disable\""
    exit 1
fi

echo "=== Running Database Migrations ==="
echo "Database: $DB_URL"
echo ""

cd "$(dirname "$0")/../.."

# Run migrations
sh ./tools/db/migrate.sh up "$DB_URL"

echo ""
echo "Migrations completed"
