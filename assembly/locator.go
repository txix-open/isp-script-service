// nolint:ireturn
package assembly

import (
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/log"
	"isp-script-service/conf"
	"isp-script-service/controller"
	"isp-script-service/domain"
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

func (l Locator) Handler(cfg conf.Remote) (*grpc.Mux, error) {
	scripts, err := MergeScripts(cfg.Scripts)
	if err != nil {
		return nil, errors.WithMessage(err, "merge scripts")
	}

	router := repository.NewRouter(l.grpcCli)
	scriptService := service.NewScript(router, l.logger, scripts, cfg.SharedScript, cfg.ScriptExecutionTimeoutMs)
	scriptController := controller.NewScript(scriptService)

	c := routes.Controllers{
		Script: scriptController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	handler := routes.Handler(mapper, c)
	return handler, nil
}

func MergeScripts(allScripts conf.AllScripts) ([]conf.ScriptDefinition, error) {
	uniqueScripts := make(map[string]conf.ScriptDefinition)

	for _, script := range allScripts.CommonScripts {
		if _, exists := uniqueScripts[script.Id]; exists {
			return nil, domain.ErrDuplicateScriptId
		}
		uniqueScripts[script.Id] = script
	}

	for _, customScript := range allScripts.CustomScripts {
		if _, exists := uniqueScripts[customScript.Id]; exists {
			return nil, domain.ErrDuplicateScriptId
		}
		uniqueScripts[customScript.Id] = customScript
	}

	mergedScripts := make([]conf.ScriptDefinition, 0, len(uniqueScripts))
	for _, script := range uniqueScripts {
		mergedScripts = append(mergedScripts, script)
	}

	return mergedScripts, nil
}
