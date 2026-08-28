package task

import (
	"context"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type NotificationSenderTask struct {
	commonTask
	maxAttempts uint32
}

func NewNotificationSenderTask(logger *slog.Logger, botAPI *telegram.Bot, service service.BotService, maxAttempts uint32) *NotificationSenderTask {
	return &NotificationSenderTask{
		commonTask: commonTask{
			logger: logger.With(
				slog.String("component", "NotificationSenderTask"),
			),
			botAPI:  botAPI,
			service: service,
		},
		maxAttempts: maxAttempts,
	}
}

func (t *NotificationSenderTask) Run(ctx context.Context) error {
	logger := t.logger.With(
		slog.String("operation", "Run"),
	)

	logger.InfoContext(ctx, "operation started")

	notifications, err := t.service.GetPendingNotifications(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	for _, notification := range notifications {
		if err := t.processNotification(ctx, notification); err != nil {
			logger.ErrorContext(ctx, "unable to process notification",
				slog.String("key", notification.Key.String()),
				slog.Uint64("userId", uint64(notification.UserID)),
				slog.Any("error", err),
			)
		}
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("notificationsCount", len(notifications)))
	return nil
}

func (t *NotificationSenderTask) processNotification(ctx context.Context, notification domain.Notification) error {
	status := domain.NotificationStatusSuccess

	_, err := t.botAPI.SendMessage(ctx, &telegram.SendMessageParams{
		ChatID:    int64(notification.UserID),
		Text:      notification.Text,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		status = domain.NotificationStatusPending
		if notification.Attempts+1 >= t.maxAttempts {
			status = domain.NotificationStatusFailed
		}
	}

	return t.service.UpdateNotificationStatus(ctx, service.UpdateNotificationStatusParams{
		Key:    notification.Key,
		UserID: notification.UserID,
		Status: status,
	})
}
