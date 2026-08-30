package sqlite

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"

	"github.com/jmoiron/sqlx"
)

type SQLiteNotificationRepository struct {
	logger *slog.Logger
	db     sqlx.ExtContext
}

func NewNotificationRepository(logger *slog.Logger, db sqlx.ExtContext) *SQLiteNotificationRepository {
	return &SQLiteNotificationRepository{
		logger: logger.With(
			slog.String("component", "SQLiteNotificationRepository"),
		),
		db: db,
	}
}

func (r *SQLiteNotificationRepository) Get(ctx context.Context, key domain.NotificationKey, userID domain.TelegramID) (*domain.Notification, error) {
	logger := r.logger.With(
		slog.String("operation", "Get"),
		slog.String("key", key.String()),
		slog.Uint64("userId", uint64(userID)),
	)

	const query = `SELECT id, key, user_id, text, status, attempts, created_at, updated_at FROM notifications WHERE key = ? AND user_id = ?`
	queryArgs := []any{key, userID}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var notification domain.Notification
	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).StructScan(&notification); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityNotFound) {
			logger.DebugContext(ctx, "entity not found")
			return nil, err
		}
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entity found", slog.Uint64("notificationId", uint64(notification.ID)))
	return &notification, nil
}

func (r *SQLiteNotificationRepository) List(ctx context.Context, params repository.NotificationListParams) ([]domain.Notification, error) {
	logger := r.logger.With(
		slog.String("operation", "List"),
		slog.Any("status", params.Status),
		slog.Uint64("limit", params.Limit),
		slog.String("order", string(params.Order)),
	)

	query := `SELECT id, key, user_id, text, status, attempts, created_at, updated_at FROM notifications WHERE status = ?`
	queryArgs := []any{params.Status}

	if params.Order == repository.SortOrderDesc {
		query += ` ORDER BY created_at DESC`
	} else {
		query += ` ORDER BY created_at ASC`
	}

	if params.Limit > 0 {
		query += ` LIMIT ?`
		queryArgs = append(queryArgs, params.Limit)
	}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	var notifications []domain.Notification
	if err := sqlx.SelectContext(ctx, r.db, &notifications, query, queryArgs...); err != nil {
		err = translateError(err)
		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return nil, err
	}

	logger.DebugContext(ctx, "entities found", slog.Int("count", len(notifications)))
	return notifications, nil
}

func (r *SQLiteNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	logger := r.logger.With(
		slog.String("operation", "Create"),
		slog.String("key", notification.Key.String()),
		slog.Uint64("userId", uint64(notification.UserID)),
	)

	const query = `INSERT INTO notifications (key, user_id, text, status, attempts, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	queryArgs := []any{
		notification.Key,
		notification.UserID,
		notification.Text,
		notification.Status,
		notification.Attempts,
		notification.CreatedAt,
		notification.UpdatedAt,
	}

	logger.DebugContext(ctx, "executing query", slog.String("query", query), slog.Any("queryArgs", queryArgs))

	if err := r.db.QueryRowxContext(ctx, query, queryArgs...).Scan(&notification.ID); err != nil {
		err = translateError(err)
		if errors.Is(err, repository.ErrEntityAlreadyExists) {
			logger.DebugContext(ctx, "entity already exists")
			return err
		}

		logger.ErrorContext(ctx, "query failed", slog.Any("error", err))
		return err
	}

	logger.DebugContext(ctx, "entity created", slog.Uint64("notificationId", uint64(notification.ID)))
	return nil
}

func (r *SQLiteNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	logger := r.logger.With(
		slog.String("operation", "Update"),
		slog.Uint64("notificationId", uint64(notification.ID)),
		slog.String("notificationKey", notification.Key.String()),
		slog.Uint64("userId", uint64(notification.UserID)),
	)

	const query = `UPDATE notifications SET status = ?, attempts = ?, updated_at = ? WHERE id = ?`
	queryArgs := []any{
		notification.Status,
		notification.Attempts,
		notification.UpdatedAt,
		notification.ID,
	}

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
