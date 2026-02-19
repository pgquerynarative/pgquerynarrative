# API examples

Base URL: `http://localhost:8080` (or set `PGQUERYNARRATIVE_PORT` if different). All examples use `application/json`.

## Run query

```bash
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category",
    "limit": 10
  }'
```

Response includes `columns`, `rows`, `row_count`, `execution_time_ms`, optional `chart_suggestions` and `period_comparison` (when the result has a time column and measures). See [Period comparison](../features/period-comparison.md).

## Save query

```bash
curl -X POST http://localhost:8080/api/v1/queries/saved \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Top Products",
    "sql": "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category",
    "tags": ["sales", "top"]
  }'
```

## Generate report

Requires a configured LLM. See [LLM setup](../getting-started/llm-setup.md) and [Configuration – LLM](../configuration.md#llm).

```bash
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category"
  }'
```

Response includes `narrative`, `metrics` (aggregates, top_categories, time_series, data_quality, perf_suggestions, and when applicable `period_current_label` / `period_previous_label`).

### Single-period query (no “previous period” in narrative)

When the result has no time-series comparison (e.g. one aggregate over “last 30 days”), the narrative will not mention “previous period” or “same period last year”:

```bash
curl -s -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT SUM(total_amount) AS total_fares, COUNT(*) AS trips FROM demo.sales"}' \
  | jq '.narrative'
```

### Get report by ID

```bash
curl -s http://localhost:8080/api/v1/reports/REPORT_UUID | jq .
```

### List reports

```bash
curl -s "http://localhost:8080/api/v1/reports?limit=5&offset=0" | jq '.items[] | {id, created_at}'
```

**See also:** [API reference](README.md), [Configuration](../configuration.md), [Period comparison](../features/period-comparison.md)
