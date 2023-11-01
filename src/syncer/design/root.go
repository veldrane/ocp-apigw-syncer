package design

import . "goa.design/goa/v3/dsl"

var DefaultStatus = ResultType("application/vnd.invoices.defaultredirect", func() {
	Attributes(func() {
		Attribute("Location", String)
	})
})

var _ = Service("root", func() {

	Description("Service provide redirect to swagger-ui")
	HTTP(func() {
		Path("/")
	})
	Meta("swagger:generate", "false")
	Method("default", func() {

		Description("Return default redirect")
		Result(func() {
			Attribute("Location", String, "Default location string")
		})

		HTTP(func() {
			GET("/")
			Response(func() {
				Description("Return redirect to swagger-ui")
				Code(StatusFound)
				Header("Location")
			})
		})
	})
})
