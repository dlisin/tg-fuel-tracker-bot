package repository

import (
	"context"
	"errors"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)

	GetByTelegramID(ctx context.Context, telegramID model.TelegramID) (*model.User, error)
}
