package command

import (
	"context"
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/service"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util"
)

type addCommandHandler struct {
	users   *service.UserService
	refuels *service.RefuelService
}

func (h *addCommandHandler) Process(ctx context.Context, msg *telegram.Message) (telegram.Chattable, error) {
	user, err := h.users.GetOrCreateUser(ctx, msg.From.ID)
	if err != nil {
		return createMessage(msg.Chat.ID, "⚠️ Не удалось загрузить профиль. Зарегистрируйтесь, выполнив команду: /start"), err
	}

	odometer, liters, totalPrice, err := h.parseCmdArgs(msg)
	if err != nil {
		return createMessage(msg.Chat.ID, "⚠️ Ошибка: "+err.Error()), nil
	}

	refuel, err := h.refuels.AddRefuel(ctx, user.ID, odometer, liters, totalPrice)
	if err != nil {
		return createMessage(msg.Chat.ID, "❌ Не удалось сохранить заправку. Попробуйте позже"), err
	}

	return createMessage(msg.Chat.ID, fmt.Sprintf("⛽ Заправка добавлена: пробег %dкм, %.2fл, цена/л: %.2f%s",
		refuel.Odometer, refuel.Liters, refuel.PricePerLiter, user.Currency)), nil
}

func (h *addCommandHandler) parseCmdArgs(msg *telegram.Message) (odometer int64, liters float64, totalPrice float64, err error) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 3 {
		return 0, 0, 0, fmt.Errorf("укажите: <пробег> <литры> <сумма чека>")
	}

	odometer, err = util.ParseInt64(args[0])
	if err != nil || odometer < 0 {
		return 0, 0, 0, fmt.Errorf("пробег должен быть целым числом ≥ 0")
	}

	liters, err = util.ParseFloat64(args[1])
	if err != nil || liters <= 0 {
		return 0, 0, 0, fmt.Errorf("литры должны быть числом > 0")
	}

	totalPrice, err = util.ParseFloat64(args[2])
	if err != nil || totalPrice <= 0 {
		return 0, 0, 0, fmt.Errorf("сумма чека должна быть числом > 0")
	}

	return odometer, liters, totalPrice, nil
}
