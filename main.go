package main

import (
	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/shutdown"
	"isp-script-service/assembly"
	"isp-script-service/conf"
	"isp-script-service/routes"
)

var version = "1.0.0"

// @title isp-script-service
// @version 1.0.0
// @description Сервис для обработки и выполнения JavaScript скриптов

// @license.name GNU GPL v3.0

// @host localhost:9000
// @BasePath /api/script

//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json

func main() {
	boot := bootstrap.New(version, conf.Remote{}, routes.EndpointDescriptors())
	app := boot.App
	logger := app.Logger()

	assembly, err := assembly.New(boot)
	if err != nil {
		logger.Fatal(app.Context(), err)
	}
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})

	err = app.Run()
	if err != nil {
		app.Shutdown()
		logger.Fatal(app.Context(), err)
	}
}
