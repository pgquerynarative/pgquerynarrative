#!/bin/bash
set -e

echo "=== Testing API Endpoints ==="
echo ""

BASE_URL="${BASE_URL:-http://localhost:8080}"

# Test query execution
echo "1. Testing query execution..."
curl -X POST "$BASE_URL/api/v1/queries/run" \
    -H "Content-Type: application/json" \
    -d '{
        "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category",
        "limit": 10
    }' | jq '.' || echo "Query test failed"

echo ""
echo "2. Testing report generation (requires Ollama)..."
curl -X POST "$BASE_URL/api/v1/reports/generate" \
    -H "Content-Type: application/json" \
    -d '{
        "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category LIMIT 5"
    }' | jq '.' || echo "Report generation test failed (Ollama may not be running)"

echo ""
echo "API tests completed"
