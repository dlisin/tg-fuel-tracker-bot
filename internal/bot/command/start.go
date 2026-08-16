package command

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type startCommand struct {
	commonCommand
}

func NewStartCommand(cfg config.BotConfig, botAPI *telegram.BotAPI, service service.BotService) Handler {
	return &startCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *startCommand) Process(_ context.Context, msg *telegram.Message) error {
	_, err := h.botAPI.Request(telegram.NewSetMyCommandsWithScope(telegram.NewBotCommandScopeChat(msg.Chat.ID),
		telegram.BotCommand{
			Command:     "/start",
			Description: "помощь",
		},
		telegram.BotCommand{
			Command:     "/refuel_add",
			Description: "добавить заправку",
		},
		telegram.BotCommand{
			Command:     "/refuel_delete",
			Description: "удалить заправку",
		},
		telegram.BotCommand{
			Command:     "/refuel_list",
			Description: "показать заправки за указанный период",
		},
		telegram.BotCommand{
			Command:     "/refuel_stats",
			Description: "показать статистику за указанный период",
		},
	))
	if err != nil {
		return err
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/start.tmpl", nil)
}
