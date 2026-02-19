package queryrunner_test

import (
	"strings"
	"testing"

	"github.com/pgquerynarrative/pgquerynarrative/app/errors"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
)

func TestValidator(t *testing.T) {
	validator := queryrunner.NewValidator([]string{"demo"}, 10000)

	tests := []struct {
		name    string
		sql     string
		wantErr error
	}{
		{"valid_select", "SELECT * FROM demo.sales", nil},
		{"valid_with", "WITH cte AS (SELECT * FROM demo.sales) SELECT * FROM cte", nil},
		{"too_long", "SELECT * FROM demo.sales WHERE '" + strings.Repeat("a", 20000) + "' = 'x'", errors.ErrQueryTooLong},
		{"non_select", "UPDATE demo.sales SET quantity = 1", errors.ErrOnlySelectAllowed},
		{"disallowed", "SELECT * FROM demo.sales; DROP TABLE demo.sales", errors.ErrMultipleStatements},
		{"schema_not_allowed", "SELECT * FROM public.users", errors.ErrSchemaNotAllowed},
		{"alias_column_allowed", "WITH cte AS (SELECT region AS r FROM demo.sales) SELECT c.r FROM cte c", nil},
		{"complex_cte_with_aliases", `WITH regional_totals AS (SELECT region, product_category, SUM(total_amount) AS revenue FROM demo.sales WHERE date >= CURRENT_DATE - INTERVAL '365 days' GROUP BY region, product_category),
  regional_ranked AS (SELECT region, product_category, revenue, RANK() OVER (PARTITION BY region ORDER BY revenue DESC) AS cat_rank_in_region FROM regional_totals),
  sales_rep_perf AS (SELECT sales_rep, region, SUM(total_amount) AS rep_revenue FROM demo.sales GROUP BY sales_rep, region)
SELECT rr.region, rr.revenue, sr.rep_revenue
FROM regional_ranked rr
LEFT JOIN LATERAL (SELECT sales_rep, rep_revenue FROM sales_rep_perf WHERE sales_rep_perf.region = rr.region ORDER BY rep_revenue DESC LIMIT 1) sr ON true
WHERE rr.cat_rank_in_region <= 3`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.sql)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && err == nil {
				t.Fatalf("expected error %v", tt.wantErr)
			}
		})
	}
}
