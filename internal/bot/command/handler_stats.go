package command

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type statsCommand struct {
	commonCommand
}

func NewStatsCommand(cfg *config.Config, botAPI *telegram.BotAPI, refuelRepository repository.RefuelRepository) Handler {
	return &statsCommand{
		commonCommand: commonCommand{
			cfg:              cfg,
			botAPI:           botAPI,
			refuelRepository: refuelRepository,
		},
	}
}

func (h statsCommand) Process(ctx context.Context, msg *telegram.Message) error {

	cmdArgs, err := parseListCommandArgs(msg.CommandArguments())
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	refuels, err := h.refuelRepository.List(ctx, model.TelegramID(msg.From.ID), cmdArgs.Period)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить данные. Попробуйте позже")
	}

	if len(refuels) < 2 {
		return h.sendMessage(msg.Chat.ID, "ℹ️ Недостаточно данных. Нужны минимум две записи в выбранном периоде")
	}

	stats := model.CreateRefuelStats(refuels)

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/stats.tmpl", struct {
		Params *listCommandArgs
		Stats  *model.RefuelStats
		Config *config.Config
	}{
		Params: cmdArgs,
		Stats:  stats,
		Config: h.cfg,
	})
}
