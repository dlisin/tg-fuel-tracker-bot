package repository

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

type RefuelRepository interface {
	Create(ctx context.Context, refuel *model.Refuel) error
}
