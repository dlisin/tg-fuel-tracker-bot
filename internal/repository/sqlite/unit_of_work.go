package sqlite

import (
	"context"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"
	"github.com/jmoiron/sqlx"
)

type SQLiteUnitOfWork struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewUnitOfWork(logger *slog.Logger, db *sqlx.DB) *SQLiteUnitOfWork {
	return &SQLiteUnitOfWork{
		logger: logger,
		db:     db,
	}
}

func (u *SQLiteUnitOfWork) Begin(ctx context.Context) (repository.Transaction, error) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &SQLiteTransaction{
		logger: u.logger,
		tx:     tx,
	}, nil
}

type SQLiteTransaction struct {
	logger *slog.Logger
	tx     *sqlx.Tx
}

func (t *SQLiteTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLiteTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLiteTransaction) CarRepository() repository.CarRepository {
	return NewCarRepository(t.logger, t.tx)
}

func (t *SQLiteTransaction) UserCarRepository() repository.UserCarRepository {
	return NewUserCarRepository(t.logger, t.tx)
}

func (t *SQLiteTransaction) RefuelRepository() repository.RefuelRepository {
	return NewRefuelRepository(t.logger, t.tx)
}

func (t *SQLiteTransaction) UserCarInviteRepository() repository.UserCarInviteRepository {
	return NewSQLiteUserCarInviteRepository(t.logger, t.tx)
}
