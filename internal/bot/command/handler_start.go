package command

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

const helpText = `Добро пожаловать в Топливный бот 🚗

Доступные команды:
/start — помощь
/add <пробег> <литры> <сумма чека> — добавить заправку
/stats [*] — показать статистика за период`

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
	return h.sendMessage(msg.Chat.ID, helpText)
}
