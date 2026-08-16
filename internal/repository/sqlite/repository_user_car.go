package sqlite

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"

	"github.com/jmoiron/sqlx"
)

type SQLiteUserCarRepository struct {
	logger *slog.Logger
	db     sqlx.ExtContext
}

func NewUserCarRepository(logger *slog.Logger, db sqlx.ExtContext) *SQLiteUserCarRepository {
	return &SQLiteUserCarRepository{
		logger: logger.With(
			slog.String("component", "SQLiteUserCarRepository"),
		),
		db: db,
	}
}

func (r *SQLiteUserCarRepository) Get(ctx context.Context, userID domain.TelegramID, carID domain.CarID) (*domain.UserCar, error) {
	logger := r.logger.With(
		slog.String("operation", "Get"),
		slog.Uint64("userId", uint64(userID)),
		slog.Uint64("carId", uint64(carID)),
	)

	const query = `SELECT id, user_id, car_id, is_owner, created_at FROM user_cars WHERE user_id = ? AND car_id = ?`
	queryArgs := []any{userID, carID}
	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var userCar domain.UserCar
	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).StructScan(&userCar); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityNotFound) {
			logger.DebugContext(ctx, "entity not found")
			return nil, err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entity found", slog.Uint64("userCarId", userCar.ID))
	return &userCar, nil
}

func (r *SQLiteUserCarRepository) List(ctx context.Context, userID domain.TelegramID) ([]domain.Car, error) {
	logger := r.logger.With(
		slog.String("operation", "List"),
		slog.Uint64("userId", uint64(userID)),
	)

	const query = `
		SELECT
			c.id,
			c.reg_number,
			c.fuel_type,
			c.odometer,
			c.created_at,
			c.updated_at
		FROM user_cars AS uc
		INNER JOIN cars AS c ON c.id = uc.car_id
		WHERE uc.user_id = ?
		ORDER BY uc.created_at
	`
	queryArgs := []any{userID}
	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var cars []domain.Car
	if err := sqlx.SelectContext(ctx, r.db, &cars, query, queryArgs...); err != nil {
		err = translateError(err)
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entities found", slog.Int("count", len(cars)))
	return cars, nil
}

func (r *SQLiteUserCarRepository) Create(ctx context.Context, userCar *domain.UserCar) error {
	logger := r.logger.With(
		slog.String("operation", "Create"),
		slog.Uint64("userId", uint64(userCar.UserID)),
		slog.Uint64("carId", uint64(userCar.CarID)),
		slog.Bool("isOwner", userCar.IsOwner),
	)

	const query = `INSERT INTO user_cars (user_id, car_id, is_owner, created_at) VALUES (?, ?, ?, ?) RETURNING id`
	queryArgs := []any{userCar.UserID, userCar.CarID, userCar.IsOwner, userCar.CreatedAt}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).Scan(&userCar.ID); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityAlreadyExists) {
			logger.DebugContext(ctx, "entity already exists")
			return err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	logger.DebugContext(ctx, "entity created", slog.Uint64("userCarId", userCar.ID))
	return nil
}

func (r *SQLiteUserCarRepository) Delete(ctx context.Context, userCar *domain.UserCar) error {
	logger := r.logger.With(
		slog.String("operation", "Delete"),
		slog.Uint64("userCarId", userCar.ID),
		slog.Uint64("userId", uint64(userCar.UserID)),
		slog.Uint64("carId", uint64(userCar.CarID)),
	)

	const query = `DELETE FROM user_cars WHERE id = ?`
	queryArgs := []any{userCar.ID}

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
