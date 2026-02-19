# Configuration

Environment variables only. Sensible defaults for local use.

## Environment variables

### Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_DEBUG` | (empty) | `1` or `true` = verbose logging (queries, report generation) |

### Server

| Variable | Default | Description |
|----------|---------|-------------|
| `PGQUERYNARRATIVE_HOST` | `0.0.0.0` | Server bind address |
| `PGQUERYNARRATIVE_PORT` | `8080` | Server port |
| `PGQUERYNARRATIVE_READ_TIMEOUT` | `15s` | Request read timeout |
| `PGQUERYNARRATIVE_WRITE_TIMEOUT` | `60s` | Response write timeout |

### Database

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_IMAGE` | `postgres:18-alpine` | Docker PostgreSQL image (use `postgres:17-alpine` or `postgres:16-alpine` for older versions) |
| `DATABASE_HOST` | `localhost` | Database host |
| `DATABASE_PORT` | `5432` | Database port |
| `DATABASE_NAME` | `pgquerynarrative` | Database name |
| `DATABASE_USER` | `pgquerynarrative_app` | Application user |
| `DATABASE_PASSWORD` | `pgquerynarrative_app` | Application password |
| `DATABASE_READONLY_USER` | `pgquerynarrative_readonly` | Read-only user |
| `DATABASE_READONLY_PASSWORD` | `pgquerynarrative_readonly` | Read-only password |
| `DATABASE_SSL_MODE` | `disable` | SSL mode (disable/require/verify-full) |
| `DATABASE_MAX_CONNECTIONS` | `10` | Max connection pool size |
| `QUERY_TIMEOUT` | `30s` | Query execution timeout |

### LLM

Required for report generation. See [LLM setup](getting-started/llm-setup.md) for providers and MCP.

| Variable | Default | Description |
|----------|---------|-------------|
| `LLM_PROVIDER` | `ollama` | LLM provider (ollama/gemini/claude/openai/groq) |
| `LLM_MODEL` | `llama3.2` | Model name |
| `LLM_BASE_URL` | `http://localhost:11434` | LLM API base URL. Use `http://host.docker.internal:11434` when running in Docker. |
| `LLM_API_KEY` | `` | API key (for cloud providers) |

**Ollama (local)**

```bash
export LLM_PROVIDER=ollama
export LLM_MODEL=llama3.2
export LLM_BASE_URL=http://localhost:11434
make start-docker   # or make start-local
```

**Gemini** – API key from [Google AI Studio](https://aistudio.google.com/app/apikey):

```bash
export LLM_PROVIDER=gemini
export LLM_MODEL=gemini-2.0-flash
export LLM_API_KEY=your_api_key_here
make start-docker   # or make start-local
```

On 404 try `gemini-1.5-flash`. On 429 the client retries with backoff. Do not commit API keys; use `.env` (gitignored) or shell.

**Claude** – API key from [Anthropic Console](https://console.anthropic.com/):

```bash
export LLM_PROVIDER=claude
export LLM_MODEL=claude-3-5-sonnet-20241022
export LLM_API_KEY=your_api_key_here
make start-docker   # or make start-local
```

See [Anthropic models](https://docs.anthropic.com/en/docs/about-claude/models). On 429 the client retries.

**OpenAI (GPT)** – API key from [OpenAI API keys](https://platform.openai.com/api-keys):

```bash
export LLM_PROVIDER=openai
export LLM_MODEL=gpt-4o-mini
export LLM_API_KEY=your_api_key_here
make start-docker   # or make start-local
```

Other models: `gpt-4o`, `gpt-4-turbo`, etc. On 429 the client retries. Do not commit API keys.

**Groq** – API key from [Groq Console](https://console.groq.com/keys):

```bash
export LLM_PROVIDER=groq
export LLM_MODEL=llama-3.3-70b-versatile
export LLM_API_KEY=your_api_key_here
make start-docker   # or make start-local
```

Other models: `llama-3.1-8b-instant`, `mixtral-8x7b-32768`. On 429 the client retries.

**MCP (Claude desktop / Cursor)**

Connect Claude desktop or Cursor via the Model Context Protocol so they can run queries and generate reports as tools. Configure by **editing the MCP config file** (no “add server” UI).

1. Start PgQueryNarrative (e.g. `make start-local`). Build the MCP server: `make build-mcp` (produces `bin/mcp-server`).

2. **Edit the MCP config file** for your client:

   **Claude desktop**
   - Config file:
     - **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
     - **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`
     - **Linux:** `~/.config/Claude/claude_desktop_config.json`
   - If the file has no `mcpServers` key yet, use the structure below. If it already has `mcpServers`, add only the `pgquerynarrative` entry inside that object.

   **Cursor**
   - Edit MCP config (e.g. Cursor Settings → MCP, or the config file Cursor uses for MCP). Add the same `pgquerynarrative` entry under `mcpServers`.

3. **Add this under `mcpServers`** (replace the path with the full path to your `bin/mcp-server`):

   ```json
   "pgquerynarrative": {
     "command": "/FULL/PATH/TO/PgQueryNarrative/bin/mcp-server"
   }
   ```

If app is not at `http://localhost:8080`, add `"env": { "PGQUERYNARRATIVE_URL": "http://localhost:8080" }`. See `config/mcp-example.json`. (4) Restart client. Tools: `run_query`, `generate_report`, `list_saved_queries`, `get_report`, `list_reports`.

### Metrics and period comparison

When query results include a date/time column and at least one numeric measure, the app computes period-over-period comparison (e.g. latest period vs previous). The trend label (up / down / flat) uses a configurable threshold.

| Variable | Default | Description |
|----------|---------|-------------|
| `PERIOD_TREND_THRESHOLD_PERCENT` | `0.5` | Minimum absolute % change to label as "up" or "down"; below this, trend is "flat". |

**Examples:** `0.25` = more sensitive (smaller changes show as up/down). `1` = only changes ≥ 1% are up/down. See [Period comparison](features/period-comparison.md) for usage and testing.

### Security

| Variable | Default | Description |
|----------|---------|-------------|
| `SECURITY_AUTH_ENABLED` | `false` | Enable auth (future) |

## Loading config

**Env:** `export PGQUERYNARRATIVE_PORT=8081` then `make start-docker` or `make start-local`.

**.env:** Create `.env` in project root (gitignored), then `export $(cat .env | xargs)` before starting. Do not commit secrets.

**Docker Compose:** Set `environment` under `app` in `docker-compose.yml`.

**Systemd:** Example unit `/etc/systemd/system/pgquerynarrative.service`:

```ini
[Unit]
Description=PgQueryNarrative Server
After=network.target postgresql.service

[Service]
Type=simple
User=pgquerynarrative
Environment="PGQUERYNARRATIVE_PORT=8080"
Environment="DATABASE_HOST=localhost"
ExecStart=/usr/local/bin/pgquerynarrative
Restart=always

[Install]
WantedBy=multi-user.target
```

## Production

- Change default passwords; use SSL for DB (`DATABASE_SSL_MODE=require`); use secrets management. Recommended: `QUERY_TIMEOUT=60s`, `DATABASE_MAX_CONNECTIONS=50`, `SECURITY_AUTH_ENABLED=true` when available.

Config is validated on startup; invalid values cause clear startup errors. If config does not apply: check `env | grep PGQUERYNARRATIVE`, `.env` format, or Compose `environment`. For other issues see [Troubleshooting](reference/troubleshooting.md).

**See also:** [Installation](getting-started/installation.md), [Troubleshooting](reference/troubleshooting.md), [Period comparison](features/period-comparison.md)
