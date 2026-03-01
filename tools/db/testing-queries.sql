-- Testing queries for PgQueryNarrative (demo schema only)
-- Use in Query Runner or Generate Report. Limit is applied by the app (default 100).

-- 1. Sales by category (good for bar/pie charts)
SELECT product_category, SUM(total_amount) AS total
FROM demo.sales
GROUP BY product_category
ORDER BY total DESC;

-- 2. Daily sales over time (good for line chart)
SELECT date, SUM(total_amount) AS daily_total
FROM demo.sales
GROUP BY date
ORDER BY date;

-- 3. Top categories by region
SELECT region, product_category, SUM(quantity) AS qty
FROM demo.sales
GROUP BY region, product_category
ORDER BY region, qty DESC;

-- 4. Sales rep performance
SELECT sales_rep, COUNT(*) AS orders, SUM(total_amount) AS revenue
FROM demo.sales
GROUP BY sales_rep
ORDER BY revenue DESC;

-- 5. Category totals with average order value
SELECT product_category,
       SUM(total_amount) AS total,
       COUNT(*) AS orders,
       ROUND(AVG(total_amount), 2) AS avg_order
FROM demo.sales
GROUP BY product_category
ORDER BY total DESC;

-- 6. Regional revenue
SELECT region, SUM(total_amount) AS revenue
FROM demo.sales
GROUP BY region
ORDER BY revenue DESC;

-- 7. Simple row sample (for table view)
SELECT date, product_category, product_name, quantity, total_amount, region
FROM demo.sales
ORDER BY date DESC;

-- 8. Monthly rollup (time series for reports)
SELECT date_trunc('month', date)::date AS month,
       SUM(total_amount) AS monthly_total,
       SUM(quantity) AS units_sold
FROM demo.sales
GROUP BY date_trunc('month', date)
ORDER BY month;
