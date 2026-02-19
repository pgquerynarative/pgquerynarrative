# Troubleshooting

Common issues and how to resolve them.

## Environment and dependencies

| Issue | Solution |
|-------|----------|
| **Docker not found** | Install [Docker Desktop](https://www.docker.com/products/docker-desktop). Verify with `docker info`. |
| **make: command not found** | Install Make (e.g. `brew install make` on macOS). |
| **Port 8080 in use** | Use another port: `PGQUERYNARRATIVE_PORT=8081 make start-docker` or `make start-local`. |

## Database

| Issue | Solution |
|-------|----------|
| **PostgreSQL connection refused** | **Docker:** run `make start-docker`. **Local:** start PostgreSQL first (e.g. `brew services start postgresql@18`), then `make start-local`. |
| **Role does not exist / permission denied** | Run `make db-init` then `make migrate`, or run `make start-local` once. If `demo.sales` is denied: `psql -d pgquerynarrative -U pgquerynarrative_app -c "GRANT SELECT ON demo.sales TO pgquerynarrative_readonly;"` |

## Reports and LLM

| Issue | Solution |
|-------|----------|
| **Failed to parse narrative JSON** | LLM output may be truncated. Ensure Ollama is running and the model is pulled (`ollama serve`, `ollama pull llama3.2`). Try a larger model or restart Ollama. |
| **Report generation fails or times out** | See [LLM setup](../getting-started/llm-setup.md). With the app in Docker and Ollama on the host, set `LLM_BASE_URL=http://host.docker.internal:11434`. |
| **Narrative shows wrong number scale (e.g. 848M instead of 84.8M)** | The prompt and sample data use comma-separated thousands and correct magnitude. Ensure you are on the latest version. |

## Features

| Issue | Solution |
|-------|----------|
| **Period comparison ("Vs previous period") not shown** | The query must include a date/time column and at least one numeric measure, with at least two result rows. See [Period comparison](../features/period-comparison.md). |

## See also

- [Configuration](../configuration.md)
- [Installation](../getting-started/installation.md)
- [LLM setup](../getting-started/llm-setup.md)
