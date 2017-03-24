package controllers

import (
	"fmt"
	"net/http"

	"github.com/MartinSahlen/installer/app"
	"github.com/MartinSahlen/installer/brew"
	"github.com/MartinSahlen/installer/install"
	"github.com/goadesign/goa"
)

// AppController implements the app resource.
type AppController struct {
	*goa.Controller
	db *brew.DB
}

// NewAppController creates a app controller.
func NewAppController(service *goa.Service, db *brew.DB) *AppController {
	return &AppController{Controller: service.NewController("AppController"), db: db}
}

// Get runs the get action.
func (c *AppController) Get(ctx *app.GetAppContext) error {
	// AppController_Get: start_implement

	// Put your logic here
	installScript, err := install.GenerateInstallScript(c.db, ctx.ConfigID)

	if err != nil {
		ctx.WriteHeader(http.StatusInternalServerError)
		ctx.Write(nil)
	}

	ctx.ResponseData.Header().Set("Content-Type", "application/zip")
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Installer.zip"))
	err = install.ArchiveInstallApp("./Installer", installScript, ctx.ResponseData.ResponseWriter)
	if err != nil {
		ctx.ResponseData.WriteHeader(http.StatusInternalServerError)
		ctx.ResponseData.Write(nil)
	}
	return nil
}
