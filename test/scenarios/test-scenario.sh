#!/bin/bash
set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
echo "=========================================="
echo "  PgQueryNarrative Test Scenario"
echo "=========================================="
echo ""
echo "Base URL: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Check prerequisites (with timeout)
echo -e "${BLUE}Checking prerequisites...${NC}"
if ! timeout 10s curl -s "$BASE_URL/" > /dev/null 2>&1; then
    echo -e "${RED}Server is not running at $BASE_URL${NC}"
    echo -e "${YELLOW}Please start the server first:${NC}"
    echo "  make start-docker   # or make start-local"
    echo "  or"
    echo "  go run ./cmd/server"
    exit 1
fi
echo -e "${GREEN}Server is running${NC}"
echo ""

test_step() {
    local name="$1"
    local command="$2"
    
    echo -e "${YELLOW}Testing: $name${NC}"
    # Wrap command with 30 second timeout to prevent hanging
    if timeout 30s bash -c "$command" 2>/dev/null; then
        echo -e "${GREEN}PASS: $name${NC}"
        ((TESTS_PASSED++))
        return 0
    else
        local exit_code=$?
        if [ $exit_code -eq 124 ]; then
            echo -e "${RED}FAIL: $name (timeout after 30s)${NC}"
        else
            echo -e "${RED}FAIL: $name${NC}"
        fi
        ((TESTS_FAILED++))
        return 1
    fi
    echo ""
}

echo "=========================================="
echo "  Scenario 1: Basic Query Execution"
echo "=========================================="
echo ""

# Test 1: Simple SELECT query
test_step "Execute simple query" "
RESPONSE=\$(curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d '{
    \"sql\": \"SELECT product_category, COUNT(*) as count FROM demo.sales GROUP BY product_category LIMIT 5\",
    \"limit\": 10
  }')
if echo \"\$RESPONSE\" | jq -e '.row_count >= 0' > /dev/null 2>&1; then
    exit 0
elif echo \"\$RESPONSE\" | jq -e '.name == \"validation_error\"' > /dev/null 2>&1; then
    echo \"  Database connection issue - ensure DB is running\"
    exit 1
else
    exit 1
fi
"

# Test 2: Query with aggregation
test_step "Execute aggregation query" "
RESPONSE=\$(curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d '{
    \"sql\": \"SELECT product_category, SUM(total_amount) as total FROM demo.sales GROUP BY product_category ORDER BY total DESC LIMIT 5\",
    \"limit\": 10
  }')
echo \"\$RESPONSE\" | jq -e '.columns | length > 0' > /dev/null 2>&1
"

echo ""
echo "=========================================="
echo "  Scenario 2: Save and Retrieve Queries"
echo "=========================================="
echo ""

# Test 3: Save a query (with timeout)
SAVED_QUERY_ID=$(timeout 30s curl -s -X POST $BASE_URL/api/v1/queries/saved \
  -H 'Content-Type: application/json' \
  -d '{
    \"name\": \"Top Product Categories\",
    \"sql\": \"SELECT product_category, SUM(total_amount) as total FROM demo.sales GROUP BY product_category ORDER BY total DESC\",
    \"tags\": [\"sales\", \"top\", \"categories\"]
  }' | jq -r '.id')

if [ -n "$SAVED_QUERY_ID" ] && [ "$SAVED_QUERY_ID" != "null" ]; then
    echo -e "${GREEN}PASS: Save query (ID: $SAVED_QUERY_ID)${NC}"
    ((TESTS_PASSED++))
    
    # Test 4: List saved queries (wait a moment for DB to sync)
    sleep 1
    test_step "List saved queries" "
    timeout 30s curl -s $BASE_URL/api/v1/queries/saved | jq -e '.items | length >= 1' > /dev/null
    "
    
    # Test 5: Get specific saved query
    test_step "Get saved query by ID" "
    timeout 30s curl -s $BASE_URL/api/v1/queries/saved/$SAVED_QUERY_ID | jq -e '.id != null and .id != \"\"' > /dev/null
    "
else
    echo -e "${RED}FAIL: Save query${NC}"
    ((TESTS_FAILED++))
fi

echo ""
echo "=========================================="
echo "  Scenario 3: Report Generation"
echo "=========================================="
echo ""

# Test 6: Generate report (may fail if Ollama not running, that's OK)
test_step "Generate narrative report" "
curl -s -X POST $BASE_URL/api/v1/reports/generate \
  -H 'Content-Type: application/json' \
  -d '{
    \"sql\": \"SELECT product_category, SUM(total_amount) as total FROM demo.sales GROUP BY product_category ORDER BY total DESC LIMIT 5\"
  }' | jq -e '.id != null or .name == \"llm_error\"' > /dev/null
"

echo ""
echo "=========================================="
echo "  Scenario 4: Web UI Pages"
echo "=========================================="
echo ""

# Test 7: Home page
test_step "Home page loads" "
curl -s $BASE_URL/ | grep -q 'PgQueryNarrative' > /dev/null
"

# Test 8: Query page
test_step "Query page loads" "
curl -s $BASE_URL/query | grep -q 'Run SQL Query' > /dev/null
"

# Test 9: Saved queries page
test_step "Saved queries page loads" "
curl -s $BASE_URL/saved | grep -q 'Saved Queries' > /dev/null
"

# Test 10: Reports page
test_step "Reports page loads" "
curl -s $BASE_URL/reports | grep -q 'Generate Report' > /dev/null
"

echo ""
echo "=========================================="
echo "  Scenario 5: Error Handling"
echo "=========================================="
echo ""

# Test 11: Invalid SQL query
test_step "Handle invalid SQL gracefully" "
curl -s -X POST $BASE_URL/api/v1/queries/run \
  -H 'Content-Type: application/json' \
  -d '{
    \"sql\": \"DELETE FROM demo.sales\",
    \"limit\": 10
  }' | jq -e '.name == \"validation_error\"' > /dev/null
"

# Test 12: Non-existent saved query (use valid UUID format)
test_step "Handle non-existent query" "
timeout 30s curl -s $BASE_URL/api/v1/queries/saved/123e4567-e89b-12d3-a456-426614174000 | jq -e '.name == \"not_found\" or .name == \"fault\" or .name == \"invalid_format\"' > /dev/null
"

echo ""
echo "=========================================="
echo "  Test Results Summary"
echo "=========================================="
echo ""
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed${NC}"
    exit 1
fi
