package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("config", func() {
	Action("create", func() {
		Routing(
			POST("config"),
		)
		Description("Create a config that contains a list of apps")
		Payload(configPayload)
		Response(Created, savedConfigPayload)
		Response(BadRequest)
	})

	Action("get", func() {
		Routing(
			GET("config/:configID"),
		)
		Params(func() {
			Param("configID", String)
		})
		Description("Get a config that contains a list of apps")
		Response(OK, savedConfigPayload)
	})

	Action("update", func() {
		Routing(
			PUT("config/:configID"),
		)
		Description("updatet a config that contains a list of apps")
		Params(func() {
			Param("configID", String)
		})
		Payload(configPayload)
		Response(OK)
	})

	Action("delete", func() {
		Routing(
			DELETE("config/:configID"),
		)
		Description("Delete a config that contains a list of apps")
		Params(func() {
			Param("configID", String)
		})
		Payload(configPayload)
		Response(OK, configPayload)
	})
})

var savedConfigPayload = MediaType("application/vnd.savedconfiglist", func() {
	Description("Payload containg a list of apps to install")
	TypeName("savedConfig")
	Attributes(func() {
		Attribute("dependencies", ArrayOf(configType))
		Attribute("id", String)
	})

	View("default", func() {
		Attribute("dependencies")
		Attribute("id")
	})
})

var configPayload = MediaType("application/vnd.configlist", func() {
	Description("Payload containg a list of apps to install")
	TypeName("config")
	Attributes(func() {
		Attribute("config", ArrayOf(configType))
	})

	View("default", func() {
		Attribute("config")
	})
	Required("config")
})

var configType = Type("configType", func() {
	Description("A config that is an app to install")
	Attribute("type", func() {
		Enum("BREW", "BREW_CASK")
	})
	Attribute("fullName", String)
	Attribute("name", String)
	Required("type", "fullName", "name")
})
