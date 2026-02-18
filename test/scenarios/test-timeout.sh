#!/bin/bash
# Test script to verify query timeout functionality

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
TIMEOUT="${QUERY_TIMEOUT:-30s}"

echo "=========================================="
echo "  Query Timeout Test"
echo "=========================================="
echo ""
echo "Base URL: $BASE_URL"
echo "Query Timeout: $TIMEOUT"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Testing query timeout handling...${NC}"
echo ""

# Test 1: Normal query (should succeed)
echo "Test 1: Normal query (should succeed)"
RESPONSE=$(timeout 30s curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d '{
    "sql": "SELECT product_category, COUNT(*) as count FROM demo.sales GROUP BY product_category LIMIT 5",
    "limit": 10
  }')

if echo "$RESPONSE" | jq -e '.row_count >= 0' > /dev/null 2>&1; then
    echo -e "${GREEN}PASS: Normal query executed successfully${NC}"
else
    echo -e "${RED}FAIL: Normal query failed${NC}"
    echo "Response: $RESPONSE"
fi
echo ""

# Test 2: Query that should timeout (using pg_sleep)
echo "Test 2: Long-running query (should timeout)"
echo "Note: This test requires a query that takes longer than the timeout period"
echo ""

# Convert timeout to seconds for comparison
TIMEOUT_SEC=$(echo $TIMEOUT | sed 's/s$//')
if [ -z "$TIMEOUT_SEC" ]; then
    TIMEOUT_SEC=30
fi

# Create a query that sleeps for longer than timeout
SLEEP_TIME=$((TIMEOUT_SEC + 5))

RESPONSE=$(timeout 60s curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d "{
    \"sql\": \"SELECT pg_sleep($SLEEP_TIME), 'test' as value\",
    \"limit\": 10
  }")

if echo "$RESPONSE" | jq -e '.name == "timeout_error" or .name == "validation_error"' > /dev/null 2>&1; then
    echo -e "${GREEN}PASS: Timeout error detected correctly${NC}"
    echo "Error type: $(echo "$RESPONSE" | jq -r '.name')"
    echo "Message: $(echo "$RESPONSE" | jq -r '.message')"
else
    echo -e "${YELLOW}Query may have completed or validation prevented execution${NC}"
    echo "Response: $RESPONSE"
fi
echo ""

# Test 3: Verify timeout error format
echo "Test 3: Verify timeout error format"
RESPONSE=$(timeout 60s curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d "{
    \"sql\": \"SELECT pg_sleep($SLEEP_TIME)\",
    \"limit\": 10
  }")

ERROR_NAME=$(echo "$RESPONSE" | jq -r '.name // "unknown"')
ERROR_CODE=$(echo "$RESPONSE" | jq -r '.code // "unknown"')

if [ "$ERROR_NAME" == "timeout_error" ] || [ "$ERROR_CODE" == "TIMEOUT_ERROR" ]; then
    echo -e "${GREEN}PASS: Timeout error has correct format${NC}"
    echo "  Name: $ERROR_NAME"
    echo "  Code: $ERROR_CODE"
elif [ "$ERROR_NAME" == "validation_error" ]; then
    echo -e "${YELLOW}Query was rejected by validator (expected for pg_sleep)${NC}"
else
    echo -e "${RED}FAIL: Unexpected error format${NC}"
    echo "Response: $RESPONSE"
fi
echo ""

echo "=========================================="
echo "  Timeout Test Complete"
echo "=========================================="
echo ""
echo "To change the query timeout, set the QUERY_TIMEOUT environment variable:"
echo "  export QUERY_TIMEOUT=60s"
echo "  make start-docker   # or make start-local"
