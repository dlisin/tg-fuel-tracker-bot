package sqlite

import (
	"database/sql"
	"errors"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"

	"github.com/mattn/go-sqlite3"
)

func translateError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrEntityNotFound
	}

	var sqliteErr sqlite3.Error
	if !errors.As(err, &sqliteErr) {
		return err
	}

	switch {
	case errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey):
		return repository.ErrEntityAlreadyExists
	case errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique):
		return repository.ErrEntityAlreadyExists
	default:
		return err
	}
}
