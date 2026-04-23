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

	Method("questions", func() {
		Description("Suggest natural-language questions based on the current schema for discovery and onboarding.")
		Payload(func() {
			Attribute("connection_id", String, "Optional connection ID; defaults to server default connection")
			Attribute("limit", Int32, "Max questions to return", func() {
				Default(8)
				Minimum(1)
				Maximum(20)
			})
		})
		Result(SuggestedQuestionsResult)
		HTTP(func() {
			GET("/api/v1/suggestions/questions")
			Params(func() {
				Param("connection_id")
				Param("limit")
			})
			Response(StatusOK)
		})
	})

	Method("similar", func() {
		Description("Return saved queries semantically similar to the given text (embedding-based). Requires embeddings to be enabled.")
		Payload(func() {
			Attribute("text", String, "Natural-language or SQL text to find similar queries for.")
			Attribute("limit", Int32, "Max number of similar queries to return", func() {
				Default(5)
				Minimum(1)
				Maximum(20)
			})
		})
		Result(SuggestedQueriesResult)
		HTTP(func() {
			GET("/api/v1/suggestions/similar")
			Params(func() {
				Param("text")
				Param("limit")
			})
			Response(StatusOK)
		})
	})

	Method("ask", func() {
		Description("Natural language to SQL and report in one step: ask a question, get a generated SELECT, run it, and return the narrative report. Requires LLM.")
		Payload(func() {
			Attribute("question", String, "Natural-language question (e.g. 'What were top 5 products by revenue last month?')", func() {
				MinLength(1)
				MaxLength(1000)
			})
			Attribute("connection_id", String, "Optional connection ID; defaults to server default connection")
			Required("question")
		})
		Result(AskResult)
		Error("validation_error", ValidationError)
		Error("llm_error", LLMError)
		HTTP(func() {
			POST("/api/v1/suggestions/ask")
			Response(StatusOK)
			Response(StatusBadRequest, "validation_error")
			Response(StatusInternalServerError, "llm_error")
		})
	})

	Method("explain", func() {
		Description("Explain a SQL query in plain English (one or two sentences). Requires LLM.")
		Payload(func() {
			Attribute("sql", String, "Read-only SQL to explain (SELECT or WITH).", func() {
				MinLength(1)
				MaxLength(10000)
			})
			Required("sql")
		})
		Result(ExplainResult)
		Error("validation_error", ValidationError)
		Error("llm_error", LLMError)
		HTTP(func() {
			POST("/api/v1/suggestions/explain")
			Response(StatusOK)
			Response(StatusBadRequest, "validation_error")
			Response(StatusInternalServerError, "llm_error")
		})
	})
})

// ExplainResult is the result of the suggestions explain method.
var ExplainResult = Type("ExplainResult", func() {
	Attribute("sql", String, "The SQL that was explained")
	Attribute("explanation", String, "Plain-English explanation (one or two sentences)")
	Required("sql", "explanation")
})

// AskResult is the result of the suggestions ask method (NL → SQL → report).
var AskResult = Type("AskResult", func() {
	Attribute("question", String, "The original natural-language question")
	Attribute("sql", String, "The generated and executed SQL")
	Attribute("report", Report, "The narrative report from the query result")
	Required("question", "sql", "report")
})

// SuggestedQueriesResult is the result of the suggestions queries method.
var SuggestedQueriesResult = Type("SuggestedQueriesResult", func() {
	Attribute("suggestions", ArrayOf(QuerySuggestion), "Suggested SQL and metadata")
	Required("suggestions")
})

// SuggestedQuestionsResult is schema-driven natural-language discovery prompts.
var SuggestedQuestionsResult = Type("SuggestedQuestionsResult", func() {
	Attribute("questions", ArrayOf(String))
	Required("questions")
})

// QuerySuggestion is one suggested query (sql, title, source).
var QuerySuggestion = Type("QuerySuggestion", func() {
	Attribute("sql", String, "Suggested SQL (use with run_query or refine)")
	Attribute("title", String, "Short label for the suggestion")
	Attribute("source", String, "Where the suggestion came from: curated or saved")
	Required("sql", "title", "source")
})
