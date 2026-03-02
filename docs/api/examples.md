# API examples

Base URL: `http://localhost:8080` (or set port via [Configuration](../configuration.md)). Full endpoint list: [API reference](README.md).

**Example — run query:**

```bash
curl -s -X POST http://localhost:8080/api/v1/queries/run -H "Content-Type: application/json" -d '{"sql":"SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category", "limit":10}'
```

Response: `columns`, `rows`, `row_count`, `execution_time_ms`, optional `chart_suggestions`, `period_comparison` (when result has a time column and measures).

**See also:** [API reference](README.md) · [Configuration](../configuration.md) · [Documentation index](../README.md)
