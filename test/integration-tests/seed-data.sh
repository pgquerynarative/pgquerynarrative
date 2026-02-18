#!/bin/bash
set -e

if [ -z "$DB_URL" ]; then
    echo "ERROR: DB_URL environment variable not set"
    exit 1
fi

echo "=== Seeding Test Data ==="
echo "Database: $DB_URL"
echo ""

cd "$(dirname "$0")/../.."

# Run seed script
psql "$DB_URL" -f ./tools/db/seed.sql

echo ""
echo "Test data seeded"
