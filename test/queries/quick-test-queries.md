# Quick Test Queries

Use these queries to test PgQueryNarrative functionality.

## Via Web UI

1. Open http://localhost:8080/query
2. Paste any query below
3. Click "Execute Query"

## Via API

```bash
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "YOUR_QUERY_HERE",
    "limit": 100
  }'
```

---

## Test Queries

### 1. Basic Category List
```sql
SELECT DISTINCT product_category 
FROM demo.sales 
ORDER BY product_category;
```

### 2. Sales Summary by Category
```sql
SELECT 
    product_category,
    SUM(total_amount) as total_sales,
    COUNT(*) as transaction_count,
    AVG(total_amount) as avg_transaction
FROM demo.sales
GROUP BY product_category
ORDER BY total_sales DESC;
```

### 3. Top 10 Products
```sql
SELECT 
    product_name,
    product_category,
    SUM(total_amount) as revenue,
    COUNT(*) as sales_count
FROM demo.sales
GROUP BY product_name, product_category
ORDER BY revenue DESC
LIMIT 10;
```

### 4. Monthly Sales Trend
```sql
SELECT 
    DATE_TRUNC('month', date)::date AS month,
    SUM(total_amount) AS monthly_total,
    COUNT(*) AS transaction_count
FROM demo.sales
GROUP BY DATE_TRUNC('month', date)
ORDER BY month
LIMIT 12;
```
**Use this to test "Vs previous period":** Run it in the UI or via API; you should see a **Vs previous period** block comparing the last month to the previous month (e.g. `monthly_total` and `transaction_count` with % change and trend).

### 5. Category Performance Analysis
```sql
SELECT 
    product_category,
    SUM(total_amount) as total,
    COUNT(*) as count,
    AVG(total_amount) as avg_sale,
    MIN(total_amount) as min_sale,
    MAX(total_amount) as max_sale
FROM demo.sales
GROUP BY product_category
ORDER BY total DESC;
```

### 6. Regional Sales Breakdown
```sql
SELECT 
    region,
    product_category,
    SUM(total_amount) as regional_sales,
    COUNT(*) as sales_count
FROM demo.sales
GROUP BY region, product_category
ORDER BY regional_sales DESC
LIMIT 20;
```

### 7. Sales Rep Performance
```sql
SELECT 
    sales_rep,
    COUNT(*) as sales_count,
    SUM(total_amount) as total_revenue,
    AVG(total_amount) as avg_sale
FROM demo.sales
GROUP BY sales_rep
ORDER BY total_revenue DESC
LIMIT 10;
```

### 8. Recent Sales (Last 30 Days)
```sql
SELECT 
    sale_date,
    product_category,
    SUM(total_amount) as daily_total
FROM demo.sales
WHERE sale_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY sale_date, product_category
ORDER BY sale_date DESC, daily_total DESC;
```

### 9. High-Value Transactions
```sql
SELECT 
    sale_date,
    product_name,
    product_category,
    total_amount,
    region
FROM demo.sales
WHERE total_amount > 1000
ORDER BY total_amount DESC
LIMIT 20;
```

### 10. Category Comparison
```sql
SELECT 
    product_category,
    SUM(total_amount) as total_revenue,
    COUNT(*) as transaction_count,
    SUM(quantity) as total_quantity,
    AVG(unit_price) as avg_unit_price
FROM demo.sales
GROUP BY product_category
ORDER BY total_revenue DESC;
```

---

## Test Report Generation

Use any of the above queries with the Reports endpoint:

```bash
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) as total FROM demo.sales GROUP BY product_category ORDER BY total DESC"
  }'
```

Or via Web UI: http://localhost:8080/reports

---

## Save a Query

```bash
curl -X POST http://localhost:8080/api/v1/queries/saved \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Top Categories",
    "sql": "SELECT product_category, SUM(total_amount) as total FROM demo.sales GROUP BY product_category ORDER BY total DESC",
    "tags": ["sales", "top", "categories"]
  }'
```

---

## Test Error Handling

Try these to test validation:

```sql
-- This should fail (DELETE not allowed)
DELETE FROM demo.sales;

-- This should fail (UPDATE not allowed)
UPDATE demo.sales SET total_amount = 0;

-- This should fail (INSERT not allowed)
INSERT INTO demo.sales VALUES (...);
```

---

## Notes

- All queries use the `demo.sales` table
- Queries are read-only (SELECT only)
- Results are limited to prevent large result sets
- Use the Web UI for easier testing: http://localhost:8080/query
