# Testing

Unit and integration tests for PgQueryNarrative. All unit tests live under `test/unit/` and can be run by package or by test name.

## Running tests

**All unit tests (recommended):**

```bash
make test-unit
```

Or explicitly:

```bash
go test ./test/unit/... ./cmd/server/... -v
```

**By package:**

```bash
go test ./test/unit/app/charts/... -v       # Chart suggestions
go test ./test/unit/app/metrics/... -v      # Metrics, time series, data quality
go test ./test/unit/app/queryrunner/... -v  # Query validation
go test ./test/unit/app/story/... -v        # Narrative sanitizer
go test ./test/unit/app/service/... -v      # Reports service, API conversion
go test ./test/unit/web/... -v              # Report HTML formatting
```

**Single test:**

```bash
go test ./test/unit/app/service/... -run TestBuildPerfSuggestions_LimitApplied -v
```

**Integration tests** (require Docker):

```bash
go test ./test/integration/... -v
```

**End-to-end tests:**

```bash
go test ./test/e2e/... -v
```

## Test layout

| Package | What is tested |
|---------|-----------------|
| `test/unit/app/charts` | Chart suggestions (bar, line, pie, area, table) |
| `test/unit/app/metrics` | Period comparison, trend, anomalies, data quality, std dev, period labels |
| `test/unit/app/queryrunner` | SQL validation (schema, SELECT-only) |
| `test/unit/app/story` | Narrative sanitizer (no fabricated "previous period") |
| `test/unit/app/service` | Perf suggestions, metrics-to-API conversion |
| `test/unit/web` | Report HTML (charts, data quality, perf, narrative sections) |
| `cmd/server` | Request logging middleware |

See `test/unit/README.md` in the repository for more detail.

## QA checklist

For manual or automated QA, the following areas should be verified:

- **Chart suggestions:** API returns suggestions (bar, line, pie, area, table); UI dropdown and report section show them.
- **Period comparison:** Time-series queries show "Vs previous period," period labels, trend summary, moving average, anomalies, next-period forecast where applicable.
- **Advanced analytics:** Data quality table (nulls, distinct, rows); performance suggestions when limit applied or query slow; std dev in aggregates when applicable.
- **Narrative:** Single-period queries do not mention "previous period" or "same period last year"; time-series reports include period comparison in takeaways.
- **API:** `POST /reports/generate`, `GET /reports/{id}`, `GET /reports` return expected structure; metrics include `data_quality`, `perf_suggestions`, `time_series` with optional fields.
- **Errors:** Invalid SQL → 400; report not found → 404; LLM failure → 500 with clear message.

Example API checks:

```bash
curl -s -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT SUM(total_amount) AS total_fares, COUNT(*) AS trips FROM demo.sales"}'

curl -s "http://localhost:8080/api/v1/reports?limit=5&offset=0"
```

## See also

- [API examples](../api/examples.md)
- [Development setup](setup.md)
- [Troubleshooting](../reference/troubleshooting.md)
