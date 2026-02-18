#!/bin/sh
set -e

CMD="${1:-up}"
DB_URL="${2:-}"

if [ -z "$DB_URL" ]; then
  echo "Usage: ./tools/db/migrate.sh up|down|version <database_url>"
  exit 1
fi

MIGRATE_PKG="github.com/golang-migrate/migrate/v4/cmd/migrate@latest"

case "$CMD" in
  up)
    go run -tags 'postgres' "$MIGRATE_PKG" -path ./app/db/migrations -database "$DB_URL" up
    ;;
  down)
    go run -tags 'postgres' "$MIGRATE_PKG" -path ./app/db/migrations -database "$DB_URL" down
    ;;
  version)
    go run -tags 'postgres' "$MIGRATE_PKG" -path ./app/db/migrations -database "$DB_URL" version
    ;;
  *)
    echo "Unknown command: $CMD"
    echo "Usage: ./tools/db/migrate.sh up|down|version <database_url>"
    exit 1
    ;;
esac
