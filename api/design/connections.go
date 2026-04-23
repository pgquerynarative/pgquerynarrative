package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("connections", func() {
	Description("Available read-only data connections")

	Method("list", func() {
		Description("List configured data connections (id and display name only)")
		Result(ConnectionListResult)
		HTTP(func() {
			GET("/api/v1/connections")
			Response(StatusOK)
		})
	})
})

var ConnectionInfo = Type("ConnectionInfo", func() {
	Attribute("id", String)
	Attribute("name", String)
	Required("id", "name")
})

var ConnectionListResult = Type("ConnectionListResult", func() {
	Attribute("items", ArrayOf(ConnectionInfo))
	Required("items")
})
