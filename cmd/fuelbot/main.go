package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dlisin/tg-fuel-tracker-bot/internal"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	app, err := internal.NewApp()
	if err != nil {
		slog.Error("unable to create application", slog.Any("error", err))
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		slog.Error("application failed", slog.Any("error", err))
		os.Exit(1)
	}
}
