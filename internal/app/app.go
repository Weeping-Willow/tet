package app

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/api"
	"github.com/Weeping-Willow/tet/internal/config"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type App struct {
	globalContext context.Context
	api           *api.API
}

func MustNew(ctx context.Context) *App {
	cfg, err := config.New()
	if err != nil {
		utils.LoggerFromContext(ctx).Error(errors.Wrap(err, "load config").Error(), nil)
		panic(err)
	}

	apiServer := api.New(ctx, cfg)

	return &App{
		api:           apiServer,
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
