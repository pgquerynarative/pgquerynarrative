# PgQueryNarrative

Turn SQL query results into business narratives with AI.

- Run read-only SQL, get metrics and insights
- Generate narrative reports via LLM (Ollama, Gemini, Claude)
- Save queries and share reports

## Quick Start

**Docker** (PostgreSQL + app in containers):

```bash
make start-docker
```

**Local PostgreSQL** (app on host; Postgres must be running):

```bash
make start-local
```

Then open **http://localhost:8080** or call the API:

```bash
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category", "limit": 10}'
```

## Requirements

- **Docker** (for `make start-docker`) or **PostgreSQL 16+** (for `make start-local`)
- **Go 1.24+** (for building)

## Commands

| Action | Command |
|--------|---------|
| Start (Docker) | `make start-docker` |
| Start (local) | `make start-local` |
| Stop | `make stop` |
| Build | `make build` |
| Test | `make test` |
| CLI | `make cli CMD='query "…"'` |

## Documentation

| Topic | Location |
|-------|----------|
| **Index** | [docs/README.md](docs/README.md) |
| Getting started | [docs/getting-started/](docs/getting-started/) (installation, quick start, LLM setup) |
| Configuration | [docs/configuration.md](docs/configuration.md) |
| API reference | [docs/api/README.md](docs/api/README.md) |
| API examples | [docs/api/examples.md](docs/api/examples.md) |
| CLI usage | [docs/usage/cli-usage.md](docs/usage/cli-usage.md) |
| Period comparison | [docs/features/period-comparison.md](docs/features/period-comparison.md) |
| Troubleshooting | [docs/troubleshooting.md](docs/troubleshooting.md) |
| Contributing | [.github/CONTRIBUTING.md](.github/CONTRIBUTING.md) |

## License

MIT. See [LICENSE](LICENSE).
