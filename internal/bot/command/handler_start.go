package command

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/service"
)

const helpStartText = `Добро пожаловать в Топливный бот 🚗

Доступные команды:
/start — помощь
/add <пробег> <литры> <сумма_чека> — добавить заправку
/stats [<месяц>|<год>|*] — показать статистика за период`

type startCommandHandler struct {
	users *service.UserService
}

func (h *startCommandHandler) Process(ctx context.Context, msg *telegram.Message) (telegram.Chattable, error) {
	_, err := h.users.GetOrCreateUser(ctx, msg.From.ID)
	if err != nil {
		return createMessage(msg.Chat.ID, "❌ Не удалось сохранить профиль пользователя. Попробуйте позже"), err
	}

	return createMessage(msg.Chat.ID, helpStartText), nil
}
