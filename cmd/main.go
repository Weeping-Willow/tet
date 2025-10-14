package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weeping-Willow/tet/internal/app"
	"github.com/Weeping-Willow/tet/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Info("shutting down gracefully")
		cancel()
	}()

	ctx = utils.ContextWithLogging(ctx, log)

	app.MustNew(ctx)
}
