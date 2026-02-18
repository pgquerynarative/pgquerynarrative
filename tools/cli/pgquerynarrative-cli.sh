#!/bin/sh
# PgQueryNarrative CLI - Terminal interface for running queries and viewing results

set -e

API_URL="${PGQUERYNARRATIVE_API_URL:-http://app:8080}"
FORMAT="${PGQUERYNARRATIVE_FORMAT:-table}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo ""
    echo "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo "${CYAN}  $1${NC}"
    echo "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

print_success() {
    echo "${GREEN}✅ $1${NC}"
}

print_error() {
    echo "${RED}❌ $1${NC}" >&2
}

print_info() {
    echo "${BLUE}ℹ️  $1${NC}"
}

# Check if API is available
check_api() {
    if ! curl -s -f "${API_URL}/api/v1/queries/saved" >/dev/null 2>&1; then
        print_error "Cannot connect to API at ${API_URL}"
        print_info "Make sure the app container is running: docker compose up -d"
        exit 1
    fi
}

# Format JSON as table
format_table() {
    jq -r '
        if .rows then
            # Query results
            if (.rows | length) > 0 then
                (.rows[0] | keys | @tsv),
                (.rows[] | [.[]] | @tsv)
            else
                "No results"
            end
        elif .items then
            # List of items
            if (.items | length) > 0 then
                (.items[0] | keys | @tsv),
                (.items[] | [.[]] | @tsv)
            else
                "No items"
            end
        else
            # Single object
            to_entries | map([.key, .value] | @tsv) | .[]
        end
    ' 2>/dev/null || echo "$1"
}

# Format JSON as JSON (pretty)
format_json() {
    jq '.' 2>/dev/null || echo "$1"
}

# Run a query
run_query() {
    SQL="$1"
    LIMIT="${2:-100}"
    
    print_header "Running Query"
    echo "${YELLOW}SQL:${NC} ${SQL}"
    echo ""
    
    RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/queries/run" \
        -H "Content-Type: application/json" \
        -d "{\"sql\": \"${SQL}\", \"limit\": ${LIMIT}}")
    
    if echo "$RESPONSE" | jq -e '.error' >/dev/null 2>&1; then
        ERROR_MSG=$(echo "$RESPONSE" | jq -r '.error // .message // "Unknown error"')
        print_error "Query failed: ${ERROR_MSG}"
        return 1
    fi
    
    print_success "Query executed successfully"
    echo ""
    
    if [ "$FORMAT" = "table" ]; then
        echo "$RESPONSE" | format_table
    else
        echo "$RESPONSE" | format_json
    fi
}

# List saved queries
list_queries() {
    print_header "Saved Queries"
    
    RESPONSE=$(curl -s "${API_URL}/api/v1/queries/saved")
    
    if echo "$RESPONSE" | jq -e '.items' >/dev/null 2>&1; then
        COUNT=$(echo "$RESPONSE" | jq '.items | length')
        if [ "$COUNT" -eq 0 ]; then
            print_info "No saved queries found"
        else
            echo "$RESPONSE" | jq -r '.items[] | "\(.id)\t\(.name)\t\(.sql)"' | awk -F'\t' '{printf "%-40s %-30s %s\n", $1, $2, $3}'
        fi
    else
        print_error "Failed to fetch saved queries"
    fi
}

# Get a saved query
get_query() {
    QUERY_ID="$1"
    
    print_header "Query Details"
    
    RESPONSE=$(curl -s "${API_URL}/api/v1/queries/saved/${QUERY_ID}")
    
    if echo "$RESPONSE" | jq -e '.error' >/dev/null 2>&1; then
        print_error "Query not found: ${QUERY_ID}"
        return 1
    fi
    
    echo "$RESPONSE" | format_json
}

# Save a query
save_query() {
    NAME="$1"
    SQL="$2"
    TAGS="${3:-}"
    
    print_header "Saving Query"
    
    if [ -n "$TAGS" ]; then
        TAGS_JSON=$(echo "$TAGS" | tr ',' '\n' | jq -R . | jq -s .)
        PAYLOAD="{\"name\": \"${NAME}\", \"sql\": \"${SQL}\", \"tags\": ${TAGS_JSON}}"
    else
        PAYLOAD="{\"name\": \"${NAME}\", \"sql\": \"${SQL}\"}"
    fi
    
    RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/queries/saved" \
        -H "Content-Type: application/json" \
        -d "$PAYLOAD")
    
    if echo "$RESPONSE" | jq -e '.id' >/dev/null 2>&1; then
        QUERY_ID=$(echo "$RESPONSE" | jq -r '.id')
        print_success "Query saved with ID: ${QUERY_ID}"
        echo "$RESPONSE" | format_json
    else
        print_error "Failed to save query"
        echo "$RESPONSE" | format_json
        return 1
    fi
}

# Generate a report
generate_report() {
    SQL="$1"
    
    print_header "Generating Narrative Report"
    echo "${YELLOW}SQL:${NC} ${SQL}"
    echo ""
    print_info "This may take a moment (requires LLM)..."
    echo ""
    
    RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/reports/generate" \
        -H "Content-Type: application/json" \
        -d "{\"sql\": \"${SQL}\"}")
    
    if echo "$RESPONSE" | jq -e '.error' >/dev/null 2>&1; then
        ERROR_MSG=$(echo "$RESPONSE" | jq -r '.error // .message // "Unknown error"')
        print_error "Report generation failed: ${ERROR_MSG}"
        return 1
    fi
    
    print_success "Report generated successfully"
    echo ""
    
    if [ "$FORMAT" = "table" ]; then
        echo "$RESPONSE" | jq -r '.narrative // .summary // "No narrative available"'
    else
        echo "$RESPONSE" | format_json
    fi
}

# Show help
show_help() {
    cat <<EOF
${CYAN}PgQueryNarrative CLI${NC}

Usage: $0 <command> [options]

Commands:
  query <sql> [limit]          Run a SQL query
  list                         List all saved queries
  get <id>                     Get a saved query by ID
  save <name> <sql> [tags]     Save a query (tags: comma-separated)
  report <sql>                 Generate a narrative report
  help                         Show this help message

Environment Variables:
  PGQUERYNARRATIVE_API_URL     API base URL (default: http://app:8080)
  PGQUERYNARRATIVE_FORMAT      Output format: table or json (default: table)

Examples:
  # Run a query
  $0 query "SELECT * FROM demo.sales LIMIT 5"
  
  # Run with custom limit
  $0 query "SELECT * FROM demo.sales" 10
  
  # List saved queries
  $0 list
  
  # Save a query
  $0 save "Top Products" "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category" "sales,top"
  
  # Generate a report
  $0 report "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"
  
  # Use JSON output
  PGQUERYNARRATIVE_FORMAT=json $0 query "SELECT * FROM demo.sales LIMIT 3"

EOF
}

# Main
check_api

case "${1:-help}" in
    query)
        if [ -z "$2" ]; then
            print_error "SQL query required"
            show_help
            exit 1
        fi
        run_query "$2" "${3:-100}"
        ;;
    list)
        list_queries
        ;;
    get)
        if [ -z "$2" ]; then
            print_error "Query ID required"
            show_help
            exit 1
        fi
        get_query "$2"
        ;;
    save)
        if [ -z "$2" ] || [ -z "$3" ]; then
            print_error "Name and SQL required"
            show_help
            exit 1
        fi
        save_query "$2" "$3" "${4:-}"
        ;;
    report)
        if [ -z "$2" ]; then
            print_error "SQL query required"
            show_help
            exit 1
        fi
        generate_report "$2"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
