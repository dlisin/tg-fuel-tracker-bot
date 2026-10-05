package task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	telegram "github.com/go-telegram/bot"
)

type MonthlyStatsTask struct {
	commonTask
	yearlyStatsEnabled bool
}

func NewMonthlyStatsTask(logger *slog.Logger, botAPI *telegram.Bot, service service.BotService, yearlyStatsEnabled bool) *MonthlyStatsTask {
	return &MonthlyStatsTask{
		commonTask: commonTask{
			logger: logger.With(
				slog.String("component", "MonthlyStatsTask"),
			),
			botAPI:  botAPI,
			service: service,
		},
		yearlyStatsEnabled: yearlyStatsEnabled,
	}
}

func (t *MonthlyStatsTask) Run(ctx context.Context) error {
	logger := t.logger.With(
		slog.String("operation", "Run"),
	)

	logger.InfoContext(ctx, "operation started")

	now := time.Now()
	if shouldSkipMonthlyStats(now, t.yearlyStatsEnabled) {
		logger.InfoContext(ctx, "operation completed", slog.String("reason", "yearly stats replace December monthly stats"))
		return nil
	}

	from, to := previousMonthPeriod(now)
	cars, err := t.service.GetAllCars(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	var taskErrors []error

	for _, car := range cars {
		if err := t.processCar(ctx, car, from, to); err != nil {
			logger.ErrorContext(ctx, "unable to process car",
				slog.Uint64("carId", uint64(car.ID)),
				slog.String("regNumber", car.RegNumber.String()),
				slog.Any("error", err),
			)

			taskErrors = append(taskErrors, fmt.Errorf("unable to process car %d: %w", car.ID, err))
		}
	}

	if err := errors.Join(taskErrors...); err != nil {
		logger.ErrorContext(ctx, "operation failed",
			slog.Int("carsCount", len(cars)),
			slog.Int("errorsCount", len(taskErrors)),
			slog.Any("error", err),
		)
		return err
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("carsCount", len(cars)))
	return nil
}

func (t *MonthlyStatsTask) processCar(ctx context.Context, car domain.Car, from time.Time, to time.Time) error {
	stats, err := t.service.GetRefuelStatsForPeriod(ctx, car.CreatedBy, service.GetRefuelsForPeriodParams{
		RegNumber: car.RegNumber,
		From:      from,
		To:        to,
	})
	if err != nil {
		if errors.Is(err, service.ErrStatsNotEnoughRefuels) {
			t.logger.DebugContext(ctx, "car has no refuels for period", slog.Uint64("carId", uint64(car.ID)))
			return nil
		}

		return err
	}

	text, err := t.renderTemplate("task/monthly_stats.jet", jet.VarMap{}.
		Set("Label", getLabel(from)).
		Set("Car", car).
		Set("Stats", stats),
	)
	if err != nil {
		return err
	}

	return t.service.AddNotification(ctx, service.AddNotificationParams{
		Key:    domain.NewNotificationKey("monthly-stats", fmt.Sprint(car.ID), from.Format("2006-01")),
		UserID: car.CreatedBy,
		Text:   text,
	})
}

func previousMonthPeriod(now time.Time) (time.Time, time.Time) {
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return currentMonth.AddDate(0, -1, 0), currentMonth.Add(-time.Nanosecond)
}

func shouldSkipMonthlyStats(now time.Time, yearlyStatsEnabled bool) bool {
	return yearlyStatsEnabled && now.Month() == time.January
}

func getLabel(date time.Time) string {
	return fmt.Sprintf("за %s %d", getMonthName(date.Month()), date.Year())
}

func getMonthName(month time.Month) string {
	months := [...]string{
		"январь",
		"февраль",
		"март",
		"апрель",
		"май",
		"июнь",
		"июль",
		"август",
		"сентябрь",
		"октябрь",
		"ноябрь",
		"декабрь",
	}

	return months[month-1]
}
