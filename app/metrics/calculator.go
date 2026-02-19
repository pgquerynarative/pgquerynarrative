package metrics

import (
	"fmt"
	"math"
	"sort"
	"strconv"
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
		DataQuality:   make(map[string]ColumnQuality),
	}

	if len(rows) == 0 {
		return metrics
	}

	// Data quality (nulls, distinct) for all columns
	metrics.calculateDataQuality(columns, rows)

	// Calculate aggregates for numeric/measure columns (includes std dev)
	for i, profile := range profiles {
		if profile.IsMeasure && profile.Type == ColumnTypeNumeric {
			agg := calculateAggregates(rows, i)
			metrics.Aggregates[profile.Name] = agg
		}
	}

	// Calculate top categories for grouped data
	metrics.calculateTopCategories(columns, rows, profiles)

	// Calculate time-series metrics (includes predictive next-period forecast)
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

	vals := make([]float64, 0, len(rows))
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
		vals = append(vals, val)

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
		_, std := MeanAndStd(vals)
		if std > 0 {
			agg.StdDev = &std
		}
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

	const maxPeriods = 24
	const movingAvgWindow = 3
	const anomalySigma = 2.0
	const trendPeriods = 6

	// Build ordered values per measure (one float per period)
	measureValues := make(map[string][]float64)
	for _, measureIdx := range measureCols {
		measureName := columns[measureIdx]
		vals := make([]float64, 0, len(periods))
		for _, p := range periods {
			vals = append(vals, periodTotals[p][measureIdx])
		}
		measureValues[measureName] = vals
	}

	// Calculate metrics for each measure (including advanced)
	for _, measureIdx := range measureCols {
		measureName := columns[measureIdx]
		values := measureValues[measureName]

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

		ts := TimeSeriesMetric{
			CurrentPeriod:    current,
			PreviousPeriod:   &previous,
			Change:           &change,
			ChangePercentage: changePct,
			Trend:            trend,
		}

		// Advanced: last N periods as Periods (newest last)
		nPeriods := len(periods)
		if nPeriods > maxPeriods {
			nPeriods = maxPeriods
		}
		start := len(periods) - nPeriods
		if start < 0 {
			start = 0
		}
		ts.Periods = make([]PeriodPoint, 0, nPeriods)
		for i := start; i < len(periods); i++ {
			ts.Periods = append(ts.Periods, PeriodPoint{
				Label: periods[i],
				Value: values[i],
			})
		}

		// Moving average (SMA) for latest period
		if len(values) >= movingAvgWindow {
			var sum float64
			for j := len(values) - movingAvgWindow; j < len(values); j++ {
				sum += values[j]
			}
			ma := sum / float64(movingAvgWindow)
			ts.MovingAverage = &ma
		}

		// Anomaly detection: z-score; flag points beyond anomalySigma std devs
		if len(values) >= 3 {
			mean, std := MeanAndStd(values)
			if std > 0 {
				for i, v := range values {
					z := (v - mean) / std
					if math.Abs(z) >= anomalySigma {
						reason := "High"
						if z < 0 {
							reason = "Low"
						}
						ts.Anomalies = append(ts.Anomalies, AnomalyPoint{
							PeriodLabel: periods[i],
							Value:       v,
							Reason:      formatAnomalyReason(reason, z, anomalySigma),
						})
					}
				}
			}
		}

		// Trend summary: linear regression over last trendPeriods (or all)
		nTrend := trendPeriods
		if len(values) < nTrend {
			nTrend = len(values)
		}
		if nTrend >= 2 {
			trendStart := len(values) - nTrend
			if trendStart < 0 {
				trendStart = 0
				nTrend = len(values)
			}
			slope, _ := LinearRegression(values[trendStart:])
			avgVal, _ := MeanAndStd(values[trendStart:])
			direction := "stable"
			if avgVal != 0 && math.Abs(slope) > 1e-10 {
				pctSlope := (slope / math.Abs(avgVal)) * 100
				if pctSlope > trendThresholdPercent {
					direction = "increasing"
				} else if pctSlope < -trendThresholdPercent {
					direction = "decreasing"
				}
			}
			summary := formatTrendSummary(direction, slope, avgVal, nTrend)
			ts.TrendSummary = &TrendSummary{
				Direction:   direction,
				Slope:       slope,
				PeriodsUsed: nTrend,
				Summary:     summary,
			}
			// Simple predictive: next period ≈ last value + slope
			if len(values) > 0 {
				forecast := values[len(values)-1] + slope
				ts.NextPeriodForecast = &forecast
			}
		}

		m.TimeSeries[measureName] = ts
	}
}

// calculateDataQuality fills DataQuality for each column (nulls, distinct count).
func (m *Metrics) calculateDataQuality(columns []string, rows [][]interface{}) {
	if len(rows) == 0 {
		return
	}
	total := len(rows)
	for colIdx, colName := range columns {
		nullCount := 0
		seen := make(map[string]struct{})
		for _, row := range rows {
			if colIdx >= len(row) {
				continue
			}
			v := row[colIdx]
			if v == nil {
				nullCount++
				continue
			}
			key := fmt.Sprintf("%v", v)
			seen[key] = struct{}{}
		}
		nullPct := 0.0
		if total > 0 {
			nullPct = float64(nullCount) / float64(total) * 100
		}
		m.DataQuality[colName] = ColumnQuality{
			NullCount:     nullCount,
			DistinctCount: len(seen),
			TotalRows:     total,
			NullPct:       nullPct,
		}
	}
}

// MeanAndStd returns mean and population standard deviation. Exported for testing.
func MeanAndStd(values []float64) (mean, std float64) {
	if len(values) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))
	var sqSum float64
	for _, v := range values {
		d := v - mean
		sqSum += d * d
	}
	variance := sqSum / float64(len(values))
	if variance <= 0 {
		return mean, 0
	}
	return mean, math.Sqrt(variance)
}

func formatAnomalyReason(highLow string, z, sigma float64) string {
	zAbs := math.Abs(z)
	if zAbs < 10 {
		return highLow + ": " + formatOneDecimal(zAbs) + "σ from mean (threshold " + formatOneDecimal(sigma) + "σ)"
	}
	return highLow + ": extreme value (" + formatOneDecimal(zAbs) + "σ from mean)"
}

func formatOneDecimal(x float64) string {
	return strconv.FormatFloat(x, 'f', 1, 64)
}

// LinearRegression returns slope and intercept for y = slope*x + intercept (x = 0,1,2,...). Exported for testing.
func LinearRegression(y []float64) (slope, intercept float64) {
	n := float64(len(y))
	if n < 2 {
		return 0, 0
	}
	var sumX, sumY, sumXY, sumX2 float64
	for i, v := range y {
		x := float64(i)
		sumX += x
		sumY += v
		sumXY += x * v
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	if math.Abs(denom) < 1e-20 {
		return 0, sumY / n
	}
	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n
	return slope, intercept
}

func formatTrendSummary(direction string, slope, avgVal float64, periodsUsed int) string {
	if periodsUsed == 0 {
		return ""
	}
	if math.Abs(avgVal) < 1e-10 {
		return direction + " over last " + strconv.Itoa(periodsUsed) + " periods (absolute change ~" + formatOneDecimal(slope) + " per period)."
	}
	pctPerPeriod := (slope / math.Abs(avgVal)) * 100
	switch direction {
	case "increasing":
		return "Increasing ~" + formatOneDecimal(pctPerPeriod) + "% per period over last " + strconv.Itoa(periodsUsed) + " periods."
	case "decreasing":
		return "Decreasing ~" + formatOneDecimal(-pctPerPeriod) + "% per period over last " + strconv.Itoa(periodsUsed) + " periods."
	default:
		return "Stable over last " + strconv.Itoa(periodsUsed) + " periods."
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
