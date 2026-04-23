package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("schedules", func() {
	Description("Scheduled report generation and delivery")

	Method("list", func() {
		Result(ScheduleListResult)
		HTTP(func() {
			GET("/api/v1/schedules")
			Response(StatusOK)
		})
	})

	Method("create", func() {
		Payload(ScheduleInput)
		Result(Schedule)
		Error("validation_error", ValidationError)
		HTTP(func() {
			POST("/api/v1/schedules")
			Response(StatusOK)
			Response(StatusBadRequest, "validation_error")
		})
	})

	Method("update", func() {
		Payload(func() {
			Attribute("id", String, func() { Format(FormatUUID) })
			Extend(ScheduleInput)
			Required("id")
		})
		Result(Schedule)
		Error("not_found", NotFoundError)
		Error("validation_error", ValidationError)
		HTTP(func() {
			PUT("/api/v1/schedules/{id}")
			Response(StatusOK)
			Response(StatusNotFound, "not_found")
			Response(StatusBadRequest, "validation_error")
		})
	})

	Method("delete", func() {
		Payload(func() {
			Attribute("id", String, func() { Format(FormatUUID) })
			Required("id")
		})
		Result(Empty)
		HTTP(func() {
			DELETE("/api/v1/schedules/{id}")
			Response(StatusNoContent)
		})
	})

	Method("run_now", func() {
		Payload(func() {
			Attribute("id", String, func() { Format(FormatUUID) })
			Required("id")
		})
		Result(ScheduleRunResult)
		Error("not_found", NotFoundError)
		HTTP(func() {
			POST("/api/v1/schedules/{id}/run")
			Response(StatusOK)
			Response(StatusNotFound, "not_found")
		})
	})
})

var ScheduleInput = Type("ScheduleInput", func() {
	Attribute("name", String, func() { MinLength(1); MaxLength(200) })
	Attribute("saved_query_id", String, func() { Format(FormatUUID) })
	Attribute("sql", String, func() { MaxLength(10000) })
	Attribute("connection_id", String)
	Attribute("cron_expr", String, "Use @every <duration> format (e.g. @every 6h)")
	Attribute("destination_type", String, "webhook|log")
	Attribute("destination_target", String, "Webhook URL or log channel name")
	Attribute("enabled", Boolean)
	Required("name", "cron_expr", "destination_type", "destination_target")
})

var Schedule = Type("Schedule", func() {
	Attribute("id", String, func() { Format(FormatUUID) })
	Attribute("name", String)
	Attribute("saved_query_id", String, func() { Format(FormatUUID) })
	Attribute("sql", String)
	Attribute("connection_id", String)
	Attribute("cron_expr", String)
	Attribute("destination_type", String)
	Attribute("destination_target", String)
	Attribute("enabled", Boolean)
	Attribute("last_run_at", String, func() { Format(FormatDateTime) })
	Attribute("last_status", String)
	Attribute("last_error", String)
	Attribute("next_run_at", String, func() { Format(FormatDateTime) })
	Attribute("created_at", String, func() { Format(FormatDateTime) })
	Attribute("updated_at", String, func() { Format(FormatDateTime) })
	Required("id", "name", "connection_id", "cron_expr", "destination_type", "destination_target", "enabled", "created_at", "updated_at")
})

var ScheduleListResult = Type("ScheduleListResult", func() {
	Attribute("items", ArrayOf(Schedule))
	Required("items")
})

var ScheduleRunResult = Type("ScheduleRunResult", func() {
	Attribute("schedule", Schedule)
	Attribute("report_id", String, func() { Format(FormatUUID) })
	Attribute("delivered", Boolean)
	Required("schedule", "delivered")
})
