# PgQueryNarrative

Turn SQL query results into business narratives with AI. Run read-only SQL against PostgreSQL, get metrics and chart suggestions, and generate narrative reports using your choice of LLM (Ollama, OpenAI, Claude, Gemini, Groq).

## Quick start

**With Docker** (PostgreSQL and app in containers):

```bash
make start-docker
```

**With local PostgreSQL** (app on host; PostgreSQL must be running):

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

- **Docker** (for `make start-docker`), or **PostgreSQL 16+** and **Go 1.24+** (for `make start-local` and building)

## Commands

| Action | Command |
|--------|--------|
| Start (Docker) | `make start-docker` |
| Start (local) | `make start-local` |
| Stop | `make stop` |
| Build | `make build` |
| Test | `make test` |
| CLI | `make cli CMD='query "SELECT * FROM demo.sales LIMIT 5"'` |

## Project structure

| Path | Purpose |
|------|---------|
| `cmd/server` | Application entrypoint |
| `app/` | Core logic (config, DB, query runner, metrics, LLM, narrative, service) |
| `api/design/` | API design (Goa); generated code in `api/gen/` |
| `web/` | Web UI and static assets |
| `docs/` | User and developer documentation |
| `test/unit/` | Unit tests by feature |
| `changelog/` | Release history |

## Documentation

Full documentation is in the **[docs](docs/README.md)** directory:

- **Getting started:** [Installation](docs/getting-started/installation.md), [Quick start](docs/getting-started/quickstart.md), [LLM setup](docs/getting-started/llm-setup.md)
- **User guides:** [Configuration](docs/configuration.md), [CLI usage](docs/usage/cli-usage.md)
- **API:** [Reference](docs/api/README.md), [Examples](docs/api/examples.md)
- **Features:** [Period comparison](docs/features/period-comparison.md)
- **Reference:** [Troubleshooting](docs/reference/troubleshooting.md), [Docker resources](docs/reference/docker-resources.md), [PostgreSQL extension](docs/reference/postgres-extension.md)
- **Development:** [Setup](docs/development/setup.md), [Testing](docs/development/testing.md)

## License

MIT. See [LICENSE](LICENSE).
