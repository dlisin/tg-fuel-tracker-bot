package command

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type refuelListCommand struct {
	commonCommand
}

func NewRefuelListCommand(cfg config.BotConfig, botAPI *telegram.BotAPI, service service.BotService) Handler {
	return &refuelListCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *refuelListCommand) Process(ctx context.Context, msg *telegram.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseListCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, err.Error())
	}

	refuels, err := h.service.GetRefuelsForPeriod(ctx, userID, service.GetRefuelsForPeriodParams{
		RegNumber: car.RegNumber,
		From:      cmdArgs.From,
		To:        cmdArgs.To,
	})
	if err != nil {
		return h.sendMessage(msg.Chat.ID, h.handleServiceError(err).Error())
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/refuel_list.tmpl", struct {
		Params  *listCommandArgs
		Car     *domain.Car
		Refuels []domain.Refuel
		Config  config.BotConfig
	}{
		Params:  cmdArgs,
		Car:     car,
		Refuels: refuels,
		Config:  h.cfg,
	})
}
