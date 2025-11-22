package command

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type listCommand struct {
	commonCommand
}

func NewListCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &listCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h listCommand) Process(ctx context.Context, msg *telegram.Message) error {
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

		_ = h.sendMessageFromTemplate(msg.Chat.ID, "templates/list.tmpl", struct {
			Params  *listCommandArgs
			Refuels []model.Refuel
			Config  *config.Config
		}{
			Params:  cmdArgs,
			Refuels: refuels,
			Config:  h.cfg,
		})

		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить данные. Попробуйте позже")
		return err
	}

	return nil
}
