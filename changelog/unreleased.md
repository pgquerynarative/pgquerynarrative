## [Unreleased]

All items below are **Release 1** scope. The next version (e.g. v1.0.0) will ship Release 1.

### Added
- **Period comparison:** Automatic period-over-period (e.g. this month vs last month) when query results have a date/time column and numeric measures; derived % change and trend (up/down/flat); "Vs previous period" block in query results and report UI with optional period labels
- **Run query API:** `period_comparison` array and optional `period_current_label` / `period_previous_label` on run response
- **Reports API:** `metrics.period_current_label` and `metrics.period_previous_label` when time series is present
- **Configurable trend threshold:** `PERIOD_TREND_THRESHOLD_PERCENT` (default 0.5) for when to label change as up/down vs flat
- **Narrative:** LLM prompt rule to include at least one takeaway on period-over-period change when time series metrics exist
- **Metrics:** Support for PostgreSQL `NUMERIC` (pgtype.Numeric) in column profiling and aggregation so measures like `SUM(...)` appear in period comparison
- Groq LLM provider: `LLM_PROVIDER=groq`, `LLM_MODEL`, `LLM_API_KEY` (OpenAI-compatible API; e.g. llama-3.3-70b-versatile)
- OpenAI (GPT) LLM provider: `LLM_PROVIDER=openai`, `LLM_MODEL`, `LLM_API_KEY` (Chat Completions API)
- Claude LLM provider: `LLM_PROVIDER=claude`, `LLM_MODEL`, `LLM_API_KEY` (Anthropic Messages API)
- Gemini LLM provider: `LLM_PROVIDER=gemini`, `LLM_MODEL`, `LLM_API_KEY` for report generation
- MCP server (`cmd/mcp-server`): tools for Claude desktop / Cursor (run query, generate report, list saved/reports); `config/mcp-example.json`, docs
- **MCP schema, context, and query suggestions:** Schema API `GET /api/v1/schema` (queryable tables/columns from `information_schema`); suggestions API `GET /api/v1/suggestions/queries?intent=...&limit=...` (curated examples + saved-query match by intent); MCP tools `get_schema`, `get_context` (schema + saved queries merged), `suggest_queries`; `app/catalog`, `app/suggestions`; design and testing in `docs/development/mcp-schema-context-design.md`
- Report UI: show LLM provider and model; improved report card layout and CSS
- PostgreSQL extension for calling PgQueryNarrative from SQL
- CLI tool for Docker-only usage
- API documentation, contributing guidelines, security policy
- Security scanning (secret scan, CodeQL, govulncheck, gosec)
- **Chart suggestions:** By data structure (time series → line/area; category+value → bar/pie; table); suggestion buttons and chart-type dropdown built from API on query page; area chart support; report page shows suggested charts; unit tests (app/charts/suggester_test.go)
- **Advanced metrics:** Richer time-series (last N periods, 3-period moving average); anomaly detection (z-score, configurable threshold); trend analysis (linear regression over last 6 periods, direction and summary); report API (`periods`, `moving_average`, `anomalies`, `trend_summary`); report UI (trend summary, anomalies list, period history table); unit tests (app/metrics/calculator_test.go)

### Changed
- **Documentation:** Reorganized into `docs/api/`, `docs/usage/`, `docs/features/`; added Period comparison doc and Metrics section in configuration (`PERIOD_TREND_THRESHOLD_PERCENT`); API examples in `docs/api/examples.md`, CLI usage in `docs/usage/cli-usage.md`; docs index and cross-links updated
- Documentation: single generic LLM setup guide (Ollama, Gemini, Claude, OpenAI, Groq, MCP); docs shortened and standardized
- Go 1.23 → 1.24; PostgreSQL 18 as default (16, 17, 18 supported)
- Docker: postgres:18-alpine, memory limits

### Fixed
- E2E migration: roles created in 000001 so 000003 GRANT succeeds; migration permission errors
- CLI shell and argument passing for Alpine/Docker
- Postgres init script role creation order
- **Narrative number scale:** LLM prompt now formats sample data with comma-separated thousands and instructs the model to preserve exact magnitude when citing metrics (avoids e.g. 848M instead of 84.8M)

### Security
- Secret scanning, dependency vulnerability scanning, CodeQL, gosec
