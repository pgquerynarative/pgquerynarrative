package design

import (
	. "goa.design/goa/v3/dsl"
)

var RunQueryPayload = Type("RunQueryPayload", func() {
	Attribute("sql", String, "SQL query to execute", func() {
		MinLength(1)
		MaxLength(10000)
		Pattern("^[^;]+$")
	})
	Attribute("limit", Int32, "Maximum number of rows to return", func() {
		Default(100)
		Minimum(1)
		Maximum(1000)
	})
	Attribute("connection_id", String, "Optional connection ID; defaults to server default connection")
	Required("sql")
})

var RunQueryResult = Type("RunQueryResult", func() {
	Attribute("columns", ArrayOf(ColumnInfo))
	Attribute("rows", ArrayOf(ArrayOf(Any)))
	Attribute("row_count", Int32)
	Attribute("execution_time_ms", Int64)
	Attribute("limit", Int32)
	Attribute("chart_suggestions", ArrayOf(ChartSuggestion), "Suggested chart types based on result shape")
	Attribute("period_comparison", ArrayOf(PeriodComparisonItem), "Period-over-period comparison when result has time + measure columns")
	Attribute("period_current_label", String, "Label for current period when period_comparison is present (e.g. date or month)")
	Attribute("period_previous_label", String, "Label for previous period")
	Required("columns", "rows", "row_count", "execution_time_ms", "limit")
})

// PeriodComparisonItem is one measure's current vs previous period (e.g. this month vs last month).
var PeriodComparisonItem = Type("PeriodComparisonItem", func() {
	Attribute("measure", String, "Measure column name")
	Attribute("current", Float64, "Current period value")
	Attribute("previous", Float64, "Previous period value")
	Attribute("change", Float64, "Absolute change (current - previous)")
	Attribute("change_percentage", Float64, "Percent change vs previous")
	Attribute("trend", String, "up, down, or flat")
	Required("measure", "current", "trend")
})

// ChartSuggestion describes a chart type suggested from the query result shape.
var ChartSuggestion = Type("ChartSuggestion", func() {
	Attribute("chart_type", String, "Chart type identifier: bar, line, pie, area, table")
	Attribute("label", String, "Human-readable label")
	Attribute("reason", String, "Why this chart fits the data")
	Required("chart_type", "label", "reason")
})

var ColumnInfo = Type("ColumnInfo", func() {
	Attribute("name", String)
	Attribute("type", String)
	Required("name", "type")
})

var ValidationError = Type("ValidationError", func() {
	Attribute("name", String)
	Attribute("message", String)
	Attribute("code", String)
	Required("name", "message")
})

var NotFoundError = Type("NotFoundError", func() {
	Attribute("name", String)
	Attribute("message", String)
	Attribute("code", String)
	Required("name", "message")
})

var SaveQueryPayload = Type("SaveQueryPayload", func() {
	Attribute("name", String, func() {
		MinLength(1)
		MaxLength(200)
	})
	Attribute("sql", String, func() {
		MinLength(1)
		MaxLength(10000)
	})
	Attribute("description", String, func() {
		MaxLength(500)
	})
	Attribute("tags", ArrayOf(String))
	Attribute("connection_id", String, "Optional connection ID; defaults to server default connection")
	Required("name", "sql")
})

var SavedQuery = Type("SavedQuery", func() {
	Attribute("id", String, func() {
		Format(FormatUUID)
	})
	Attribute("name", String)
	Attribute("sql", String)
	Attribute("description", String)
	Attribute("tags", ArrayOf(String))
	Attribute("connection_id", String)
	Attribute("created_at", String, func() {
		Format(FormatDateTime)
	})
	Attribute("updated_at", String, func() {
		Format(FormatDateTime)
	})
	Required("id", "name", "sql", "connection_id", "created_at", "updated_at")
})

var SavedQueryList = Type("SavedQueryList", func() {
	Attribute("items", ArrayOf(SavedQuery))
	Attribute("limit", Int32)
	Attribute("offset", Int32)
	Required("items", "limit", "offset")
})
