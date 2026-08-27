package bot

import (
	"context"
	"fmt"
	"log/slog"

	telegram "github.com/go-telegram/bot"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/scheduler"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type Bot struct {
	logger          *slog.Logger
	cfg             config.BotConfig
	commandRegistry *CommandRegistry
	taskRegistry    *TaskRegistry
}

func New(logger *slog.Logger, cfg config.BotConfig, service service.BotService, scheduler scheduler.Scheduler) *Bot {
	return &Bot{
		logger:          logger,
		cfg:             cfg,
		commandRegistry: NewCommandRegistry(logger, service),
		taskRegistry:    NewTaskRegistry(logger, cfg.Tasks, service, scheduler),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	botAPI, err := telegram.New(b.cfg.Token, telegram.WithDebug())
	if err != nil {
		return fmt.Errorf("unable to create bot API: %w", err)
	}

	if err := b.commandRegistry.Register(botAPI); err != nil {
		return fmt.Errorf("unable to register commands: %w", err)
	}

	if err := b.taskRegistry.Register(botAPI); err != nil {
		return fmt.Errorf("unable to register tasks: %w", err)
	}

	botAPI.Start(ctx)

	return nil
}
