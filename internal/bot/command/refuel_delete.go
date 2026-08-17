package command

import (
	"context"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type refuelDeleteCommand struct {
	commonCommand
}

func NewRefuelDeleteCommand(cfg config.BotConfig, botAPI *telegram.Bot, service service.BotService) Handler {
	return &refuelDeleteCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *refuelDeleteCommand) Process(ctx context.Context, msg *models.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseRefuelDeleteCommandArgs(
		parseCommandArgs(msg.Text),
	)
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, err.Error())
	}

	refuel, err := h.service.DeleteRefuel(ctx, userID, service.DeleteRefuelParams{
		RegNumber: car.RegNumber,
		Odometer:  cmdArgs.Odometer,
	})
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, h.handleServiceError(err).Error())
	}

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "templates/refuel_delete.tmpl",
		struct {
			Car    *domain.Car
			Refuel *domain.Refuel
			Config config.BotConfig
		}{
			Car:    car,
			Refuel: refuel,
			Config: h.cfg,
		},
	)
}
