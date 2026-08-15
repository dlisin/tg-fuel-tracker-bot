package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/storage/sqlite"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"
)

type Storage interface {
	Open(ctx context.Context) error
	Close() error

	UnitOfWork() repository.UnitOfWork
}

func New(logger *slog.Logger, cfg config.StorageConfig) (Storage, error) {
	switch cfg.Provider {
	case config.StorageProviderSQLite:
		return sqlite.New(logger, cfg.SQLite), nil

	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Provider)
	}
}
