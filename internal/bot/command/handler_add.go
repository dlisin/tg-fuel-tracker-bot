package command

import (
	"context"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/sliceutils"
)

type addCommand struct {
	commonCommand
}

func NewAddCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &addCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h *addCommand) Process(ctx context.Context, msg *telegram.Message) error {
	var prevRefuel, newRefuel *model.Refuel

	err := repository.WithTransaction(ctx, h.uow, func(ctx context.Context, tx repository.Transaction) error {
		userID := model.TelegramID(msg.From.ID)
		prevRefuels, err := tx.RefuelRepository().List(ctx, userID, repository.RefuelFilter{Limit: 1})
		if err != nil {
			return err
		}
		prevRefuel = sliceutils.First(prevRefuels)

		cmdArgs, err := parseAddCommandArgs(msg.CommandArguments(), prevRefuel)
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
			return nil
		}

		newRefuel, err = tx.RefuelRepository().Create(ctx, &model.Refuel{
			UserID:        userID,
			Odometer:      cmdArgs.Odometer,
			Liters:        cmdArgs.Liters,
			PriceTotal:    cmdArgs.TotalPrice,
			PricePerLiter: cmdArgs.TotalPrice / cmdArgs.Liters,
			CreatedAt:     time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось сохранить заправку. Попробуйте позже")
		return err
	}

	if newRefuel != nil {
		var stats *model.RefuelStats
		if prevRefuel != nil {
			stats = model.CreateRefuelStats([]model.Refuel{*prevRefuel, *newRefuel})
		}

		_ = h.sendMessageFromTemplate(msg.Chat.ID, "templates/add.tmpl", struct {
			Refuel *model.Refuel
			Stats  *model.RefuelStats
			Config *config.Config
		}{
			Refuel: newRefuel,
			Stats:  stats,
			Config: h.cfg,
		})
	}

	return nil
}
