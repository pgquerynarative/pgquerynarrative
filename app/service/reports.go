package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
	"github.com/pgquerynarrative/pgquerynarrative/app/apilog"
	"github.com/pgquerynarrative/pgquerynarrative/app/charts"
	"github.com/pgquerynarrative/pgquerynarrative/app/config"
	"github.com/pgquerynarrative/pgquerynarrative/app/debuglog"
	"github.com/pgquerynarrative/pgquerynarrative/app/embedding"
	"github.com/pgquerynarrative/pgquerynarrative/app/llm"
	"github.com/pgquerynarrative/pgquerynarrative/app/metrics"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
	"github.com/pgquerynarrative/pgquerynarrative/app/story"
)

type ReportsService struct {
	readOnlyPool   *pgxpool.Pool
	appPool        *pgxpool.Pool
	runner         *queryrunner.Runner
	llmClient      llm.Client
	generator      *story.Generator
	metricsOpts    *metrics.Options
	embedder       embedding.Embedder
	embeddingStore *embedding.Store
}

func NewReportsService(readOnlyPool, appPool *pgxpool.Pool, runner *queryrunner.Runner, llmClient llm.Client, metricsCfg config.MetricsConfig) *ReportsService {
	opts := metricsOptionsFromConfig(metricsCfg)
	return &ReportsService{
		readOnlyPool: readOnlyPool,
		appPool:      appPool,
		runner:       runner,
		llmClient:    llmClient,
		generator:    story.NewGenerator(llmClient),
		metricsOpts:  opts,
	}
}

// NewReportsServiceWithRAG is like NewReportsService but enables RAG: similar past
// queries are retrieved and added to the narrative prompt when generating reports.
func NewReportsServiceWithRAG(readOnlyPool, appPool *pgxpool.Pool, runner *queryrunner.Runner, llmClient llm.Client, metricsCfg config.MetricsConfig, embedder embedding.Embedder, embeddingStore *embedding.Store) *ReportsService {
	opts := metricsOptionsFromConfig(metricsCfg)
	return &ReportsService{
		readOnlyPool:   readOnlyPool,
		appPool:        appPool,
		runner:         runner,
		llmClient:      llmClient,
		generator:      story.NewGenerator(llmClient),
		metricsOpts:    opts,
		embedder:       embedder,
		embeddingStore: embeddingStore,
	}
}

func (s *ReportsService) Generate(ctx context.Context, payload *reports.GenerateReportPayload) (*reports.Report, error) {
	debuglog.Log("report generation started")
	// Execute query
	queryResult, err := s.runner.Run(ctx, payload.SQL, 1000)
	if err != nil {
		kind, userMsg := ClassifyRunError(err)
		if kind == RunErrorTimeout {
			apilog.ValidationError("generate", "timeout_error", err.Error())
			return nil, &reports.ValidationError{Name: "timeout_error", Message: userMsg, Code: strPtr("TIMEOUT_ERROR")}
		}
		apilog.ValidationError("generate", "validation_error", err.Error())
		return nil, &reports.ValidationError{Name: "validation_error", Message: userMsg, Code: strPtr("VALIDATION_ERROR")}
	}

	// Extract column names and types
	columnNames := make([]string, len(queryResult.Columns))
	columnTypes := make([]string, len(queryResult.Columns))
	for i, col := range queryResult.Columns {
		columnNames[i] = col.Name
		columnTypes[i] = col.Type
	}

	// Chart suggestions from result shape
	chartSuggestions := suggestToReports(charts.Suggest(columnNames, columnTypes, queryResult.Rows))

	// Profile columns
	profiles := metrics.ProfileColumns(columnNames, queryResult.Rows)

	// Calculate metrics
	calcMetrics := metrics.CalculateMetrics(columnNames, queryResult.Rows, profiles, s.metricsOpts)
	calcMetrics.PerfSuggestions = BuildPerfSuggestions(queryResult)

	// Optional RAG: retrieve similar past queries and add to prompt context
	var similarContext string
	if s.embedder != nil && s.embeddingStore != nil {
		if vec, err := s.embedder.Embed(ctx, payload.SQL); err == nil {
			if similar, err := s.embeddingStore.FindSimilar(ctx, vec, 3); err == nil && len(similar) > 0 {
				const maxSQLLen = 200
				var b strings.Builder
				for _, q := range similar {
					b.WriteString("- ")
					b.WriteString(q.Name)
					b.WriteString(": ")
					sql := q.SQL
					if len(sql) > maxSQLLen {
						sql = sql[:maxSQLLen] + "..."
					}
					b.WriteString(sql)
					b.WriteString("\n")
				}
				similarContext = strings.TrimSpace(b.String())
			}
		}
	}

	// Generate narrative
	debuglog.Log("calling LLM for narrative generation")
	narrative, err := s.generator.Generate(ctx, payload.SQL, columnNames, queryResult.Rows, calcMetrics, similarContext)
	if err != nil {
		llmMsg := err.Error()
		apilog.LLMError(llmMsg)
		// Fallback keeps report generation and Ask usable even when the model is slow
		// or temporarily unavailable.
		narrative = buildFallbackNarrative(queryResult.RowCount, calcMetrics.PerfSuggestions)
	}

	// Convert metrics to API format
	metricsData := ConvertMetrics(calcMetrics)

	// Store report in database
	debuglog.Log("storing report in database")
	reportID, err := s.storeReport(ctx, payload, narrative, calcMetrics, queryResult)
	if err != nil {
		return nil, err
	}
	apilog.Request("generate", "report_id="+reportID)

	// Convert narrative to API format
	narrativeData := &reports.NarrativeContent{
		Headline:        narrative.Headline,
		Takeaways:       narrative.Takeaways,
		Drivers:         narrative.Drivers,
		Limitations:     narrative.Limitations,
		Recommendations: narrative.Recommendations,
	}

	return &reports.Report{
		ID:               reportID,
		SavedQueryID:     payload.SavedQueryID,
		SQL:              payload.SQL,
		Narrative:        narrativeData,
		Metrics:          metricsData,
		ChartSuggestions: chartSuggestions,
		CreatedAt:        time.Now().Format(time.RFC3339),
		LlmModel:         s.llmClient.Name(),
		LlmProvider:      s.llmClient.Name(),
	}, nil
}

func buildFallbackNarrative(rowCount int, perfSuggestions []string) *story.NarrativeContent {
	n := &story.NarrativeContent{
		Headline: "Report generated without LLM narrative",
		Takeaways: []string{
			"Query executed successfully and returned " + strconv.Itoa(rowCount) + " rows.",
		},
		Limitations: []string{
			"Natural-language narrative is unavailable right now; showing metrics and raw results instead.",
		},
	}
	if len(perfSuggestions) > 0 {
		n.Recommendations = append(n.Recommendations, perfSuggestions...)
	}
	return n
}

// suggestToReports converts charts.Suggestion slice to reports API type.
func suggestToReports(in []charts.Suggestion) []*reports.ChartSuggestion {
	if len(in) == 0 {
		return nil
	}
	out := make([]*reports.ChartSuggestion, len(in))
	for i := range in {
		out[i] = &reports.ChartSuggestion{
			ChartType: in[i].ChartType,
			Label:     in[i].Label,
			Reason:    in[i].Reason,
		}
	}
	return out
}

func (s *ReportsService) storeReport(ctx context.Context, payload *reports.GenerateReportPayload, narrative *story.NarrativeContent, calcMetrics *metrics.Metrics, queryResult *queryrunner.Result) (string, error) {
	narrativeJSON, _ := json.Marshal(narrative)
	metricsJSON, _ := json.Marshal(calcMetrics)
	statsJSON, _ := json.Marshal(map[string]interface{}{
		"execution_time_ms": queryResult.ExecutionTimeMs,
		"row_count":         queryResult.RowCount,
	})

	var reportID string
	err := s.appPool.QueryRow(ctx, `
		INSERT INTO app.reports (
			saved_query_id, sql, narrative_md, narrative_json, metrics, stats,
			llm_model, llm_provider, success
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, payload.SavedQueryID, payload.SQL, narrative.Headline, narrativeJSON, metricsJSON, statsJSON,
		s.llmClient.Name(), s.llmClient.Name(), true).Scan(&reportID)

	return reportID, err
}

func (s *ReportsService) Get(ctx context.Context, payload *reports.GetPayload) (*reports.Report, error) {
	row := s.appPool.QueryRow(ctx, `
		SELECT id, saved_query_id, sql, narrative_json, metrics, created_at, llm_model, llm_provider
		FROM app.reports
		WHERE id = $1
	`, payload.ID)

	var report reports.Report
	var savedQueryID sql.NullString
	var narrativeJSON []byte
	var metricsJSON []byte
	var createdAt time.Time

	err := row.Scan(&report.ID, &savedQueryID, &report.SQL, &narrativeJSON, &metricsJSON, &createdAt, &report.LlmModel, &report.LlmProvider)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &reports.NotFoundError{
				Name:    "not_found",
				Message: "report not found",
				Code:    strPtr("NOT_FOUND"),
			}
		}
		return nil, err
	}

	if savedQueryID.Valid {
		report.SavedQueryID = &savedQueryID.String
	}

	var narrative story.NarrativeContent
	if err := json.Unmarshal(narrativeJSON, &narrative); err == nil {
		report.Narrative = &reports.NarrativeContent{
			Headline:        narrative.Headline,
			Takeaways:       narrative.Takeaways,
			Drivers:         narrative.Drivers,
			Limitations:     narrative.Limitations,
			Recommendations: narrative.Recommendations,
		}
	}

	var calcMetrics metrics.Metrics
	if err := json.Unmarshal(metricsJSON, &calcMetrics); err == nil {
		report.Metrics = ConvertMetrics(&calcMetrics)
	}

	report.CreatedAt = createdAt.Format(time.RFC3339)

	return &report, nil
}

func (s *ReportsService) List(ctx context.Context, payload *reports.ListPayload) (*reports.ReportList, error) {
	limit := int(payload.Limit)
	offset := int(payload.Offset)

	var rows pgx.Rows
	var err error

	if payload.SavedQueryID != nil && *payload.SavedQueryID != "" {
		rows, err = s.appPool.Query(ctx, `
			SELECT id, saved_query_id, sql, narrative_json, metrics, created_at, llm_model, llm_provider
			FROM app.reports
			WHERE saved_query_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, *payload.SavedQueryID, limit, offset)
	} else {
		rows, err = s.appPool.Query(ctx, `
			SELECT id, saved_query_id, sql, narrative_json, metrics, created_at, llm_model, llm_provider
			FROM app.reports
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*reports.Report, 0, limit)
	for rows.Next() {
		var report reports.Report
		var savedQueryID sql.NullString
		var narrativeJSON []byte
		var metricsJSON []byte
		var createdAt time.Time

		if err := rows.Scan(&report.ID, &savedQueryID, &report.SQL, &narrativeJSON, &metricsJSON, &createdAt, &report.LlmModel, &report.LlmProvider); err != nil {
			return nil, err
		}

		if savedQueryID.Valid {
			report.SavedQueryID = &savedQueryID.String
		}

		var narrative story.NarrativeContent
		if err := json.Unmarshal(narrativeJSON, &narrative); err == nil {
			report.Narrative = &reports.NarrativeContent{
				Headline:        narrative.Headline,
				Takeaways:       narrative.Takeaways,
				Drivers:         narrative.Drivers,
				Limitations:     narrative.Limitations,
				Recommendations: narrative.Recommendations,
			}
		}

		var calcMetrics metrics.Metrics
		if err := json.Unmarshal(metricsJSON, &calcMetrics); err == nil {
			report.Metrics = ConvertMetrics(&calcMetrics)
		}

		report.CreatedAt = createdAt.Format(time.RFC3339)
		items = append(items, &report)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &reports.ReportList{
		Items:  items,
		Limit:  payload.Limit,
		Offset: payload.Offset,
	}, nil
}

// BuildPerfSuggestions returns performance suggestions from query result. Exported for testing.
func BuildPerfSuggestions(r *queryrunner.Result) []string {
	var suggestions []string
	if r.ExecutionTimeMs > 2000 {
		suggestions = append(suggestions, "Query took over 2s; consider adding filters or indexes.")
	}
	if r.RowCount >= 1000 {
		suggestions = append(suggestions, "Result set is large (limit applied); consider narrowing date range or dimensions.")
	}
	return suggestions
}

// ConvertMetrics converts app metrics to API type. Exported for testing.
func ConvertMetrics(m *metrics.Metrics) *reports.MetricsData {
	aggregates := make(map[string]*reports.AggregateData, len(m.Aggregates))
	for col, agg := range m.Aggregates {
		count := int32(agg.Count)
		ad := &reports.AggregateData{
			Sum:   agg.Sum,
			Avg:   agg.Avg,
			Min:   agg.Min,
			Max:   agg.Max,
			Count: &count,
		}
		if agg.StdDev != nil {
			ad.StdDev = agg.StdDev
		}
		aggregates[col] = ad
	}

	topCategories := make(map[string][]*reports.TopCategoryData, len(m.TopCategories))
	for col, cats := range m.TopCategories {
		categoryData := make([]*reports.TopCategoryData, len(cats))
		for i, cat := range cats {
			categoryData[i] = &reports.TopCategoryData{
				Category:   cat.Category,
				Value:      cat.Value,
				Percentage: cat.Percentage,
			}
		}
		topCategories[col] = categoryData
	}

	timeSeries := make(map[string]*reports.TimeSeriesData, len(m.TimeSeries))
	for col, ts := range m.TimeSeries {
		tsData := &reports.TimeSeriesData{
			CurrentPeriod: ts.CurrentPeriod,
			Trend:         ts.Trend,
		}
		if ts.PreviousPeriod != nil {
			tsData.PreviousPeriod = ts.PreviousPeriod
		}
		if ts.Change != nil {
			tsData.Change = ts.Change
		}
		if ts.ChangePercentage != nil {
			tsData.ChangePercentage = ts.ChangePercentage
		}
		if len(ts.Periods) > 0 {
			tsData.Periods = make([]*reports.PeriodPointData, len(ts.Periods))
			for i := range ts.Periods {
				tsData.Periods[i] = &reports.PeriodPointData{
					Label: ts.Periods[i].Label,
					Value: ts.Periods[i].Value,
				}
			}
		}
		if ts.MovingAverage != nil {
			tsData.MovingAverage = ts.MovingAverage
		}
		if len(ts.Anomalies) > 0 {
			tsData.Anomalies = make([]*reports.AnomalyPointData, len(ts.Anomalies))
			for i := range ts.Anomalies {
				tsData.Anomalies[i] = &reports.AnomalyPointData{
					PeriodLabel: ts.Anomalies[i].PeriodLabel,
					Value:       ts.Anomalies[i].Value,
					Reason:      ts.Anomalies[i].Reason,
				}
			}
		}
		if ts.TrendSummary != nil {
			pu := int32(ts.TrendSummary.PeriodsUsed)
			tsData.TrendSummary = &reports.TrendSummaryData{
				Direction:   ts.TrendSummary.Direction,
				Summary:     ts.TrendSummary.Summary,
				Slope:       &ts.TrendSummary.Slope,
				PeriodsUsed: &pu,
			}
		}
		if ts.NextPeriodForecast != nil {
			tsData.NextPeriodForecast = ts.NextPeriodForecast
		}
		if ts.ForecastCILower != nil {
			tsData.ForecastCiLower = ts.ForecastCILower
		}
		if ts.ForecastCIUpper != nil {
			tsData.ForecastCiUpper = ts.ForecastCIUpper
		}
		if ts.PredictiveSummary != "" {
			tsData.PredictiveSummary = &ts.PredictiveSummary
		}
		if ts.ExponentialSmoothForecast != nil {
			tsData.ExponentialSmoothForecast = ts.ExponentialSmoothForecast
		}
		if ts.HoltForecast != nil {
			tsData.HoltForecast = ts.HoltForecast
		}
		if ts.SeasonalPeriod != 0 {
			sp := int32(ts.SeasonalPeriod)
			tsData.SeasonalPeriod = &sp
		}
		if ts.SeasonallyAdjustedForecast != nil {
			tsData.SeasonallyAdjustedForecast = ts.SeasonallyAdjustedForecast
		}
		timeSeries[col] = tsData
	}

	correlations := make([]*reports.CorrelationPairData, len(m.Correlations))
	for i := range m.Correlations {
		c := &m.Correlations[i]
		correlations[i] = &reports.CorrelationPairData{
			ColumnA:  c.ColumnA,
			ColumnB:  c.ColumnB,
			Pearson:  c.Pearson,
			Spearman: c.Spearman,
		}
	}

	dataQuality := make(map[string]*reports.ColumnQualityData, len(m.DataQuality))
	for col, q := range m.DataQuality {
		dataQuality[col] = &reports.ColumnQualityData{
			NullCount:     int32(q.NullCount),
			DistinctCount: int32(q.DistinctCount),
			TotalRows:     int32(q.TotalRows),
			NullPct:       q.NullPct,
		}
	}

	cohorts := make([]*reports.CohortMetricData, 0)
	if len(m.Cohorts) > 0 {
		cohorts = make([]*reports.CohortMetricData, len(m.Cohorts))
		for i := range m.Cohorts {
			co := &m.Cohorts[i]
			periods := make([]*reports.CohortPeriodPointData, len(co.Periods))
			for j := range co.Periods {
				periods[j] = &reports.CohortPeriodPointData{
					PeriodLabel: co.Periods[j].PeriodLabel,
					Value:       co.Periods[j].Value,
				}
			}
			cohorts[i] = &reports.CohortMetricData{
				CohortLabel:  co.CohortLabel,
				Periods:      periods,
				RetentionPct: co.RetentionPct,
			}
		}
	}
	out := &reports.MetricsData{
		Aggregates:      aggregates,
		TopCategories:   topCategories,
		TimeSeries:      timeSeries,
		Correlations:    correlations,
		Cohorts:         cohorts,
		DataQuality:     dataQuality,
		PerfSuggestions: m.PerfSuggestions,
	}
	if m.CurrentPeriodLabel != "" {
		out.PeriodCurrentLabel = &m.CurrentPeriodLabel
	}
	if m.PreviousPeriodLabel != "" {
		out.PeriodPreviousLabel = &m.PreviousPeriodLabel
	}
	return out
}
