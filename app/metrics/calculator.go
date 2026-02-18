package metrics

import (
	"math"
	"sort"
	"time"
)

// CalculateMetrics computes all metrics from query results.
// trendThresholdPercent is the minimum absolute % change to label as "up" or "down" (vs "flat"); 0 or negative uses default 0.5.
func CalculateMetrics(columns []string, rows [][]interface{}, profiles []ColumnProfile, trendThresholdPercent float64) *Metrics {
	metrics := &Metrics{
		Profiles:      profiles,
		Aggregates:    make(map[string]Aggregates),
		TopCategories: make(map[string][]TopCategory),
		TimeSeries:    make(map[string]TimeSeriesMetric),
	}

	if len(rows) == 0 {
		return metrics
	}

	// Calculate aggregates for numeric/measure columns
	for i, profile := range profiles {
		if profile.IsMeasure && profile.Type == ColumnTypeNumeric {
			agg := calculateAggregates(rows, i)
			metrics.Aggregates[profile.Name] = agg
		}
	}

	// Calculate top categories for grouped data
	metrics.calculateTopCategories(columns, rows, profiles)

	// Calculate time-series metrics
	if trendThresholdPercent <= 0 {
		trendThresholdPercent = 0.5
	}
	metrics.calculateTimeSeries(columns, rows, profiles, trendThresholdPercent)

	return metrics
}

// calculateAggregates computes sum, avg, min, max, count for a numeric column
func calculateAggregates(rows [][]interface{}, colIndex int) Aggregates {
	var sum float64
	var count int
	var min, max *float64

	for _, row := range rows {
		if colIndex >= len(row) {
			continue
		}

		val, ok := GetNumericValue(row[colIndex])
		if !ok {
			continue
		}

		sum += val
		count++

		if min == nil || val < *min {
			min = &val
		}
		if max == nil || val > *max {
			max = &val
		}
	}

	agg := Aggregates{
		Sum:   &sum,
		Min:   min,
		Max:   max,
		Count: count,
	}

	if count > 0 {
		avg := sum / float64(count)
		agg.Avg = &avg
	}

	return agg
}

// calculateTopCategories identifies top categories by measure value
func (m *Metrics) calculateTopCategories(columns []string, rows [][]interface{}, profiles []ColumnProfile) {
	// Find dimension columns (categories) and measure columns
	dimensionCols := make([]int, 0, len(profiles))
	measureCols := make([]int, 0, len(profiles))

	for i, profile := range profiles {
		if profile.IsDimension {
			dimensionCols = append(dimensionCols, i)
		}
		if profile.IsMeasure {
			measureCols = append(measureCols, i)
		}
	}

	// For each measure, find top categories
	for _, measureIdx := range measureCols {
		measureName := columns[measureIdx]

		// Group by first dimension column
		if len(dimensionCols) == 0 {
			continue
		}

		dimensionIdx := dimensionCols[0]
		categoryMap := make(map[string]float64)

		// Aggregate values by category
		for _, row := range rows {
			if measureIdx >= len(row) || dimensionIdx >= len(row) {
				continue
			}

			category := getStringValue(row[dimensionIdx])
			value, ok := GetNumericValue(row[measureIdx])
			if !ok {
				continue
			}

			categoryMap[category] += value
		}

		// Convert to slice and sort
		categories := make([]TopCategory, 0, len(categoryMap))
		var total float64
		for cat, val := range categoryMap {
			total += val
			categories = append(categories, TopCategory{
				Category: cat,
				Value:    val,
			})
		}

		// Sort by value descending
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Value > categories[j].Value
		})

		// Calculate percentages and limit to top 10
		topN := 10
		if len(categories) < topN {
			topN = len(categories)
		}

		for i := 0; i < topN; i++ {
			if total > 0 {
				categories[i].Percentage = (categories[i].Value / total) * 100
			}
		}

		m.TopCategories[measureName] = categories[:topN]
	}
}

func (m *Metrics) calculateTimeSeries(columns []string, rows [][]interface{}, profiles []ColumnProfile, trendThresholdPercent float64) {
	var timeColIdx = -1
	measureCols := make([]int, 0, len(profiles))

	for i, profile := range profiles {
		if profile.IsTimeSeries {
			timeColIdx = i
		}
		if profile.IsMeasure {
			measureCols = append(measureCols, i)
		}
	}

	if timeColIdx == -1 || len(measureCols) == 0 {
		return
	}

	if len(rows) < 2 {
		return
	}

	// Group by period and calculate totals
	periodTotals := make(map[string]map[int]float64) // period -> measureCol -> total
	periodTimes := make(map[string]time.Time)        // period -> parsed time for sorting

	for _, row := range rows {
		if timeColIdx >= len(row) {
			continue
		}

		periodStr := getStringValue(row[timeColIdx])
		if periodStr == "" {
			continue
		}

		// Try to parse as time for proper sorting
		var periodTime time.Time
		var err error
		timeFormats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range timeFormats {
			if t, parseErr := time.Parse(format, periodStr); parseErr == nil {
				periodTime = t
				err = nil
				break
			}
		}
		if err != nil {
			// Use string as-is if parsing fails
			periodTime = time.Time{}
		}

		if periodTotals[periodStr] == nil {
			periodTotals[periodStr] = make(map[int]float64)
		}
		if !periodTime.IsZero() {
			periodTimes[periodStr] = periodTime
		}

		for _, measureIdx := range measureCols {
			if measureIdx >= len(row) {
				continue
			}

			value, ok := GetNumericValue(row[measureIdx])
			if ok {
				periodTotals[periodStr][measureIdx] += value
			}
		}
	}

	// Get periods sorted by time if available, otherwise by string
	periods := make([]string, 0, len(periodTotals))
	for p := range periodTotals {
		periods = append(periods, p)
	}

	if len(periods) < 2 {
		return
	}

	// Sort periods by time if available, otherwise by string
	sort.Slice(periods, func(i, j int) bool {
		ti, hasTimeI := periodTimes[periods[i]]
		tj, hasTimeJ := periodTimes[periods[j]]
		if hasTimeI && hasTimeJ {
			return ti.Before(tj)
		}
		if hasTimeI {
			return true
		}
		if hasTimeJ {
			return false
		}
		return periods[i] < periods[j]
	})

	// Use last two periods for comparison
	currentPeriod := periods[len(periods)-1]
	previousPeriod := periods[len(periods)-2]
	m.CurrentPeriodLabel = currentPeriod
	m.PreviousPeriodLabel = previousPeriod

	// Calculate metrics for each measure
	for _, measureIdx := range measureCols {
		measureName := columns[measureIdx]

		current := periodTotals[currentPeriod][measureIdx]
		previous := periodTotals[previousPeriod][measureIdx]

		change := current - previous
		var changePct *float64
		if previous != 0 {
			pct := (change / previous) * 100
			changePct = &pct
		}

		trend := "flat"
		if changePct != nil {
			if math.Abs(*changePct) < trendThresholdPercent {
				trend = "flat"
			} else if *changePct > 0 {
				trend = "up"
			} else {
				trend = "down"
			}
		}

		m.TimeSeries[measureName] = TimeSeriesMetric{
			CurrentPeriod:    current,
			PreviousPeriod:   &previous,
			Change:           &change,
			ChangePercentage: changePct,
			Trend:            trend,
		}
	}
}

// getStringValue converts an interface{} to string
func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return ""
	}
}
