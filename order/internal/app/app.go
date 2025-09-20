package app

import (
	"context"
	"fmt"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/config"
	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/closer"
	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initDB,
		a.initMigrations,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP Order server listening on %s", a.diContainer.HTTPServer(ctx).Addr))
	return a.diContainer.HTTPServer(ctx).ListenAndServe()
}

func (a *App) initDB(ctx context.Context) error {
	_ = a.diContainer.DBPool(ctx)
	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	return a.diContainer.runMigrations(ctx)
}

func (a *App) initHTTPServer(ctx context.Context) error {
	_ = a.diContainer.HTTPServer(ctx)
	return nil
}
