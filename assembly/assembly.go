package assembly

import (
	"context"

	stdgrpc "google.golang.org/grpc"

	"github.com/integration-system/isp-kit/app"
	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/client"
	"github.com/integration-system/isp-kit/log"
	"github.com/pkg/errors"
	"isp-script-service/conf"
)

const GrpcMaxRecvMsgSize = 64 * 1024 * 1024

type Assembly struct {
	boot   *bootstrap.Bootstrap
	server *grpc.Server
	logger *log.Adapter
	router *client.Client
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	server := grpc.NewServer(stdgrpc.MaxRecvMsgSize(GrpcMaxRecvMsgSize))
	router, err := client.Default()
	if err != nil {
		return nil, errors.WithMessage(err, "create repository client")
	}
	return &Assembly{
		boot:   boot,
		server: server,
		logger: boot.App.Logger(),
		router: router,
	}, nil
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) error {
	var (
		newCfg  conf.Remote
		prevCfg conf.Remote
	)
	err := a.boot.RemoteConfig.Upgrade(remoteConfig, &newCfg, &prevCfg)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "upgrade remote config"))
	}

	a.logger.SetLevel(newCfg.LogLevel)

	locator := NewLocator(a.logger, a.router)
	handler := locator.Handler(newCfg)

	a.server.Upgrade(handler)

	return nil
}

func (a *Assembly) Runners() []app.Runner {
	eventHandler := cluster.NewEventHandler().
		RequireModule("repository", a.router).
		RemoteConfigReceiver(a)
	return []app.Runner{
		app.RunnerFunc(func(ctx context.Context) error {
			return a.server.ListenAndServe(a.boot.BindingAddress)
		}),
		app.RunnerFunc(func(ctx context.Context) error {
			return a.boot.ClusterCli.Run(ctx, eventHandler)
		}),
	}
}

func (a *Assembly) Closers() []app.Closer {
	return []app.Closer{
		a.boot.ClusterCli,
		app.CloserFunc(func() error {
			a.server.Shutdown()
			return nil
		}),
		a.router,
	}
}
