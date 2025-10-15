package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Weeping-Willow/tet/internal/config"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
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
	logger := utils.LoggerFromContext(a.globalCtx)
	logger.Info("Starting API server")

	mux := a.newHandler()

	server := &http.Server{
		Addr:    ":" + a.cfg.App.PortHTTP,
		Handler: mux,
	}

	done := make(chan error, 1)
	go func() {
		done <- server.ListenAndServe()
	}()

	select {
	case <-a.globalCtx.Done():
		logger.Info("Shutting down API server")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		if err := server.Shutdown(ctxShutdown); err != nil {
			return errors.Wrap(err, "shutdown http server")
		}

		return nil
	case err := <-done:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "start server")
		}

		return nil
	}
}

func (a *API) newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.Handle("/api/v1/rates/latest", a.latestExchangeRatesHandler())
	mux.Handle("/api/v1/rates/history", a.currencyExchangeRateHistoryHandler())

	return a.loggingMiddleware(mux)
}
