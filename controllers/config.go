package controllers

import (
	"github.com/MartinSahlen/installer/app"
	"github.com/MartinSahlen/installer/brew"
	"github.com/goadesign/goa"
)

// ConfigController implements the config resource.
type ConfigController struct {
	*goa.Controller
	db *brew.DB
}

// NewConfigController creates a config controller.
func NewConfigController(service *goa.Service, db *brew.DB) *ConfigController {
	return &ConfigController{Controller: service.NewController("ConfigController"), db: db}
}

// Create runs the create action.
func (c *ConfigController) Create(ctx *app.CreateConfigContext) error {

	deps := []brew.Dependency{}

	for _, d := range ctx.Payload.Config {
		dep := brew.Dependency{
			Name:     d.Name,
			FullName: d.FullName,
		}
		if d.Type == "BREW" {
			dep.Type = brew.Brew
		}
		if d.Type == "BREW_CASK" {
			dep.Type = brew.BrewCask
		}
		deps = append(deps, dep)
	}

	set, err := c.db.SaveDependencies(deps)

	if err != nil {
		ctx.WriteHeader(500)
		return nil
	}
	return ctx.Created(&app.SavedConfig{
		Dependencies: ctx.Payload.Config,
		ID:           &set.ID,
	})
}

// Delete runs the delete action.
func (c *ConfigController) Delete(ctx *app.DeleteConfigContext) error {
	// ConfigController_Delete: start_implement

	// Put your logic here

	// ConfigController_Delete: end_implement
	res := &app.Config{}
	return ctx.OK(res)
}

// Get runs the get action.
func (c *ConfigController) Get(ctx *app.GetConfigContext) error {
	// ConfigController_Get: start_implement
	deps, err := c.db.GetDependenciesForID(ctx.ConfigID)
	if err != nil {
		ctx.WriteHeader(500)
		return nil
	}
	d := app.SavedConfig{
		ID:           &ctx.ConfigID,
		Dependencies: []*app.ConfigType{},
	}

	for _, dd := range deps {
		var t string = ""
		if dd.Type == brew.Brew {
			t = "BREW"
		}
		if dd.Type == brew.BrewCask {
			t = "BREW_CASK"
		}
		d.Dependencies = append(d.Dependencies, &app.ConfigType{
			Name:     dd.Name,
			Type:     t,
			FullName: dd.FullName,
		})
	}

	// Put your logic here
	return ctx.OK(&d)
	// ConfigController_Get: end_implement
}

// Update runs the update action.
func (c *ConfigController) Update(ctx *app.UpdateConfigContext) error {
	// ConfigController_Update: start_implement

	// Put your logic here

	// ConfigController_Update: end_implement
	return nil
}
