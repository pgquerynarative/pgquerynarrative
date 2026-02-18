#!/bin/sh
set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

make docker-up
make db-init

export DB_URL="postgres://postgres:postgres@localhost:5432/pgquerynarrative?sslmode=disable"
make migrate
make seed
make generate

PORT="${PGQUERYNARRATIVE_PORT:-8080}"
if lsof -i ":${PORT}" >/dev/null 2>&1; then
  PORT=8081
  echo "Port 8080 is in use; starting server on ${PORT}."
fi

PGQUERYNARRATIVE_PORT="${PORT}" make run
