package sqlite

import (
	"context"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	repository "github.com/dlisin/tg-fuel-tracker-bot/internal/repository"
	sqliteRepository "github.com/dlisin/tg-fuel-tracker-bot/internal/repository/sqlite"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	logger *slog.Logger
	config config.SQLiteStorageConfig

	db  *sqlx.DB
	uow *sqliteRepository.SQLiteUnitOfWork
}

func New(logger *slog.Logger, cfg config.SQLiteStorageConfig) *SQLiteStorage {
	logger = logger.With(
		slog.String("component", "SQLiteStorage"),
	)

	logger.Debug("component initialized")
	return &SQLiteStorage{
		logger: logger,
		config: cfg,
	}
}

func (s *SQLiteStorage) Open(ctx context.Context) error {
	logger := s.logger.With(
		slog.String("operation", "Open"),
		slog.String("path", s.config.Path),
	)

	logger.InfoContext(ctx, "operation started")
	err := func() error {
		db, err := sqlx.Open("sqlite3", s.config.Path)
		if err != nil {
			return err
		}

		db.SetMaxOpenConns(s.config.MaxConnections)
		db.SetMaxIdleConns(s.config.MaxConnections)

		logger.DebugContext(ctx, "checking database connection")
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			return err
		}
		logger.DebugContext(ctx, "database connection is up")

		logger.DebugContext(ctx, "database migrations started")
		if err := migrate(db.DB); err != nil {
			_ = db.Close()
			return err
		}
		logger.DebugContext(ctx, "database migrations completed")

		s.db = db
		s.uow = sqliteRepository.NewUnitOfWork(s.logger, db)

		return nil
	}()
	if err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	logger.InfoContext(ctx, "operation completed")
	return nil
}

func (s *SQLiteStorage) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *SQLiteStorage) UnitOfWork() repository.UnitOfWork {
	return s.uow
}
