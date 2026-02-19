package metrics

// ColumnType represents the type of a column in query results
type ColumnType string

const (
	ColumnTypeNumeric ColumnType = "numeric"
	ColumnTypeDate    ColumnType = "date"
	ColumnTypeText    ColumnType = "text"
	ColumnTypeBoolean ColumnType = "boolean"
)

// ColumnProfile contains profiling information about a column
type ColumnProfile struct {
	Name         string
	Type         ColumnType
	IsMeasure    bool // True if this column contains measurable values (sum, avg, etc.)
	IsDimension  bool // True if this column contains categorical/grouping data
	IsTimeSeries bool // True if this column contains date/time data
}

// Aggregates contains calculated aggregate metrics (including std dev for stats).
type Aggregates struct {
	Sum    *float64
	Avg    *float64
	Min    *float64
	Max    *float64
	Count  int
	StdDev *float64 // Standard deviation (stats)
}

// TopCategory represents a top category with its value and contribution
type TopCategory struct {
	Category   string
	Value      float64
	Percentage float64
}

// PeriodPoint is a single (label, value) point in a time series.
type PeriodPoint struct {
	Label string
	Value float64
}

// AnomalyPoint marks a period value that was flagged as anomalous.
type AnomalyPoint struct {
	PeriodLabel string
	Value       float64
	Reason      string // e.g. "High: 2.3σ above mean"
}

// TrendSummary describes trend over multiple periods (e.g. linear regression).
type TrendSummary struct {
	Direction   string  // "increasing", "decreasing", "stable"
	Slope       float64 // change per period (e.g. % or absolute)
	PeriodsUsed int
	Summary     string // human-readable, e.g. "Increasing ~5.2% per period over last 6 periods"
}

// TimeSeriesMetric contains time-series analysis results.
// Advanced metrics: Periods (last N), MovingAverage, Anomalies, TrendSummary.
type TimeSeriesMetric struct {
	CurrentPeriod    float64
	PreviousPeriod   *float64
	Change           *float64
	ChangePercentage *float64
	Trend            string // "up", "down", "flat"

	// Advanced: multi-period series (newest last)
	Periods []PeriodPoint
	// Moving average for current/latest period (e.g. 3-period SMA)
	MovingAverage *float64
	// Points flagged as statistical anomalies
	Anomalies []AnomalyPoint
	// Trend over full series or last N periods
	TrendSummary *TrendSummary
	// NextPeriodForecast is a simple predictive value: last period + slope (optional).
	NextPeriodForecast *float64
}

// ColumnQuality holds data-quality metrics for a column (nulls, distinct ratio).
type ColumnQuality struct {
	NullCount     int     // Number of nulls
	DistinctCount int     // Number of distinct non-null values
	TotalRows     int     // Total rows
	NullPct       float64 // Null percentage (0–100)
}

// Metrics contains all calculated metrics for query results.
// PerfSuggestions is set by the service layer from query execution context (timing, row count).
type Metrics struct {
	Profiles            []ColumnProfile
	Aggregates          map[string]Aggregates       // Keyed by column name
	TopCategories       map[string][]TopCategory    // Keyed by measure column, contains top categories
	TimeSeries          map[string]TimeSeriesMetric // Keyed by measure column
	CurrentPeriodLabel  string                      // Label for current period when time series present
	PreviousPeriodLabel string                      // Label for previous period
	DataQuality         map[string]ColumnQuality    // Per-column data quality (nulls, distinct)
	PerfSuggestions     []string                    // Performance suggestions (set by service from execution time/row count)
}
