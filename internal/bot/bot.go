package bot

import (
	"context"
	"fmt"
	"log/slog"

	telegram "github.com/go-telegram/bot"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type Bot struct {
	logger          *slog.Logger
	botAPI          *telegram.Bot
	commandRegistry *CommandRegistry
}

func New(logger *slog.Logger, cfg config.BotConfig, service service.BotService) (*Bot, error) {
	botAPI, err := telegram.New(cfg.Token, telegram.WithDebug())
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}

	return &Bot{
		logger:          logger,
		botAPI:          botAPI,
		commandRegistry: NewCommandRegistry(logger, cfg, botAPI, service),
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.commandRegistry.Register()
	b.botAPI.Start(ctx)

	return nil
}
