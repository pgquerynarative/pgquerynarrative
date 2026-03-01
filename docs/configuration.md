# Configuration

PgQueryNarrative is configured via **environment variables** only. Sensible defaults apply for local use. Invalid config produces clear startup errors (see [Troubleshooting](reference/troubleshooting.md)).

## Loading config

| Method | Usage |
|--------|--------|
| **Env** | `export PGQUERYNARRATIVE_PORT=8081` then start. |
| **.env** | Create `.env` in project root (gitignored); `export $(cat .env | xargs)` before starting. Do not commit secrets. |
| **Docker Compose** | Set `environment` under `app` in [docker-compose.yml](../docker-compose.yml). |

---

## Logging

| Variable | Default | Description |
|---------|---------|-------------|
| `LOG_DEBUG` | (empty) | `1` or `true` = verbose logging (query execution, report generation). |

---

## Server

| Variable | Default | Description |
|---------|---------|-------------|
| `PGQUERYNARRATIVE_HOST` | `0.0.0.0` | Bind address. |
| `PGQUERYNARRATIVE_PORT` | `8080` | Server port. |
| `PGQUERYNARRATIVE_READ_TIMEOUT` | `15s` | Request read timeout. |
| `PGQUERYNARRATIVE_WRITE_TIMEOUT` | `60s` | Response write timeout. |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown timeout. |
| `CORS_ORIGINS` | (empty) | Comma-separated origins for CORS; when set, `Access-Control-Allow-Origin` is sent for matching request origins. Empty = same-origin only. |

---

## Database

| Variable | Default | Description |
|---------|---------|-------------|
| `POSTGRES_IMAGE` | `postgres:18-alpine` | Docker Postgres image (root Compose). |
| `DATABASE_HOST` | `localhost` | Database host. |
| `DATABASE_PORT` | `5432` | Database port. |
| `DATABASE_NAME` | `pgquerynarrative` | Database name. |
| `DATABASE_USER` | `pgquerynarrative_app` | Application user (migrations, saved_queries, reports). |
| `DATABASE_PASSWORD` | `pgquerynarrative_app` | Application password. |
| `DATABASE_READONLY_USER` | `pgquerynarrative_readonly` | Read-only user (query execution). |
| `DATABASE_READONLY_PASSWORD` | `pgquerynarrative_readonly` | Read-only password. |
| `DATABASE_SSL_MODE` | `disable` | SSL mode: `disable` \| `require` \| `verify-full`. |
| `DATABASE_MAX_CONNECTIONS` | `10` | Max connection pool size. |
| `QUERY_TIMEOUT` | `30s` | Query execution timeout. |
| `DATABASE_ALLOWED_SCHEMAS` | `public,demo` | Comma-separated schemas queries may access. Use your own schemas here (e.g. `public` or `public,analytics`). |

---

## LLM {#llm}

Required for [report generation](api/README.md#reports). See [LLM setup](getting-started/llm-setup.md).

| Variable | Default | Description |
|---------|---------|-------------|
| `LLM_PROVIDER` | `ollama` | `ollama` \| `gemini` \| `claude` \| `openai` \| `groq`. |
| `LLM_MODEL` | `llama3.2` | Model name. |
| `LLM_BASE_URL` | `http://localhost:11434` | LLM API base URL. Docker with Ollama on host: `http://host.docker.internal:11434`. |
| `LLM_API_KEY` | (empty) | API key (required for cloud providers). |

---

## Embeddings (optional) {#embeddings}

Used for [GET /api/v1/suggestions/similar](api/README.md#suggestions) and RAG context in report generation. When not set, those features are disabled.

| Variable | Default | Description |
|---------|---------|-------------|
| `EMBEDDING_BASE_URL` | (empty) | Embedding API URL. If empty and `LLM_PROVIDER=ollama`, defaults to `LLM_BASE_URL`. |
| `EMBEDDING_MODEL` | `nomic-embed-text` | Embedding model (e.g. Ollama `nomic-embed-text`). |

Ollama: `ollama pull nomic-embed-text`. See [Semantic search (pgvector)](reference/semantic-search-pgvector.md).

---

## MCP (Claude desktop / Cursor) {#mcp-claude-desktop--cursor}

For [LLM setup – MCP](getting-started/llm-setup.md#mcp-claude-desktop--cursor):

1. Build: `make build-mcp` → `bin/mcp-server`.
2. Edit MCP config:
   - **Claude:** macOS `~/Library/Application Support/Claude/claude_desktop_config.json`; Windows `%APPDATA%\Claude\`; Linux `~/.config/Claude/`.
   - **Cursor:** `.cursor/mcp.json` in the project root (or Settings → MCP). See `config/mcp-example.json` for the template.
3. Add under `mcpServers` (replace path):
   ```json
   "pgquerynarrative": {
     "command": "/path/to/pgquerynarrative/bin/mcp-server"
   }
   ```
   If the app is not at http://localhost:8080, set `"env": { "PGQUERYNARRATIVE_URL": "http://localhost:PORT" }`. If the app has auth enabled (`SECURITY_AUTH_ENABLED=true`), set `"env": { "PGQUERYNARRATIVE_API_KEY": "your-secret-key" }` (same value as `SECURITY_API_KEY`). See `config/mcp-example.json`.
4. Restart the client. Available tools: `run_query`, `generate_report`, `list_saved_queries`, `get_report`, `list_reports`, `get_schema`, `get_context`, `suggest_queries`, `list_schemas`, `ask_question`, `explain_sql`. See [How to use MCP tools](getting-started/llm-setup.md#how-to-use-mcp-tools-in-cursor--claude) below.

---

## Metrics

| Variable | Default | Description |
|---------|---------|-------------|
| `PERIOD_TREND_THRESHOLD_PERCENT` | `0.5` | Min % change to label trend "up"/"down"; below = "flat". |

---

## Security

| Variable | Default | Description |
|---------|---------|-------------|
| `SECURITY_AUTH_ENABLED` | `false` | When true, `/api/*` and `/web/reports/export*` require `Authorization: Bearer <SECURITY_API_KEY>`. `/health` and `/ready` are never protected. |
| `SECURITY_API_KEY` | (empty) | Bearer token for API auth. Required when `SECURITY_AUTH_ENABLED` is true. |
| `SECURITY_RATE_LIMIT_RPM` | `0` | Max requests per minute per client IP (0 = disabled). |
| `SECURITY_RATE_LIMIT_BURST` | `0` | Burst size for rate limiter (0 = 2× RPM). |

---

## Production

- Change default passwords; use secrets management.
- Use SSL for DB: `DATABASE_SSL_MODE=require`.
- Recommended: `QUERY_TIMEOUT=60s`, `DATABASE_MAX_CONNECTIONS=50`.

See [Deployment](reference/deployment.md) and [Operations](reference/operations.md).

---

## See also

- [Installation](getting-started/installation.md) · [API reference](api/README.md) · [Deployment](reference/deployment.md) · [Documentation index](README.md)
