package sqlite

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"

	"github.com/jmoiron/sqlx"
)

type SQLiteRefuelRepository struct {
	logger *slog.Logger
	db     sqlx.ExtContext
}

func NewRefuelRepository(logger *slog.Logger, db sqlx.ExtContext) *SQLiteRefuelRepository {
	return &SQLiteRefuelRepository{
		logger: logger.With(
			slog.String("component", "SQLiteRefuelRepository"),
		),
		db: db,
	}
}

func (r *SQLiteRefuelRepository) Get(ctx context.Context, carID domain.CarID, odometer domain.Mileage) (*domain.Refuel, error) {
	logger := r.logger.With(
		slog.String("operation", "Get"),
		slog.Uint64("carId", uint64(carID)),
		slog.Uint64("odometer", uint64(odometer)),
	)

	const query = `SELECT id, car_id, odometer, liters, price_per_liter, price_total, created_at, created_by FROM refuels WHERE car_id = ? AND odometer = ?`
	queryArgs := []any{carID, odometer}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var refuel domain.Refuel
	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).StructScan(&refuel); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityNotFound) {
			logger.DebugContext(ctx, "entity not found")
			return nil, err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entity found", slog.Int64("refuelId", refuel.ID))
	return &refuel, nil
}

func (r *SQLiteRefuelRepository) List(ctx context.Context, carID domain.CarID, params repository.RefuelListParams) ([]domain.Refuel, error) {
	logger := r.logger.With(
		slog.String("operation", "List"),
		slog.Uint64("carId", uint64(carID)),
		slog.Time("from", params.From),
		slog.Time("to", params.To),
		slog.Int("limit", params.Limit),
	)

	query := `SELECT id, car_id, odometer, liters, price_per_liter, price_total, created_at, created_by FROM refuels WHERE car_id = ?`
	queryArgs := []any{carID}

	if !params.From.IsZero() {
		query += ` AND created_at >= ?`
		queryArgs = append(queryArgs, params.From)
	}

	if !params.To.IsZero() {
		query += ` AND created_at <= ?`
		queryArgs = append(queryArgs, params.To)
	}

	if params.Order == repository.SortOrderDesc {
		query += ` ORDER BY odometer DESC`
	} else {
		query += ` ORDER BY odometer ASC`
	}

	if params.Limit > 0 {
		query += ` LIMIT ?`
		queryArgs = append(queryArgs, params.Limit)
	}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var refuels []domain.Refuel
	if err := sqlx.SelectContext(ctx, r.db, &refuels, query, queryArgs...); err != nil {
		err = translateError(err)
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entities found", slog.Int("count", len(refuels)))
	return refuels, nil
}

func (r *SQLiteRefuelRepository) Create(ctx context.Context, refuel *domain.Refuel) error {
	logger := r.logger.With(
		slog.String("operation", "Create"),
		slog.Uint64("carId", uint64(refuel.CarID)),
		slog.Uint64("odometer", uint64(refuel.Odometer)),
	)

	const query = `INSERT INTO refuels (car_id, odometer, liters, price_per_liter, price_total, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	queryArgs := []any{refuel.CarID, refuel.Odometer, refuel.Liters, refuel.PricePerLiter, refuel.PriceTotal, refuel.CreatedAt, refuel.CreatedBy}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).Scan(&refuel.ID); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityAlreadyExists) {
			logger.DebugContext(ctx, "entity already exists")
			return err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	logger.DebugContext(ctx, "entity created", slog.Int64("refuelId", refuel.ID))
	return nil
}

func (r *SQLiteRefuelRepository) Delete(ctx context.Context, refuel *domain.Refuel) error {
	logger := r.logger.With(
		slog.String("operation", "Delete"),
		slog.Int64("refuelId", refuel.ID),
		slog.Uint64("carId", uint64(refuel.CarID)),
		slog.Uint64("odometer", uint64(refuel.Odometer)),
	)

	const query = `DELETE FROM refuels WHERE id = ?`
	queryArgs := []any{refuel.ID}

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
