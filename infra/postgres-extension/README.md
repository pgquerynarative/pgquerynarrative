# PgQueryNarrative PostgreSQL Extension

This extension allows you to call PgQueryNarrative functions directly from within PostgreSQL SQL queries.

## Installation

### Option 1: Using http Extension (Recommended)

1. Install the http extension:
```sql
CREATE EXTENSION IF NOT EXISTS http;
```

2. Install PgQueryNarrative extension:
```sql
CREATE EXTENSION pgquerynarrative;
```

3. Apply the full implementation:
```sql
\i infra/postgres-extension/pgquerynarrative--1.0--with-http.sql
```

### Option 2: Basic Installation (Placeholder)

If you don't have the http extension, you can install the basic version:

```sql
CREATE EXTENSION pgquerynarrative;
```

This provides the function signatures but requires PL/Python or another method for HTTP calls.

## Configuration

Set the PgQueryNarrative API URL:

```sql
SELECT pgquerynarrative_set_api_url('http://localhost:8080');
```

Get the current API URL:

```sql
SELECT pgquerynarrative_get_api_url();
```

## Usage

### Run a Query

Execute a SQL query through PgQueryNarrative:

```sql
SELECT pgquerynarrative_run_query(
    'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category',
    10
);
```

### Generate a Narrative Report

Generate an AI-powered narrative from query results:

```sql
SELECT pgquerynarrative_generate_report(
    'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category'
);
```

### List Saved Queries

Retrieve saved queries:

```sql
SELECT pgquerynarrative_list_saved(50, 0);
```

## Example: Using in a View

```sql
CREATE VIEW sales_summary AS
SELECT 
    product_category,
    total_amount,
    pgquerynarrative_generate_report(
        'SELECT product_category, SUM(total_amount) FROM demo.sales WHERE product_category = ''' || product_category || ''' GROUP BY product_category'
    ) as narrative
FROM (
    SELECT product_category, SUM(total_amount) as total_amount
    FROM demo.sales
    GROUP BY product_category
) sub;
```

## Requirements

- PostgreSQL 12+
- http extension (for full functionality) or PL/Python
- PgQueryNarrative service running and accessible

## Security Considerations

- The extension makes HTTP calls to the PgQueryNarrative API
- Ensure proper network security between PostgreSQL and the API
- Consider using SSL/TLS for API connections
- Review and adjust function permissions as needed

## Troubleshooting

If you get errors about the http extension:

1. Install the http extension: `CREATE EXTENSION http;`
2. Or use PL/Python to implement HTTP calls
3. Or use the basic placeholder version for testing

## Limitations

- HTTP calls are synchronous and will block the query
- Timeout handling depends on PostgreSQL's http extension settings
- Large result sets may impact performance
