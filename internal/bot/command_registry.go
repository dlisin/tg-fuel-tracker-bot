package bot

import (
	"context"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/command"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandRegistry struct {
	logger   *slog.Logger
	handlers map[string]command.Handler
}

func NewCommandRegistry(logger *slog.Logger, cfg config.BotConfig, botAPI *telegram.BotAPI, service service.BotService) *CommandRegistry {
	return &CommandRegistry{
		logger: logger.With(
			slog.String("component", "CommandRegistry"),
		),
		handlers: map[string]command.Handler{
			"start": command.NewStartCommand(cfg, botAPI, service),
			// "car-add":       command.NewCarAddCommand(cfg, botAPI, service),
			"refuel-add":    command.NewRefuelAddCommand(cfg, botAPI, service),
			"refuel-delete": command.NewRefuelDeleteCommand(cfg, botAPI, service),
			"refuel-list":   command.NewRefuelListCommand(cfg, botAPI, service),
			"refuel-stats":  command.NewRefuelStatsCommand(cfg, botAPI, service),
		},
	}
}

func (r *CommandRegistry) ProcessUpdate(ctx context.Context, update telegram.Update) {
	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return
	}

	logger := r.logger.With(
		slog.String("operation", "ProcessUpdate"),
		slog.String("command", msg.Command()),
		slog.Int64("chatId", msg.Chat.ID),
		slog.Int64("userId", msg.From.ID),
	)

	logger.InfoContext(ctx, "operation started")

	handler, ok := r.handlers[msg.Command()]
	if !ok {
		logger.WarnContext(ctx, "operation aborted", slog.String("error", "unsupported command"))
		return
	}

	if err := handler.Process(ctx, msg); err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return
	}

	logger.InfoContext(ctx, "operation completed")
}
