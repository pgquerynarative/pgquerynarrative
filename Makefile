.PHONY: setup tidy generate build build-mcp run test test-unit test-integration test-e2e lint fmt migrate seed dev dev-stop dev-watch dev-build dev-teardown docker-up docker-down docker-logs db-init start start-docker start-local stop cli cli-shell changelog

GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOA ?= goa

DB_URL ?= postgres://pgquerynarrative_app:pgquerynarrative_app@localhost:5432/pgquerynarrative?sslmode=disable

# ============================================================================
# SIMPLIFIED COMMANDS - Start here!
# ============================================================================

# Start: choose Docker or local PostgreSQL explicitly
start:
	@echo "Choose how to run PgQueryNarrative:"
	@echo ""
	@echo "  make start-docker   Use Docker (PostgreSQL + app in containers)"
	@echo "  make start-local    Use local PostgreSQL (app runs on host)"
	@echo ""
	@echo "Then open http://localhost:8080"
	@exit 1

# Start with Docker (PostgreSQL + app in containers)
start-docker:
	@$(MAKE) docker-start

# Start with local PostgreSQL (app runs on host; requires Postgres already running)
start-local:
	@$(MAKE) local-start

# Stop everything
stop:
	@echo "🛑 Stopping PgQueryNarrative..."
	@if docker ps | grep -q pgquerynarrative; then \
		docker compose down; \
	else \
		echo "No Docker containers running"; \
	fi
	@echo "✅ Stopped"

# ============================================================================
# Docker-based startup
# ============================================================================

docker-start:
	@echo "📦 Setting up with Docker..."
	@echo ""
	@echo "Step 1: Starting PostgreSQL..."
	@docker compose up -d postgres || (echo "❌ Failed to start PostgreSQL. Is Docker running?" && exit 1)
	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then \
			echo "✅ PostgreSQL is ready!"; \
			break; \
		fi; \
		sleep 1; \
		timeout=$$((timeout - 1)); \
	done; \
	if [ $$timeout -eq 0 ]; then \
		echo "❌ PostgreSQL failed to start"; \
		exit 1; \
	fi
	@echo ""
	@echo "Step 2: Setting up database..."
	@$(MAKE) db-init || true
	@echo ""
	@echo "Step 3: Running migrations..."
	@export DB_URL="postgres://postgres:postgres@localhost:5432/pgquerynarrative?sslmode=disable"; \
	$(MAKE) migrate || true
	@echo ""
	@echo "Step 4: Seeding demo data..."
	@export DB_URL="postgres://postgres:postgres@localhost:5432/pgquerynarrative?sslmode=disable"; \
	$(MAKE) seed || true
	@echo ""
	@echo "Step 5: Building application..."
	@$(MAKE) generate build
	@echo ""
	@echo "Step 6: Starting application..."
	@echo "✅ PgQueryNarrative is starting!"
	@echo ""
	@echo "🌐 API will be available at: http://localhost:8080"
	@echo "📊 Try: curl http://localhost:8080/api/v1/queries/saved"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo ""
	@docker compose up app

# ============================================================================
# Local PostgreSQL startup (no Docker)
# ============================================================================

# URL for local PostgreSQL (app user). Set DATABASE_URL or DB_URL to override.
LOCAL_DB_URL ?= postgres://pgquerynarrative_app:pgquerynarrative_app@localhost:5432/pgquerynarrative?sslmode=disable

local-start:
	@echo "💻 Setting up with local PostgreSQL..."
	@echo ""
	@echo "Step 1: Checking PostgreSQL connection..."
	@pg_isready -h localhost -p 5432 >/dev/null 2>&1 || (echo "❌ PostgreSQL not running. Start it with: brew services start postgresql@18" && exit 1)
	@echo "✅ PostgreSQL is ready!"
	@echo ""
	@echo "Step 2: Setting up database..."
	@$(MAKE) db-init || true
	@echo ""
	@echo "Step 3: Running migrations..."
	@DB_URL="$${DB_URL:-$${DATABASE_URL:-$(LOCAL_DB_URL)}}"; $(MAKE) migrate
	@echo ""
	@echo "Step 4: Seeding demo data..."
	@DB_URL="$${DB_URL:-$${DATABASE_URL:-$(LOCAL_DB_URL)}}"; $(MAKE) seed
	@echo ""
	@echo "Step 5: Building application..."
	@$(MAKE) generate build
	@echo ""
	@echo "Step 6: Starting application..."
	@echo "✅ PgQueryNarrative is starting!"
	@echo ""
	@echo "🌐 API will be available at: http://localhost:8080"
	@echo "📊 Try: curl http://localhost:8080/api/v1/queries/saved"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo ""
	@$(MAKE) run

# ============================================================================
# Development commands
# ============================================================================

setup:
	@echo "📥 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependencies installed"

tidy:
	$(GO) mod tidy

generate:
	@echo "🔧 Generating API code..."
	@if ! command -v goa >/dev/null 2>&1; then \
		echo "Installing Goa..."; \
		$(GO) install goa.design/goa/v3/cmd/goa@latest; \
	fi
	$(GO) generate ./...
	$(GOA) gen github.com/pgquerynarrative/pgquerynarrative/api/design
	@sh ./tools/fix-gen-metrics-validator.sh 2>/dev/null || true
	@sh ./tools/copy-gen-to-api-gen.sh 2>/dev/null || true
	@echo "✅ Code generated"

build:
	@echo "🔨 Building application..."
	$(GO) build -o bin/server ./cmd/server
	@echo "✅ Build complete: bin/server"

# Build MCP server for Claude / Cursor (stdio transport)
build-mcp:
	$(GO) build -o bin/mcp-server ./cmd/mcp-server
	@echo "✅ MCP server: bin/mcp-server"

run:
	$(GO) run ./cmd/server

# ============================================================================
# Testing
# ============================================================================

test: test-unit test-integration

test-unit:
	@echo "🧪 Running unit tests..."
	$(GO) test ./app/... ./cmd/... -v

test-integration:
	@echo "🧪 Running integration tests..."
	$(GO) test ./test/integration/... -v

test-e2e:
	@echo "🧪 Running E2E tests..."
	$(GO) test ./test/e2e/... -v

# ============================================================================
# Code quality
# ============================================================================

lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt:
	@echo "🎨 Formatting code..."
	$(GO) fmt ./...
	@echo "✅ Code formatted"

# ============================================================================
# Database operations
# ============================================================================

migrate:
	@DB_URL="$${DB_URL:-$${DATABASE_URL:-$(LOCAL_DB_URL)}}"; \
	if [ -z "$$DB_URL" ]; then \
		echo "❌ DB_URL or DATABASE_URL not set. Using default..."; \
		DB_URL="$(LOCAL_DB_URL)"; \
	fi; \
	sh ./tools/db/migrate.sh up "$$DB_URL"

seed:
	@DB_URL="$${DB_URL:-$${DATABASE_URL:-$(LOCAL_DB_URL)}}"; \
	if [ -z "$$DB_URL" ]; then \
		echo "❌ DB_URL or DATABASE_URL not set. Using default..."; \
		DB_URL="$(LOCAL_DB_URL)"; \
	fi; \
	psql "$$DB_URL" -f ./tools/db/seed.sql || echo "⚠️  Seed data already exists or database not accessible"

db-init:
	@echo "🗄️  Initializing database..."
	@if docker ps | grep -q pgquerynarrative-postgres; then \
		sh ./tools/db/init.sh; \
	else \
		$(MAKE) local-db-init || true; \
	fi

# Local PostgreSQL: create database and roles (no Docker). Uses default connection
# (e.g. current user on macOS Homebrew). Requires superuser.
local-db-init:
	@echo "🗄️  Creating database and roles (local PostgreSQL)..."
	@createdb pgquerynarrative 2>/dev/null || true
	@psql -d pgquerynarrative -f infra/postgres-init/00-init.sql && echo "✅ Database and roles ready" || (echo "⚠️  Run as PostgreSQL superuser (e.g. your macOS user). If roles exist, run: make migrate seed"; exit 0)

# ============================================================================
# Docker Compose commands
# ============================================================================

dev:
	sh ./tools/dev/dev.sh

dev-stop:
	docker compose down

docker-up:
	sh ./tools/docker/up.sh

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

dev-watch:
	@echo "👀 Starting in watch mode (auto-reload on changes)..."
	@echo "This will build and start services with file watching."
	@echo ""
	docker compose up --build --watch

dev-build:
	@echo "🔨 Rebuilding app container..."
	$(GO) build -o bin/server ./cmd/server
	docker compose up --build -d app

dev-teardown:
	@echo "🧹 Tearing down development environment..."
	docker compose down -v
	rm -rf infra/data
	@echo "✅ Development environment reset complete"

# ============================================================================
# CLI Commands (Docker-only, no browser needed)
# ============================================================================

cli:
	@echo "💻 Running CLI command..."
	@sh ./tools/docker/docker-cli.sh $(CMD)

cli-shell:
	@echo "💻 Starting interactive CLI shell..."
	@echo "Type 'pgquerynarrative help' for commands"
	@echo "Or use 'pqn' as alias"
	@echo ""
	@docker compose run --rm -it --entrypoint /bin/sh cli -l

# ============================================================================
# Changelog
# ============================================================================

changelog:
	@echo "📝 Building CHANGELOG.md from changelog/..."
	@sh ./tools/changelog/build.sh

# ============================================================================
# PostgreSQL Extension
# ============================================================================

install-extension:
	@echo "📦 Installing PgQueryNarrative PostgreSQL extension..."
	@sh ./tools/db/install-extension.sh

install-extension-docker:
	@echo "📦 Installing PgQueryNarrative extension in Docker PostgreSQL..."
	@docker compose exec postgres psql -U postgres -d pgquerynarrative -f /docker-entrypoint-initdb.d/extension/pgquerynarrative--1.0.sql || \
		docker compose exec postgres psql -U pgquerynarrative_app -d pgquerynarrative -f /docker-entrypoint-initdb.d/extension/pgquerynarrative--1.0.sql || \
		echo "⚠ Extension installation requires manual setup. See infra/postgres-extension/README.md"
