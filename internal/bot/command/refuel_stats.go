package command

import (
	"context"

	"github.com/CloudyKit/jet/v6"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RefuelStatsCommand struct {
	commonCommand
}

func NewRefuelStatsCommand(botAPI *telegram.Bot, service service.BotService) *RefuelStatsCommand {
	return &RefuelStatsCommand{
		commonCommand: commonCommand{
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h RefuelStatsCommand) Process(ctx context.Context, msg *models.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseListCommandArgs(parseCommandArgs(msg.Text))
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(ctx, msg.Chat.ID, err.Error())
	}

	stats, err := h.service.GetRefuelStatsForPeriod(ctx, userID, service.GetRefuelsForPeriodParams{
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
	variables.Set("Stats", stats)

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "command/refuel_stats.jet", variables)
}
