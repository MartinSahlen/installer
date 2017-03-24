package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("app", func() {
	Action("get", func() {
		Routing(
			GET("install/:configID"),
		)
		Description("Get the app installer for your selection of apps")
		Params(func() {
			Param("configID", String, "Config ID")
		})
		Response(OK)
	})
})
