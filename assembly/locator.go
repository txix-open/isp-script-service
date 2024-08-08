// nolint:ireturn
package assembly

import (
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/log"
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

func (l Locator) Handler(cfg conf.Remote) *grpc.Mux {
	router := repository.NewRouter(l.grpcCli)
	scriptService := service.NewScript(router, l.logger, cfg.Scripts, cfg.SharedScript, cfg.ScriptExecutionTimeoutMs)
	scriptController := controller.NewScript(scriptService)

	c := routes.Controllers{
		Script: scriptController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	handler := routes.Handler(mapper, c)
	return handler
}
