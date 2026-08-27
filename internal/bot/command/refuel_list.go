package command

import (
	"context"

	"github.com/CloudyKit/jet/v6"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RefuelListCommand struct {
	commonCommand
}

func NewRefuelListCommand(botAPI *telegram.Bot, service service.BotService) *RefuelListCommand {
	return &RefuelListCommand{
		commonCommand: commonCommand{
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *RefuelListCommand) Process(ctx context.Context, msg *models.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseListCommandArgs(parseCommandArgs(msg.Text))
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, err.Error())
	}

	refuels, err := h.service.GetRefuelsForPeriod(ctx, userID, service.GetRefuelsForPeriodParams{
		RegNumber: car.RegNumber,
		From:      cmdArgs.From,
		To:        cmdArgs.To,
	})
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, h.handleServiceError(err).Error())
	}

	variables := jet.VarMap{}
	variables.Set("Params", cmdArgs)
	variables.Set("Car", car)
	variables.Set("Refuels", refuels)

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "command/refuel_list.jet", variables)
}
