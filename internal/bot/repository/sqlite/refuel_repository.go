package sqlite

import (
	"context"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/jmoiron/sqlx"
)

type refuelRepository struct {
	tx *sqlx.Tx
}

func (r *refuelRepository) Create(ctx context.Context, refuel *model.Refuel) error {
	err := r.tx.QueryRowxContext(ctx,
		"INSERT INTO refuels (user_id, odometer, liters, price_per_liter, price_total) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		refuel.UserID, refuel.Odometer, refuel.Liters, refuel.PricePerLiter, refuel.PriceTotal).Scan(&refuel.ID)
	if err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *refuelRepository) wrapError(err error) error {
	return err
}
