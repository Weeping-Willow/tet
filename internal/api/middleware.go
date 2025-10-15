package api

import (
	"net/http"
	"time"

	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/google/uuid"
)

func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := utils.LoggerFromContext(a.globalCtx).With("request_id", uuid.New().String()).With("route", r.URL.Path)
		r = r.WithContext(utils.ContextWithLogging(r.Context(), logger))

		startTime := time.Now()

		logger.Info("incoming request")
		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		logger.Info("request completed", "duration_ms", duration.Milliseconds())
	})
}
