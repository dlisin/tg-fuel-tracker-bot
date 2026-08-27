package command

import (
	"context"

	"github.com/CloudyKit/jet/v6"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type RefuelAddCommand struct {
	commonCommand
}

func NewRefuelAddCommand(botAPI *telegram.Bot, service service.BotService) *RefuelAddCommand {
	return &RefuelAddCommand{
		commonCommand: commonCommand{
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *RefuelAddCommand) Process(ctx context.Context, msg *models.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseRefuelAddCommandArgs(parseCommandArgs(msg.Text))
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, err.Error())
	}

	refuel, err := h.service.AddRefuel(ctx, userID, service.AddRefuelParams{
		RegNumber:  car.RegNumber,
		Odometer:   cmdArgs.Odometer,
		Liters:     cmdArgs.Liters,
		PriceTotal: cmdArgs.TotalPrice,
	})
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, h.handleServiceError(err).Error())
	}

	stats, _ := h.service.GetLatestRefuelStats(ctx, userID, car.RegNumber)

	variables := jet.VarMap{}
	variables.Set("Car", car)
	variables.Set("Refuel", refuel)
	variables.Set("Stats", stats)

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "command/refuel_add.jet", variables)
}
