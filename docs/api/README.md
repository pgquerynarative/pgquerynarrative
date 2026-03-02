# API reference

REST API base: `http://localhost:8080/api/v1` (override with [Configuration](../configuration.md) server port). All request/response bodies are JSON. When `SECURITY_AUTH_ENABLED` is true, send `Authorization: Bearer <SECURITY_API_KEY>`; otherwise requests are unauthenticated. Rate limiting (when `SECURITY_RATE_LIMIT_RPM` > 0) returns 429. **OpenAPI 3:** `api/gen/http/openapi3.json` and `api/gen/http/openapi3.yaml` for codegen and API tooling.

## Queries

| Method | Path | Description |
|--------|------|-------------|
| POST | `/queries/run` | Body: `{"sql":"...", "limit": 100}`. Run read-only SQL. Returns `columns`, `rows`, `row_count`, `execution_time_ms`, optional `period_comparison`, `chart_suggestions`. |
| POST | `/queries/saved` | Body: `{"name","sql","tags"}`. Save query. |
| GET | `/queries/saved` | Query: `limit`, `offset`, `tags`. List saved queries. |
| GET | `/queries/saved/{id}` | Get saved query by ID. |
| DELETE | `/queries/saved/{id}` | Delete saved query. |

## Schema

| Method | Path | Description |
|--------|------|-------------|
| GET | `/schema` | Allowed schemas, tables, columns. Used by MCP `get_schema`. |

## Suggestions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/suggestions/queries` | Query: `intent`, `limit` (default 5). Suggested SQL (curated + saved-query match). |
| GET | `/suggestions/similar` | Query: `text` (required), `limit` (default 5). Saved queries semantically similar to text (embedding-based). Requires [embeddings](../reference/semantic-search-pgvector.md) enabled. |
| POST | `/suggestions/ask` | Body: `{"question":"..."}`. Natural language → SQL → run → narrative report. Requires [LLM](../getting-started/llm-setup.md). |
| POST | `/suggestions/explain` | Body: `{"sql":"..."}`. Plain-English explanation of SQL. Requires [LLM](../getting-started/llm-setup.md). |

## Reports

| Method | Path | Description |
|--------|------|-------------|
| POST | `/reports/generate` | Body: `{"sql":"...", "saved_query_id": "uuid"}`. Generate report (requires [LLM](../getting-started/llm-setup.md)). Returns `narrative`, `metrics`. |
| GET | `/reports/{id}` | Get report. Metrics: `time_series`, `correlations`, `cohorts` (when [cohort shape](../configuration.md#cohort-analysis) present), `data_quality`, `perf_suggestions`, etc. |
| GET | `/reports` | Query: `limit`, `offset`, `saved_query_id`. List reports. |

## Errors

Response JSON: `{"name","message","code"}`. Codes: `VALIDATION_ERROR`, `TIMEOUT_ERROR`, `NOT_FOUND`, `LLM_ERROR`, `UNAUTHORIZED` (401 when auth enabled and missing/invalid Bearer), `RATE_LIMIT_EXCEEDED` (429 when rate limit applies).

## See also

- [API examples](examples.md) — Single run-query cURL example
- [Configuration](../configuration.md) — Environment variables
- [Deployment](../reference/deployment.md) — Running the API
- [Embedded integration](../getting-started/embedded.md) — Library and middleware (different path prefix)
- [Documentation index](../README.md)
