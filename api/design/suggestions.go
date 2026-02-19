package design

import (
	. "goa.design/goa/v3/dsl"
)

// Suggestions service provides query suggestions (curated examples and saved-query matches).
var _ = Service("suggestions", func() {
	Description("Query suggestions for MCP and other clients (curated examples and saved-query match by intent)")

	Method("queries", func() {
		Description("Return suggested SQL queries: curated examples plus saved queries matching optional intent.")
		Payload(func() {
			Attribute("intent", String, "Optional natural-language intent (e.g. 'sales by region'); used to match saved queries.")
			Attribute("limit", Int32, "Max number of suggestions to return", func() {
				Default(5)
				Minimum(1)
				Maximum(20)
			})
		})
		Result(SuggestedQueriesResult)
		HTTP(func() {
			GET("/api/v1/suggestions/queries")
			Params(func() {
				Param("intent")
				Param("limit")
			})
			Response(StatusOK)
		})
	})
})

// SuggestedQueriesResult is the result of the suggestions queries method.
var SuggestedQueriesResult = Type("SuggestedQueriesResult", func() {
	Attribute("suggestions", ArrayOf(QuerySuggestion), "Suggested SQL and metadata")
	Required("suggestions")
})

// QuerySuggestion is one suggested query (sql, title, source).
var QuerySuggestion = Type("QuerySuggestion", func() {
	Attribute("sql", String, "Suggested SQL (use with run_query or refine)")
	Attribute("title", String, "Short label for the suggestion")
	Attribute("source", String, "Where the suggestion came from: curated or saved")
	Required("sql", "title", "source")
})
