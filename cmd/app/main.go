package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weeping-Willow/tet/internal/app"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/spf13/cobra"
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

	appInstance := app.MustNew(ctx)

	cmd := cobra.Command{}
	cmd.AddCommand(appInstance.NewServerCommand())
	cmd.AddCommand(appInstance.NewFetchCommand())

	if err := cmd.Execute(); err != nil {
		log.Error("failed to execute command", slog.String("error", err.Error()))
	}
}
