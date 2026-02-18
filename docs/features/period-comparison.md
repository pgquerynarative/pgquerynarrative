# Period comparison

When query results have a **date/time column** and at least one **numeric measure**, the app automatically computes **period-over-period** comparison (e.g. latest period vs previous) and shows derived fields: current value, previous value, absolute change, **% change**, and trend (up / down / flat).

## Where it appears

- **Run Query (UI):** A “Vs previous period” block above the results table, with period labels (e.g. `2026-01-01 vs 2025-12-01`) when available.
- **Reports (UI):** Same “Vs previous period” section in the report view, with period labels in the heading when the report was generated from a time-series query.
- **API:** Run-query response includes `period_comparison` (array of measure, current, previous, change, change_percentage, trend) and optional `period_current_label` / `period_previous_label`. Report response includes `metrics.time_series` and `metrics.period_current_label` / `metrics.period_previous_label` when applicable.

## Requirements

- At least one column detected as **date/time** (e.g. `DATE_TRUNC('month', date)::date`).
- At least one **numeric measure** (e.g. `SUM(total_amount)`, `COUNT(*)`).
- At least **two rows** (two periods); comparison uses the last row vs the second-to-last.

## Example query

```sql
SELECT DATE_TRUNC('month', date)::date AS month,
       SUM(total_amount) AS monthly_total,
       COUNT(*) AS transaction_count
FROM demo.sales
GROUP BY DATE_TRUNC('month', date)
ORDER BY month
LIMIT 12
```

## Testing in the UI

1. **Query results:** Open [Run Query](http://localhost:8080/query), paste the query above, run it. Confirm the “Vs previous period” block and period labels.
2. **Reports:** Open [Reports](http://localhost:8080/reports), generate a report with the same SQL. Confirm “Vs previous period” in the report and that **Key takeaways** can mention change vs previous period (narrative rule).

## Configurable trend threshold

The threshold for labeling a change as “up” or “down” (vs “flat”) is configurable. See [Configuration – Metrics and period comparison](../configuration.md#metrics-and-period-comparison). Example:

```bash
# Only changes ≥ 1% are up/down
export PERIOD_TREND_THRESHOLD_PERCENT=1
make start-local
```

**See also:** [Configuration](../configuration.md), [API reference](../api/README.md), [API examples](../api/examples.md)
