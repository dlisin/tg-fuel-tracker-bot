package bot

import (
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/task"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/scheduler"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
)

const yearlyStatsNotificationSenderSchedule = "*/5 23 31 12 *"

type TaskRegistry struct {
	logger    *slog.Logger
	cfg       config.BotTasksConfig
	service   service.BotService
	scheduler scheduler.Scheduler
}

func NewTaskRegistry(logger *slog.Logger, cfg config.BotTasksConfig, service service.BotService, scheduler scheduler.Scheduler) *TaskRegistry {
	return &TaskRegistry{
		logger: logger.With(
			slog.String("component", "TaskRegistry"),
		),
		cfg:       cfg,
		service:   service,
		scheduler: scheduler,
	}
}

func (r *TaskRegistry) Register(botAPI *telegram.Bot) error {
	notificationSenderTask := task.NewNotificationSenderTask(r.logger, botAPI, r.service, r.cfg.NotificationSender.MaxAttempts)
	if err := r.scheduler.Schedule("notification-sender", r.cfg.NotificationSender.Schedule, notificationSenderTask); err != nil {
		return fmt.Errorf("unable to schedule notification-sender task: %w", err)
	}

	if r.cfg.MonthlyStats.Enabled {
		t := task.NewMonthlyStatsTask(r.logger, botAPI, r.service, r.cfg.YearlyStats.Enabled)
		if err := r.scheduler.Schedule("monthly-stats", r.cfg.MonthlyStats.Schedule, t); err != nil {
			return fmt.Errorf("unable to schedule monthly-stats task: %w", err)
		}
	}

	if r.cfg.YearlyStats.Enabled {
		t := task.NewYearlyStatsTask(r.logger, botAPI, r.service)
		if err := r.scheduler.Schedule("yearly-stats", r.cfg.YearlyStats.Schedule, t); err != nil {
			return fmt.Errorf("unable to schedule yearly-stats task: %w", err)
		}

		if err := r.scheduler.Schedule("yearly-stats-notification-sender", yearlyStatsNotificationSenderSchedule, notificationSenderTask); err != nil {
			return fmt.Errorf("unable to schedule yearly-stats-notification-sender task: %w", err)
		}
	}

	return nil
}
