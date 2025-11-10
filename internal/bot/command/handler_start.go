package command

import (
	"context"
	"errors"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type startCommand struct {
	commonCommand
}

func NewStartCommand(cfg *config.Config, botAPI *telegram.BotAPI, uow repository.UnitOfWork) Handler {
	return &startCommand{
		commonCommand: commonCommand{
			cfg:    cfg,
			botAPI: botAPI,
			uow:    uow,
		},
	}
}

func (h *startCommand) Process(ctx context.Context, msg *telegram.Message) error {
	err := repository.WithTransaction(ctx, h.uow, func(ctx context.Context, tx repository.Transaction) error {
		_, err := tx.UserRepository().Create(ctx, &model.User{
			TelegramID: msg.From.ID,
			CreatedAt:  time.Now(),
		})
		if err != nil && !errors.Is(err, repository.ErrUserAlreadyExists) {
			return err
		}
		return nil
	})
	if err != nil {
		_ = h.sendMessage(msg.Chat.ID, "❌ Не удалось сохранить профиль пользователя. Попробуйте позже")
		return err
	}

	return h.sendHelpMessage(msg.Chat.ID)
}
