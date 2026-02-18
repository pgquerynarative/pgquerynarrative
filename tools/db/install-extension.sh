#!/bin/sh
# Script to install PgQueryNarrative PostgreSQL extension

set -e

DB_NAME="${PGQUERYNARRATIVE_DB:-pgquerynarrative}"
DB_USER="${PGQUERYNARRATIVE_USER:-pgquerynarrative_app}"
DB_HOST="${PGHOST:-localhost}"
DB_PORT="${PGPORT:-5432}"

EXTENSION_DIR="infra/postgres-extension"

echo "Installing PgQueryNarrative PostgreSQL extension..."
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo "Host: $DB_HOST:$DB_PORT"
echo ""

# Check if http extension is available
echo "Checking for http extension..."
if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 1 FROM pg_extension WHERE extname = 'http'" | grep -q 1; then
    echo "✓ http extension found"
    HAS_HTTP=true
else
    echo "⚠ http extension not found"
    echo "  Installing http extension..."
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "CREATE EXTENSION IF NOT EXISTS http;" 2>/dev/null; then
        echo "✓ http extension installed"
        HAS_HTTP=true
    else
        echo "⚠ Could not install http extension (may require superuser)"
        echo "  Will use basic version without HTTP support"
        HAS_HTTP=false
    fi
fi

# Install the extension
echo ""
echo "Installing PgQueryNarrative extension..."
if [ "$HAS_HTTP" = "true" ]; then
    # Install basic extension first
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$EXTENSION_DIR/pgquerynarrative--1.0.sql"
    
    # Then apply the HTTP-enabled version
    echo "Applying HTTP-enabled functions..."
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$EXTENSION_DIR/pgquerynarrative--1.0--with-http.sql"
else
    # Install basic version only
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$EXTENSION_DIR/pgquerynarrative--1.0.sql"
fi

echo ""
echo "✓ PgQueryNarrative extension installed successfully!"
echo ""
echo "Usage examples:"
echo "  SELECT pgquerynarrative_get_api_url();"
echo "  SELECT pgquerynarrative_set_api_url('http://localhost:8080');"
echo "  SELECT pgquerynarrative_run_query('SELECT COUNT(*) FROM demo.sales', 10);"
echo "  SELECT pgquerynarrative_list_saved(50, 0);"
