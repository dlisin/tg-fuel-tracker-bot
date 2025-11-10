package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type userRepository struct {
	tx *sqlx.Tx
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.tx.QueryRowxContext(ctx,
		"INSERT INTO users (telegram_id, created_at) VALUES ($1, $2) RETURNING id",
		user.TelegramID, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return nil, r.wrapError(err)
	}

	log.Printf("User created: %+v\n", user)
	return user, nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID model.TelegramID) (*model.User, error) {
	var user = &model.User{}
	err := r.tx.QueryRowxContext(ctx,
		"SELECT id, telegram_id, created_at FROM users WHERE telegram_id = $1",
		telegramID).StructScan(user)
	if err != nil {
		return nil, r.wrapError(err)
	}

	log.Printf("User found: %+v\n", user)
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
