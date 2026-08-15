package command

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type refuelAddCommand struct {
	commonCommand
}

func NewRefuelAddCommand(cfg config.BotConfig, botAPI *telegram.BotAPI, service service.BotService) Handler {
	return &refuelAddCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *refuelAddCommand) Process(ctx context.Context, msg *telegram.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseRefuelAddCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, err.Error())
	}

	refuel, err := h.service.AddRefuel(ctx, userID, service.AddRefuelParams{
		RegNumber:  car.RegNumber,
		Odometer:   cmdArgs.Odometer,
		Liters:     cmdArgs.Liters,
		PriceTotal: cmdArgs.TotalPrice,
	})
	if err != nil {
		return h.sendMessage(msg.Chat.ID, h.handleServiceError(err).Error())
	}

	stats, _ := h.service.GetLatestRefuelStats(ctx, userID, car.RegNumber)

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/refuel_add.tmpl", struct {
		Car    *domain.Car
		Refuel *domain.Refuel
		Stats  *service.RefuelStats
		Config config.BotConfig
	}{
		Car:    car,
		Refuel: refuel,
		Stats:  stats,
		Config: h.cfg,
	})
}
