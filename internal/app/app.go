package app

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/api"
	"github.com/Weeping-Willow/tet/internal/config"
	"github.com/Weeping-Willow/tet/internal/rates"
	"github.com/Weeping-Willow/tet/internal/storage"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type App struct {
	api         *api.API
	rateService rates.Service

	globalContext context.Context
}

func MustNew(ctx context.Context) *App {
	cfg, err := config.New()
	if err != nil {
		utils.LoggerFromContext(ctx).Error(errors.Wrap(err, "load config").Error())
		panic(err)
	}

	db, err := config.NewDb(cfg)
	if err != nil {
		utils.LoggerFromContext(ctx).Error(errors.Wrap(err, "connect to database").Error())
		panic(err)
	}

	err = storage.Migrate(db)
	if err != nil {
		utils.LoggerFromContext(ctx).Error(errors.Wrap(err, "migrate database").Error())
		panic(err)
	}

	rateFetcher := rates.NewEcbRssFetcher(cfg)
	rateRepo := storage.NewRates(db)
	rateService := rates.NewService(rateRepo, rateFetcher)

	apiServer := api.New(ctx, rateService, cfg)

	return &App{
		api:         apiServer,
		rateService: rateService,

		globalContext: ctx,
	}
}

func (a *App) NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the API server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := a.api.Start(); err != nil {
				utils.LoggerFromContext(a.globalContext).Error(errors.Wrap(err, "run API server").Error())
				return
			}
		},
	}
}

func (a *App) NewFetchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Fetch and update currency exchange rates",
		Run: func(cmd *cobra.Command, args []string) {
			if err := a.rateService.UpdateRates(a.globalContext); err != nil {
				utils.LoggerFromContext(a.globalContext).Error(errors.Wrap(err, "update rates").Error())
				return
			}
		},
	}
}
