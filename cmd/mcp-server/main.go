// Package main runs an MCP server that exposes PgQueryNarrative as tools.
// Connect Claude (desktop or Cursor) to this server to run queries and generate reports via MCP.
//
// Prerequisites: PgQueryNarrative HTTP server must be running (e.g. make start-local).
// Set PGQUERYNARRATIVE_URL if the server is not at http://localhost:8080.
//
// Run: go run ./cmd/mcp-server
// Or build and run: go build -o bin/mcp-server ./cmd/mcp-server && ./bin/mcp-server
//
// Add to Claude desktop: MCP server command = path to this binary (stdio transport).
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultBaseURL = "http://localhost:8080"

func main() {
	baseURL := os.Getenv("PGQUERYNARRATIVE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	client := &apiClient{baseURL: baseURL, http: &http.Client{Timeout: 60 * time.Second}}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pgquerynarrative",
		Version: "1.0.0",
	}, nil)

	// run_query: execute a read-only SQL query against the demo schema
	mcp.AddTool(server, &mcp.Tool{
		Name:        "run_query",
		Description: "Run a read-only SQL query against the PgQueryNarrative database (demo schema). Returns columns and rows as JSON.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input RunQueryInput) (*mcp.CallToolResult, any, error) {
		limit := input.Limit
		if limit <= 0 {
			limit = 100
		}
		body, err := client.post(ctx, "/api/v1/queries/run", map[string]any{"sql": input.SQL, "limit": limit})
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// generate_report: run a query and generate an LLM narrative report
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_report",
		Description: "Run a SQL query and generate a narrative report (headline, takeaways, drivers, limitations, recommendations) using the configured LLM.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GenerateReportInput) (*mcp.CallToolResult, any, error) {
		body, err := client.post(ctx, "/api/v1/reports/generate", map[string]any{"sql": input.SQL})
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// list_saved_queries: list saved queries
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_saved_queries",
		Description: "List saved queries (optional limit and offset).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListSavedQueriesInput) (*mcp.CallToolResult, any, error) {
		limit, offset := input.Limit, input.Offset
		if limit <= 0 {
			limit = 20
		}
		path := fmt.Sprintf("/api/v1/queries/saved?limit=%d&offset=%d", limit, offset)
		body, err := client.get(ctx, path)
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// get_report: fetch a report by ID
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_report",
		Description: "Get a report by its ID (UUID).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetReportInput) (*mcp.CallToolResult, any, error) {
		body, err := client.get(ctx, "/api/v1/reports/"+input.ID)
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// list_reports: list reports with optional limit/offset
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_reports",
		Description: "List generated reports (optional limit and offset).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListReportsInput) (*mcp.CallToolResult, any, error) {
		limit, offset := input.Limit, input.Offset
		if limit <= 0 {
			limit = 20
		}
		path := fmt.Sprintf("/api/v1/reports?limit=%d&offset=%d", limit, offset)
		body, err := client.get(ctx, path)
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// get_schema: return database schema (allowed schemas, tables, columns) for querying
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_schema",
		Description: "Returns the database schema available for querying (allowed schemas, tables, columns). Use this to see what tables and columns you can use in run_query.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetSchemaInput) (*mcp.CallToolResult, any, error) {
		body, err := client.get(ctx, "/api/v1/schema")
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	// get_context: combined schema + saved queries for full context
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_context",
		Description: "Returns combined context: schema (tables, columns) plus a list of saved queries (name, sql, description). Use this to understand the data model and existing saved queries before suggesting or running SQL.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetContextInput) (*mcp.CallToolResult, any, error) {
		savedLimit, savedOffset := input.SavedLimit, input.SavedOffset
		if savedLimit <= 0 {
			savedLimit = 20
		}
		schemaBody, err1 := client.get(ctx, "/api/v1/schema")
		savedPath := fmt.Sprintf("/api/v1/queries/saved?limit=%d&offset=%d", savedLimit, savedOffset)
		savedBody, err2 := client.get(ctx, savedPath)
		out := buildContextResult(schemaBody, savedBody, err1, err2)
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: out}}}, nil, nil
	})

	// suggest_queries: suggest SQL based on optional intent (curated + saved-query match)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "suggest_queries",
		Description: "Suggests SQL queries based on optional intent (e.g. 'sales by category'). Returns curated examples and saved queries that match. Use the suggested SQL with run_query or refine before running.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SuggestQueriesInput) (*mcp.CallToolResult, any, error) {
		path := "/api/v1/suggestions/queries?"
		if input.Intent != "" {
			path += "intent=" + url.QueryEscape(input.Intent) + "&"
		}
		if input.Limit > 0 {
			path += fmt.Sprintf("limit=%d", input.Limit)
		} else {
			path += "limit=5"
		}
		body, err := client.get(ctx, path)
		if err != nil {
			return toolError(err), nil, nil
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: body}}}, nil, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "mcp-server: %v\n", err)
		os.Exit(1)
	}
}

type RunQueryInput struct {
	SQL   string `json:"sql" jsonschema:"Read-only SQL query e.g. SELECT from demo.sales"`
	Limit int    `json:"limit" jsonschema:"Max rows to return"`
}

type GenerateReportInput struct {
	SQL string `json:"sql" jsonschema:"SQL query for the report e.g. SELECT category SUM amount FROM demo.sales GROUP BY category"`
}

type ListSavedQueriesInput struct {
	Limit  int `json:"limit" jsonschema:"Max items to return"`
	Offset int `json:"offset" jsonschema:"Offset for pagination"`
}

type GetReportInput struct {
	ID string `json:"id" jsonschema:"Report UUID"`
}

type ListReportsInput struct {
	Limit  int `json:"limit" jsonschema:"Max items to return"`
	Offset int `json:"offset" jsonschema:"Offset for pagination"`
}

type GetSchemaInput struct{}

type GetContextInput struct {
	SavedLimit  int `json:"saved_limit" jsonschema:"Max saved queries to include (default 20)"`
	SavedOffset int `json:"saved_offset" jsonschema:"Offset for saved queries (default 0)"`
}

type SuggestQueriesInput struct {
	Intent string `json:"intent" jsonschema:"Optional natural-language intent to match saved queries (e.g. sales by region)"`
	Limit  int    `json:"limit" jsonschema:"Max suggestions to return (default 5)"`
}

type apiClient struct {
	baseURL string
	http    *http.Client
}

func (c *apiClient) post(ctx context.Context, path string, body map[string]any) (string, error) {
	enc, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(enc))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *apiClient) get(ctx context.Context, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return "", err
	}
	return c.do(req)
}

func (c *apiClient) do(req *http.Request) (string, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API %s %s: %d %s", req.Method, req.URL.Path, resp.StatusCode, string(b))
	}
	return string(b), nil
}

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
	}
}

// buildContextResult merges schema and saved-queries API responses into one text block.
// On partial failure, includes what succeeded and notes the error.
func buildContextResult(schemaBody, savedBody string, schemaErr, savedErr error) string {
	var b bytes.Buffer
	b.WriteString("=== Schema (queryable tables and columns) ===\n")
	if schemaErr != nil {
		b.WriteString("(schema unavailable: " + schemaErr.Error() + ")\n")
	} else {
		b.WriteString(schemaBody)
	}
	b.WriteString("\n\n=== Saved queries ===\n")
	if savedErr != nil {
		b.WriteString("(saved queries unavailable: " + savedErr.Error() + ")\n")
	} else {
		b.WriteString(savedBody)
	}
	return b.String()
}
