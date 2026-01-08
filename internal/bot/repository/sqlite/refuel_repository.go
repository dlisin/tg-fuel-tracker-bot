package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/jmoiron/sqlx"
)

type refuelRepository struct {
	db *sqlx.DB
}

func (r *refuelRepository) Create(ctx context.Context, refuel *model.Refuel) (*model.Refuel, error) {
	query := "INSERT INTO refuels (user_id, odometer, liters, price_per_liter, price_total) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	queryArgs := []interface{}{refuel.UserID, refuel.Odometer, refuel.Liters, refuel.PricePerLiter, refuel.PriceTotal}

	err := r.db.QueryRowxContext(ctx, query, queryArgs...).Scan(&refuel.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Refuel created: %+v\n", refuel)
	return refuel, nil
}

func (r *refuelRepository) Delete(ctx context.Context, refuel *model.Refuel) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refuels WHERE id = $1", refuel.ID)
	if err != nil {
		return err
	}

	log.Printf("Refuel deleted: %+v\n", refuel)
	return nil
}

func (r *refuelRepository) GetByOdometer(ctx context.Context, userID model.TelegramID, odometer int64) (*model.Refuel, error) {
	var query string
	var queryArgs []interface{}

	if odometer == 0 { // search for refuel with max odometer value
		query = "SELECT id, user_id, odometer, liters, price_per_liter, price_total, created_at FROM refuels WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1"
		queryArgs = []interface{}{userID}
	} else { // search for refuel with specified odometer value
		query = "SELECT id, user_id, odometer, liters, price_per_liter, price_total, created_at FROM refuels WHERE user_id = $1 and odometer = $2"
		queryArgs = []interface{}{userID, odometer}
	}

	var refuel model.Refuel
	err := r.db.GetContext(ctx, &refuel, query, queryArgs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	log.Printf("Refuel found: %+v\n", refuel)
	return &refuel, nil
}

func (r *refuelRepository) List(ctx context.Context, userID model.TelegramID, period model.Range[time.Time]) ([]model.Refuel, error) {
	query := "SELECT id, user_id, odometer, liters, price_per_liter, price_total, created_at FROM refuels WHERE user_id = $1"
	queryArgs := []interface{}{userID}

	if !period.Start.IsZero() {
		queryArgs = append(queryArgs, period.Start)
		query += fmt.Sprintf(" AND created_at >= $%d", len(queryArgs))
	}

	if !period.End.IsZero() {
		queryArgs = append(queryArgs, period.End)
		query += fmt.Sprintf(" AND created_at <= $%d", len(queryArgs))

	}

	query += " ORDER BY created_at DESC"

	var refuels []model.Refuel
	err := r.db.SelectContext(ctx, &refuels, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	log.Printf("Refuels found: %+v\n", refuels)
	return refuels, nil
}
