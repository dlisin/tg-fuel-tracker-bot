package sqlite

import (
	"context"
	"fmt"
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	"github.com/jmoiron/sqlx"
)

type refuelRepository struct {
	tx *sqlx.Tx
}

func (r *refuelRepository) Create(ctx context.Context, refuel *model.Refuel) (*model.Refuel, error) {
	err := r.tx.QueryRowxContext(ctx,
		"INSERT INTO refuels (user_id, odometer, liters, price_per_liter, price_total) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		refuel.UserID, refuel.Odometer, refuel.Liters, refuel.PricePerLiter, refuel.PriceTotal).Scan(&refuel.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Refuel created: %+v\n", refuel)
	return refuel, nil
}

func (r *refuelRepository) List(ctx context.Context, userID model.TelegramID, filter repository.RefuelFilter) ([]model.Refuel, error) {
	query := "SELECT id, user_id, odometer, liters, price_per_liter, price_total, created_at FROM refuels WHERE user_id = $1"
	queryArgs := []interface{}{userID}

	if filter.CreatedAt != nil {
		if !filter.CreatedAt.Start.IsZero() {
			queryArgs = append(queryArgs, filter.CreatedAt.Start)
			query += fmt.Sprintf(" AND created_at >= $%d", len(queryArgs))
		}

		if !filter.CreatedAt.End.IsZero() {
			queryArgs = append(queryArgs, filter.CreatedAt.End)
			query += fmt.Sprintf(" AND created_at <= $%d", len(queryArgs))

		}
	}

	query += " ORDER BY created_at ASC"

	if filter.Limit > 0 {
		queryArgs = append(queryArgs, filter.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(queryArgs))
	}

	refuels := []model.Refuel{}
	err := r.tx.SelectContext(ctx, &refuels, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	log.Printf("Refuels found:\n\tFilter: %+v\n\tResult: %+v\n", filter, refuels)
	return refuels, nil
}
