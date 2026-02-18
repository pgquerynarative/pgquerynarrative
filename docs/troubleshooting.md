# Troubleshooting

**Docker not found:** Install [Docker Desktop](https://www.docker.com/products/docker-desktop); run `docker info` to confirm.

**PostgreSQL connection refused:** Docker: `make start-docker`. Local: start Postgres first (e.g. `brew services start postgresql@18`), then `make start-local`.

**role does not exist / permission denied:** Run `make local-db-init` then `make migrate`, or `make start-local` once. If `demo.sales` denied: `psql -d pgquerynarrative -U pgquerynarrative_app -c "GRANT SELECT ON demo.sales TO pgquerynarrative_readonly;"`

**Port 8080 in use:** `PGQUERYNARRATIVE_PORT=8081 make start-docker` (or `make start-local`).

**make: command not found:** Install Make (e.g. `brew install make`).

**failed to parse narrative JSON:** LLM output truncated. Ensure Ollama running and model pulled (`ollama serve`, `ollama pull llama3.2`); try larger model or restart Ollama.

**Report generation fails or times out:** [LLM setup](getting-started/llm-setup.md). With app in Docker and Ollama on host: `LLM_BASE_URL=http://host.docker.internal:11434`.

**Period comparison ("Vs previous period") not shown:** Query must include a date/time column and at least one numeric measure, with at least two result rows. See [Period comparison](features/period-comparison.md).

**See also:** [Configuration](configuration.md), [Installation](getting-started/installation.md)
