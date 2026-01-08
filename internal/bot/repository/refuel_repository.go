package repository

import (
	"context"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

type RefuelRepository interface {
	// Create refuel record
	Create(ctx context.Context, refuel *model.Refuel) (*model.Refuel, error)

	// Delete refuel record
	Delete(ctx context.Context, refuel *model.Refuel) error

	// GetByOdometer return refuel record for specified userID and odometer
	// If odometer is 0 - then return refuel record with max odometer value
	GetByOdometer(ctx context.Context, userID model.TelegramID, odometer int64) (*model.Refuel, error)

	// List return refuel records for specified date range
	// If date range is not specified - then return all refuel records
	List(ctx context.Context, userID model.TelegramID, period model.Range[time.Time]) ([]model.Refuel, error)
}
