# API reference

REST API base: `http://localhost:8080/api/v1`. No authentication in current version.

## Queries

| Method | Path | Description |
|--------|------|-------------|
| POST | `/queries/run` | Body: `{"sql":"...", "limit": 100}`. Run read-only SQL. Returns optional `period_comparison` and `chart_suggestions`. |
| POST | `/queries/saved` | Body: `{"name","sql","tags"}`. Save query. |
| GET | `/queries/saved` | Query: `limit`, `offset`, `tags`. List saved. |
| GET | `/queries/saved/{id}` | Get saved query. |
| DELETE | `/queries/saved/{id}` | Delete saved query. |

## Reports

| Method | Path | Description |
|--------|------|-------------|
| POST | `/reports/generate` | Body: `{"sql":"...", "saved_query_id": "uuid"}`. Generate report (requires LLM). |
| GET | `/reports/{id}` | Get report. Metrics may include `time_series`, `period_current_label`, `period_previous_label`. |
| GET | `/reports` | Query: `limit`, `offset`, `saved_query_id`. List reports. |

## Errors

JSON: `{"name","message","code"}`. Codes: `VALIDATION_ERROR`, `TIMEOUT_ERROR`, `NOT_FOUND`, `LLM_ERROR`. No rate limiting in current version.

**OpenAPI:** `gen/http/openapi.json`, `openapi.yaml`, `openapi3.json`, `openapi3.yaml` (in repo). **Examples:** [API examples](examples.md).
