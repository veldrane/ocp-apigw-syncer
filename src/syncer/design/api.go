package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("syncer", func() {
	Title("Syncer")
	Description("Service achieving atomic synchronization for jwt token in ng-plus distributed environment")
	Server("syncer", func() {
		Host("", func() {
			URI("http://localhost:8080")
		})
	})
})
