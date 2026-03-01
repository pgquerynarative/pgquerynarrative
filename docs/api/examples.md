# API examples

Base URL: `http://localhost:8080` (or set `PGQUERYNARRATIVE_PORT` via [Configuration](../configuration.md)). All examples use `Content-Type: application/json`. Full endpoint list: [API reference](README.md).

## Run query

```bash
curl -X POST http://localhost:8080/api/v1/queries/run -H "Content-Type: application/json" -d '{"sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category", "limit": 10}'
```

Response: `columns`, `rows`, `row_count`, `execution_time_ms`, optional `chart_suggestions`, `period_comparison` (when result has a time column and measures).

## Save query

```bash
curl -X POST http://localhost:8080/api/v1/queries/saved -H "Content-Type: application/json" -d '{"name": "Top Products", "sql": "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category", "tags": ["sales", "top"]}'
```

## Generate report

Requires a configured [LLM](../getting-started/llm-setup.md). See [Configuration – LLM](../configuration.md#llm).

```bash
curl -X POST http://localhost:8080/api/v1/reports/generate -H "Content-Type: application/json" -d '{"sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category"}'
```

Response: `narrative`, `metrics` (aggregates, time_series, data_quality, perf_suggestions, etc.).

## Get / list reports

```bash
curl -s "http://localhost:8080/api/v1/reports?limit=5&offset=0" | jq '.items[] | {id, created_at}'
```

## See also

- [API reference](README.md) · [LLM setup](../getting-started/llm-setup.md) · [Configuration](../configuration.md) · [Documentation index](../README.md)
