# PostgreSQL extension

Call PgQueryNarrative from SQL: run queries, generate reports, list saved queries. Requires PgQueryNarrative service and PostgreSQL 12+.

## Install

**1. HTTP extension (recommended):** `CREATE EXTENSION IF NOT EXISTS http;` (may need superuser).

**2. PgQueryNarrative extension:** `make install-extension` or copy `infra/postgres-extension/*.control` and `*.sql` to `$(pg_config --sharedir)/extension/`, then `CREATE EXTENSION pgquerynarrative;`. If using http: `\i infra/postgres-extension/pgquerynarrative--1.0--with-http.sql`.

## Config

```sql
SELECT pgquerynarrative_set_api_url('http://localhost:8080');
SELECT pgquerynarrative_get_api_url();
```

## Functions

| Function | Description |
|----------|-------------|
| `pgquerynarrative_run_query(query_sql, row_limit)` | Run query; returns JSON. Default limit 100. |
| `pgquerynarrative_generate_report(query_sql)` | Generate narrative report (JSON). |
| `pgquerynarrative_list_saved(limit, offset)` | List saved queries (JSON). |

**Example**

```sql
SELECT pgquerynarrative_run_query('SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category', 10);
SELECT pgquerynarrative_generate_report('SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category');
```

## Troubleshooting

- **Extension not found:** Check files in extension dir; PostgreSQL version; `CREATE EXTENSION` privilege.
- **Placeholder/empty result:** Install http extension and apply `pgquerynarrative--1.0--with-http.sql`.
- **Connection errors:** Verify service running, API URL, network/firewall.

**See also:** [API reference](api/README.md), [Configuration](configuration.md)
