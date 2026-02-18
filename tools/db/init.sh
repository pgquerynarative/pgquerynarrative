#!/bin/sh
set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

docker compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1 || {
  echo "Postgres is not ready. Run: make docker-up"
  exit 1
}

docker compose exec -T postgres psql -U postgres -d postgres -tc \
  "SELECT 1 FROM pg_database WHERE datname='pgquerynarrative'" | grep -q 1 || \
  docker compose exec -T postgres createdb -U postgres pgquerynarrative

if [ -f "$ROOT_DIR/tools/db/create-roles.sql" ]; then
  docker compose exec -T postgres psql -U postgres -d pgquerynarrative < "$ROOT_DIR/tools/db/create-roles.sql"
else
  echo "Note: create-roles.sql not found, roles should be created by init script"
fi
