// Package charts provides chart-type suggestions based on query result shape.
package charts

import (
	"fmt"
	"strings"

	"github.com/pgquerynarrative/pgquerynarrative/app/metrics"
)

// Suggestion holds a single chart suggestion (chart_type, label, reason).
type Suggestion struct {
	ChartType string
	Label     string
	Reason    string
}

// Suggest returns chart suggestions based on column names, types, and row data.
// Uses metrics.ProfileColumns to infer dimensions, measures, and time series,
// then suggests bar, line, pie, or table as appropriate.
func Suggest(columnNames []string, columnTypes []string, rows [][]interface{}) []Suggestion {
	if len(columnNames) == 0 {
		return nil
	}
	profiles := metrics.ProfileColumns(columnNames, rows)

	var dateCols, numericCols, textCols []int
	for i, p := range profiles {
		switch p.Type {
		case metrics.ColumnTypeDate:
			dateCols = append(dateCols, i)
		case metrics.ColumnTypeNumeric:
			numericCols = append(numericCols, i)
		case metrics.ColumnTypeText:
			textCols = append(textCols, i)
		}
	}

	var out []Suggestion
	seen := make(map[string]bool)
	add := func(chartType, label, reason string) {
		if seen[chartType] {
			return
		}
		seen[chartType] = true
		out = append(out, Suggestion{ChartType: chartType, Label: label, Reason: reason})
	}

	// Time series: date + at least one numeric
	if len(dateCols) > 0 && len(numericCols) > 0 {
		add("line", "Line chart", "Date/time column with numeric values suits a time series line chart.")
	}

	// Category + value: text/dimension + numeric
	if len(textCols) > 0 && len(numericCols) > 0 {
		add("bar", "Bar chart", "Category column with numeric values suits a bar chart.")
		// Pie only if few categories (sample first column for distinct count)
		distinct := distinctCount(rows, textCols[0])
		if distinct >= 2 && distinct <= 12 {
			add("pie", "Pie chart", "Few categories with a single value column suit a pie chart.")
		}
	}

	// Multiple numeric series (e.g. multiple metrics)
	if len(numericCols) > 1 {
		add("line", "Line chart (multi-series)", "Multiple numeric columns can be shown as series on a line chart.")
	}

	// Table is always useful for raw data
	add("table", "Table", "Tabular view of the result set.")

	return out
}

func distinctCount(rows [][]interface{}, colIndex int) int {
	m := make(map[string]bool)
	for _, row := range rows {
		if colIndex >= len(row) {
			continue
		}
		v := row[colIndex]
		var key string
		if v == nil {
			key = ""
		} else {
			key = strings.TrimSpace(stringify(v))
		}
		m[key] = true
	}
	return len(m)
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return fmt.Sprint(v)
	}
}
