package design

import (
	. "goa.design/goa/v3/dsl"
)

// Schema service exposes the database catalog (allowed schemas, tables, columns)
// for use by MCP and other clients to discover queryable objects.
var _ = Service("schema", func() {
	Description("Database schema and catalog for allowed queryable objects")

	Method("get", func() {
		Description("Return the list of allowed schemas with their tables and columns (from information_schema, read-only).")
		Payload(func() {
			Attribute("connection_id", String)
		})
		Result(SchemaResult)
		HTTP(func() {
			GET("/api/v1/schema")
			Params(func() {
				Param("connection_id")
			})
			Response(StatusOK)
		})
	})
})

// SchemaResult is the full catalog of allowed schemas and their tables/columns.
var SchemaResult = Type("SchemaResult", func() {
	Attribute("schemas", ArrayOf(SchemaInfo), "Allowed schemas with tables and columns")
	Required("schemas")
})

// SchemaInfo describes one schema and its tables.
var SchemaInfo = Type("SchemaInfo", func() {
	Attribute("name", String, "Schema name (e.g. demo)")
	Attribute("tables", ArrayOf(TableInfo), "Tables in this schema")
	Required("name", "tables")
})

// TableInfo describes one table and its columns.
var TableInfo = Type("TableInfo", func() {
	Attribute("name", String, "Table name")
	Attribute("columns", ArrayOf(ColumnInfo), "Columns (name and type)")
	Required("name", "columns")
})
