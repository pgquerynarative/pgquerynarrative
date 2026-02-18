#!/bin/bash
# Comprehensive Business Scenario Demo
# Q4 Sales Performance Review for TechRetail Inc.

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  COMPREHENSIVE BUSINESS SCENARIO: Q4 SALES PERFORMANCE REVIEW"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Company: TechRetail Inc."
echo "Situation: CEO needs Q4 performance review for board meeting"
echo "Challenge: Synthesize multiple data points into clear business insights"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}ANALYSIS 1: Overall Q4 vs Q3 Performance Comparison${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY1='SELECT CASE WHEN EXTRACT(QUARTER FROM sale_date) = 3 THEN '\''Q3'\'' WHEN EXTRACT(QUARTER FROM sale_date) = 4 THEN '\''Q4'\'' END as quarter, SUM(total_amount) as total_revenue, COUNT(*) as transaction_count, AVG(total_amount) as avg_transaction_value FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) IN (3, 4) GROUP BY quarter ORDER BY quarter'
RESULT1=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY1\", \"limit\": 10}")
echo "$RESULT1" | jq -r '
  "Quarter | Revenue      | Transactions | Avg Transaction
  --------|--------------|--------------|----------------" ,
  (.rows[] | "\(.[0] // "N/A") | $\(.[1] // 0) | \(.[2] // 0) | $\(.[3] // 0)")
' 2>/dev/null || echo "$RESULT1" | jq '.row_count' && echo " rows returned"
echo ""

echo -e "${BLUE}ANALYSIS 2: Category Performance (Q4)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY2='SELECT product_category, SUM(total_amount) as q4_revenue, COUNT(*) as q4_transactions, AVG(total_amount) as avg_sale FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) = 4 GROUP BY product_category ORDER BY q4_revenue DESC'
RESULT2=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY2\", \"limit\": 10}")
echo "$RESULT2" | jq -r '
  "Category        | Revenue      | Transactions | Avg Sale
  ----------------|--------------|--------------|----------" ,
  (.rows[] | "\(.[0] // "N/A") | $\(.[1] // 0) | \(.[2] // 0) | $\(.[3] // 0)")
' 2>/dev/null || echo "$RESULT2" | jq '.row_count' && echo " categories analyzed"
echo ""

echo -e "${BLUE}ANALYSIS 3: Top 10 Products (Q4)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY3='SELECT product_name, product_category, SUM(total_amount) as total_revenue, SUM(quantity) as total_units_sold FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) = 4 GROUP BY product_name, product_category ORDER BY total_revenue DESC LIMIT 10'
RESULT3=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY3\", \"limit\": 10}")
echo "$RESULT3" | jq -r '
  "Product                    | Category      | Revenue      | Units
  ----------------------------|---------------|--------------|-------" ,
  (.rows[] | "\(.[0] // "N/A") | \(.[1] // "N/A") | $\(.[2] // 0) | \(.[3] // 0)")
' 2>/dev/null || echo "$RESULT3" | jq '.row_count' && echo " top products"
echo ""

echo -e "${BLUE}🌍 ANALYSIS 4: Regional Performance (Q4)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY4='SELECT region, SUM(total_amount) as regional_revenue, COUNT(*) as transaction_count FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) = 4 GROUP BY region ORDER BY regional_revenue DESC'
RESULT4=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY4\", \"limit\": 10}")
echo "$RESULT4" | jq -r '
  "Region  | Revenue      | Transactions
  --------|--------------|--------------" ,
  (.rows[] | "\(.[0] // "N/A") | $\(.[1] // 0) | \(.[2] // 0)")
' 2>/dev/null || echo "$RESULT4" | jq '.row_count' && echo " regions analyzed"
echo ""

echo -e "${BLUE}👥 ANALYSIS 5: Sales Team Performance (Q4)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY5='SELECT sales_rep, COUNT(*) as sales_count, SUM(total_amount) as total_revenue, AVG(total_amount) as avg_deal_size FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) = 4 GROUP BY sales_rep ORDER BY total_revenue DESC LIMIT 10'
RESULT5=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY5\", \"limit\": 10}")
echo "$RESULT5" | jq -r '
  "Sales Rep | Sales Count | Revenue      | Avg Deal
  ----------|-------------|--------------|----------" ,
  (.rows[] | "\(.[0] // "N/A") | \(.[1] // 0) | $\(.[2] // 0) | $\(.[3] // 0)")
' 2>/dev/null || echo "$RESULT5" | jq '.row_count' && echo " sales reps analyzed"
echo ""

echo -e "${BLUE}📅 ANALYSIS 6: Monthly Trend Analysis (Q4)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
QUERY6='SELECT DATE_TRUNC('\''month'\'', sale_date) as month, SUM(total_amount) as monthly_revenue, COUNT(*) as transaction_count FROM demo.sales WHERE EXTRACT(QUARTER FROM sale_date) = 4 GROUP BY DATE_TRUNC('\''month'\'', sale_date) ORDER BY month'
RESULT6=$(timeout 15s curl -s -X POST "$BASE_URL/api/v1/queries/run" -H "Content-Type: application/json" -d "{\"sql\": \"$QUERY6\", \"limit\": 10}")
echo "$RESULT6" | jq -r '
  "Month       | Revenue      | Transactions
  ------------|--------------|--------------" ,
  (.rows[] | "\(.[0] // "N/A") | $\(.[1] // 0) | \(.[2] // 0)")
' 2>/dev/null || echo "$RESULT6" | jq '.row_count' && echo " months analyzed"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}COMPREHENSIVE ANALYSIS COMPLETE${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Summary:"
echo "  • 6 comprehensive analyses completed"
echo "  • Q4 performance reviewed across multiple dimensions"
echo "  • Data ready for narrative generation"
echo ""
echo "Next Step: Generate AI narratives for each analysis"
echo "   Visit: http://localhost:8080/reports"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
