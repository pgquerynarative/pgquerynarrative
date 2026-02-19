# PostgreSQL extension

Call PgQueryNarrative from SQL: run queries, generate reports, and list saved queries. Requires the PgQueryNarrative service to be running and PostgreSQL 12+.

## Install

1. **HTTP extension (recommended):**  
   `CREATE EXTENSION IF NOT EXISTS http;` (may require superuser).

2. **PgQueryNarrative extension:**  
   Run `make install-extension`, or copy `infra/postgres-extension/*.control` and `*.sql` to `$(pg_config --sharedir)/extension/`, then:
   ```sql
   CREATE EXTENSION pgquerynarrative;
   ```
   If using the HTTP extension:  
   `\i infra/postgres-extension/pgquerynarrative--1.0--with-http.sql`

## Configuration

```sql
SELECT pgquerynarrative_set_api_url('http://localhost:8080');
SELECT pgquerynarrative_get_api_url();
```

## Functions

| Function | Description |
|----------|-------------|
| `pgquerynarrative_run_query(query_sql, row_limit)` | Run a read-only query; returns JSON. Default limit 100. |
| `pgquerynarrative_generate_report(query_sql)` | Generate a narrative report (JSON). |
| `pgquerynarrative_list_saved(limit, offset)` | List saved queries (JSON). |

## Example

```sql
SELECT pgquerynarrative_run_query(
  'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category',
  10
);
SELECT pgquerynarrative_generate_report(
  'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category'
);
```

## Troubleshooting

- **Extension not found:** Check that files are in the extension directory, PostgreSQL version is supported, and you have `CREATE EXTENSION` privilege.
- **Placeholder or empty result:** Install the HTTP extension and apply `pgquerynarrative--1.0--with-http.sql`.
- **Connection errors:** Verify the PgQueryNarrative service is running, the API URL is correct, and network/firewall allow the connection.

## See also

- [API reference](../api/README.md)
- [Configuration](../configuration.md)
