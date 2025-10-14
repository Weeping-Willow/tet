package api

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/utils"
)

type API struct {
	globalCtx context.Context
}

func New(ctx context.Context) *API {
	return &API{
		globalCtx: ctx,
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
