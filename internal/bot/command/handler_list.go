package command

import (
	"context"
	"fmt"
	"time"

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

		if len(refuels) == 0 {
			_ = h.sendMessage(msg.Chat.ID, "ℹ️ Нет заправок в выбранном периоде. Используйте /add чтобы добавить первую")
		}

		text := fmt.Sprintf("📝 *Заправки %s:*\n\n", cmdArgs.Label)
		for _, refuel := range refuels {
			text += fmt.Sprintf("*%d*. %s, пробег %dкм, %.2fл, цена/л: %.2f%s\n\n",
				refuel.ID,
				refuel.CreatedAt.Format(time.DateOnly),
				refuel.Odometer,
				refuel.Liters,
				refuel.PricePerLiter,
				h.cfg.DefaultCurrency)
		}

		_ = h.sendMessage(msg.Chat.ID, text)

		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось загрузить данные. Попробуйте позже")
		return err
	}

	return nil
}
