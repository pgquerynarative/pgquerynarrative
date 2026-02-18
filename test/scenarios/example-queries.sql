-- Example SQL Queries for Testing PgQueryNarrative
-- These queries work with the demo.sales table

-- 1. Basic Query: List all product categories
SELECT DISTINCT product_category 
FROM demo.sales 
ORDER BY product_category;

-- 2. Aggregation: Total sales by category
SELECT 
    product_category,
    SUM(total_amount) as total_sales,
    COUNT(*) as transaction_count,
    AVG(total_amount) as avg_transaction
FROM demo.sales
GROUP BY product_category
ORDER BY total_sales DESC;

-- 3. Time Series: Sales by month
SELECT 
    DATE_TRUNC('month', sale_date) as month,
    SUM(total_amount) as monthly_total,
    COUNT(*) as transaction_count
FROM demo.sales
GROUP BY DATE_TRUNC('month', sale_date)
ORDER BY month DESC
LIMIT 12;

-- 4. Top Products: Best performing products
SELECT 
    product_name,
    product_category,
    SUM(total_amount) as total_revenue,
    COUNT(*) as sales_count
FROM demo.sales
GROUP BY product_name, product_category
ORDER BY total_revenue DESC
LIMIT 10;

-- 5. Customer Analysis: Top customers
SELECT 
    customer_id,
    COUNT(*) as purchase_count,
    SUM(total_amount) as total_spent,
    AVG(total_amount) as avg_purchase
FROM demo.sales
GROUP BY customer_id
ORDER BY total_spent DESC
LIMIT 10;

-- 6. Category Comparison: Compare categories
SELECT 
    product_category,
    SUM(total_amount) as total,
    COUNT(*) as count,
    MIN(total_amount) as min_sale,
    MAX(total_amount) as max_sale,
    AVG(total_amount) as avg_sale
FROM demo.sales
GROUP BY product_category
ORDER BY total DESC;

-- 7. Recent Sales: Last 30 days
SELECT 
    sale_date,
    product_category,
    SUM(total_amount) as daily_total
FROM demo.sales
WHERE sale_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY sale_date, product_category
ORDER BY sale_date DESC, daily_total DESC;

-- 8. Sales Trends: Month-over-month growth
SELECT 
    DATE_TRUNC('month', sale_date) as month,
    SUM(total_amount) as monthly_total,
    LAG(SUM(total_amount)) OVER (ORDER BY DATE_TRUNC('month', sale_date)) as previous_month,
    SUM(total_amount) - LAG(SUM(total_amount)) OVER (ORDER BY DATE_TRUNC('month', sale_date)) as change
FROM demo.sales
GROUP BY DATE_TRUNC('month', sale_date)
ORDER BY month DESC
LIMIT 12;
