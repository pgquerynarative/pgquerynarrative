# Quick start

## Prerequisites

- Docker, or PostgreSQL 16+ and Go 1.24+

## Docker

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make start-docker
```

App: `http://localhost:8080`

## Local PostgreSQL

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
pg_isready   # ensure Postgres is running
make start-local
```

## First steps

- **Web:** http://localhost:8080
- **CLI:** `make cli CMD='query "SELECT * FROM demo.sales LIMIT 5"'` ([CLI usage](../usage/cli-usage.md))
- **API:** [API examples](../api/examples.md)

**See also:** [Installation](installation.md), [LLM setup](llm-setup.md) (reports, MCP), [Troubleshooting](../troubleshooting.md), [Documentation index](../README.md)
