# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Entries are managed in `changelog/` — see [changelog/README.md](changelog/README.md) for how to add and release.

## [Unreleased]

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
- Report UI: show LLM provider and model; improved report card layout and CSS
- PostgreSQL extension for calling PgQueryNarrative from SQL
- CLI tool for Docker-only usage
- API documentation, contributing guidelines, security policy
- Security scanning (secret scan, CodeQL, govulncheck, gosec)

### Changed
- **Documentation:** Reorganized into `docs/api/`, `docs/usage/`, `docs/features/`; added [Period comparison](docs/features/period-comparison.md), Metrics section in configuration (`PERIOD_TREND_THRESHOLD_PERCENT`); API examples moved to `docs/api/examples.md`, CLI usage to `docs/usage/cli-usage.md`; docs index and cross-links updated
- Documentation: single generic LLM setup guide (Ollama, Gemini, Claude, OpenAI, Groq, MCP); docs shortened and standardized
- Go 1.23 → 1.24; PostgreSQL 18 as default (16, 17, 18 supported)
- Docker: postgres:18-alpine, memory limits

### Fixed
- E2E migration: roles created in 000001 so 000003 GRANT succeeds; migration permission errors
- CLI shell and argument passing for Alpine/Docker
- Postgres init script role creation order

### Security
- Secret scanning, dependency vulnerability scanning, CodeQL, gosec

## [0.1.0] - 2026-02-17

### Added
- Initial release
- Query execution engine
- Query validation and security
- Metrics calculation
- LLM integration for narrative generation
- Web UI
- RESTful API
- Docker support
- Database migrations
- Demo data seeding

[Unreleased]: https://github.com/pgquerynarrative/pgquerynarrative/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/pgquerynarrative/pgquerynarrative/releases/tag/v0.1.0
