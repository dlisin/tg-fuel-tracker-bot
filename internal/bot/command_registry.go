package bot

import (
	"context"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/command"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CommandRegistry struct {
	logger   *slog.Logger
	botAPI   *telegram.Bot
	handlers map[string]command.Handler
}

func NewCommandRegistry(logger *slog.Logger, cfg config.BotConfig, botAPI *telegram.Bot, service service.BotService) *CommandRegistry {
	return &CommandRegistry{
		logger: logger.With(
			slog.String("component", "CommandRegistry"),
		),
		botAPI: botAPI,
		handlers: map[string]command.Handler{
			"start": command.NewStartCommand(cfg, botAPI, service),
			// "car-add":       command.NewCarAddCommand(cfg, botAPI, service),
			"refuel_add":    command.NewRefuelAddCommand(cfg, botAPI, service),
			"refuel_delete": command.NewRefuelDeleteCommand(cfg, botAPI, service),
			"refuel_list":   command.NewRefuelListCommand(cfg, botAPI, service),
			"refuel_stats":  command.NewRefuelStatsCommand(cfg, botAPI, service),
		},
	}
}

func (r *CommandRegistry) Register() {
	for commandName, handler := range r.handlers {
		r.botAPI.RegisterHandler(
			telegram.HandlerTypeMessageText,
			commandName,
			telegram.MatchTypeCommandStartOnly,
			r.commandHandler(commandName, handler),
		)
	}
}

func (r *CommandRegistry) commandHandler(commandName string, handler command.Handler) telegram.HandlerFunc {
	return func(ctx context.Context, _ *telegram.Bot, update *models.Update) {
		msg := update.Message
		if msg == nil || msg.From == nil {
			return
		}

		logger := r.logger.With(
			slog.String("operation", "ProcessCommand"),
			slog.String("command", commandName),
			slog.Int64("chatId", msg.Chat.ID),
			slog.Int64("userId", msg.From.ID),
		)

		logger.InfoContext(ctx, "operation started")

		if err := handler.Process(ctx, msg); err != nil {
			logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
			return
		}

		logger.InfoContext(ctx, "operation completed")
	}
}
