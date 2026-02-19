package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
)

type Handlers struct {
	queriesEndpoints *queries.Endpoints
	reportsEndpoints *reports.Endpoints
}

func NewHandlers(queriesEndpoints *queries.Endpoints, reportsEndpoints *reports.Endpoints) *Handlers {
	return &Handlers{
		queriesEndpoints: queriesEndpoints,
		reportsEndpoints: reportsEndpoints,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	// Check API health
	apiStatus := "operational"
	apiStatusClass := "status-ok"
	apiMessage := "All systems operational"

	// Try to check if API is responding
	ctx := r.Context()
	_, err := h.queriesEndpoints.ListSaved(ctx, &queries.ListSavedPayload{Limit: 1, Offset: 0})
	if err != nil {
		apiStatus = "degraded"
		apiStatusClass = "status-warning"
		apiMessage = "API may be experiencing issues"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>PgQueryNarrative</title>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<script src="/static/js/htmx-config.js"></script>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body data-page="home">
	<nav class="navbar">
		<div class="container navbar-inner">
			<a href="/" class="navbar-brand">
				<span class="navbar-logo">PgQueryNarrative</span>
				<span class="navbar-tagline">Turn SQL into narratives</span>
			</a>
			<div class="navbar-links">
				<a href="/" class="nav-link">Home</a>
				<a href="/query" class="nav-link">Run Query</a>
				<a href="/saved" class="nav-link">Saved</a>
				<a href="/reports" class="nav-link">Reports</a>
			</div>
		</div>
	</nav>
	<main class="container">
		<div class="status-badge %s">
			<span class="status-indicator"></span>
			<span><strong>System Status:</strong> %s - %s</span>
		</div>
		
		<div class="hero">
			<h2>Welcome to PgQueryNarrative</h2>
			<p>Execute SQL queries and generate AI-powered business narratives from your PostgreSQL data</p>
		</div>
		
		<div class="quick-stats">
			<div class="stat-card">
				<div class="stat-value">✓</div>
				<div class="stat-label">Server Running</div>
			</div>
			<div class="stat-card">
				<div class="stat-value">API</div>
				<div class="stat-label">RESTful API Active</div>
			</div>
			<div class="stat-card">
				<div class="stat-value">AI</div>
				<div class="stat-label">Narrative Generation Ready</div>
			</div>
		</div>
		
		<div class="features">
			<div class="feature-card">
				<div class="feature-icon">🔍</div>
				<h3>Run Queries</h3>
				<p>Execute SQL safely with read-only access and validation. View results as a table or chart.</p>
				<a href="/query" class="btn btn-primary">Run Query</a>
			</div>
			<div class="feature-card">
				<div class="feature-icon">💾</div>
				<h3>Saved Queries</h3>
				<p>Store and reuse your frequently used queries across sessions.</p>
				<a href="/saved" class="btn btn-secondary">View Saved</a>
			</div>
			<div class="feature-card">
				<div class="feature-icon">📊</div>
				<h3>Reports</h3>
				<p>Generate AI-powered business narratives and chart suggestions from query results.</p>
				<a href="/reports" class="btn btn-secondary">Create Report</a>
			</div>
		</div>
		
		<div class="actions">
			<a href="/query" class="btn btn-primary">Get Started</a>
		</div>
	</main>
	<footer>
		<div class="container">
			<p>&copy; 2026 PgQueryNarrative &middot; <a href="/api/v1/queries/saved">API</a></p>
		</div>
	</footer>
</body>
</html>`, apiStatusClass, apiStatus, apiMessage)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, html)
}

func (h *Handlers) QueryPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>Run Query - PgQueryNarrative</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>
	<script src="/static/js/htmx-config.js"></script>
	<script src="/static/js/charts.js"></script>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body data-page="query">
	<nav class="navbar">
		<div class="container navbar-inner">
			<a href="/" class="navbar-brand">
				<span class="navbar-logo">PgQueryNarrative</span>
				<span class="navbar-tagline">Turn SQL into narratives</span>
			</a>
			<div class="navbar-links">
				<a href="/" class="nav-link">Home</a>
				<a href="/query" class="nav-link">Run Query</a>
				<a href="/saved" class="nav-link">Saved</a>
				<a href="/reports" class="nav-link">Reports</a>
			</div>
		</div>
	</nav>
	<main class="container">
		<div class="query-page">
			<h2 class="page-title">Run SQL Query</h2>
			<p class="page-subtitle">Execute read-only SQL against the demo schema. Results can be charted.</p>
			<form hx-post="/web/query/run" hx-target="#results" hx-swap="innerHTML" hx-indicator="#loading" class="query-form">
				<div class="form-group">
					<label for="sql">SQL Query</label>
					<textarea id="sql" name="sql" rows="10" required placeholder="SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category"></textarea>
				</div>
				<div class="form-group">
					<label for="limit">Result Limit</label>
					<input type="number" id="limit" name="limit" value="100" min="1" max="1000"/>
				</div>
				<button type="submit" class="btn btn-primary">Execute Query</button>
			</form>
			<div id="loading" class="loading" style="display:none;">Loading...</div>
			<div id="results" class="results"></div>
		</div>
	</main>
	<footer>
		<div class="container">
			<p>&copy; 2026 PgQueryNarrative</p>
		</div>
	</footer>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html)
}

func (h *Handlers) SavedQueries(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>Saved Queries - PgQueryNarrative</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<script src="/static/js/htmx-config.js"></script>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body data-page="saved">
	<nav class="navbar">
		<div class="container navbar-inner">
			<a href="/" class="navbar-brand">
				<span class="navbar-logo">PgQueryNarrative</span>
				<span class="navbar-tagline">Turn SQL into narratives</span>
			</a>
			<div class="navbar-links">
				<a href="/" class="nav-link">Home</a>
				<a href="/query" class="nav-link">Run Query</a>
				<a href="/saved" class="nav-link">Saved</a>
				<a href="/reports" class="nav-link">Reports</a>
			</div>
		</div>
	</nav>
	<main class="container">
		<div class="saved-queries">
			<h2 class="page-title">Saved Queries</h2>
			<p class="page-subtitle">Reuse and manage your saved queries.</p>
			<div hx-get="/api/v1/queries/saved" hx-trigger="load" hx-swap="innerHTML" id="queries-list">
				<p>Loading...</p>
			</div>
		</div>
	</main>
	<footer>
		<div class="container">
			<p>&copy; 2026 PgQueryNarrative</p>
		</div>
	</footer>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html)
}

func (h *Handlers) Reports(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>Reports - PgQueryNarrative</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<script src="/static/js/htmx-config.js"></script>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body data-page="reports">
	<nav class="navbar">
		<div class="container navbar-inner">
			<a href="/" class="navbar-brand">
				<span class="navbar-logo">PgQueryNarrative</span>
				<span class="navbar-tagline">Turn SQL into narratives</span>
			</a>
			<div class="navbar-links">
				<a href="/" class="nav-link">Home</a>
				<a href="/query" class="nav-link">Run Query</a>
				<a href="/saved" class="nav-link">Saved</a>
				<a href="/reports" class="nav-link">Reports</a>
			</div>
		</div>
	</nav>
	<main class="container">
		<div class="reports-page">
			<h2 class="page-title">Generate Report</h2>
			<p class="page-subtitle">Run a query and generate an AI narrative with suggested charts.</p>
			<form hx-post="/web/reports/generate" hx-target="#report-result" hx-swap="innerHTML" class="query-form">
				<div class="form-group">
					<label for="sql">SQL Query</label>
					<textarea id="sql" name="sql" rows="10" required placeholder="SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category"></textarea>
				</div>
				<button type="submit" class="btn btn-primary">Generate Narrative Report</button>
			</form>
			<div id="report-result" class="results"></div>
		</div>
	</main>
	<footer>
		<div class="container">
			<p>&copy; 2026 PgQueryNarrative</p>
		</div>
	</footer>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html)
}

func FormatQueryResultsHTML(result *queries.RunQueryResult) string {
	rows := len(result.Rows)
	cols := len(result.Columns)
	estSize := 512 + rows*cols*32
	if estSize < 4096 {
		estSize = 4096
	}
	var sb strings.Builder
	sb.Grow(estSize)
	sb.WriteString("<div class='query-results'>")
	sb.WriteString("<h3>Query Results</h3>")
	sb.WriteString("<div class='result-info'>")
	sb.WriteString("<p><strong>Rows:</strong> ")
	sb.WriteString(strconv.Itoa(int(result.RowCount)))
	sb.WriteString("</p>")
	if result.ExecutionTimeMs > 0 {
		execTime := float64(result.ExecutionTimeMs) / 1000.0
		if execTime < 1.0 {
			sb.WriteString("<p><strong>Execution Time:</strong> ")
			sb.WriteString(strconv.Itoa(int(result.ExecutionTimeMs)))
			sb.WriteString("ms</p>")
		} else {
			sb.WriteString("<p><strong>Execution Time:</strong> ")
			sb.WriteString(fmt.Sprintf("%.2f", execTime))
			sb.WriteString("s</p>")
		}
	}
	if result.Limit > 0 {
		sb.WriteString("<p><strong>Limit:</strong> ")
		sb.WriteString(strconv.Itoa(int(result.Limit)))
		sb.WriteString(" rows</p>")
	}
	sb.WriteString("</div>")

	if len(result.PeriodComparison) > 0 {
		titleAttr := "Compares the latest period in your result to the period before it."
		sb.WriteString("<div class='period-comparison' title='")
		sb.WriteString(template.HTMLEscapeString(titleAttr))
		sb.WriteString("'><h4>Vs previous period")
		if result.PeriodCurrentLabel != nil && result.PeriodPreviousLabel != nil && *result.PeriodCurrentLabel != "" && *result.PeriodPreviousLabel != "" {
			sb.WriteString(" <span class='period-labels'>(")
			sb.WriteString(template.HTMLEscapeString(*result.PeriodCurrentLabel))
			sb.WriteString(" vs ")
			sb.WriteString(template.HTMLEscapeString(*result.PeriodPreviousLabel))
			sb.WriteString(")</span>")
		}
		sb.WriteString("</h4><p class='period-comparison-hint'>Latest period in result vs the one before it.</p><ul class='period-comparison-list'>")
		for _, p := range result.PeriodComparison {
			if p == nil {
				continue
			}
			sb.WriteString("<li class='period-item'><strong>")
			sb.WriteString(template.HTMLEscapeString(p.Measure))
			sb.WriteString("</strong>: ")
			sb.WriteString(formatFloatWithCommas(p.Current))
			if p.Previous != nil {
				sb.WriteString(" (prev: ")
				sb.WriteString(formatFloatWithCommas(*p.Previous))
				sb.WriteString(")")
			}
			if p.ChangePercentage != nil {
				sb.WriteString(" <span class='period-change trend-")
				sb.WriteString(template.HTMLEscapeString(p.Trend))
				sb.WriteString("'>")
				if *p.ChangePercentage >= 0 {
					sb.WriteString("+")
				}
				sb.WriteString(fmt.Sprintf("%.1f%%", *p.ChangePercentage))
				sb.WriteString("</span>")
			}
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul></div>")
	}

	if len(result.ChartSuggestions) > 0 {
		sb.WriteString("<div class='chart-suggestions'><h4>Suggested charts</h4><ul class='suggestion-list' id='chart-suggestions-list'>")
		for _, s := range result.ChartSuggestions {
			if s == nil {
				continue
			}
			sb.WriteString("<li><button type='button' class='chart-type-btn' data-chart-type='")
			sb.WriteString(template.HTMLEscapeString(s.ChartType))
			sb.WriteString("' title='")
			sb.WriteString(template.HTMLEscapeString(s.Reason))
			sb.WriteString("'>")
			sb.WriteString(template.HTMLEscapeString(s.Label))
			sb.WriteString("</button></li>")
		}
		sb.WriteString("</ul></div>")
	}

	chartData := chartDataFromResult(result)
	if chartData != nil {
		jsonBytes, _ := json.Marshal(chartData)
		chartSelectOptions := buildChartSelectOptions(result.ChartSuggestions)
		sb.WriteString("<div class='chart-area' data-chart-data='")
		sb.WriteString(template.HTMLEscapeString(string(jsonBytes)))
		sb.WriteString("'><div class='chart-toolbar'><span class='chart-label'>Chart:</span> <select id='chart-type-select' class='chart-select'><option value=''>—</option>")
		sb.WriteString(chartSelectOptions)
		sb.WriteString("</select></div><div class='chart-canvas-wrap'><canvas id='result-chart' width='400' height='250'></canvas></div></div>")
	}

	if rows == 0 {
		sb.WriteString("<div class='no-results'><p>No results returned.</p><p class='hint'>Try adjusting your query or check if the data exists.</p></div>")
		sb.WriteString("</div>")
		return sb.String()
	}

	sb.WriteString("<div class='table-container'>")
	sb.WriteString("<table class='results-table'><thead><tr>")
	for _, col := range result.Columns {
		sb.WriteString("<th>")
		sb.WriteString(template.HTMLEscapeString(col.Name))
		if col.Type != "" {
			sb.WriteString(" <span class='col-type'>(")
			sb.WriteString(template.HTMLEscapeString(col.Type))
			sb.WriteString(")</span>")
		}
		sb.WriteString("</th>")
	}
	sb.WriteString("</tr></thead><tbody>")
	for _, row := range result.Rows {
		sb.WriteString("<tr>")
		for _, val := range row {
			sb.WriteString("<td>")
			sb.WriteString(template.HTMLEscapeString(formatValue(val)))
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>")
	}
	sb.WriteString("</tbody></table>")
	sb.WriteString("</div>")
	sb.WriteString("</div>")
	return sb.String()
}

// buildChartSelectOptions returns <option> HTML for the chart-type dropdown,
// built from API suggestions (excluding table). Uses defaults (bar, line, pie) if none.
func buildChartSelectOptions(suggestions []*queries.ChartSuggestion) string {
	defaults := []struct{ value, label string }{
		{"bar", "Bar"},
		{"line", "Line"},
		{"pie", "Pie"},
	}
	seen := make(map[string]bool)
	var opts []struct{ value, label string }
	for _, s := range suggestions {
		if s == nil || s.ChartType == "table" {
			continue
		}
		if seen[s.ChartType] {
			continue
		}
		seen[s.ChartType] = true
		label := s.Label
		if label == "" {
			label = s.ChartType
		}
		opts = append(opts, struct{ value, label string }{s.ChartType, label})
	}
	if len(opts) == 0 {
		opts = defaults
	}
	var sb strings.Builder
	sb.Grow(len(opts) * 64)
	for _, o := range opts {
		sb.WriteString("<option value='")
		sb.WriteString(template.HTMLEscapeString(o.value))
		sb.WriteString("'>")
		sb.WriteString(template.HTMLEscapeString(o.label))
		sb.WriteString("</option>")
	}
	return sb.String()
}

// chartDataForJS is the struct serialized to data-chart-data for the chart UI.
type chartDataForJS struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

// chartDataFromResult builds JSON-serializable chart data from a query result (max 200 rows).
func chartDataFromResult(result *queries.RunQueryResult) *chartDataForJS {
	if result == nil || len(result.Columns) == 0 || len(result.Rows) == 0 {
		return nil
	}
	maxRows := 200
	if len(result.Rows) < maxRows {
		maxRows = len(result.Rows)
	}
	cols := make([]string, len(result.Columns))
	for i, c := range result.Columns {
		cols[i] = c.Name
	}
	rows := make([][]interface{}, maxRows)
	for i := 0; i < maxRows; i++ {
		row := make([]interface{}, len(result.Rows[i]))
		for j, val := range result.Rows[i] {
			row[j] = rowValueToScalar(val)
		}
		rows[i] = row
	}
	return &chartDataForJS{Columns: cols, Rows: rows}
}

func rowValueToScalar(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	if n, ok := val.(pgtype.Numeric); ok {
		if !n.Valid || n.Int == nil {
			return nil
		}
		f, _ := n.Int.Float64()
		exp := int(n.Exp)
		for exp < 0 {
			f /= 10
			exp++
		}
		for exp > 0 {
			f *= 10
			exp--
		}
		return f
	}
	switch v := val.(type) {
	case float64, int, int32, int64:
		return v
	case float32:
		return float64(v)
	case string:
		return v
	default:
		return fmt.Sprint(val)
	}
}

// formatFloat formats a float64 for display (e.g. in period comparison).
func formatFloat(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// formatFloatWithCommas formats a float64 with thousands separators (e.g. 45793291.51 -> "45,793,291.51").
func formatFloatWithCommas(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	idx := strings.Index(s, ".")
	if idx == -1 {
		idx = len(s)
	}
	integerPart := s[:idx]
	var b strings.Builder
	for i, c := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteRune(c)
	}
	if idx < len(s) {
		b.WriteString(s[idx:])
	}
	return b.String()
}

func formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	// Handle PostgreSQL numeric types
	if numeric, ok := val.(pgtype.Numeric); ok {
		if !numeric.Valid {
			return "NULL"
		}
		// Convert pgtype.Numeric to string representation
		if numeric.Int != nil {
			// For decimal numerics, convert using Exp
			if numeric.Exp != 0 {
				// Convert big.Int to float64 and apply exponent
				valFloat, _ := numeric.Int.Float64()
				exp := int(numeric.Exp)
				for exp < 0 {
					valFloat = valFloat / 10.0
					exp++
				}
				for exp > 0 {
					valFloat = valFloat * 10.0
					exp--
				}
				return fmt.Sprintf("%.2f", valFloat)
			}
			// For integer-like numerics
			return numeric.Int.String()
		}
		return "0"
	}

	switch v := val.(type) {
	case string:
		return template.HTMLEscapeString(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return fmt.Sprintf("%.2f", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// For complex types, try to convert via fmt.Sprintf
		return template.HTMLEscapeString(fmt.Sprintf("%v", v))
	}
}

func (h *Handlers) RunQuery(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			errorHTML := "<div class='error-message'><strong>Error:</strong> An unexpected error occurred. Please try again.</div>"
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, errorHTML)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		errorHTML := "<div class='error-message'><strong>Error:</strong> Failed to parse form data. Please try again.</div>"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errorHTML)
		return
	}

	sql := strings.TrimSpace(r.FormValue("sql"))
	if sql == "" {
		errorHTML := "<div class='error-message'><strong>Error:</strong> SQL query cannot be empty.</div>"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errorHTML)
		return
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limitStr = "100"
	}
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	ctx := r.Context()
	payload := &queries.RunQueryPayload{
		SQL:   sql,
		Limit: int32(limit),
	}

	result, err := h.queriesEndpoints.Run(ctx, payload)
	if err != nil {
		errorHTML := "<div class='error-message'>"
		errorHTML += "<strong>Error:</strong> "
		if validationErr, ok := err.(*queries.ValidationError); ok {
			errorHTML += template.HTMLEscapeString(validationErr.Message)
			if validationErr.Code != nil {
				errorHTML += " <span class='error-code'>(" + template.HTMLEscapeString(*validationErr.Code) + ")</span>"
			}
		} else {
			errorHTML += template.HTMLEscapeString(err.Error())
		}
		errorHTML += "</div>"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errorHTML)
		return
	}

	queryResult, ok := result.(*queries.RunQueryResult)
	if !ok {
		errorHTML := "<div class='error-message'><strong>Error:</strong> Invalid response type. Please try again.</div>"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, errorHTML)
		return
	}

	html := FormatQueryResultsHTML(queryResult)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, html)
}

// GenerateReport handles form submission from web UI and converts to API call
func (h *Handlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	sql := r.FormValue("sql")

	// Call the reports service
	ctx := r.Context()
	payload := &reports.GenerateReportPayload{
		SQL: sql,
	}

	report, err := h.reportsEndpoints.Generate(ctx, payload)
	if err != nil {
		// Format error for HTML display
		errorHTML := "<div class='error-message'><strong>Error:</strong> "
		if validationErr, ok := err.(*reports.ValidationError); ok {
			errorHTML += template.HTMLEscapeString(validationErr.Message)
		} else if llmErr, ok := err.(*reports.LLMError); ok {
			errorHTML += template.HTMLEscapeString(llmErr.Message)
		} else {
			errorHTML += template.HTMLEscapeString(err.Error())
		}
		errorHTML += "</div>"
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, errorHTML)
		return
	}

	// Type assert the report
	reportResult, ok := report.(*reports.Report)
	if !ok {
		errorHTML := "<div class='error-message'><strong>Error:</strong> Invalid response type</div>"
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, errorHTML)
		return
	}

	// Format report as HTML
	html := FormatReportHTML(reportResult)
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html)
}

// FormatReportHTML formats report results as HTML with model badge and improved layout.
func FormatReportHTML(report *reports.Report) string {
	var sb strings.Builder
	sb.Grow(2048)
	sb.WriteString("<div class=\"report-results\">")
	// Report id and model on one line
	sb.WriteString("<div class=\"report-meta\">")
	sb.WriteString("<span class=\"report-id\">Report: ")
	sb.WriteString(template.HTMLEscapeString(report.ID))
	sb.WriteString("</span>")
	if report.LlmProvider != "" || report.LlmModel != "" {
		modelLabel := report.LlmProvider
		if report.LlmModel != "" {
			if modelLabel != "" {
				modelLabel += " / "
			}
			modelLabel += report.LlmModel
		}
		sb.WriteString("<span class=\"report-model\">Model: ")
		sb.WriteString(template.HTMLEscapeString(modelLabel))
		sb.WriteString("</span>")
	}
	sb.WriteString("</div>")

	if len(report.ChartSuggestions) > 0 {
		sb.WriteString("<div class=\"chart-suggestions\"><h4 class=\"report-section-title\">Suggested charts</h4><ul class=\"suggestion-list report-list\">")
		for _, s := range report.ChartSuggestions {
			if s == nil {
				continue
			}
			sb.WriteString("<li><strong>")
			sb.WriteString(template.HTMLEscapeString(s.Label))
			sb.WriteString("</strong> — ")
			sb.WriteString(template.HTMLEscapeString(s.Reason))
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul></div>")
	}

	if report.Metrics != nil && len(report.Metrics.TimeSeries) > 0 {
		sb.WriteString("<div class=\"period-comparison\" title=\"Compares the latest period in the result to the period before it.\"><h4 class=\"report-section-title\">Vs previous period")
		if report.Metrics.PeriodCurrentLabel != nil && report.Metrics.PeriodPreviousLabel != nil && *report.Metrics.PeriodCurrentLabel != "" && *report.Metrics.PeriodPreviousLabel != "" {
			sb.WriteString(" <span class=\"period-labels\">(")
			sb.WriteString(template.HTMLEscapeString(*report.Metrics.PeriodCurrentLabel))
			sb.WriteString(" vs ")
			sb.WriteString(template.HTMLEscapeString(*report.Metrics.PeriodPreviousLabel))
			sb.WriteString(")</span>")
		}
		sb.WriteString("</h4><p class=\"period-comparison-hint\">Latest period vs the one before it.</p><ul class=\"period-comparison-list\">")
		measures := make([]string, 0, len(report.Metrics.TimeSeries))
		for m := range report.Metrics.TimeSeries {
			measures = append(measures, m)
		}
		sort.Strings(measures)
		for _, measure := range measures {
			ts := report.Metrics.TimeSeries[measure]
			if ts == nil {
				continue
			}
			sb.WriteString("<li class=\"period-item\"><strong>")
			sb.WriteString(template.HTMLEscapeString(measure))
			sb.WriteString("</strong>: ")
			sb.WriteString(formatFloatWithCommas(ts.CurrentPeriod))
			if ts.PreviousPeriod != nil {
				sb.WriteString(" (prev: ")
				sb.WriteString(formatFloatWithCommas(*ts.PreviousPeriod))
				sb.WriteString(")")
			}
			if ts.ChangePercentage != nil {
				sb.WriteString(" <span class=\"period-change trend-")
				sb.WriteString(template.HTMLEscapeString(ts.Trend))
				sb.WriteString("\">")
				if *ts.ChangePercentage >= 0 {
					sb.WriteString("+")
				}
				sb.WriteString(fmt.Sprintf("%.1f%%", *ts.ChangePercentage))
				sb.WriteString("</span>")
			}
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul></div>")

		// Advanced metrics: trend summary, anomalies, period history per measure
		sb.WriteString("<div class=\"advanced-metrics\">")
		for _, measure := range measures {
			ts := report.Metrics.TimeSeries[measure]
			if ts == nil {
				continue
			}
			hasAdvanced := (ts.TrendSummary != nil && ts.TrendSummary.Summary != "") ||
				len(ts.Anomalies) > 0 || len(ts.Periods) > 0
			if !hasAdvanced {
				continue
			}
			sb.WriteString("<div class=\"advanced-metric-measure\"><h5 class=\"report-measure-title\">")
			sb.WriteString(template.HTMLEscapeString(measure))
			sb.WriteString("</h5>")

			if ts.TrendSummary != nil && ts.TrendSummary.Summary != "" {
				sb.WriteString("<p class=\"trend-summary\"><span class=\"trend-badge trend-")
				sb.WriteString(template.HTMLEscapeString(ts.TrendSummary.Direction))
				sb.WriteString("\">")
				sb.WriteString(template.HTMLEscapeString(ts.TrendSummary.Direction))
				sb.WriteString("</span> ")
				sb.WriteString(template.HTMLEscapeString(ts.TrendSummary.Summary))
				if ts.MovingAverage != nil {
					sb.WriteString(" Moving avg (latest): ")
					sb.WriteString(formatFloatWithCommas(*ts.MovingAverage))
					sb.WriteString(".")
				}
				sb.WriteString("</p>")
			}
			if ts.NextPeriodForecast != nil {
				sb.WriteString("<p class=\"next-period-forecast\"><strong>Next period forecast:</strong> ")
				sb.WriteString(formatFloatWithCommas(*ts.NextPeriodForecast))
				sb.WriteString("</p>")
			}

			if len(ts.Anomalies) > 0 {
				sb.WriteString("<div class=\"anomalies-list\"><strong>Anomalies:</strong> <ul class=\"anomaly-list\">")
				for _, a := range ts.Anomalies {
					if a == nil {
						continue
					}
					sb.WriteString("<li><span class=\"anomaly-period\">")
					sb.WriteString(template.HTMLEscapeString(a.PeriodLabel))
					sb.WriteString("</span> ")
					sb.WriteString(formatFloatWithCommas(a.Value))
					sb.WriteString(" — ")
					sb.WriteString(template.HTMLEscapeString(a.Reason))
					sb.WriteString("</li>")
				}
				sb.WriteString("</ul></div>")
			}

			if len(ts.Periods) > 0 {
				sb.WriteString("<details class=\"period-history-details\"><summary>Period history (")
				sb.WriteString(strconv.Itoa(len(ts.Periods)))
				sb.WriteString(" periods)</summary><table class=\"period-history-table\"><thead><tr><th>Period</th><th>Value</th></tr></thead><tbody>")
				// Show last 15 to keep UI compact
				start := 0
				if len(ts.Periods) > 15 {
					start = len(ts.Periods) - 15
				}
				for i := start; i < len(ts.Periods); i++ {
					p := ts.Periods[i]
					if p == nil {
						continue
					}
					sb.WriteString("<tr><td>")
					sb.WriteString(template.HTMLEscapeString(p.Label))
					sb.WriteString("</td><td>")
					sb.WriteString(formatFloatWithCommas(p.Value))
					sb.WriteString("</td></tr>")
				}
				sb.WriteString("</tbody></table></details>")
			}
			sb.WriteString("</div>")
		}
		sb.WriteString("</div>")
	}

	if report.Metrics != nil && len(report.Metrics.PerfSuggestions) > 0 {
		sb.WriteString("<div class=\"perf-suggestions\"><h4 class=\"report-section-title\">Performance suggestions</h4><ul class=\"report-list\">")
		for _, s := range report.Metrics.PerfSuggestions {
			sb.WriteString("<li>")
			sb.WriteString(template.HTMLEscapeString(s))
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul></div>")
	}

	if report.Metrics != nil && len(report.Metrics.DataQuality) > 0 {
		sb.WriteString("<div class=\"data-quality\"><h4 class=\"report-section-title\">Data quality</h4><table class=\"data-quality-table\"><thead><tr><th>Column</th><th>Nulls</th><th>Null %</th><th>Distinct</th><th>Rows</th></tr></thead><tbody>")
		dqCols := make([]string, 0, len(report.Metrics.DataQuality))
		for c := range report.Metrics.DataQuality {
			dqCols = append(dqCols, c)
		}
		sort.Strings(dqCols)
		for _, col := range dqCols {
			q := report.Metrics.DataQuality[col]
			if q == nil {
				continue
			}
			sb.WriteString("<tr><td>")
			sb.WriteString(template.HTMLEscapeString(col))
			sb.WriteString("</td><td>")
			sb.WriteString(strconv.Itoa(int(q.NullCount)))
			sb.WriteString("</td><td>")
			sb.WriteString(fmt.Sprintf("%.1f", q.NullPct))
			sb.WriteString("%</td><td>")
			sb.WriteString(strconv.Itoa(int(q.DistinctCount)))
			sb.WriteString("</td><td>")
			sb.WriteString(strconv.Itoa(int(q.TotalRows)))
			sb.WriteString("</td></tr>")
		}
		sb.WriteString("</tbody></table></div>")
	}

	if report.Narrative != nil {
		sb.WriteString("<div class=\"report-narrative narrative-content\">")
		if report.Narrative.Headline != "" {
			sb.WriteString("<h3 class=\"report-headline\">")
			sb.WriteString(template.HTMLEscapeString(report.Narrative.Headline))
			sb.WriteString("</h3>")
		}
		if len(report.Narrative.Takeaways) > 0 {
			sb.WriteString("<h4 class=\"report-section-title\">Key takeaways</h4><ul class=\"report-takeaways\">")
			for _, takeaway := range report.Narrative.Takeaways {
				sb.WriteString("<li>")
				sb.WriteString(template.HTMLEscapeString(takeaway))
				sb.WriteString("</li>")
			}
			sb.WriteString("</ul>")
		}
		if len(report.Narrative.Drivers) > 0 {
			sb.WriteString("<h4 class=\"report-section-title\">Drivers</h4><ul class=\"report-list\">")
			for _, d := range report.Narrative.Drivers {
				sb.WriteString("<li>")
				sb.WriteString(template.HTMLEscapeString(d))
				sb.WriteString("</li>")
			}
			sb.WriteString("</ul>")
		}
		if len(report.Narrative.Limitations) > 0 {
			sb.WriteString("<h4 class=\"report-section-title\">Limitations</h4><ul class=\"report-list\">")
			for _, l := range report.Narrative.Limitations {
				sb.WriteString("<li>")
				sb.WriteString(template.HTMLEscapeString(l))
				sb.WriteString("</li>")
			}
			sb.WriteString("</ul>")
		}
		if len(report.Narrative.Recommendations) > 0 {
			sb.WriteString("<h4 class=\"report-section-title\">Recommendations</h4><ul class=\"report-list\">")
			for _, r := range report.Narrative.Recommendations {
				sb.WriteString("<li>")
				sb.WriteString(template.HTMLEscapeString(r))
				sb.WriteString("</li>")
			}
			sb.WriteString("</ul>")
		}
		sb.WriteString("</div>")
	}

	sb.WriteString("</div>")
	return sb.String()
}
