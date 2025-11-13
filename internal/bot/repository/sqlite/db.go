package sqlite

import (
	"context"
	"embed"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

//go:embed migrations/*.sql
var embedMigrations embed.FS

type sqliteUnitOfWork struct {
	db *sqlx.DB
}

type sqliteTransaction struct {
	tx *sqlx.Tx
}

func NewSQLiteDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, cfg.Path)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func NewUnitOfWork(db *sqlx.DB) repository.UnitOfWork {
	return &sqliteUnitOfWork{
		db: db,
	}
}

func (u *sqliteUnitOfWork) Begin(ctx context.Context) (repository.Transaction, error) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &sqliteTransaction{tx: tx}, nil
}

func (t *sqliteTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *sqliteTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *sqliteTransaction) RefuelRepository() repository.RefuelRepository {
	return &refuelRepository{t.tx}
}

func runMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driverName); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
