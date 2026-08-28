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
}

func NewMonthlyStatsTask(logger *slog.Logger, botAPI *telegram.Bot, service service.BotService) *MonthlyStatsTask {
	return &MonthlyStatsTask{
		commonTask: commonTask{
			logger: logger.With(
				slog.String("component", "MonthlyStatsTask"),
			),
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (t *MonthlyStatsTask) Run(ctx context.Context) error {
	logger := t.logger.With(
		slog.String("operation", "Run"),
	)

	logger.InfoContext(ctx, "operation started")

	from, to := previousMonthPeriod(time.Now())
	cars, err := t.service.GetAllCars(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
		return err
	}

	for _, car := range cars {
		if err := t.processCar(ctx, car, from, to); err != nil {
			logger.ErrorContext(ctx, "unable to process car",
				slog.Uint64("carId", uint64(car.ID)),
				slog.String("regNumber", car.RegNumber.String()),
				slog.Any("error", err),
			)
		}
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

func getLabel(date time.Time) string {
	months := [...]string{
		"Январь",
		"Февраль",
		"Март",
		"Апрель",
		"Май",
		"Июнь",
		"Июль",
		"Август",
		"Сентябрь",
		"Октябрь",
		"Ноябрь",
		"Декабрь",
	}

	return fmt.Sprintf("за %s %d", months[date.Month()-1], date.Year())
}
