package utils

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func ContextWithLogging(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}
