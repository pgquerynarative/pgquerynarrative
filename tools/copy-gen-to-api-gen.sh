#!/bin/sh
# Copy generated code from gen/ to api/gen/ and fix import paths.
# Goa generates into gen/; the app imports api/gen. This script keeps api/gen in sync.
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

[ -d gen ] || exit 0

cp gen/queries/*.go api/gen/queries/
cp gen/reports/*.go api/gen/reports/
cp -r gen/http/* api/gen/http/

# Fix imports: gen/ -> api/gen/ in api/gen/http so server and app use the same types
for f in api/gen/http/queries/server/*.go api/gen/http/queries/client/*.go \
         api/gen/http/reports/server/*.go api/gen/http/reports/client/*.go; do
	[ -f "$f" ] || continue
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/queries|github.com/pgquerynarrative/pgquerynarrative/api/gen/queries|g' "$f"
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/reports|github.com/pgquerynarrative/pgquerynarrative/api/gen/reports|g' "$f"
	rm -f "$f.bak"
done
for f in api/gen/http/cli/pgquerynarrative/cli.go; do
	[ -f "$f" ] || continue
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/http/|github.com/pgquerynarrative/pgquerynarrative/api/gen/http/|g' "$f"
	rm -f "$f.bak"
done

# fix-gen-metrics-validator.sh runs on gen/ before this copy, so api/gen gets the patched file
echo "Synced gen/ to api/gen/ (imports fixed)"
