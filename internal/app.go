package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/logger"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/scheduler"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/storage"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg       config.Config
	logger    *slog.Logger
	storage   storage.Storage
	scheduler scheduler.Scheduler
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

	appScheduler, err := scheduler.New(appLogger, cfg.Scheduler)
	if err != nil {
		return nil, fmt.Errorf("unable to create scheduler: %w", err)
	}

	return &App{
		cfg:       *cfg,
		logger:    appLogger,
		storage:   appStorage,
		scheduler: appScheduler,
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

	service := service.NewBotService(a.logger, a.storage.UnitOfWork())
	telegramBot := bot.New(a.logger, a.cfg.Bot, service, a.scheduler)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := telegramBot.Run(ctx); err != nil {
			return fmt.Errorf("unable to run bot: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		if err := a.scheduler.Run(ctx); err != nil {
			return fmt.Errorf("unable to run scheduler: %w", err)
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}
