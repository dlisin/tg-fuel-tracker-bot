package command

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type deleteCommand struct {
	commonCommand
}

func NewDeleteCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &deleteCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h *deleteCommand) Process(ctx context.Context, msg *telegram.Message) error {
	err := repository.WithTransaction(ctx, h.uow, func(ctx context.Context, tx repository.Transaction) error {
		userID := model.TelegramID(msg.From.ID)

		cmdArgs, err := parseDeleteCommandArgs(msg.CommandArguments())
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
			return nil
		}

		refuel, err := tx.RefuelRepository().GetByOdometer(ctx, userID, cmdArgs.Odometer)
		if err != nil {
			return err
		}

		if refuel != nil {
			err = tx.RefuelRepository().Delete(ctx, refuel)
			if err != nil {
				return err
			}

			_ = h.sendMessageFromTemplate(msg.Chat.ID, "templates/delete.tmpl", struct {
				Refuel *model.Refuel
				Stats  *model.RefuelStats
				Config *config.Config
			}{
				Refuel: refuel,
				Config: h.cfg,
			})

			return nil
		} else {
			_ = h.sendMessage(msg.Chat.ID, fmt.Sprintf("⚠️ Нет заправки с указанным пробегом: %d", cmdArgs.Odometer))
			return nil
		}
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось удалить заправку. Попробуйте позже")
		return err
	}

	return nil
}
