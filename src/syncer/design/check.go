package design

import (
	. "goa.design/goa/v3/dsl"
)

var NginxInstance = Type("app/vnd.nginx", func() {
	Description("Stored nginx instance")
	TypeName("NginxInstance")
	Attribute("hostname", String)
	Attribute("address", String)
	Attribute("port", String)
})

var NginxInstancies = ArrayOf(NginxInstance)

var ResultOfNginxInstance = ResultType("app/vnd.nginxs", func() {
	Reference(NginxInstance)
})

var ResultOfSync = ResultType("app/vnd.sync", func() {
	Attribute("status", String, func() {
		Default("synced")
	})
})

var _ = Service("checker", func() {
	Description("Check the replicas of nginx instancies")

	Error("NotFound", func() {
		Description("Notfound is the error returned by the service methods when the id of the stuff is not found.")
	})

	Error("InternalError", func() {
		Description("Internal Server Error")
	})

	Error("Unauthorized", func() {
		Description("Unauthorized")
	})

	Method("get", func() {

		Description("Get last full report")
		Result(ResultOfSync)

		HTTP(func() {
			GET("/v1/synced")
			Response(StatusOK)
		})

	})

})
