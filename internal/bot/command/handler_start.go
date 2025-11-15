package command

import (
	"context"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type startCommand struct {
	commonCommand
}

func NewStartCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &startCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h *startCommand) Process(_ context.Context, msg *telegram.Message) error {
	_, err := h.botAPI.Send(telegram.NewSetMyCommandsWithScope(telegram.NewBotCommandScopeChat(msg.Chat.ID),
		telegram.BotCommand{
			Command:     "/start",
			Description: "помощь",
		},
		telegram.BotCommand{
			Command:     "/add",
			Description: "добавить заправку",
		},
		telegram.BotCommand{
			Command:     "/list",
			Description: "показать заправки за указанный период",
		},
		telegram.BotCommand{
			Command:     "/stats",
			Description: "показать статистику за указанный период",
		},
	))
	if err != nil {
		log.Println("Failed to update bot menu: ", err)
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/start.tmpl", nil)
}
