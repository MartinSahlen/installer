package main

import (
	"os"

	"github.com/MartinSahlen/installer/app"
	"github.com/MartinSahlen/installer/brew"
	"github.com/MartinSahlen/installer/controllers"
	"github.com/go-kit/kit/log"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/logging/kit"
	"github.com/goadesign/goa/middleware"
)

func main() {
	service := goa.New("installer")
	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)
	service.WithLogger(goakit.New(logger))

	// Setup basic middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Setup database connection
	db, err := brew.NewDB()
	if err != nil {
		panic(err)
	}

	configController := controllers.NewConfigController(service, db)
	app.MountConfigController(service, configController)

	appController := controllers.NewAppController(service, db)
	app.MountAppController(service, appController)

	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError(err.Error())
	}
}
