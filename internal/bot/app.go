package bot

import (
	"context"
	"fmt"
	"log/slog"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type App struct {
	logger          *slog.Logger
	botAPI          *telegram.BotAPI
	commandRegistry *CommandRegistry
}

func NewApp(logger *slog.Logger, cfg config.BotConfig, service service.BotService) (*App, error) {
	botAPI, err := telegram.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}
	logger.Info("telegram bot authorized", slog.String("username", botAPI.Self.UserName))
	botAPI.Debug = true

	commandRegistry := NewCommandRegistry(logger, cfg, botAPI, service)

	return &App{
		logger:          logger,
		botAPI:          botAPI,
		commandRegistry: commandRegistry,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	updates := a.botAPI.GetUpdatesChan(telegram.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})

	for {
		select {
		case <-ctx.Done():
			return nil

		case update, ok := <-updates:
			if !ok {
				return nil
			}

			a.commandRegistry.ProcessUpdate(ctx, update)
		}
	}
}
