# Development setup

## Prerequisites

- Go 1.24+, PostgreSQL 16+ (or Docker), Git, Make

## Setup

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make setup
make generate
```

**Database:** Docker: `docker compose up -d postgres` then `make db-init && make migrate && make seed`. Local: `make db-init && make migrate && make seed` (Postgres running).

**Test:** `make test`

## Run locally

```bash
make run
# or
go run ./cmd/server
```

App: `http://localhost:8080`. Verbose logging: `LOG_DEBUG=1 make run`.

## Workflow

1. Branch: `git checkout -b feature/name`
2. Code, test (`make test`), lint (`make lint`), format (`make fmt`)
3. Commit: Conventional Commits (e.g. `feat: add X`)
4. After changing `api/design/*.go`: `make generate`

**Migrations:** Add `00000N_name.up.sql` and `00000N_name.down.sql` in `app/db/migrations/`; test with `make migrate`.

## Commands

| Command | Purpose |
|---------|---------|
| `make fmt` | Format code |
| `make lint` | Lint |
| `make test-unit` | Unit tests |
| `make test-integration` | Integration tests (Docker) |
| `make test-e2e` | E2E tests |
| `make build` | Build `bin/server` |

**See also:** [Contributing](../../.github/CONTRIBUTING.md), [Testing](testing.md), [Troubleshooting](../reference/troubleshooting.md)
