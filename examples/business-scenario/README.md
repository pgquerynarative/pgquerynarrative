# Business Scenario: Sales Performance Analysis

## Overview

This example demonstrates a real-world use case of PgQueryNarrative: **A CEO needs comprehensive sales performance insights for a board meeting**.

The scenario shows how PgQueryNarrative transforms raw SQL query results into actionable business insights automatically.

---

## Scenario Description

**Company:** TechRetail Inc.  
**Situation:** CEO needs Q4 sales performance review for board meeting tomorrow  
**Challenge:** Data is scattered across multiple queries, needs to be synthesized into clear insights  
**Solution:** Use PgQueryNarrative to generate evidence-based narratives from sales data

### Business Questions Answered

1. **Overall Performance:** How do sales compare across categories?
2. **Top Products:** What are the best-selling products?
3. **Regional Analysis:** Which regions are performing best?
4. **Sales Team:** Which sales reps are top performers?
5. **Category Comparison:** How do categories compare in detail?

---

## Quick Start

### Prerequisites

1. **Server Running:** PgQueryNarrative server must be running
   ```bash
   # If not running, start it:
   make start-docker   # or make start-local
   # or
   go run ./cmd/server
   ```

2. **Database Ready:** Ensure database has demo data
   ```bash
   # If needed, seed data:
   export PGQUERYNARRATIVE_SEED=true
   make seed
   ```

### Run the Scenario

**Option 1: Automated Script (Recommended)**
```bash
cd examples/business-scenario
bash run-business-scenario.sh
```

**Option 2: Manual Execution**
Follow the step-by-step guide below.

---

## Step-by-Step User Guide

### Step 1: Verify Server is Running

```bash
# Check if server is responding
curl http://localhost:8080/

# Should return HTML page
```

### Step 2: Run Analysis Queries

#### Analysis 1: Category Performance

**Via Web UI:**
1. Open http://localhost:8080/query
2. Paste this query:
   ```sql
   SELECT product_category, SUM(total_amount) as total_revenue, 
          COUNT(*) as transaction_count, AVG(total_amount) as avg_sale 
   FROM demo.sales 
   GROUP BY product_category 
   ORDER BY total_revenue DESC
   ```
3. Click "Execute Query"
4. Review results showing revenue by category

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) as total_revenue, COUNT(*) as transaction_count, AVG(total_amount) as avg_sale FROM demo.sales GROUP BY product_category ORDER BY total_revenue DESC",
    "limit": 10
  }'
```

#### Analysis 2: Top Products

**Query:**
```sql
SELECT product_name, product_category, SUM(total_amount) as revenue, 
       SUM(quantity) as units_sold 
FROM demo.sales 
GROUP BY product_name, product_category 
ORDER BY revenue DESC 
LIMIT 5
```

**What to Look For:**
- Which products generate the most revenue?
- Which products sell the most units?
- Are there patterns across categories?

#### Analysis 3: Regional Performance

**Query:**
```sql
SELECT region, SUM(total_amount) as revenue, COUNT(*) as transactions 
FROM demo.sales 
GROUP BY region 
ORDER BY revenue DESC
```

**What to Look For:**
- Which regions are top performers?
- Which regions need attention?
- Regional distribution of sales

#### Analysis 4: Sales Team Performance

**Query:**
```sql
SELECT sales_rep, COUNT(*) as sales_count, SUM(total_amount) as revenue, 
       AVG(total_amount) as avg_deal 
FROM demo.sales 
GROUP BY sales_rep 
ORDER BY revenue DESC 
LIMIT 5
```

**What to Look For:**
- Top performing sales reps
- Average deal sizes
- Sales volume vs revenue correlation

#### Analysis 5: Category Comparison

**Query:**
```sql
SELECT product_category, SUM(total_amount) as total, 
       MIN(total_amount) as min_sale, MAX(total_amount) as max_sale, 
       AVG(total_amount) as avg_sale 
FROM demo.sales 
GROUP BY product_category 
ORDER BY total DESC
```

**What to Look For:**
- Revenue ranges per category
- Average transaction values
- Category diversity

### Step 3: Generate Narrative Reports

For each query, you can generate an AI-powered narrative:

**Via Web UI:**
1. Go to http://localhost:8080/reports
2. Paste your query
3. Click "Generate Narrative Report"
4. Review the generated insights

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) as total_revenue FROM demo.sales GROUP BY product_category ORDER BY total_revenue DESC"
  }'
```

**Note:** Narrative generation requires an LLM (Ollama, Groq, or Gemini). If not configured, you'll get an error but can still view query results.

### Step 4: Save Important Queries

Save queries for future use:

**Via Web UI:**
1. Go to http://localhost:8080/query
2. After running a query, you can save it (feature in development)

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/queries/saved \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Category Performance Analysis",
    "sql": "SELECT product_category, SUM(total_amount) as total_revenue FROM demo.sales GROUP BY product_category ORDER BY total_revenue DESC",
    "tags": ["sales", "category", "performance"]
  }'
```

---

## Expected Results

After running all analyses, you should see:

### Category Performance
- Furniture, Office Supplies, and Clothing leading with ~$22M each
- Accessories and Electronics at ~$11M each

### Top Products
- Beta and Delta products performing well across categories
- Revenue ranging from $4.4M to $4.6M for top products

### Regional Performance
- West, East, and South regions performing strongly (~$22M each)
- North and Central regions at ~$11M (growth opportunity)

### Sales Team
- Top 3 reps generating ~$22M each
- Consistent average deal sizes (~$2,800)

### Category Comparison
- Consistent pricing across categories
- Wide product range indicating diverse portfolio

---

## Business Insights

Based on the analysis results:

1. **Regional Expansion:** Focus on North and Central regions (2x growth potential)
2. **Product Strategy:** Leverage Beta and Delta brand strength
3. **Sales Training:** Document best practices from top 3 performers
4. **Inventory Management:** Maintain strong stock for top 5 products

---

## Troubleshooting

### Server Not Responding
```bash
# Check if server is running
curl http://localhost:8080/

# If not, start it:
make start-docker   # or make start-local
```

### Database Connection Error
```bash
# Check database status
docker compose ps postgres

# If not running:
docker compose up -d postgres

# Wait for it to be ready:
docker compose exec postgres pg_isready -U postgres
```

### No Data in Results
```bash
# Seed demo data:
export PGQUERYNARRATIVE_SEED=true
make seed
```

### LLM Not Available
- Narrative generation requires Ollama, Groq, or Gemini
- Query execution works without LLM
- Install Ollama: https://ollama.ai
- Or configure Groq/Gemini API keys

---

## Files in This Example

- **README.md** - This file (user guide)
- **business-scenario.md** - Detailed scenario description
- **run-business-scenario.sh** - Automated execution script
- **SCENARIO_RESULTS.md** - Expected results and insights

---

## Learning Outcomes

After completing this scenario, you'll understand:

1. How to execute SQL queries safely
2. How to analyze multi-dimensional data
3. How to generate business insights automatically
4. How to use PgQueryNarrative for real business problems
5. How to interpret query results for decision-making

---

## Next Steps

1. **Modify Queries:** Try different aggregations and filters
2. **Create Custom Reports:** Generate narratives for your own queries
3. **Save Queries:** Build a library of useful analyses
4. **Explore Web UI:** Use the interactive interface at http://localhost:8080

---

## 💡 Tips

- **Start Simple:** Begin with basic queries, then add complexity
- **Use Web UI:** Easier to experiment and see results
- **Save Queries:** Build a reusable query library
- **Generate Narratives:** Use AI to get insights you might miss
- **Combine Analyses:** Run multiple queries for comprehensive view

---

## 📞 Need Help?

- Check main README: `/README.md`
- Review API docs: `/docs/api/README.md`
- Test queries: `/test/queries/quick-test-queries.md`

---

