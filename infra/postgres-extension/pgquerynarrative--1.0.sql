-- PostgreSQL extension for PgQueryNarrative
-- Provides functions to interact with PgQueryNarrative API from within PostgreSQL

-- Configuration: API URL (defaults to http://localhost:8080)
CREATE OR REPLACE FUNCTION pgquerynarrative_set_api_url(url TEXT)
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM set_config('pgquerynarrative.api_url', url, false);
END;
$$;

-- Get current API URL
CREATE OR REPLACE FUNCTION pgquerynarrative_get_api_url()
RETURNS TEXT
LANGUAGE plpgsql
STABLE
AS $$
BEGIN
    RETURN COALESCE(
        current_setting('pgquerynarrative.api_url', true),
        'http://localhost:8080'
    );
END;
$$;

-- Execute a query via PgQueryNarrative API
-- Returns JSON with query results
CREATE OR REPLACE FUNCTION pgquerynarrative_run_query(
    query_sql TEXT,
    row_limit INTEGER DEFAULT 100
)
RETURNS JSON
LANGUAGE plpgsql
AS $$
DECLARE
    api_url TEXT;
    response JSON;
    http_response TEXT;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    
    -- Use PostgreSQL's http extension if available
    -- Otherwise, this requires PL/Python or similar
    -- For now, we'll provide a template that can be extended
    
    -- Note: This requires the http extension or PL/Python
    -- Example with http extension:
    -- SELECT content::json FROM http((
    --     'POST',
    --     api_url || '/api/v1/queries/run',
    --     ARRAY[
    --         http_header('Content-Type', 'application/json')
    --     ],
    --     'application/json',
    --     json_build_object('sql', query_sql, 'limit', row_limit)::text
    -- )::http_request);
    
    -- For basic implementation, return a placeholder
    -- Full implementation requires http extension or PL/Python
    RAISE NOTICE 'PgQueryNarrative API URL: %', api_url;
    RAISE NOTICE 'Query: %', query_sql;
    RAISE NOTICE 'Limit: %', row_limit;
    
    -- Return a JSON structure indicating the function was called
    RETURN json_build_object(
        'status', 'pending',
        'message', 'HTTP extension or PL/Python required for full functionality',
        'api_url', api_url,
        'query', query_sql,
        'limit', row_limit
    );
END;
$$;

-- Generate a narrative report from a query
CREATE OR REPLACE FUNCTION pgquerynarrative_generate_report(
    query_sql TEXT
)
RETURNS JSON
LANGUAGE plpgsql
AS $$
DECLARE
    api_url TEXT;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    
    RAISE NOTICE 'PgQueryNarrative API URL: %', api_url;
    RAISE NOTICE 'Generating report for query: %', query_sql;
    
    -- Return a JSON structure indicating the function was called
    RETURN json_build_object(
        'status', 'pending',
        'message', 'HTTP extension or PL/Python required for full functionality',
        'api_url', api_url,
        'query', query_sql
    );
END;
$$;

-- List saved queries
CREATE OR REPLACE FUNCTION pgquerynarrative_list_saved(
    query_limit INTEGER DEFAULT 50,
    query_offset INTEGER DEFAULT 0
)
RETURNS JSON
LANGUAGE plpgsql
AS $$
DECLARE
    api_url TEXT;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    
    RAISE NOTICE 'PgQueryNarrative API URL: %', api_url;
    RAISE NOTICE 'Listing saved queries (limit: %, offset: %)', query_limit, query_offset;
    
    RETURN json_build_object(
        'status', 'pending',
        'message', 'HTTP extension or PL/Python required for full functionality',
        'api_url', api_url,
        'limit', query_limit,
        'offset', query_offset
    );
END;
$$;

-- Grant execute permissions to public (adjust as needed for security)
GRANT EXECUTE ON FUNCTION pgquerynarrative_set_api_url(TEXT) TO PUBLIC;
GRANT EXECUTE ON FUNCTION pgquerynarrative_get_api_url() TO PUBLIC;
GRANT EXECUTE ON FUNCTION pgquerynarrative_run_query(TEXT, INTEGER) TO PUBLIC;
GRANT EXECUTE ON FUNCTION pgquerynarrative_generate_report(TEXT) TO PUBLIC;
GRANT EXECUTE ON FUNCTION pgquerynarrative_list_saved(INTEGER, INTEGER) TO PUBLIC;
