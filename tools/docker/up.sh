#!/bin/sh
set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

if docker info >/dev/null 2>&1; then
  docker compose up -d
  exit 0
fi

if [ -z "${DOCKER_API_VERSION:-}" ]; then
  export DOCKER_API_VERSION=1.41
fi

if docker info >/dev/null 2>&1; then
  docker compose up -d
  exit 0
fi

echo "Docker is not reachable. Start Docker Desktop and retry."
echo "If needed: export DOCKER_API_VERSION=1.41"
exit 1
