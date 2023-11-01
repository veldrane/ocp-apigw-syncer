package design

import . "goa.design/goa/v3/dsl"

var _ = Service("swagger", func() {
	Files("/swagger-ui/index.html", "public/swagger/index.html", func() {
		Description("Serve home page.")
		Docs(func() {
			Description("Additional documentation")
			URL("https://goa.design")
		})
	})
	Files("/swagger-ui/{*path}", "public/swagger", func() {
		Description("Serve static content.")
	})
	Meta("swagger:generate", "false")
})
