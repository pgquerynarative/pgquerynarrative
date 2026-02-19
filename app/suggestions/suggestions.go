// Package suggestions provides query suggestion logic: curated example queries
// and matching saved queries by intent (keyword/substring).
package suggestions

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	suggestions "github.com/pgquerynarrative/pgquerynarrative/gen/suggestions"
)

// Curated example queries for the demo schema (demo.sales).
// These are always included when limit allows, so the AI has starting points.
var curatedQueries = []*suggestions.QuerySuggestion{
	{
		SQL:    "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category ORDER BY total DESC",
		Title:  "Sales by category (aggregate)",
		Source: "curated",
	},
	{
		SQL:    "SELECT date, SUM(total_amount) AS daily_total FROM demo.sales GROUP BY date ORDER BY date",
		Title:  "Daily sales over time",
		Source: "curated",
	},
	{
		SQL:    "SELECT region, product_category, SUM(quantity) AS qty FROM demo.sales GROUP BY region, product_category ORDER BY region, qty DESC",
		Title:  "Quantity by region and category",
		Source: "curated",
	},
}

// Suggester returns suggested SQL (curated + saved-query matches by intent).
type Suggester struct {
	appPool *pgxpool.Pool
}

// NewSuggester creates a suggester that uses the app pool to list saved queries
// for intent matching.
func NewSuggester(appPool *pgxpool.Pool) *Suggester {
	return &Suggester{appPool: appPool}
}

// Queries implements the suggestions service: returns curated examples plus
// saved queries matching the optional intent, up to limit.
func (s *Suggester) Queries(ctx context.Context, payload *suggestions.QueriesPayload) (*suggestions.SuggestedQueriesResult, error) {
	limit := int(payload.Limit)
	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	var out []*suggestions.QuerySuggestion

	// Add curated (up to limit)
	for _, c := range curatedQueries {
		if len(out) >= limit {
			break
		}
		out = append(out, &suggestions.QuerySuggestion{
			SQL:    c.SQL,
			Title:  c.Title,
			Source: c.Source,
		})
	}

	// If intent provided, add saved queries that match (name, description, or sql)
	intent := ""
	if payload.Intent != nil {
		intent = strings.TrimSpace(*payload.Intent)
	}
	if intent != "" && len(out) < limit {
		saved, err := s.matchSavedQueries(ctx, intent, limit-len(out))
		if err == nil {
			for _, q := range saved {
				out = append(out, q)
				if len(out) >= limit {
					break
				}
			}
		}
	}

	return &suggestions.SuggestedQueriesResult{Suggestions: out}, nil
}

// matchSavedQueries returns saved queries whose name, description, or sql
// contain the intent string (case-insensitive substring), up to max.
// Uses position() so intent is not interpreted as a LIKE pattern.
func (s *Suggester) matchSavedQueries(ctx context.Context, intent string, max int) ([]*suggestions.QuerySuggestion, error) {
	lowerIntent := strings.ToLower(intent)
	rows, err := s.appPool.Query(ctx, `
		SELECT name, sql, COALESCE(description, '')
		FROM app.saved_queries
		WHERE position($1 in lower(name)) > 0
		   OR position($1 in lower(COALESCE(description, ''))) > 0
		   OR position($1 in lower(sql)) > 0
		ORDER BY updated_at DESC
		LIMIT $2
	`, lowerIntent, max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*suggestions.QuerySuggestion
	for rows.Next() {
		var name, sql, desc string
		if err := rows.Scan(&name, &sql, &desc); err != nil {
			return nil, err
		}
		title := name
		if desc != "" {
			title = name + ": " + truncate(desc, 60)
		}
		result = append(result, &suggestions.QuerySuggestion{
			SQL:    sql,
			Title:  title,
			Source: "saved",
		})
	}
	return result, rows.Err()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
