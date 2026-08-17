package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/logger"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/storage"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type App struct {
	cfg     config.Config
	logger  *slog.Logger
	storage storage.Storage
}

func NewApp() (*App, error) {
	bootstrapLogger := slog.Default()

	cfg, err := config.Load(bootstrapLogger)
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	appLogger, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("unable to create logger: %w", err)
	}

	appStorage, err := storage.New(appLogger, cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("unable to create storage: %w", err)
	}

	return &App{
		cfg:     *cfg,
		logger:  appLogger,
		storage: appStorage,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.storage.Open(ctx); err != nil {
		return fmt.Errorf("unable to open storage: %w", err)
	}

	defer func() {
		if err := a.storage.Close(); err != nil {
			a.logger.ErrorContext(context.WithoutCancel(ctx), "unable to unable to close storage", slog.Any("error", err))
		}
	}()

	botService := service.NewBotService(a.logger, a.storage.UnitOfWork())

	telegramBot, err := bot.New(a.logger, a.cfg.Bot, botService)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}

	if err := telegramBot.Run(ctx); err != nil {
		return fmt.Errorf("unable to run bot: %w", err)
	}

	return nil
}
