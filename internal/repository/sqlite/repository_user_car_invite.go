package sqlite

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type SQLiteUserCarInviteRepository struct {
	logger *slog.Logger
	db     sqlx.ExtContext
}

func NewUserCarInviteRepository(logger *slog.Logger, db sqlx.ExtContext) *SQLiteUserCarRepository {
	return &SQLiteUserCarRepository{
		logger: logger.With(
			slog.String("component", "SQLiteUserCarInviteRepository"),
		),
		db: db,
	}
}
