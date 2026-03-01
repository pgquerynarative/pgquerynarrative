# Testing

Unit, integration, and E2E tests. Unit tests live in `test/unit/`; run by package or by test name. See [Development setup](setup.md) for prerequisites and commands.

## Running tests

**All unit tests:**

```bash
make test-unit
```

Or: `go test ./test/unit/... ./cmd/server/... ./pkg/narrative/... -v`

**By package:**

```bash
go test ./test/unit/app/catalog/... -v
go test ./test/unit/app/charts/... -v
go test ./test/unit/app/metrics/... -v
go test ./test/unit/app/llm/... -v
go test ./test/unit/app/queryrunner/... -v
go test ./test/unit/app/service/... -v
go test ./test/unit/app/story/... -v
go test ./test/unit/app/suggestions/... -v
go test ./test/unit/web/... -v
go test ./cmd/server/... -v
go test ./pkg/narrative/... -v
```

**Single test:** `go test ./test/unit/app/service/... -run TestBuildPerfSuggestions_LimitApplied -v`

**Integration tests** (require Docker): `go test ./test/integration/... -v`

**E2E tests:** `go test ./test/e2e/... -v`

## Test layout

| Package | What is tested |
|---------|-----------------|
| `test/unit/app/catalog` | Schema/catalog loader |
| `test/unit/app/charts` | Chart suggestions |
| `test/unit/app/metrics` | Period comparison, trend, anomalies, data quality |
| `test/unit/app/llm` | Narrative prompt builder |
| `test/unit/app/queryrunner` | SQL validation (schema, SELECT-only, disallowed keywords) |
| `test/unit/app/story` | Narrative sanitizer |
| `test/unit/app/service` | Perf suggestions, metrics-to-API conversion |
| `test/unit/app/suggestions` | Query suggestions (curated, limit) |
| `test/unit/web` | Report HTML/export |
| `cmd/server` | Request logging middleware |
| `pkg/narrative` | Client, run-query options, validation |

**Integration** (`test/integration/...`): query runner; schema and suggestions against real Postgres; reports List/Get.

**E2E** (`test/e2e/...`): full HTTP API against real Postgres (queries, schema, suggestions, reports).

## QA checklist

**Queries:** Run query; period comparison and chart suggestions in response; save, list, get, delete saved queries; invalid SQL → 400.

**Schema:** `GET /api/v1/schema` returns allowed schemas, tables, columns.

**Suggestions:** `GET /api/v1/suggestions/queries` (curated + intent match); `GET /api/v1/suggestions/similar` (embeddings); `POST /api/v1/suggestions/ask` (NL→SQL→report); `POST /api/v1/suggestions/explain` (plain-English SQL explanation).

**Reports:** Generate, get, list; metrics (period comparison, time-series, anomalies, trend, data quality, perf); narrative content; errors (not found → 404, LLM failure → 500).

**Export:** HTML and PDF at `/web/reports/export?id=...` and `/web/reports/export/pdf?id=...`.

**Probes:** `GET /health`, `GET /ready`, `GET /metrics`, `GET /version` return 200.

**UI:** Query Runner schema browser, query suggestions card, shortcuts (Ctrl+E, Ctrl+Enter); Reports export buttons.

Example API checks: [API examples](../api/examples.md).

## Testing auth, rate limit, and audit

You must start the server with auth and rate limiting enabled; otherwise unauthenticated requests succeed and rate limits never trigger. Defaults: `SECURITY_AUTH_ENABLED=false`, `SECURITY_RATE_LIMIT_RPM=0` (disabled).

Start the server with security enabled and a low rate limit so you can trigger 429 quickly:

```bash
SECURITY_AUTH_ENABLED=true SECURITY_API_KEY=test-key SECURITY_RATE_LIMIT_RPM=5 make run
```

**Auth** — Without a token, protected endpoints return **401**. With the correct Bearer token they return **200**.

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/api/v1/queries/saved
curl -s -o /dev/null -w "%{http_code}\n" -H "Authorization: Bearer test-key" http://localhost:8080/api/v1/queries/saved
```

**Health/ready** — Always unprotected; returns **200**.

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/health
```

**Rate limit** — With `SECURITY_RATE_LIMIT_RPM=5`, send more than 5 requests per minute from the same machine; the 6th and later return **429** (first five return **200**).

```bash
for i in $(seq 8); do curl -s -o /dev/null -w "%{http_code}\n" -H "Authorization: Bearer test-key" http://localhost:8080/api/v1/queries/saved; done
```

**Audit log** — Check `app.audit_logs` for `API_REQUEST`, `AUTH_FAILURE`, `RATE_LIMIT_EXCEEDED`.

```bash
psql -d pgquerynarrative -c "SELECT event_type, details, user_id, ip_address, created_at FROM app.audit_logs ORDER BY created_at DESC LIMIT 10;"
```

## See also

- [Development setup](setup.md) · [API examples](../api/examples.md) · [Documentation index](../README.md)
