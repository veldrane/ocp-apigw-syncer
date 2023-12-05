package design

import (
	. "goa.design/goa/v3/dsl"
)

var ResultOfHealth = ResultType("app/vnd.health", func() {
	Attribute("status", String, func() {
		Default("Up")
	})
})

var _ = Service("health", func() {

	Method("get", func() {

		Description("Ping endpoin")
		Result(ResultOfHealth)

		HTTP(func() {
			GET("/health")
			Response(func() {
				Code(StatusOK)
			})
		})

	})
})
