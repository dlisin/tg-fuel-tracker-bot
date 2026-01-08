package command

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type deleteCommand struct {
	commonCommand
}

func NewDeleteCommand(cfg *config.Config, botAPI *telegram.BotAPI, refuelRepository repository.RefuelRepository) Handler {
	return &deleteCommand{
		commonCommand: commonCommand{
			cfg:              cfg,
			botAPI:           botAPI,
			refuelRepository: refuelRepository,
		},
	}
}

func (h *deleteCommand) Process(ctx context.Context, msg *telegram.Message) error {
	cmdArgs, err := parseDeleteCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	refuel, err := h.refuelRepository.GetByOdometer(ctx, model.TelegramID(msg.From.ID), cmdArgs.Odometer)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "❌ Не удалось удалить заправку. Попробуйте позже")
	}

	if refuel != nil {
		err = h.refuelRepository.Delete(ctx, refuel)
		if err != nil {
			return h.sendMessage(msg.Chat.ID, "❌ Не удалось удалить заправку. Попробуйте позже")
		}

		return h.sendMessageFromTemplate(msg.Chat.ID, "templates/delete.tmpl", struct {
			Refuel *model.Refuel
			Stats  *model.RefuelStats
			Config *config.Config
		}{
			Refuel: refuel,
			Config: h.cfg,
		})

	} else {
		return h.sendMessage(msg.Chat.ID, fmt.Sprintf("⚠️ Нет заправки с указанным пробегом: %d", cmdArgs.Odometer))
	}

}
