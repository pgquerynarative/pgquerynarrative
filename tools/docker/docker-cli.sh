#!/bin/sh
# Wrapper script for running CLI commands in Docker

if [ $# -eq 0 ]; then
    echo "Usage: make cli CMD='<command> [args...]'"
    echo ""
    echo "Examples:"
    echo "  make cli CMD='list'"
    echo "  make cli CMD='query \"SELECT * FROM demo.sales LIMIT 5\"'"
    echo "  make cli CMD='save \"My Query\" \"SELECT * FROM demo.sales\"'"
    echo "  make cli CMD='report \"SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category\"'"
    exit 1
fi

docker compose run --rm cli /usr/local/bin/pgquerynarrative "$@"
