package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/logger"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/storage"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	bootstrapLogger := slog.Default()

	cfg, err := config.Load(bootstrapLogger)
	if err != nil {
		bootstrapLogger.Error("unable to load config", slog.Any("error", err))
		os.Exit(1)
	}

	logger, err := logger.New(cfg.Log)
	if err != nil {
		bootstrapLogger.Error("unable to create logger", slog.Any("error", err))
		os.Exit(1)
	}

	storage, err := storage.New(logger, cfg.Storage)
	if err != nil {
		logger.Error("unable to create storage", slog.Any("error", err))
		os.Exit(1)
	}

	if err := storage.Open(ctx); err != nil {
		logger.Error("unable to open storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logger.Error("unable to close storage", slog.Any("error", err))
		}
	}()

	botService := service.NewBotService(
		logger,
		storage.UnitOfWork(),
	)

	app, err := bot.NewApp(logger, cfg.Bot, botService)
	if err != nil {
		logger.Error("unable to create bot application", slog.Any("error", err))
		os.Exit(1)
	}
	if err := app.Run(ctx); err != nil {
		logger.Error("application failed", slog.Any("error", err))
		os.Exit(1)
	}
}
