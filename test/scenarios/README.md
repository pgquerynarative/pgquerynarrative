# Test Scenarios

This directory contains test scenarios and examples for PgQueryNarrative.

## Quick Start

### Prerequisites

1. **Server running**: The application server must be running
   ```bash
   make start-docker   # or make start-local
   # or
   go run ./cmd/server
   ```

2. **Database accessible**: Ensure PostgreSQL is running and accessible
   - If using Docker: `docker compose up -d postgres`
   - If using local: Ensure PostgreSQL is running on port 5432

3. **Dependencies**: `jq` for JSON parsing (install with `brew install jq` on macOS)

### Run Automated Tests

```bash
# Simple test run
bash test/scenarios/test-scenario.sh

# With setup check
bash test/scenarios/setup-and-test.sh
```

## Test Scenarios

### 1. Automated Test Script (`test-scenario.sh`)

Comprehensive automated test suite covering:

- **Basic Query Execution**: Simple SELECT and aggregation queries
- **Save and Retrieve Queries**: Save queries and retrieve them
- **Report Generation**: Generate narrative reports (requires LLM)
- **Web UI Pages**: Verify all web pages load correctly
- **Error Handling**: Test invalid SQL and non-existent resources

**Expected Output:**
```
==========================================
  PgQueryNarrative Test Scenario
==========================================

PASS: Execute simple query
PASS: Execute aggregation query
PASS: Save query
PASS: List saved queries
...

Tests Passed: 12
Tests Failed: 0
All tests passed!
```

### 2. Example Queries (`example-queries.sql`)

8 ready-to-use SQL queries for manual testing:

1. **Basic Query**: List all product categories
2. **Aggregation**: Total sales by category
3. **Time Series**: Sales by month
4. **Top Products**: Best performing products
5. **Customer Analysis**: Top customers
6. **Category Comparison**: Compare categories
7. **Recent Sales**: Last 30 days
8. **Sales Trends**: Month-over-month growth

**Usage:**
```bash
# Copy a query and paste into the web UI at http://localhost:8080/query
# Or use with the API:
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT ... FROM demo.sales ..."}'
```

### 3. Manual Test Guide (`manual-test-guide.md`)

Step-by-step instructions for manual testing:

- **Scenario 1**: Basic Query Execution (Web UI & API)
- **Scenario 2**: Save and Reuse Query
- **Scenario 3**: Generate Narrative Report
- **Scenario 4**: Error Handling

## Example Test Run

```bash
# 1. Start the server (in one terminal)
make start-docker   # or make start-local

# 2. Run tests (in another terminal, from project root)
cd /path/to/PgQueryNarrative   # or: cd "$(git rev-parse --show-toplevel)"
bash test/scenarios/test-scenario.sh
```

## Troubleshooting

### "Server is not running"
- Start the server: `make start-docker` or `make start-local` or `go run ./cmd/server`
- Verify: `curl http://localhost:8080/`

### "Database connection refused"
- Start database: `docker compose up -d postgres`
- Wait for it to be ready: `docker compose exec postgres pg_isready`
- Run migrations: Check that `make start-docker` or `make start-local` completed successfully

### "jq: command not found"
- Install jq: `brew install jq` (macOS) or `apt-get install jq` (Linux)

### Tests fail with "validation_error"
- Database may not be accessible
- Check database connection: `psql $DB_URL`
- Ensure demo data is seeded: `export PGQUERYNARRATIVE_SEED=true && make seed`

## Customization

### Change Base URL

```bash
BASE_URL=http://localhost:9090 bash test/scenarios/test-scenario.sh
```

### Skip Specific Tests

Edit `test-scenario.sh` and comment out specific test cases.

## Next Steps

After running tests:

1. **Review Results**: Check which tests passed/failed
2. **Manual Testing**: Use `example-queries.sql` for interactive testing
3. **Web UI**: Open http://localhost:8080 and explore the interface
4. **API Testing**: Use curl or Postman with the API endpoints
