package command

import (
	"context"
	"fmt"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/sliceutils"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
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
		_ = h.sendMessage(msg.Chat.ID, fmt.Sprintf("⛽ Заправка добавлена:\n пробег %dкм, %.2fл, цена/л: %.2f%s",
			newRefuel.Odometer, newRefuel.Liters, newRefuel.PricePerLiter, h.cfg.DefaultCurrency))

		if prevRefuel != nil {
			stats := model.CreateRefuelStats([]model.Refuel{*prevRefuel, *newRefuel})
			_ = h.sendMessage(msg.Chat.ID,
				fmt.Sprintf("📊 Статистика с предыдущей заправки:\n• Пробег: %dкм\n• Средний расход: %.2fл/100км\n• Цена/л: %.2f%s → %.2f%s (%+.2f%s; %+.1f%%)",
					stats.TotalDistance, stats.FuelConsumption,
					stats.PricePerLiterFirst, h.cfg.DefaultCurrency,
					stats.PricePerLiterLast, h.cfg.DefaultCurrency,
					stats.PricePerLiterDeltaAbs, h.cfg.DefaultCurrency, stats.PricePerLiterDeltaPct))
		}
	}

	return nil
}
