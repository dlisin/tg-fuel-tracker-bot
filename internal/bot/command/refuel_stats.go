package command

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type refuelStatsCommand struct {
	commonCommand
}

func NewRefuelStatsCommand(cfg config.BotConfig, botAPI *telegram.BotAPI, service service.BotService) Handler {
	return &refuelStatsCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h refuelStatsCommand) Process(ctx context.Context, msg *telegram.Message) error {
	userID := domain.TelegramID(msg.From.ID)

	cmdArgs, err := parseListCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	car, err := h.resolveCar(ctx, userID, cmdArgs.RegNumber)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, err.Error())
	}

	stats, err := h.service.GetRefuelStatsForPeriod(ctx, userID, service.GetRefuelsForPeriodParams{
		RegNumber: car.RegNumber,
		From:      cmdArgs.From,
		To:        cmdArgs.To,
	})
	if err != nil {
		return h.sendMessage(msg.Chat.ID, h.handleServiceError(err).Error())
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/refuel_stats.tmpl", struct {
		Params *listCommandArgs
		Car    *domain.Car
		Stats  *service.RefuelStats
		Config config.BotConfig
	}{

		Params: cmdArgs,
		Car:    car,
		Stats:  stats,
		Config: h.cfg,
	})
}
