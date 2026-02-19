# MCP Schema, Context, and Query Suggestion Design

This document describes the architecture and workflow for MCP schema tools, context retrieval, and integrated query suggestion capabilities.

## Goals

- **Schema tools:** Expose allowed database schema (tables, columns) to MCP clients so AI assistants know what can be queried.
- **Context:** Provide a combined view of schema plus saved queries so the assistant has full context for refining or suggesting SQL.
- **Query suggestion:** Return suggested SQL (curated examples and/or saved-query matches) from natural-language intent.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│  MCP Server (cmd/mcp-server)                                            │
│  Tools: get_schema, get_context, suggest_queries, run_query, ...         │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │ HTTP (PGQUERYNARRATIVE_URL)
                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  PgQueryNarrative API                                                    │
│  GET /api/v1/schema           → catalog service (read-only pool)        │
│  GET /api/v1/queries/saved    → existing queries service                 │
│  GET /api/v1/suggestions/queries?intent=... → suggestions service        │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │
        ┌───────────────────────┼───────────────────────┐
        ▼                       ▼                       ▼
┌───────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ app/catalog   │     │ app/service     │     │ app/suggestions  │
│ (schema from  │     │ (queries,        │     │ (curated +       │
│ information_ │     │  reports)        │     │  saved-query     │
│ schema)       │     │                  │     │  match)          │
└───────┬───────┘     └────────┬─────────┘     └────────┬─────────┘
        │                     │                         │
        ▼                     ▼                         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  DB: read-only pool (queries, schema) │ app pool (saved_queries, reports) │
└─────────────────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| **app/catalog** | Queries PostgreSQL `information_schema` (via read-only pool) for allowed schemas; returns tables and columns. Only schemas configured as allowed (e.g. `demo`) are included. |
| **Schema API** | Goa service `schema` with one method: get schema/catalog. Used by MCP `get_schema` and by suggestions. |
| **Context** | No new backend endpoint. MCP tool `get_context` calls existing schema + list_saved_queries and merges the response. |
| **app/suggestions** | Builds suggested SQL: (1) curated example queries for demo schema, (2) optional keyword match over saved queries (name, description, sql). |
| **Suggestions API** | Goa service or method: `GET /api/v1/suggestions/queries?intent=...` returns list of suggested SQL with optional title/source. |

## Data Flow

### Schema discovery

1. Server is configured with allowed schemas (e.g. `["demo"]`), same as query validator.
2. Catalog loads table/column metadata from `information_schema.tables` and `information_schema.columns`, filtered by `table_schema IN (allowed)`.
3. Result: list of schemas, each with tables, each table with columns (name, data_type). No DDL, no row counts (kept simple).

### Context assembly

1. MCP client calls `get_context` (optional limit for saved queries).
2. MCP server calls `GET /api/v1/schema` and `GET /api/v1/queries/saved?limit=20&offset=0`.
3. Server merges into one response: schema summary + saved queries (id, name, sql, description). Returned as a single text/JSON block for the LLM.

### Query suggestion pipeline

1. MCP client calls `suggest_queries` with optional `intent` (e.g. "sales by region").
2. MCP server calls `GET /api/v1/suggestions/queries?intent=...`.
3. Backend suggestions logic:
   - Always include 2–3 curated example queries (e.g. demo.sales aggregates, time series).
   - If `intent` is provided: match saved queries by substring in name, description, or sql; return up to N (e.g. 5) suggestions.
4. Response: list of `{ "sql", "title", "source" }` (source = "curated" or "saved").
5. MCP returns this as tool result for the AI to use or present to the user.

## MCP Tool Contracts

### get_schema

- **Name:** `get_schema`
- **Description:** Returns the database schema available for querying (allowed schemas, tables, columns). Use this to see what tables and columns you can use in run_query.
- **Input:** None (optional `schemas` filter could be added later).
- **Output:** JSON or formatted text: list of schemas, each with tables and columns (name, type).
- **Errors:** API/network errors returned as tool error.

### get_context

- **Name:** `get_context`
- **Description:** Returns combined context: schema (tables, columns) plus a list of saved queries (name, sql, description). Use this to understand the data model and existing saved queries before suggesting or running SQL.
- **Input:** Optional `saved_limit` (default 20), `saved_offset` (default 0).
- **Output:** Single blob: schema summary + saved queries list.
- **Errors:** If schema or saved list fails, return partial result with error note.

### suggest_queries

- **Name:** `suggest_queries`
- **Description:** Suggests SQL queries based on optional intent (e.g. "sales by category"). Returns curated examples and, when intent is provided, saved queries that match. Use the suggested SQL with run_query or refine before running.
- **Input:** Optional `intent` (string). Optional `limit` (max suggestions, default 5).
- **Output:** List of suggestions: each with sql, title, source (curated | saved).
- **Errors:** API errors returned as tool error.

## Security and Validation

- Schema listing uses the **read-only** pool; only objects visible to the read-only user are exposed (same as run_query).
- All suggested or returned SQL remains subject to existing run_query validation (read-only, allowed schemas, no dangerous keywords).
- No PII or credentials in context; saved queries are already stored by the app.

## Scalability and Limits

- Schema: single DB; catalog result is small (one schema, few tables). No pagination for schema.
- Saved queries in context: use limit/offset to cap response size.
- Suggestions: cap curated list (e.g. 3), cap saved-query matches (e.g. 5), limit total response size.

## Implementation Phases

1. **Phase 1 – Schema:** Add `app/catalog`, Goa schema service, GET /api/v1/schema, wire in server; add MCP `get_schema`.
2. **Phase 2 – Context:** Add MCP `get_context` that calls schema + list_saved_queries and merges.
3. **Phase 3 – Suggestions:** Add `app/suggestions`, suggestions API, GET /api/v1/suggestions/queries; add MCP `suggest_queries`.
4. **Phase 4 – Tests and docs:** Unit tests for catalog and suggestions; integration/manual test for MCP tools; update MCP and API docs.

## Testing the MCP new features

### 1. Backend API (no MCP)

With the app running (`make start-local`), the MCP tools call these endpoints. Verify them first:

```bash
# Schema
curl -s http://localhost:8080/api/v1/schema | jq '.schemas[0].tables[0].name'
# Expect: "sales"

# Suggestions (no intent)
curl -s "http://localhost:8080/api/v1/suggestions/queries?limit=3" | jq '.suggestions | length'
# Expect: 3

# Suggestions with intent (matches saved queries)
curl -s "http://localhost:8080/api/v1/suggestions/queries?intent=sales&limit=5" | jq '.suggestions'
```

### 2. Cursor (or Claude desktop) – real MCP usage

1. **Start the app** in one terminal: `make start-local` (keep it running).
2. **Build the MCP server:** `make build-mcp` (produces `bin/mcp-server`).
3. **Configure MCP** in Cursor: Settings → MCP → add server, e.g.:
   - **Command:** full path to `bin/mcp-server`, e.g. `/path/to/pgquerynarative/bin/mcp-server`
   - **Env** (if app is not on default URL): `PGQUERYNARRATIVE_URL=http://localhost:8080`
4. **Use the tools in chat:** Ask the AI to:
   - “Call **get_schema**” → should return schema with `demo.sales` and columns.
   - “Call **get_context**” → schema + saved queries in one block.
   - “Call **suggest_queries** with intent ‘sales by region’” → list of suggested SQL (curated + any matching saved).

### 3. MCP Inspector (stdio, no Cursor)

Use the official Inspector to drive the MCP server over stdio and call tools with custom arguments:

1. Start the app: `make start-local` (in another terminal).
2. Install and run (Node.js required):
   ```bash
   npx @modelcontextprotocol/inspector ./bin/mcp-server
   ```
   Or with explicit env:
   ```bash
   PGQUERYNARRATIVE_URL=http://localhost:8080 npx @modelcontextprotocol/inspector ./bin/mcp-server
   ```
3. In the Inspector: open the **Tools** tab, select `get_schema`, `get_context`, or `suggest_queries`, set arguments (e.g. `{"intent":"sales"}` for suggest_queries), run and check the result.

### 4. Unit tests

Catalog and suggestions logic are covered by unit tests; run:

```bash
make test-unit
```

Relevant packages: `test/unit/app/catalog`, `test/unit/app/suggestions`.

## References

- Existing MCP tools: `cmd/mcp-server/main.go`
- Query validator allowed schemas: `app/queryrunner/validator.go`
- Demo schema: `app/db/migrations/000003_create_demo_schema.up.sql`
- STATUS: `docs/project/STATUS.md` (Phase 5 – MCP schema tools, query suggestions, context retrieval)
