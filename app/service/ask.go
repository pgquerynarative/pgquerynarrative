package service

import (
	"context"
	"strings"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
	suggestions "github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/app/catalog"
	"github.com/pgquerynarrative/pgquerynarrative/app/llm"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
)

// AskService performs the full NL→SQL→report flow: load schema, generate SQL
// from a natural-language question via LLM, validate and run the query, then
// generate the narrative report.
type AskService struct {
	catalogLoader *catalog.Loader
	llmClient     llm.Client
	validator     *queryrunner.Validator
	reportsSvc    *ReportsService
}

// NewAskService creates an AskService with the given dependencies.
func NewAskService(
	catalogLoader *catalog.Loader,
	llmClient llm.Client,
	validator *queryrunner.Validator,
	reportsSvc *ReportsService,
) *AskService {
	return &AskService{
		catalogLoader: catalogLoader,
		llmClient:     llmClient,
		validator:     validator,
		reportsSvc:    reportsSvc,
	}
}

// Ask implements the suggestions service Ask method: question → SQL → report.
func (s *AskService) Ask(ctx context.Context, payload *suggestions.AskPayload) (*suggestions.AskResult, error) {
	question := strings.TrimSpace(payload.Question)
	if question == "" {
		return nil, &suggestions.ValidationError{Name: "validation_error", Message: "question is required", Code: strPtr("VALIDATION_ERROR")}
	}

	schemaResult, err := s.catalogLoader.Load(ctx)
	if err != nil {
		return nil, &suggestions.LLMError{Name: "llm_error", Message: "failed to load schema: " + err.Error(), Code: strPtr("SCHEMA_ERROR")}
	}
	schemaText := llm.FormatSchemaForPrompt(schemaResult)
	prompt := llm.BuildNL2SQLPrompt(question, schemaText)

	response, err := s.llmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, &suggestions.LLMError{Name: "llm_error", Message: err.Error(), Code: strPtr("LLM_ERROR")}
	}
	sql := llm.ParseSQLFromResponse(response)
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, &suggestions.LLMError{Name: "llm_error", Message: "LLM did not return any SQL", Code: strPtr("LLM_ERROR")}
	}

	if err := s.validator.Validate(sql); err != nil {
		return nil, &suggestions.ValidationError{Name: "validation_error", Message: err.Error(), Code: strPtr("VALIDATION_ERROR")}
	}

	reportPayload := &reports.GenerateReportPayload{SQL: sql}
	report, err := s.reportsSvc.Generate(ctx, reportPayload)
	if err != nil {
		if ve, ok := err.(*reports.ValidationError); ok {
			return nil, &suggestions.ValidationError{Name: ve.Name, Message: ve.Message, Code: ve.Code}
		}
		if le, ok := err.(*reports.LLMError); ok {
			return nil, &suggestions.LLMError{Name: le.Name, Message: le.Message, Code: le.Code}
		}
		return nil, &suggestions.LLMError{Name: "llm_error", Message: err.Error(), Code: strPtr("REPORT_ERROR")}
	}

	return &suggestions.AskResult{
		Question: question,
		SQL:      sql,
		Report:   reportToSuggestions(report),
	}, nil
}

// Explain implements the suggestions service Explain method: SQL → plain-English explanation.
func (s *AskService) Explain(ctx context.Context, payload *suggestions.ExplainPayload) (*suggestions.ExplainResult, error) {
	sql := strings.TrimSpace(payload.SQL)
	if sql == "" {
		return nil, &suggestions.ValidationError{Name: "validation_error", Message: "sql is required", Code: strPtr("VALIDATION_ERROR")}
	}
	sql = strings.TrimSuffix(sql, ";")
	sql = strings.TrimSpace(sql)
	if err := s.validator.Validate(sql); err != nil {
		return nil, &suggestions.ValidationError{Name: "validation_error", Message: err.Error(), Code: strPtr("VALIDATION_ERROR")}
	}
	prompt := llm.BuildExplainPrompt(sql)
	response, err := s.llmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, &suggestions.LLMError{Name: "llm_error", Message: err.Error(), Code: strPtr("LLM_ERROR")}
	}
	explanation := strings.TrimSpace(response)
	if explanation == "" {
		explanation = "No explanation returned."
	}
	return &suggestions.ExplainResult{SQL: sql, Explanation: explanation}, nil
}

// reportToSuggestions copies a reports.Report into a suggestions.Report (same design, different packages).
func reportToSuggestions(r *reports.Report) *suggestions.Report {
	if r == nil {
		return nil
	}
	out := &suggestions.Report{
		ID:           r.ID,
		SavedQueryID: r.SavedQueryID,
		SQL:          r.SQL,
		CreatedAt:    r.CreatedAt,
		LlmModel:     r.LlmModel,
		LlmProvider:  r.LlmProvider,
	}
	if r.Narrative != nil {
		out.Narrative = &suggestions.NarrativeContent{
			Headline:        r.Narrative.Headline,
			Takeaways:       append([]string(nil), r.Narrative.Takeaways...),
			Drivers:         append([]string(nil), r.Narrative.Drivers...),
			Limitations:     append([]string(nil), r.Narrative.Limitations...),
			Recommendations: append([]string(nil), r.Narrative.Recommendations...),
		}
	}
	if r.Metrics != nil {
		out.Metrics = copyMetricsToSuggestions(r.Metrics)
	}
	if len(r.ChartSuggestions) > 0 {
		out.ChartSuggestions = make([]*suggestions.ChartSuggestion, len(r.ChartSuggestions))
		for i, c := range r.ChartSuggestions {
			if c != nil {
				out.ChartSuggestions[i] = &suggestions.ChartSuggestion{ChartType: c.ChartType, Label: c.Label, Reason: c.Reason}
			}
		}
	}
	return out
}

func copyMetricsToSuggestions(m *reports.MetricsData) *suggestions.MetricsData {
	if m == nil {
		return nil
	}
	out := &suggestions.MetricsData{
		PeriodCurrentLabel:  m.PeriodCurrentLabel,
		PeriodPreviousLabel: m.PeriodPreviousLabel,
		PerfSuggestions:     append([]string(nil), m.PerfSuggestions...),
	}
	if len(m.Correlations) > 0 {
		out.Correlations = make([]*suggestions.CorrelationPairData, len(m.Correlations))
		for i, c := range m.Correlations {
			if c != nil {
				out.Correlations[i] = &suggestions.CorrelationPairData{
					ColumnA:  c.ColumnA,
					ColumnB:  c.ColumnB,
					Pearson:  c.Pearson,
					Spearman: c.Spearman,
				}
			}
		}
	}
	if m.Aggregates != nil {
		out.Aggregates = make(map[string]*suggestions.AggregateData)
		for k, v := range m.Aggregates {
			if v != nil {
				out.Aggregates[k] = &suggestions.AggregateData{Sum: v.Sum, Avg: v.Avg, Min: v.Min, Max: v.Max, Count: v.Count, StdDev: v.StdDev}
			}
		}
	}
	if m.TopCategories != nil {
		out.TopCategories = make(map[string][]*suggestions.TopCategoryData)
		for k, arr := range m.TopCategories {
			for _, t := range arr {
				if t != nil {
					out.TopCategories[k] = append(out.TopCategories[k], &suggestions.TopCategoryData{Category: t.Category, Value: t.Value, Percentage: t.Percentage})
				}
			}
		}
	}
	if m.TimeSeries != nil {
		out.TimeSeries = make(map[string]*suggestions.TimeSeriesData)
		for k, v := range m.TimeSeries {
			if v != nil {
				ts := &suggestions.TimeSeriesData{
					CurrentPeriod:              v.CurrentPeriod,
					PreviousPeriod:             v.PreviousPeriod,
					Change:                     v.Change,
					ChangePercentage:           v.ChangePercentage,
					Trend:                      v.Trend,
					MovingAverage:              v.MovingAverage,
					NextPeriodForecast:         v.NextPeriodForecast,
					ForecastCiLower:            v.ForecastCiLower,
					ForecastCiUpper:            v.ForecastCiUpper,
					PredictiveSummary:          v.PredictiveSummary,
					ExponentialSmoothForecast:  v.ExponentialSmoothForecast,
					HoltForecast:               v.HoltForecast,
					SeasonalPeriod:             v.SeasonalPeriod,
					SeasonallyAdjustedForecast: v.SeasonallyAdjustedForecast,
				}
				for _, p := range v.Periods {
					if p != nil {
						ts.Periods = append(ts.Periods, &suggestions.PeriodPointData{Label: p.Label, Value: p.Value})
					}
				}
				for _, a := range v.Anomalies {
					if a != nil {
						ts.Anomalies = append(ts.Anomalies, &suggestions.AnomalyPointData{PeriodLabel: a.PeriodLabel, Value: a.Value, Reason: a.Reason})
					}
				}
				if v.TrendSummary != nil {
					ts.TrendSummary = &suggestions.TrendSummaryData{Direction: v.TrendSummary.Direction, Slope: v.TrendSummary.Slope, PeriodsUsed: v.TrendSummary.PeriodsUsed, Summary: v.TrendSummary.Summary}
				}
				out.TimeSeries[k] = ts
			}
		}
	}
	if m.DataQuality != nil {
		out.DataQuality = make(map[string]*suggestions.ColumnQualityData)
		for k, v := range m.DataQuality {
			if v != nil {
				out.DataQuality[k] = &suggestions.ColumnQualityData{NullCount: v.NullCount, DistinctCount: v.DistinctCount, TotalRows: v.TotalRows, NullPct: v.NullPct}
			}
		}
	}
	if len(m.Cohorts) > 0 {
		out.Cohorts = make([]*suggestions.CohortMetricData, len(m.Cohorts))
		for i, co := range m.Cohorts {
			if co == nil {
				continue
			}
			periods := make([]*suggestions.CohortPeriodPointData, len(co.Periods))
			for j, p := range co.Periods {
				if p != nil {
					periods[j] = &suggestions.CohortPeriodPointData{PeriodLabel: p.PeriodLabel, Value: p.Value}
				}
			}
			out.Cohorts[i] = &suggestions.CohortMetricData{
				CohortLabel:  co.CohortLabel,
				Periods:      periods,
				RetentionPct: co.RetentionPct,
			}
		}
	}
	return out
}

// SuggestionsServiceWrapper implements suggestions.Service by delegating
// Queries and Similar to Suggester and Ask to AskService.
type SuggestionsServiceWrapper struct {
	Suggester suggestionsSuggester
	AskSvc    *AskService
}

// suggestionsSuggester is the interface used by the wrapper.
type suggestionsSuggester interface {
	Queries(context.Context, *suggestions.QueriesPayload) (*suggestions.SuggestedQueriesResult, error)
	Similar(context.Context, *suggestions.SimilarPayload) (*suggestions.SuggestedQueriesResult, error)
}

// Queries delegates to the suggester.
func (w *SuggestionsServiceWrapper) Queries(ctx context.Context, p *suggestions.QueriesPayload) (*suggestions.SuggestedQueriesResult, error) {
	return w.Suggester.Queries(ctx, p)
}

// Similar delegates to the suggester.
func (w *SuggestionsServiceWrapper) Similar(ctx context.Context, p *suggestions.SimilarPayload) (*suggestions.SuggestedQueriesResult, error) {
	return w.Suggester.Similar(ctx, p)
}

// Ask delegates to the AskService.
func (w *SuggestionsServiceWrapper) Ask(ctx context.Context, p *suggestions.AskPayload) (*suggestions.AskResult, error) {
	return w.AskSvc.Ask(ctx, p)
}

// Explain delegates to the AskService.
func (w *SuggestionsServiceWrapper) Explain(ctx context.Context, p *suggestions.ExplainPayload) (*suggestions.ExplainResult, error) {
	return w.AskSvc.Explain(ctx, p)
}
