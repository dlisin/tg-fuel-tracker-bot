package command

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

const helpText = `Добро пожаловать в Топливный бот 🚗

Доступные команды:
/start — помощь
/add <пробег> <литры> <сумма чека> — добавить заправку`

// stats [<месяц>|<год>|*] — показать статистика за период

type Handler interface {
	Process(ctx context.Context, msg *telegram.Message) error
}

type commonCommand struct {
	cfg    *config.Config
	botAPI *telegram.BotAPI
	uow    repository.UnitOfWork
}

func (h *commonCommand) sendHelpMessage(chatID int64) error {
	return h.sendMessage(chatID, helpText)
}

func (h *commonCommand) sendMessage(chatID int64, msgText string) error {
	msg := telegram.NewMessage(chatID, telegram.EscapeText(telegram.ModeMarkdown, msgText))
	msg.ParseMode = telegram.ModeMarkdown

	_, err := h.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
