package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/go-co-op/gocron/v2"
)

type SchedulerImpl struct {
	logger    *slog.Logger
	cfg       config.SchedulerConfig
	scheduler gocron.Scheduler
}

func New(logger *slog.Logger, cfg config.SchedulerConfig) (*SchedulerImpl, error) {
	logger = logger.With(
		slog.String("component", "Scheduler"),
	)

	cronScheduler, err := gocron.NewScheduler(
		gocron.WithGlobalJobOptions(
			gocron.WithSingletonMode(
				gocron.LimitModeReschedule,
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create scheduler: %w", err)
	}

	scheduler := &SchedulerImpl{
		logger:    logger,
		cfg:       cfg,
		scheduler: cronScheduler,
	}

	return scheduler, nil
}

func (s *SchedulerImpl) Schedule(taskName string, taskCronExpression string, scheduledTask ScheduledTask) error {
	logger := s.logger.With(
		slog.String("operation", "Schedule"),
		slog.String("task", taskName),
		slog.String("schedule", taskCronExpression),
	)

	logger.Info("operation started")

	runner := newTaskRunner(logger, s.cfg.RetryPolicy, func(ctx context.Context) error {
		return s.runTask(ctx, taskName, scheduledTask)
	})

	job, err := s.scheduler.NewJob(gocron.CronJob(taskCronExpression, false), gocron.NewTask(runner.Run), gocron.WithName(taskName))
	if err != nil {
		err = fmt.Errorf("unable to schedule task %q: %w", taskName, err)
		logger.Error("operation failed", slog.Any("error", err))
		return err
	}

	runner.setJob(job)

	logger.Info("operation completed")
	return nil
}

func (s *SchedulerImpl) Run(ctx context.Context) error {
	logger := s.logger.With(
		slog.String("operation", "Run"),
	)

	logger.InfoContext(ctx, "operation started")

	s.scheduler.Start()
	<-ctx.Done()

	if err := s.scheduler.Shutdown(); err != nil {
		err = fmt.Errorf("unable to shutdown scheduler: %w", err)
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	logger.InfoContext(ctx, "operation completed")
	return nil
}

func (s *SchedulerImpl) runTask(ctx context.Context, name string, scheduledTask ScheduledTask) error {
	logger := s.logger.With(
		slog.String("operation", "RunTask"),
		slog.String("task", name),
	)

	logger.InfoContext(ctx, "operation started")

	if err := scheduledTask.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	logger.InfoContext(ctx, "operation completed")
	return nil
}
