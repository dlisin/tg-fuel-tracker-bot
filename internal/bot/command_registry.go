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

type Command interface {
	Process(ctx context.Context, msg *models.Message) error
}

type CommandRegistry struct {
	logger  *slog.Logger
	cfg     config.BotConfig
	service service.BotService
}

func NewCommandRegistry(logger *slog.Logger, cfg config.BotConfig, service service.BotService) *CommandRegistry {
	return &CommandRegistry{
		logger: logger.With(
			slog.String("component", "CommandRegistry"),
		),
		cfg:     cfg,
		service: service,
	}
}

func (r *CommandRegistry) Register(botAPI *telegram.Bot) error {
	r.registerCommand(botAPI, "start", command.NewStartCommand(r.cfg, botAPI, r.service))
	// r.registerCommand(botAPI, "car_add", command.NewCarAddCommand(r.cfg, botAPI, r.service)
	r.registerCommand(botAPI, "refuel_add", command.NewRefuelAddCommand(r.cfg, botAPI, r.service))
	r.registerCommand(botAPI, "refuel_delete", command.NewRefuelDeleteCommand(r.cfg, botAPI, r.service))
	r.registerCommand(botAPI, "refuel_list", command.NewRefuelListCommand(r.cfg, botAPI, r.service))
	r.registerCommand(botAPI, "refuel_stats", command.NewRefuelStatsCommand(r.cfg, botAPI, r.service))

	return nil
}

func (r *CommandRegistry) registerCommand(botAPI *telegram.Bot, commandName string, handler Command) {
	botAPI.RegisterHandler(
		telegram.HandlerTypeMessageText,
		commandName,
		telegram.MatchTypeCommandStartOnly,
		r.commandHandler(commandName, handler),
	)
}

func (r *CommandRegistry) commandHandler(commandName string, handler Command) telegram.HandlerFunc {
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
