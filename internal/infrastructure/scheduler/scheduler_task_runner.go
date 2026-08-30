package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/go-co-op/gocron/v2"
)

type taskRunner struct {
	logger *slog.Logger
	retry  config.SchedulerTaskRetryConfig
	run    func(context.Context) error

	mu       sync.Mutex
	job      gocron.Job
	attempts uint32
	timer    *time.Timer
	token    uint64
}

func newTaskRunner(logger *slog.Logger, retry config.SchedulerTaskRetryConfig, run func(context.Context) error) *taskRunner {
	return &taskRunner{
		logger: logger,
		retry:  retry,
		run:    run,
	}
}

func (r *taskRunner) setJob(job gocron.Job) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.job = job
}

func (r *taskRunner) Run(ctx context.Context) error {
	r.beginRun()

	err := r.run(ctx)
	if err == nil {
		r.completeRun()
		return nil
	}

	r.scheduleRetry(ctx)
	return err
}

func (r *taskRunner) beginRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.token++

	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
}

func (r *taskRunner) completeRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.attempts = 0
}

func (r *taskRunner) scheduleRetry(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.attempts++
	if r.retry.MaxAttempts <= 1 || r.attempts >= r.retry.MaxAttempts {
		r.attempts = 0
		return
	}

	nextAttempt := r.attempts + 1
	token := r.token

	r.timer = time.AfterFunc(r.retry.Delay, func() {
		r.runRetry(ctx, token, nextAttempt)
	})

	r.logger.WarnContext(ctx, "task retry scheduled",
		slog.Uint64("attempt", uint64(nextAttempt)),
		slog.Duration("delay", r.retry.Delay),
	)
}

func (r *taskRunner) runRetry(ctx context.Context, token uint64, attempt uint32) {
	if ctx.Err() != nil {
		return
	}

	r.mu.Lock()

	if r.token != token || r.timer == nil {
		r.mu.Unlock()
		return
	}

	r.timer = nil
	job := r.job

	r.mu.Unlock()

	if ctx.Err() != nil {
		return
	}

	if job == nil {
		r.logger.Error("unable to reschedule task retry",
			slog.Uint64("attempt", uint64(attempt)),
			slog.String("error", "job is not initialized"),
		)
		r.reset()
		return
	}

	if err := job.RunNow(); err != nil {
		r.logger.Error("unable to reschedule task retry",
			slog.Uint64("attempt", uint64(attempt)),
			slog.Any("error", err),
		)
		r.reset()
	}
}

func (r *taskRunner) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.attempts = 0
	r.token++

	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
}
