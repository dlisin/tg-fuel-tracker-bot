package repository

import (
	"context"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

type RefuelFilter struct {
	CreatedAt *model.Range[time.Time]

	Limit int
}

type RefuelRepository interface {
	Create(ctx context.Context, refuel *model.Refuel) (*model.Refuel, error)

	List(ctx context.Context, userID model.UserID, filter RefuelFilter) ([]model.Refuel, error)
}
