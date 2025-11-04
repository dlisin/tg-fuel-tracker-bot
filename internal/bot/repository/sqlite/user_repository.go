package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type userRepository struct {
	tx *sqlx.Tx
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := r.tx.QueryRowxContext(ctx,
		"INSERT INTO users (telegram_id, fuel_type, currency) VALUES ($1, $2, $3) RETURNING id",
		user.TelegramID, user.FuelType, user.Currency).Scan(&user.ID)
	if err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	err := r.tx.ExecContext(ctx)
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID model.TelegramID) (*model.User, error) {
	var user = &model.User{}
	err := r.tx.QueryRowxContext(ctx,
		"SELECT id, telegram_id, fuel_type, currency, created_at, updated_at FROM users WHERE telegram_id = $1",
		telegramID).Scan(&user.ID, &user.TelegramID, &user.FuelType, &user.Currency, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, r.wrapError(err)
	}

	return user, nil
}

func (r *userRepository) wrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrUserNotFound
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return repository.ErrUserAlreadyExists
		}
	}

	return err
}
