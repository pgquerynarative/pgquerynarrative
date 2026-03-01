# API examples

Base URL: `http://localhost:8080` (or set `PGQUERYNARRATIVE_PORT` via [Configuration](../configuration.md)). All examples use `Content-Type: application/json`. Full endpoint list: [API reference](README.md).

## Run query

```bash
curl -X POST http://localhost:8080/api/v1/queries/run -H "Content-Type: application/json" -d '{"sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category", "limit": 10}'
```

Response: `columns`, `rows`, `row_count`, `execution_time_ms`, optional `chart_suggestions`, `period_comparison` (when result has a time column and measures).

Use `GET /api/v1/schema` to see available schemas and tables (default: `public`, `demo`).

## Datasets

Import your own data with `psql` and `\copy`. Run `make migrate` so the read-only user can access `public` schema. Example: gold prices (CSV with Date, Close, High, Low, Open, Volume):

```sql
CREATE TABLE public.gold_prices (date DATE, close NUMERIC, high NUMERIC, low NUMERIC, open NUMERIC, volume BIGINT, adj_close NUMERIC, daily_return NUMERIC, ma_20 NUMERIC, ma_50 NUMERIC, ma_200 NUMERIC, volatility_20 NUMERIC, year INT, month INT, day_of_week INT, quarter INT);
\copy public.gold_prices FROM '/path/to/gold_prices_10y.csv' WITH (FORMAT csv, HEADER true, NULL '');
```

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

## Ask (NL → SQL → report)

One step: natural-language question → generated SQL → run → narrative report. Requires [LLM](../getting-started/llm-setup.md).

```bash
curl -X POST http://localhost:8080/api/v1/suggestions/ask -H "Content-Type: application/json" -d '{"question": "What were the top 5 products by revenue?"}'
```

Response: `question`, `sql`, `report` (same shape as generate).

## Explain SQL

Plain-English explanation of a SQL query. Requires [LLM](../getting-started/llm-setup.md).

```bash
curl -X POST http://localhost:8080/api/v1/suggestions/explain -H "Content-Type: application/json" -d '{"sql": "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"}'
```

Response: `sql`, `explanation`.

## Get / list reports

```bash
curl -s "http://localhost:8080/api/v1/reports?limit=5&offset=0" | jq '.items[] | {id, created_at}'
```

## See also

- [API reference](README.md) · [LLM setup](../getting-started/llm-setup.md) · [Configuration](../configuration.md) · [Documentation index](../README.md)
