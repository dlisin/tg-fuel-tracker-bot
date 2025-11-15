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

func NewStatsCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &statsCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h statsCommand) Process(ctx context.Context, msg *telegram.Message) error {
	err := repository.WithTransaction(ctx, h.uow, func(ctx context.Context, tx repository.Transaction) error {
		userID := model.TelegramID(msg.From.ID)

		cmdArgs, err := parseListCommandArgs(msg.CommandArguments())
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
			return nil
		}

		refuels, err := tx.RefuelRepository().List(ctx, userID, repository.RefuelFilter{CreatedAt: cmdArgs.Period})
		if err != nil {
			return err
		}

		if len(refuels) < 2 {
			_ = h.sendMessage(msg.Chat.ID, "ℹ️ Недостаточно данных. Нужны минимум две записи в выбранном периоде")
			return nil
		}

		stats := model.CreateRefuelStats(refuels)
		_ = h.sendMessageFromTemplate(msg.Chat.ID, "templates/stats.tmpl", struct {
			Params *listCommandArgs
			Stats  *model.RefuelStats
			Config *config.Config
		}{
			Params: cmdArgs,
			Stats:  stats,
			Config: h.cfg,
		})

		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить данные. Попробуйте позже")
		return err
	}

	return nil
}
