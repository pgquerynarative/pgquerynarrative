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

// Aggregates contains calculated aggregate metrics
type Aggregates struct {
	Sum   *float64
	Avg   *float64
	Min   *float64
	Max   *float64
	Count int
}

// TopCategory represents a top category with its value and contribution
type TopCategory struct {
	Category   string
	Value      float64
	Percentage float64
}

// TimeSeriesMetric contains time-series analysis results
type TimeSeriesMetric struct {
	CurrentPeriod    float64
	PreviousPeriod   *float64
	Change           *float64
	ChangePercentage *float64
	Trend            string // "up", "down", "flat"
}

// Metrics contains all calculated metrics for query results
type Metrics struct {
	Profiles            []ColumnProfile
	Aggregates          map[string]Aggregates       // Keyed by column name
	TopCategories       map[string][]TopCategory    // Keyed by measure column, contains top categories
	TimeSeries          map[string]TimeSeriesMetric // Keyed by measure column
	CurrentPeriodLabel  string                      // Label for current period when time series present (e.g. "2025-01-01")
	PreviousPeriodLabel string                      // Label for previous period
}
