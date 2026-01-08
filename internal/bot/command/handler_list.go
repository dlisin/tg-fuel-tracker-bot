package command

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type listCommand struct {
	commonCommand
}

func NewListCommand(cfg *config.Config, botAPI *telegram.BotAPI, refuelRepository repository.RefuelRepository) Handler {
	return &listCommand{
		commonCommand: commonCommand{
			cfg:              cfg,
			botAPI:           botAPI,
			refuelRepository: refuelRepository,
		},
	}
}

func (h listCommand) Process(ctx context.Context, msg *telegram.Message) error {
	cmdArgs, err := parseListCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	refuels, err := h.refuelRepository.List(ctx, model.TelegramID(msg.From.ID), cmdArgs.Period)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить данные. Попробуйте позже")
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/list.tmpl", struct {
		Params  *listCommandArgs
		Refuels []model.Refuel
		Config  *config.Config
	}{
		Params:  cmdArgs,
		Refuels: refuels,
		Config:  h.cfg,
	})
}
