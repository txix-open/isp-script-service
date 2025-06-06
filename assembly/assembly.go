// nolint:wrapcheck
package assembly

import (
	"context"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/app"
	"github.com/txix-open/isp-kit/bootstrap"
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/log"
	"isp-script-service/conf"
)

type Assembly struct {
	boot   *bootstrap.Bootstrap
	server *grpc.Server
	logger *log.Adapter
	router *client.Client
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	server := grpc.DefaultServer()
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
	handler, err := locator.Handler(newCfg)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "new handler"))
	}

	a.server.Upgrade(handler)

	return nil
}

func (a *Assembly) Runners() []app.Runner {
	eventHandler := cluster.NewEventHandler().
		RequireModule("router", a.router).
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
