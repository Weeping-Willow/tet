package api

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/config"
	"github.com/Weeping-Willow/tet/internal/utils"
)

type API struct {
	globalCtx context.Context
	cfg       config.Config
}

func New(ctx context.Context, cfg config.Config) *API {
	return &API{
		globalCtx: ctx,
		cfg:       cfg,
	}
}

func (a *API) Start() error {
	utils.LoggerFromContext(a.globalCtx).Info("Starting API server")

	for {
		select {
		case <-a.globalCtx.Done():
			utils.LoggerFromContext(a.globalCtx).Info("Shutting down API server")
			return nil
		}
	}
}
