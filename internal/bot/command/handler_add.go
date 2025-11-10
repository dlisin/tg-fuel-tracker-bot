package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/sliceutils"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/stringutils"
)

type addHandler struct {
	commonHandler
}

type addCommandArgs struct {
	Odometer   int64
	Liters     float64
	TotalPrice float64
}

func NewAddHandler(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &addHandler{
		commonHandler: commonHandler{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h *addHandler) Process(ctx context.Context, msg *telegram.Message) error {
	var prevRefuel, newRefuel *model.Refuel

	err := repository.WithTransaction(ctx, h.uow, func(ctx context.Context, tx repository.Transaction) error {
		user, err := tx.UserRepository().GetByTelegramID(ctx, msg.From.ID)
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Не удалось загрузить профиль пользователя. Зарегистрируйтесь, выполнив команду: /start")
			return nil
		}

		prevRefuels, err := tx.RefuelRepository().List(ctx, user.ID, repository.RefuelFilter{Limit: 1})
		if err != nil {
			return err
		}
		prevRefuel = sliceutils.First(prevRefuels)

		cmdArgs, err := h.parseCmdArgs(msg, prevRefuel)
		if err != nil {
			_ = h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
			return nil
		}

		newRefuel, err = tx.RefuelRepository().Create(ctx, &model.Refuel{
			UserID:        user.ID,
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

func (h *addHandler) parseCmdArgs(msg *telegram.Message, prevRefuel *model.Refuel) (*addCommandArgs, error) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 3 {
		return nil, fmt.Errorf("недостаточно параметров, укажите <пробег> <литры> <сумма чека>")
	}

	odometer, err := stringutils.ParseInt64(args[0])
	if err != nil || odometer < 0 {
		return nil, fmt.Errorf("пробег должен быть целым числом ≥ 0")
	}

	liters, err := stringutils.ParseFloat64(args[1])
	if err != nil || liters <= 0 {
		return nil, fmt.Errorf("литры должны быть числом > 0")
	}

	totalPrice, err := stringutils.ParseFloat64(args[2])
	if err != nil || totalPrice <= 0 {
		return nil, fmt.Errorf("сумма чека должна быть числом > 0")
	}

	if prevRefuel != nil {
		prevOdometer := prevRefuel.Odometer
		if prevOdometer >= odometer {
			return nil, fmt.Errorf("пробег должен быть больше предыдущего (%d)", prevOdometer)
		}
	}

	return &addCommandArgs{
		Odometer:   odometer,
		Liters:     liters,
		TotalPrice: totalPrice,
	}, nil
}
