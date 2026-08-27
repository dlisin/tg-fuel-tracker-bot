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
	if r.cfg.MonthlyStats.Enabled {
		task := task.NewMonthlyStatsTask(r.logger, botAPI, r.service)
		if err := r.scheduler.Schedule("monthly-stats", r.cfg.MonthlyStats.Schedule, task); err != nil {
			return fmt.Errorf("unable to schedule monthly-stats task: %w", err)
		}
	}

	return nil
}
