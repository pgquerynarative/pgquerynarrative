// Package web provides HTTP handlers and PDF export for the PgQueryNarrative UI.
// pdf.go builds well-structured PDF reports from the report type.
package web

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
)

// BuildReportPDF writes a structured PDF report to w. Uses Helvetica (ASCII-safe).
func BuildReportPDF(w io.Writer, report *reports.Report) error {
	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.SetCreator("PgQueryNarrative", false)
	pdf.SetProducer("PgQueryNarrative", false)
	pdf.SetTitle("PgQueryNarrative Report", false)
	pdf.SetMargins(40, 40, 40)
	pdf.SetAutoPageBreak(true, 28)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 10)

	// Title and meta
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 16, "PgQueryNarrative Report", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 10, "Generated: "+report.CreatedAt, "", 1, "L", false, 0, "")
	if report.LlmProvider != "" || report.LlmModel != "" {
		model := report.LlmProvider
		if report.LlmModel != "" {
			if model != "" {
				model += " / "
			}
			model += report.LlmModel
		}
		pdf.CellFormat(0, 10, "Model: "+model, "", 1, "L", false, 0, "")
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(4)

	// Query (collapsed in a box)
	sqlStr := report.SQL
	if len(sqlStr) > 400 {
		sqlStr = sqlStr[:400] + "..."
	}
	pdf.SetFont("Courier", "", 8)
	pdf.SetFillColor(245, 245, 245)
	pdf.MultiCell(0, 12, safePDFString(sqlStr), "1", "L", true)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Ln(8)

	// Narrative: headline and sections
	if report.Narrative != nil {
		if report.Narrative.Headline != "" {
			pdf.SetFont("Helvetica", "B", 14)
			pdf.MultiCell(0, 14, safePDFString(report.Narrative.Headline), "", "L", false)
			pdf.SetFont("Helvetica", "", 10)
			pdf.Ln(6)
		}
		if len(report.Narrative.Takeaways) > 0 {
			sectionTitle(pdf, "Key takeaways")
			for _, t := range report.Narrative.Takeaways {
				pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
				pdf.MultiCell(0, 10, safePDFString(t), "", "L", false)
			}
			pdf.Ln(4)
		}
		for _, title := range []string{"Drivers", "Limitations", "Recommendations"} {
			var items []string
			switch title {
			case "Drivers":
				items = report.Narrative.Drivers
			case "Limitations":
				items = report.Narrative.Limitations
			case "Recommendations":
				items = report.Narrative.Recommendations
			}
			if len(items) > 0 {
				sectionTitle(pdf, title)
				for _, s := range items {
					pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
					pdf.MultiCell(0, 10, safePDFString(s), "", "L", false)
				}
				pdf.Ln(4)
			}
		}
	}

	// Metrics bar chart (when we have time series data)
	if report.Metrics != nil && len(report.Metrics.TimeSeries) > 0 {
		drawMetricsBarChart(pdf, report.Metrics.TimeSeries)
		pdf.Ln(8)
	}

	// Chart suggestions
	if len(report.ChartSuggestions) > 0 {
		sectionTitle(pdf, "Suggested charts")
		for _, s := range report.ChartSuggestions {
			if s == nil {
				continue
			}
			pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "B", 10)
			pdf.CellFormat(0, 10, safePDFString(s.Label), "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 10)
			pdf.Ln(10)
			pdf.CellFormat(20, 10, "", "", 0, "L", false, 0, "")
			pdf.MultiCell(0, 10, safePDFString(s.Reason), "", "L", false)
		}
		pdf.Ln(4)
	}

	// Time series / period comparison
	if report.Metrics != nil && len(report.Metrics.TimeSeries) > 0 {
		sectionTitle(pdf, "Vs previous period")
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
			line := fmt.Sprintf("%s: %.2f", measure, ts.CurrentPeriod)
			if ts.PreviousPeriod != nil {
				line += fmt.Sprintf(" (prev: %.2f)", *ts.PreviousPeriod)
			}
			if ts.ChangePercentage != nil {
				line += fmt.Sprintf(" %+.1f%%", *ts.ChangePercentage)
			}
			line += " " + ts.Trend
			pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
			pdf.MultiCell(0, 10, line, "", "L", false)
		}
		pdf.Ln(4)
	}

	// Cohorts
	if report.Metrics != nil && len(report.Metrics.Cohorts) > 0 {
		sectionTitle(pdf, "Cohorts")
		for _, c := range report.Metrics.Cohorts {
			if c == nil {
				continue
			}
			pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "B", 10)
			pdf.MultiCell(0, 10, safePDFString(c.CohortLabel), "", "L", false)
			pdf.SetFont("Helvetica", "", 10)
			if c.RetentionPct != nil {
				pdf.CellFormat(20, 10, "", "", 0, "L", false, 0, "")
				pdf.MultiCell(0, 10, fmt.Sprintf("Retention: %.1f%%", *c.RetentionPct), "", "L", false)
			}
			if len(c.Periods) > 0 {
				w0, w1 := 60.0, 80.0
				pdf.SetFont("Helvetica", "B", 9)
				pdf.CellFormat(w0, 10, "Period", "1", 0, "L", true, 0, "")
				pdf.CellFormat(w1, 10, "Value", "1", 1, "R", true, 0, "")
				pdf.SetFont("Helvetica", "", 9)
				for _, p := range c.Periods {
					if p == nil {
						continue
					}
					pdf.CellFormat(w0, 10, safePDFString(p.PeriodLabel), "1", 0, "L", false, 0, "")
					pdf.CellFormat(w1, 10, fmt.Sprintf("%.2f", p.Value), "1", 1, "R", false, 0, "")
				}
				pdf.Ln(2)
			}
		}
		pdf.Ln(4)
	}

	// Data quality table
	if report.Metrics != nil && len(report.Metrics.DataQuality) > 0 {
		sectionTitle(pdf, "Data quality")
		cols := make([]string, 0, len(report.Metrics.DataQuality))
		for c := range report.Metrics.DataQuality {
			cols = append(cols, c)
		}
		sort.Strings(cols)
		// Table header
		w0, w1, w2, w3, w4 := 80.0, 45.0, 50.0, 55.0, 45.0
		pdf.SetFont("Helvetica", "B", 9)
		pdf.CellFormat(w0, 10, "Column", "1", 0, "L", true, 0, "")
		pdf.CellFormat(w1, 10, "Nulls", "1", 0, "R", true, 0, "")
		pdf.CellFormat(w2, 10, "Null %", "1", 0, "R", true, 0, "")
		pdf.CellFormat(w3, 10, "Distinct", "1", 0, "R", true, 0, "")
		pdf.CellFormat(w4, 10, "Rows", "1", 1, "R", true, 0, "")
		pdf.SetFont("Helvetica", "", 9)
		for _, col := range cols {
			q := report.Metrics.DataQuality[col]
			if q == nil {
				continue
			}
			pdf.CellFormat(w0, 10, safePDFString(col), "1", 0, "L", false, 0, "")
			pdf.CellFormat(w1, 10, fmt.Sprintf("%d", q.NullCount), "1", 0, "R", false, 0, "")
			pdf.CellFormat(w2, 10, fmt.Sprintf("%.1f", q.NullPct), "1", 0, "R", false, 0, "")
			pdf.CellFormat(w3, 10, fmt.Sprintf("%d", q.DistinctCount), "1", 0, "R", false, 0, "")
			pdf.CellFormat(w4, 10, fmt.Sprintf("%d", q.TotalRows), "1", 1, "R", false, 0, "")
		}
		pdf.Ln(4)
	}

	// Performance suggestions
	if report.Metrics != nil && len(report.Metrics.PerfSuggestions) > 0 {
		sectionTitle(pdf, "Performance suggestions")
		for _, s := range report.Metrics.PerfSuggestions {
			pdf.CellFormat(12, 10, "", "", 0, "L", false, 0, "")
			pdf.MultiCell(0, 10, safePDFString(s), "", "L", false)
		}
	}

	return pdf.Output(w)
}

func sectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "", 10)
}

// drawMetricsBarChart draws a simple bar chart from time series metrics (current period values).
func drawMetricsBarChart(pdf *gofpdf.Fpdf, timeSeries map[string]*reports.TimeSeriesData) {
	sectionTitle(pdf, "Metrics overview (chart)")
	measures := make([]string, 0, len(timeSeries))
	var maxVal float64
	for m, ts := range timeSeries {
		if ts == nil {
			continue
		}
		measures = append(measures, m)
		if ts.CurrentPeriod > maxVal {
			maxVal = ts.CurrentPeriod
		}
	}
	if len(measures) == 0 {
		return
	}
	sort.Strings(measures)
	const (
		chartLeft   = 20.0
		labelWidth  = 90.0
		barAreaW    = 280.0
		barHeight   = 14.0
		barGap      = 4.0
		barMaxWidth = barAreaW - 10
	)
	if maxVal <= 0 {
		maxVal = 1
	}
	pdf.SetFont("Helvetica", "", 9)
	for _, measure := range measures {
		ts := timeSeries[measure]
		if ts == nil {
			continue
		}
		v := ts.CurrentPeriod
		barW := (v / maxVal) * barMaxWidth
		if barW < 2 && v > 0 {
			barW = 2
		}
		label := safePDFString(measure)
		if len(label) > 18 {
			label = label[:15] + "..."
		}
		pdf.CellFormat(chartLeft, barHeight, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(labelWidth, barHeight, label, "", 0, "L", false, 0, "")
		pdf.SetFillColor(70, 130, 180) // steel blue
		pdf.Rect(chartLeft+labelWidth, pdf.GetY()+2, barW, barHeight-4, "F")
		pdf.CellFormat(barAreaW-barW, barHeight, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, barHeight, fmt.Sprintf("%.2f", v), "", 1, "R", false, 0, "")
		pdf.SetFillColor(255, 255, 255)
	}
}

// safePDFString returns a string safe for gofpdf (Helvetica is ASCII/Latin-1). Non-ASCII runes are replaced with '?'.
func safePDFString(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 128 || (r >= 160 && r <= 255) {
			b.WriteRune(r)
		} else {
			b.WriteRune('?')
		}
	}
	return b.String()
}
