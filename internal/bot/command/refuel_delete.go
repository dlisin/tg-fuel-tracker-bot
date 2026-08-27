package command

import (
	"context"

	"github.com/CloudyKit/jet/v6"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type RefuelDeleteCommand struct {
	commonCommand
}

func NewRefuelDeleteCommand(botAPI *telegram.Bot, service service.BotService) *RefuelDeleteCommand {
	return &RefuelDeleteCommand{
		commonCommand: commonCommand{
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *RefuelDeleteCommand) Process(ctx context.Context, msg *models.Message) error {
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

	variables := jet.VarMap{}
	variables.Set("Car", car)
	variables.Set("Refuel", refuel)

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "command/refuel_delete.jet", variables)
}
