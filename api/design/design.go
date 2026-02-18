package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("pgquerynarrative", func() {
	Title("PgQueryNarrative API")
	Description("Data Storyteller AI - Convert PostgreSQL SQL queries to business narratives")
	Version("v1")
	Server("pgquerynarrative", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})
