package sqlite

import (
	"embed"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

func NewRefuelRepository(db *sqlx.DB) repository.RefuelRepository {
	return &refuelRepository{
		db: db,
	}
}

func runMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect(driverName); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
