package metrics_test

import (
	"testing"

	"github.com/pgquerynarrative/pgquerynarrative/app/metrics"
)

func TestCalculateMetrics_TimeSeries_AdvancedFields(t *testing.T) {
	columns := []string{"month", "revenue"}
	profiles := []metrics.ColumnProfile{
		{Name: "month", Type: metrics.ColumnTypeDate, IsTimeSeries: true},
		{Name: "revenue", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{"2025-01", 100.0},
		{"2025-02", 105.0},
		{"2025-03", 98.0},
		{"2025-04", 110.0},
		{"2025-05", 115.0},
		{"2025-06", 120.0},
		{"2025-07", 118.0},
		{"2025-08", 125.0},
	}
	m := metrics.CalculateMetrics(columns, rows, profiles, 0.5)
	if m == nil {
		t.Fatal("expected non-nil metrics")
	}
	if len(m.TimeSeries) == 0 {
		t.Fatal("expected time series for revenue")
	}
	ts, ok := m.TimeSeries["revenue"]
	if !ok {
		t.Fatal("expected revenue in time series")
	}
	if ts.CurrentPeriod != 125.0 {
		t.Errorf("current period = %v, want 125", ts.CurrentPeriod)
	}
	if len(ts.Periods) == 0 {
		t.Error("expected periods to be populated")
	}
	if ts.MovingAverage == nil {
		t.Error("expected moving average (8 periods >= 3)")
	}
	if ts.TrendSummary == nil {
		t.Error("expected trend summary (>= 2 periods)")
	}
	if ts.TrendSummary != nil && ts.TrendSummary.Summary == "" {
		t.Error("expected non-empty trend summary")
	}
	if ts.NextPeriodForecast == nil {
		t.Error("expected next period forecast when trend summary is present")
	}
}

func TestCalculateMetrics_TimeSeries_AnomalyDetection(t *testing.T) {
	columns := []string{"month", "value"}
	profiles := []metrics.ColumnProfile{
		{Name: "month", Type: metrics.ColumnTypeDate, IsTimeSeries: true},
		{Name: "value", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{"2025-01", 10.0},
		{"2025-02", 11.0},
		{"2025-03", 10.5},
		{"2025-04", 1000.0},
		{"2025-05", 10.0},
		{"2025-06", 11.0},
	}
	m := metrics.CalculateMetrics(columns, rows, profiles, 0.5)
	ts, ok := m.TimeSeries["value"]
	if !ok {
		t.Fatal("expected time series for value")
	}
	if len(ts.Anomalies) == 0 {
		t.Error("expected at least one anomaly (1000 is outlier)")
	}
}

func TestLinearRegression(t *testing.T) {
	y := []float64{1, 3, 5, 7, 9}
	slope, intercept := metrics.LinearRegression(y)
	if slope < 1.99 || slope > 2.01 {
		t.Errorf("slope = %v, want ~2", slope)
	}
	if intercept < 0.99 || intercept > 1.01 {
		t.Errorf("intercept = %v, want ~1", intercept)
	}
}

func TestMeanAndStd(t *testing.T) {
	vals := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, std := metrics.MeanAndStd(vals)
	if mean < 4.99 || mean > 5.01 {
		t.Errorf("mean = %v, want 5", mean)
	}
	if std < 2.0 || std > 2.1 {
		t.Errorf("std = %v, want ~2", std)
	}
}

func TestCalculateMetrics_DataQualityAndStats(t *testing.T) {
	columns := []string{"a", "b"}
	profiles := []metrics.ColumnProfile{
		{Name: "a", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
		{Name: "b", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{10.0, 20.0},
		{12.0, 20.0},
		{nil, 22.0},
	}
	m := metrics.CalculateMetrics(columns, rows, profiles, 0.5)
	if len(m.DataQuality) != 2 {
		t.Errorf("data quality: got %d columns, want 2", len(m.DataQuality))
	}
	if q, ok := m.DataQuality["a"]; ok {
		if q.NullCount != 1 || q.TotalRows != 3 || q.DistinctCount != 2 {
			t.Errorf("data quality a: nulls=%d total=%d distinct=%d", q.NullCount, q.TotalRows, q.DistinctCount)
		}
	}
	if agg, ok := m.Aggregates["a"]; ok {
		if agg.StdDev == nil {
			t.Error("expected std dev for numeric measure a")
		}
	}
}

func TestCalculateMetrics_TimeSeries_PeriodLabelsAndPreviousPeriod(t *testing.T) {
	columns := []string{"month", "revenue"}
	profiles := []metrics.ColumnProfile{
		{Name: "month", Type: metrics.ColumnTypeDate, IsTimeSeries: true},
		{Name: "revenue", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{"2025-01", 80.0},
		{"2025-02", 120.0},
	}
	m := metrics.CalculateMetrics(columns, rows, profiles, 0.5)
	if len(m.TimeSeries) == 0 {
		t.Fatal("expected time series")
	}
	ts := m.TimeSeries["revenue"]
	if ts.PreviousPeriod == nil {
		t.Error("expected previous_period when two periods exist")
	}
	if *ts.PreviousPeriod != 80.0 {
		t.Errorf("previous_period = %v, want 80", *ts.PreviousPeriod)
	}
	if ts.CurrentPeriod != 120.0 {
		t.Errorf("current_period = %v, want 120", ts.CurrentPeriod)
	}
	if m.CurrentPeriodLabel == "" {
		t.Error("expected current period label")
	}
	if m.PreviousPeriodLabel == "" {
		t.Error("expected previous period label")
	}
}

func TestCalculateMetrics_NoTimeSeries(t *testing.T) {
	columns := []string{"category", "total"}
	profiles := []metrics.ColumnProfile{
		{Name: "category", Type: metrics.ColumnTypeText, IsDimension: true},
		{Name: "total", Type: metrics.ColumnTypeNumeric, IsMeasure: true},
	}
	rows := [][]interface{}{
		{"A", 10.0},
		{"B", 20.0},
	}
	m := metrics.CalculateMetrics(columns, rows, profiles, 0.5)
	if len(m.TimeSeries) != 0 {
		t.Errorf("expected no time series, got %d", len(m.TimeSeries))
	}
}
