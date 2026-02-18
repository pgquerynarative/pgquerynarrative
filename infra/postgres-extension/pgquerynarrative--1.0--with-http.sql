-- PostgreSQL extension for PgQueryNarrative (requires http extension)
-- Full implementation using PostgreSQL's http extension
-- Install with: CREATE EXTENSION http; CREATE EXTENSION pgquerynarrative;

-- Execute a query via PgQueryNarrative API using http extension
CREATE OR REPLACE FUNCTION pgquerynarrative_run_query(
    query_sql TEXT,
    row_limit INTEGER DEFAULT 100
)
RETURNS JSON
LANGUAGE plpgsql
AS $$
DECLARE
    api_url TEXT;
    request_body TEXT;
    response http_response;
    result JSON;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    request_body := json_build_object(
        'sql', query_sql,
        'limit', row_limit
    )::text;
    
    -- Make HTTP POST request to PgQueryNarrative API
    SELECT * INTO response FROM http((
        'POST',
        api_url || '/api/v1/queries/run',
        ARRAY[
            http_header('Content-Type', 'application/json')
        ],
        'application/json',
        request_body
    )::http_request);
    
    -- Parse response
    IF response.status = 200 THEN
        result := response.content::json;
    ELSE
        RAISE EXCEPTION 'PgQueryNarrative API error: % - %', response.status, response.content;
    END IF;
    
    RETURN result;
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
    request_body TEXT;
    response http_response;
    result JSON;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    request_body := json_build_object('sql', query_sql)::text;
    
    -- Make HTTP POST request to PgQueryNarrative API
    SELECT * INTO response FROM http((
        'POST',
        api_url || '/api/v1/reports/generate',
        ARRAY[
            http_header('Content-Type', 'application/json')
        ],
        'application/json',
        request_body
    )::http_request);
    
    -- Parse response
    IF response.status = 200 THEN
        result := response.content::json;
    ELSE
        RAISE EXCEPTION 'PgQueryNarrative API error: % - %', response.status, response.content;
    END IF;
    
    RETURN result;
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
    response http_response;
    result JSON;
BEGIN
    api_url := pgquerynarrative_get_api_url();
    
    -- Make HTTP GET request to PgQueryNarrative API
    SELECT * INTO response FROM http((
        'GET',
        api_url || '/api/v1/queries/saved?limit=' || query_limit || '&offset=' || query_offset,
        ARRAY[]::http_header[]
    )::http_request);
    
    -- Parse response
    IF response.status = 200 THEN
        result := response.content::json;
    ELSE
        RAISE EXCEPTION 'PgQueryNarrative API error: % - %', response.status, response.content;
    END IF;
    
    RETURN result;
END;
$$;
