// Package service provides business logic for queries and reports.
// It acts as a bridge between the API layer and the data/query execution layer.
package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/app/apilog"
	"github.com/pgquerynarrative/pgquerynarrative/app/charts"
	"github.com/pgquerynarrative/pgquerynarrative/app/config"
	"github.com/pgquerynarrative/pgquerynarrative/app/embedding"
	"github.com/pgquerynarrative/pgquerynarrative/app/format"
	"github.com/pgquerynarrative/pgquerynarrative/app/metrics"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
)

// QueriesService handles query execution and saved query management.
type QueriesService struct {
	readOnlyPool   *pgxpool.Pool       // Pool for executing queries (read-only)
	appPool        *pgxpool.Pool       // Pool for saving queries (full access)
	runner         *queryrunner.Runner // Query execution engine
	metricsOpts    *metrics.Options    // Metrics and time-series options (windows, thresholds)
	embedder       embedding.Embedder  // Optional: for storing query embeddings on save
	embeddingStore *embedding.Store    // Optional: for RAG / similar-query retrieval
	embeddingModel string              // Model name to store with embedding (e.g. nomic-embed-text)
}

var strPtr = format.StrPtr

// NewQueriesService creates a new queries service with the specified dependencies.
// metricsCfg supplies trend threshold, anomaly sigma, trend periods, and moving-average window; nil uses defaults.
func NewQueriesService(readOnlyPool, appPool *pgxpool.Pool, runner *queryrunner.Runner, metricsCfg config.MetricsConfig) *QueriesService {
	opts := metricsOptionsFromConfig(metricsCfg)
	return &QueriesService{
		readOnlyPool: readOnlyPool,
		appPool:      appPool,
		runner:       runner,
		metricsOpts:  opts,
	}
}

// NewQueriesServiceWithEmbedding is like NewQueriesService but enables storing embeddings
// when saved queries are created, for similar-query retrieval and RAG. embeddingModel
// is the name of the embedding model (e.g. nomic-embed-text).
func NewQueriesServiceWithEmbedding(readOnlyPool, appPool *pgxpool.Pool, runner *queryrunner.Runner, metricsCfg config.MetricsConfig, embedder embedding.Embedder, embeddingStore *embedding.Store, embeddingModel string) *QueriesService {
	opts := metricsOptionsFromConfig(metricsCfg)
	return &QueriesService{
		readOnlyPool:   readOnlyPool,
		appPool:        appPool,
		runner:         runner,
		metricsOpts:    opts,
		embedder:       embedder,
		embeddingStore: embeddingStore,
		embeddingModel: embeddingModel,
	}
}

// Run executes a SQL query and returns the results.
//
// The query is validated, executed with timeout protection, and results are
// limited to prevent memory exhaustion. Errors are converted to appropriate
// API error types.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - payload: Query execution request (SQL and limit)
//
// Returns:
//   - RunQueryResult with columns, rows, and metadata
//   - ValidationError if query is invalid or times out
func (s *QueriesService) Run(ctx context.Context, payload *queries.RunQueryPayload) (*queries.RunQueryResult, error) {
	result, err := s.runner.Run(ctx, payload.SQL, int(payload.Limit))
	if err != nil {
		kind, userMsg := ClassifyRunError(err)
		if kind == RunErrorTimeout {
			apilog.ValidationError("run", "timeout_error", err.Error())
			return nil, &queries.ValidationError{Name: "timeout_error", Message: userMsg, Code: strPtr("TIMEOUT_ERROR")}
		}
		apilog.ValidationError("run", "validation_error", err.Error())
		return nil, &queries.ValidationError{Name: "validation_error", Message: userMsg, Code: strPtr("VALIDATION_ERROR")}
	}

	cols := make([]*queries.ColumnInfo, 0, len(result.Columns))
	colNames := make([]string, len(result.Columns))
	colTypes := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		cols = append(cols, &queries.ColumnInfo{
			Name: col.Name,
			Type: col.Type,
		})
		colNames[i] = col.Name
		colTypes[i] = col.Type
	}

	chartSuggestions := suggestToQueries(charts.Suggest(colNames, colTypes, result.Rows))
	periodComparison, currentLabel, previousLabel := timeSeriesToPeriodComparison(colNames, result.Rows, s.metricsOpts)

	var rowCount32 int32 = math.MaxInt32
	if result.RowCount < math.MaxInt32 {
		rowCount32 = int32(result.RowCount)
	}
	res := &queries.RunQueryResult{
		Columns:          cols,
		Rows:             result.Rows,
		RowCount:         rowCount32,
		ExecutionTimeMs:  result.ExecutionTimeMs,
		Limit:            int32(result.RowLimitApplied),
		ChartSuggestions: chartSuggestions,
		PeriodComparison: periodComparison,
	}
	if currentLabel != "" {
		res.PeriodCurrentLabel = &currentLabel
	}
	if previousLabel != "" {
		res.PeriodPreviousLabel = &previousLabel
	}
	return res, nil
}

// suggestToQueries converts charts.Suggestion slice to API type.
func suggestToQueries(in []charts.Suggestion) []*queries.ChartSuggestion {
	if len(in) == 0 {
		return nil
	}
	out := make([]*queries.ChartSuggestion, len(in))
	for i := range in {
		out[i] = &queries.ChartSuggestion{
			ChartType: in[i].ChartType,
			Label:     in[i].Label,
			Reason:    in[i].Reason,
		}
	}
	return out
}

// metricsOptionsFromConfig builds metrics.Options from config. Nil or zero config uses defaults.
func metricsOptionsFromConfig(c config.MetricsConfig) *metrics.Options {
	o := &metrics.Options{
		TrendThresholdPercent:    c.TrendThresholdPercent,
		AnomalySigma:             c.AnomalySigma,
		AnomalyMethod:            c.AnomalyMethod,
		TrendPeriods:             c.TrendPeriods,
		MovingAvgWindow:          c.MovingAvgWindow,
		ConfidenceLevel:          c.ConfidenceLevel,
		MinRowsForCorrelation:    c.MinRowsForCorrelation,
		SmoothingAlpha:           c.SmoothingAlpha,
		SmoothingBeta:            c.SmoothingBeta,
		MaxSeasonalLag:           c.MaxSeasonalLag,
		MinPeriodsForSeasonality: c.MinPeriodsForSeasonality,
		MaxTimeSeriesPeriods:     c.MaxTimeSeriesPeriods,
	}
	o.ApplyDefaults()
	return o
}

// timeSeriesToPeriodComparison computes period-over-period from query result and returns API slice and period labels.
func timeSeriesToPeriodComparison(columnNames []string, rows [][]interface{}, opts *metrics.Options) ([]*queries.PeriodComparisonItem, string, string) {
	if len(rows) < 2 {
		return nil, "", ""
	}
	profiles := metrics.ProfileColumns(columnNames, rows)
	m := metrics.CalculateMetrics(columnNames, rows, profiles, opts)
	if len(m.TimeSeries) == 0 {
		return nil, "", ""
	}
	out := make([]*queries.PeriodComparisonItem, 0, len(m.TimeSeries))
	for measure, ts := range m.TimeSeries {
		item := &queries.PeriodComparisonItem{
			Measure: measure,
			Current: ts.CurrentPeriod,
			Trend:   ts.Trend,
		}
		if ts.PreviousPeriod != nil {
			item.Previous = ts.PreviousPeriod
		}
		if ts.Change != nil {
			item.Change = ts.Change
		}
		if ts.ChangePercentage != nil {
			item.ChangePercentage = ts.ChangePercentage
		}
		out = append(out, item)
	}
	return out, m.CurrentPeriodLabel, m.PreviousPeriodLabel
}

// Save stores a query for later reuse.
//
// Parameters:
//   - ctx: Context for cancellation
//   - payload: Query to save (name, SQL, description, tags)
//
// Returns:
//   - SavedQuery with generated ID and timestamps
//   - Error if database operation fails
func (s *QueriesService) Save(ctx context.Context, payload *queries.SaveQueryPayload) (*queries.SavedQuery, error) {
	row := s.appPool.QueryRow(ctx, `
		INSERT INTO app.saved_queries (name, sql, description, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, sql, description, tags, created_at, updated_at
	`, payload.Name, payload.SQL, payload.Description, payload.Tags)

	var item queries.SavedQuery
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&item.ID, &item.Name, &item.SQL, &item.Description, &item.Tags, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Optionally store embedding for similar-query retrieval and RAG
	if s.embedder != nil && s.embeddingStore != nil && s.embeddingModel != "" {
		text := item.Name
		if item.Description != nil && *item.Description != "" {
			text = text + " " + *item.Description
		}
		text = text + " " + item.SQL
		vec, err := s.embedder.Embed(ctx, text)
		if err == nil {
			_ = s.embeddingStore.Upsert(ctx, item.ID, vec, s.embeddingModel)
		}
	}

	return &item, nil
}

// ListSaved retrieves a paginated list of saved queries.
//
// Supports optional filtering by tags. Results are ordered by creation date (newest first).
//
// Parameters:
//   - ctx: Context for cancellation
//   - payload: Pagination and optional tag filter
//
// Returns:
//   - SavedQueryList with items, limit, and offset
//   - Error if database operation fails
func (s *QueriesService) ListSaved(ctx context.Context, payload *queries.ListSavedPayload) (*queries.SavedQueryList, error) {
	limit := int(payload.Limit)
	offset := int(payload.Offset)

	var rows pgx.Rows
	var err error
	if len(payload.Tags) > 0 {
		rows, err = s.appPool.Query(ctx, `
			SELECT id, name, sql, description, tags, created_at, updated_at
			FROM app.saved_queries
			WHERE tags && $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, payload.Tags, limit, offset)
	} else {
		rows, err = s.appPool.Query(ctx, `
			SELECT id, name, sql, description, tags, created_at, updated_at
			FROM app.saved_queries
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*queries.SavedQuery, 0, limit)
	for rows.Next() {
		var item queries.SavedQuery
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&item.ID, &item.Name, &item.SQL, &item.Description, &item.Tags, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, &item)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &queries.SavedQueryList{
		Items:  items,
		Limit:  payload.Limit,
		Offset: payload.Offset,
	}, nil
}

// GetSaved retrieves a saved query by ID.
//
// Parameters:
//   - ctx: Context for cancellation
//   - payload: Query ID to retrieve
//
// Returns:
//   - SavedQuery if found
//   - NotFoundError if query doesn't exist
//   - Error if database operation fails
func (s *QueriesService) GetSaved(ctx context.Context, payload *queries.GetSavedPayload) (*queries.SavedQuery, error) {
	row := s.appPool.QueryRow(ctx, `
		SELECT id, name, sql, description, tags, created_at, updated_at
		FROM app.saved_queries
		WHERE id = $1
	`, payload.ID)

	var item queries.SavedQuery
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&item.ID, &item.Name, &item.SQL, &item.Description, &item.Tags, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &queries.NotFoundError{
				Name:    "not_found",
				Message: "saved query not found",
				Code:    strPtr("NOT_FOUND"),
			}
		}
		return nil, err
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &item, nil
}

// DeleteSaved removes a saved query by ID.
//
// Parameters:
//   - ctx: Context for cancellation
//   - payload: Query ID to delete
//
// Returns:
//   - nil if deletion successful
//   - NotFoundError if query doesn't exist
//   - Error if database operation fails
func (s *QueriesService) DeleteSaved(ctx context.Context, payload *queries.DeleteSavedPayload) error {
	tag, err := s.appPool.Exec(ctx, `DELETE FROM app.saved_queries WHERE id = $1`, payload.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return &queries.NotFoundError{
			Name:    "not_found",
			Message: "saved query not found",
			Code:    strPtr("NOT_FOUND"),
		}
	}
	return nil
}
