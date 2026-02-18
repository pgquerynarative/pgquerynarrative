package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("reports", func() {
	Description("Report generation and management")

	Method("generate", func() {
		Description("Generate a narrative report from a SQL query")
		Payload(GenerateReportPayload)
		Result(Report)
		Error("validation_error", ValidationError)
		Error("llm_error", LLMError)
		HTTP(func() {
			POST("/api/v1/reports/generate")
			Response(StatusOK)
			Response(StatusBadRequest, "validation_error")
			Response(StatusInternalServerError, "llm_error")
		})
	})

	Method("get", func() {
		Description("Get a report by ID")
		Payload(func() {
			Attribute("id", String, func() {
				Format(FormatUUID)
			})
			Required("id")
		})
		Result(Report)
		Error("not_found", NotFoundError)
		HTTP(func() {
			GET("/api/v1/reports/{id}")
			Response(StatusOK)
			Response(StatusNotFound, "not_found")
		})
	})

	Method("list", func() {
		Description("List generated reports")
		Payload(func() {
			Attribute("saved_query_id", String, func() {
				Format(FormatUUID)
			})
			Attribute("limit", Int32, func() {
				Default(50)
				Minimum(1)
				Maximum(100)
			})
			Attribute("offset", Int32, func() {
				Default(0)
				Minimum(0)
			})
		})
		Result(ReportList)
		HTTP(func() {
			GET("/api/v1/reports")
			Params(func() {
				Param("saved_query_id")
				Param("limit")
				Param("offset")
			})
		})
	})
})

var GenerateReportPayload = Type("GenerateReportPayload", func() {
	Attribute("sql", String, func() {
		MinLength(1)
		MaxLength(10000)
		Pattern("^[^;]+$")
	})
	Attribute("saved_query_id", String, func() {
		Format(FormatUUID)
	})
	Required("sql")
})

var Report = Type("Report", func() {
	Attribute("id", String, func() {
		Format(FormatUUID)
	})
	Attribute("saved_query_id", String, func() {
		Format(FormatUUID)
	})
	Attribute("sql", String)
	Attribute("narrative", NarrativeContent)
	Attribute("metrics", MetricsData)
	Attribute("chart_suggestions", ArrayOf(ChartSuggestion), "Suggested chart types based on result shape")
	Attribute("created_at", String, func() {
		Format(FormatDateTime)
	})
	Attribute("llm_model", String)
	Attribute("llm_provider", String)
	Required("id", "sql", "narrative", "metrics", "created_at", "llm_model", "llm_provider")
})

var NarrativeContent = Type("NarrativeContent", func() {
	Attribute("headline", String)
	Attribute("takeaways", ArrayOf(String))
	Attribute("drivers", ArrayOf(String))
	Attribute("limitations", ArrayOf(String))
	Attribute("recommendations", ArrayOf(String))
	Required("headline", "takeaways")
})

var MetricsData = Type("MetricsData", func() {
	Attribute("aggregates", MapOf(String, AggregateData))
	Attribute("top_categories", MapOf(String, ArrayOf(TopCategoryData)))
	Attribute("time_series", MapOf(String, TimeSeriesData))
	Attribute("period_current_label", String, "Label for current period when time_series is present")
	Attribute("period_previous_label", String, "Label for previous period")
})

var AggregateData = Type("AggregateData", func() {
	Attribute("sum", Float64)
	Attribute("avg", Float64)
	Attribute("min", Float64)
	Attribute("max", Float64)
	Attribute("count", Int32)
})

var TopCategoryData = Type("TopCategoryData", func() {
	Attribute("category", String)
	Attribute("value", Float64)
	Attribute("percentage", Float64)
	Required("category", "value", "percentage")
})

var TimeSeriesData = Type("TimeSeriesData", func() {
	Attribute("current_period", Float64)
	Attribute("previous_period", Float64)
	Attribute("change", Float64)
	Attribute("change_percentage", Float64)
	Attribute("trend", String)
	Required("current_period", "trend")
})

var ReportList = Type("ReportList", func() {
	Attribute("items", ArrayOf(Report))
	Attribute("limit", Int32)
	Attribute("offset", Int32)
	Required("items", "limit", "offset")
})

var LLMError = Type("LLMError", func() {
	Attribute("name", String)
	Attribute("message", String)
	Attribute("code", String)
	Required("name", "message")
})
