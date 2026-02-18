#!/usr/bin/env sh
# Build CHANGELOG.md from changelog/unreleased.md and changelog/released/*.md
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
CHANGELOG_DIR="$ROOT_DIR/changelog"
OUTPUT="$ROOT_DIR/CHANGELOG.md"

LATEST=$(ls -1 "$CHANGELOG_DIR/released/"*.md 2>/dev/null | xargs -I{} basename {} .md | sort -V -r | head -1)
LATEST=${LATEST:-0.1.0}

{
  echo "# Changelog"
  echo ""
  echo "All notable changes to this project will be documented in this file."
  echo ""
  echo "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),"
  echo "and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)."
  echo ""
  echo "Entries are managed in \`changelog/\` — see [changelog/README.md](changelog/README.md) for how to add and release."
  echo ""
  cat "$CHANGELOG_DIR/unreleased.md"
  echo ""
  for f in $(ls -1 "$CHANGELOG_DIR/released/"*.md 2>/dev/null | sort -V -r); do
    cat "$f"
    echo ""
  done
  echo "[Unreleased]: https://github.com/pgquerynarrative/pgquerynarrative/compare/v${LATEST}...HEAD"
  for f in $(ls -1 "$CHANGELOG_DIR/released/"*.md 2>/dev/null | sort -V); do
    ver=$(basename "$f" .md)
    echo "[$ver]: https://github.com/pgquerynarrative/pgquerynarrative/releases/tag/v$ver"
  done
} > "$OUTPUT"

echo "Generated $OUTPUT"
