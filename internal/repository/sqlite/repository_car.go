package sqlite

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"

	"github.com/jmoiron/sqlx"
)

type SQLiteCarRepository struct {
	logger *slog.Logger
	db     sqlx.ExtContext
}

func NewCarRepository(logger *slog.Logger, db sqlx.ExtContext) *SQLiteCarRepository {
	return &SQLiteCarRepository{
		logger: logger.With(
			slog.String("component", "SQLiteCarRepository"),
		),
		db: db,
	}
}

func (r *SQLiteCarRepository) Get(ctx context.Context, regNum domain.RegNumber) (*domain.Car, error) {
	logger := r.logger.With(
		slog.String("operation", "Get"),
		slog.String("regNumber", string(regNum)),
	)

	const query = `SELECT id, reg_number, fuel_type, odometer, created_at, updated_at FROM cars WHERE reg_number = ?`
	queryArgs := []any{regNum}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var car domain.Car
	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).StructScan(&car); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityNotFound) {
			logger.DebugContext(ctx, "entity not found")
			return nil, err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entity found", slog.Uint64("carId", uint64(car.ID)))
	return &car, nil
}

func (r *SQLiteCarRepository) Create(ctx context.Context, car *domain.Car) error {
	logger := r.logger.With(
		slog.String("operation", "Create"),
		slog.String("regNumber", string(car.RegNumber)),
	)

	const query = `INSERT INTO cars (reg_number, fuel_type, odometer, created_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id`
	queryArgs := []any{car.RegNumber, car.FuelType, car.Odometer, car.CreatedAt, car.UpdatedAt}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	err := r.db.QueryRowxContext(ctx, query, queryArgs...).Scan(&car.ID)
	if err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityAlreadyExists) {
			logger.DebugContext(ctx, "entity already exists")
			return err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	logger.DebugContext(
		ctx,
		"entity created",
		slog.Uint64("carId", uint64(car.ID)),
	)

	return nil
}

func (r *SQLiteCarRepository) Update(ctx context.Context, car *domain.Car) error {
	logger := r.logger.With(
		slog.String("operation", "Update"),
		slog.Uint64("carId", uint64(car.ID)),
		slog.String("regNumber", string(car.RegNumber)),
	)

	const query = `UPDATE cars SET fuel_type = ?, odometer = ?, updated_at = ? WHERE id = ?`
	queryArgs := []any{car.FuelType, car.Odometer, car.UpdatedAt, car.ID}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	result, err := r.db.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		err = translateError(err)
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	if rowsAffected == 0 {
		logger.DebugContext(ctx, "entity not found")
		return repository.ErrEntityNotFound
	}

	logger.DebugContext(ctx, "entity updated")
	return nil
}

func (r *SQLiteCarRepository) Delete(ctx context.Context, car *domain.Car) error {
	logger := r.logger.With(
		slog.String("operation", "Delete"),
		slog.Uint64("carId", uint64(car.ID)),
		slog.String("regNumber", string(car.RegNumber)),
	)

	const query = `DELETE FROM cars WHERE id = ?`
	queryArgs := []any{car.ID}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	result, err := r.db.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		err = translateError(err)
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	if rowsAffected == 0 {
		logger.DebugContext(ctx, "entity not found")
		return repository.ErrEntityNotFound
	}

	logger.DebugContext(ctx, "entity deleted")
	return nil
}
