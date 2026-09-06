#!/bin/sh
set -e

DB_HOST="${DATABASE_HOST:-postgres}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-pgquerynarrative}"
DB_USER="${DATABASE_USER:-pgquerynarrative_app}"
DB_PASSWORD="${DATABASE_PASSWORD:-pgquerynarrative_app}"

export PGPASSWORD="${DB_PASSWORD}"

attempts=30
while [ $attempts -gt 0 ]; do
  if pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" >/dev/null 2>&1; then
    break
  fi
  attempts=$((attempts - 1))
  sleep 1
done

if [ $attempts -eq 0 ]; then
  echo "Postgres is not ready at ${DB_HOST}:${DB_PORT}"
  exit 1
fi

export DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DATABASE_SSL_MODE:-disable}"

# Migrations are DDL: they create extensions (pg_stat_statements, pgvector,
# hypopg) and ALTER ROLE, none of which the runtime application role is allowed
# to do — by design, since that role also executes user SQL. Running them as
# DATABASE_USER made a fresh deploy of this image fail at migration 000019 with
# "permission denied to create extension".
#
# Use DATABASE_MIGRATION_USER / DATABASE_MIGRATION_PASSWORD (or a full
# DATABASE_MIGRATION_URL) for a role that may create extensions and alter roles.
# Falling back to DATABASE_USER keeps existing deployments working, where the
# schema is already at the required version and `up` is a no-op.
MIGRATE_USER="${DATABASE_MIGRATION_USER:-$DB_USER}"
MIGRATE_PASSWORD="${DATABASE_MIGRATION_PASSWORD:-$DB_PASSWORD}"
MIGRATE_URL="${DATABASE_MIGRATION_URL:-postgres://${MIGRATE_USER}:${MIGRATE_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DATABASE_SSL_MODE:-disable}}"

if [ "${MIGRATE_USER}" = "${DB_USER}" ] && [ -z "${DATABASE_MIGRATION_URL:-}" ]; then
  echo "Running migrations as ${DB_USER}. If this is a fresh database, set" >&2
  echo "DATABASE_MIGRATION_USER to a role that may create extensions and alter roles." >&2
fi

/app/bin/migrate -path /app/app/db/migrations -database "${MIGRATE_URL}" up

if [ "${PGQUERYNARRATIVE_SEED:-false}" = "true" ]; then
  psql "${DB_URL}" -f /app/tools/db/seed.sql
fi

exec /app/bin/server
