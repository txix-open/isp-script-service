package main

import (
	"context"
	"os"

	"isp-script-service/conf"
	_ "isp-script-service/docs"
	"isp-script-service/helper"
	"isp-script-service/router"
	"isp-script-service/service"

	"github.com/integration-system/isp-lib/v2/backend"
	"github.com/integration-system/isp-lib/v2/bootstrap"
	"github.com/integration-system/isp-lib/v2/config/schema"
	"github.com/integration-system/isp-lib/v2/metric"
	"github.com/integration-system/isp-lib/v2/structure"
	log "github.com/integration-system/isp-log"
	"github.com/integration-system/isp-log/stdcodes"
	"google.golang.org/grpc"
)

var version = "1.0.0"

// @title isp-script-service
// @version 1.0.0
// @description Сервис для обработка и выполения JavaScript скриптов

// @license.name GNU GPL v3.0

// @host localhost:9000
// @BasePath /api/script

func main() {
	bootstrap.
		ServiceBootstrap(&conf.Configuration{}, &conf.RemoteConfig{}).
		OnLocalConfigLoad(onLocalConfigLoad).
		SocketConfiguration(socketConfiguration).
		DefaultRemoteConfigPath(schema.ResolveDefaultConfigPath("default_remote_config.json")).
		OnSocketErrorReceive(onRemoteErrorReceive).
		OnConfigErrorReceive(onRemoteConfigErrorReceive).
		DeclareMe(routesData).
		RequireModule("router", router.Client.ReceiveAddressList, true).
		OnRemoteConfigReceive(onRemoteConfigReceive).
		OnShutdown(onShutdown).
		Run()
}

func socketConfiguration(cfg interface{}) structure.SocketConfiguration {
	appConfig := cfg.(*conf.Configuration)

	return structure.SocketConfiguration{
		Host:   appConfig.ConfigServiceAddress.IP,
		Port:   appConfig.ConfigServiceAddress.Port,
		Secure: false,
		UrlParams: map[string]string{
			"module_name":   appConfig.ModuleName,
			"instance_uuid": appConfig.InstanceUuid,
		},
	}
}

func onShutdown(_ context.Context, _ os.Signal) {
	backend.StopGrpcServer()
}

func onRemoteConfigReceive(remoteConfig, oldRemoteConfig *conf.RemoteConfig) {
	service.Script.ReceiveConfiguration(remoteConfig.Scripts)
	metric.InitCollectors(remoteConfig.Metrics, oldRemoteConfig.Metrics)
	metric.InitHttpServer(remoteConfig.Metrics)
}

func onLocalConfigLoad(cfg *conf.Configuration) {
	const msgSize = 1024 * 1024 * 512
	metric.InitProfiling(cfg.ModuleName)
	handlers := helper.GetAllHandlers()
	service := backend.GetDefaultService(cfg.ModuleName, handlers...)
	backend.StartBackendGrpcServer(
		cfg.GrpcInnerAddress, service,
		grpc.MaxRecvMsgSize(msgSize),
		grpc.MaxSendMsgSize(msgSize),
	)
}

func routesData(localConfig interface{}) bootstrap.ModuleInfo {
	cfg := localConfig.(*conf.Configuration)

	return bootstrap.ModuleInfo{
		ModuleName:       cfg.ModuleName,
		ModuleVersion:    version,
		GrpcOuterAddress: cfg.GrpcOuterAddress,
		Handlers:         helper.GetAllHandlers(),
	}
}

func onRemoteErrorReceive(errorMessage map[string]interface{}) {
	log.WithMetadata(errorMessage).Error(stdcodes.ReceiveErrorFromConfig, "error from config service")
}

func onRemoteConfigErrorReceive(errorMessage string) {
	log.WithMetadata(map[string]interface{}{
		"message": errorMessage,
	}).Error(stdcodes.ReceiveErrorOnGettingConfigFromConfig, "error on getting remote configuration")
}
