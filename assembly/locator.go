package assembly

import (
	"github.com/integration-system/isp-kit/grpc/client"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/grpc/isp"
	"github.com/integration-system/isp-kit/log"
	"isp-script-service/conf"
	"isp-script-service/controller"
	"isp-script-service/repository"
	"isp-script-service/routes"
	"isp-script-service/service"
)

type Locator struct {
	logger  log.Logger
	grpcCli *client.Client
}

func NewLocator(logger log.Logger, grpcCli *client.Client) Locator {
	return Locator{
		logger:  logger,
		grpcCli: grpcCli,
	}
}

func (l Locator) Handler(cfg conf.Remote) isp.BackendServiceServer {
	router := repository.NewRouter(l.grpcCli)
	scriptService := service.NewScript(router, l.logger, cfg.Scripts, cfg.SharedScript, cfg.ScriptExecutionTimeoutMs)
	scriptController := controller.NewScript(scriptService)

	handler := routes.Handler(
		endpoint.DefaultWrapper(
			l.logger,
		),
		routes.Controllers{
			Script: scriptController,
		},
	)

	return handler
}
