package command

import (
	"context"
	"fmt"

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
		user, err := tx.UserRepository().GetByTelegramID(ctx, msg.From.ID)
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Не удалось загрузить профиль пользователя. Зарегистрируйтесь, выполнив команду: /start")
			return nil
		}

		cmdArgs, err := parseStatsCommandArgs(msg.CommandArguments())
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
			return nil
		}

		refuels, err := tx.RefuelRepository().List(ctx, user.ID, repository.RefuelFilter{CreatedAt: cmdArgs.Period})
		if err != nil {
			return err
		}

		if len(refuels) < 2 {
			_ = h.sendMessage(msg.Chat.ID, "ℹ️ Недостаточно данных. Нужны минимум две записи в выбранном периоде")
			return nil
		}

		stats := model.CreateRefuelStats(refuels)
		_ = h.sendMessage(msg.Chat.ID,
			fmt.Sprintf("📊 Статистика %s:\n• Пробег: %dкм\n• Средний расход: %.2fл/100км\n• Цена/л: %.2f%s → %.2f%s (%+.2f%s; %+.1f%%)",
				cmdArgs.Label, stats.TotalDistance, stats.FuelConsumption,
				stats.PricePerLiterFirst, h.cfg.DefaultCurrency,
				stats.PricePerLiterLast, h.cfg.DefaultCurrency,
				stats.PricePerLiterDeltaAbs, h.cfg.DefaultCurrency, stats.PricePerLiterDeltaPct))

		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить статистику. Попробуйте позже")
		return err
	}

	return nil
}
