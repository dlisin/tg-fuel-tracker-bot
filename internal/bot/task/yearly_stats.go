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

type YearlyStatsTask struct {
	commonTask
}

func NewYearlyStatsTask(logger *slog.Logger, botAPI *telegram.Bot, service service.BotService) *YearlyStatsTask {
	return &YearlyStatsTask{
		commonTask: commonTask{
			logger: logger.With(
				slog.String("component", "YearlyStatsTask"),
			),
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (t *YearlyStatsTask) Run(ctx context.Context) error {
	logger := t.logger.With(
		slog.String("operation", "Run"),
	)

	logger.InfoContext(ctx, "operation started")

	from, to := currentYearPeriod(time.Now())
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

func (t *YearlyStatsTask) processCar(ctx context.Context, car domain.Car, from time.Time, to time.Time) error {
	stats, err := t.service.GetYearStatsForPeriod(ctx, car.CreatedBy, service.GetRefuelsForPeriodParams{
		RegNumber: car.RegNumber,
		From:      from,
		To:        to,
	})
	if err != nil {
		if errors.Is(err, service.ErrStatsNotEnoughRefuels) {
			t.logger.DebugContext(ctx, "car has not enough refuels for year", slog.Uint64("carId", uint64(car.ID)))
			return nil
		}

		return err
	}

	text, err := t.renderTemplate("task/yearly_stats.jet", jet.VarMap{}.
		Set("Year", from.Year()).
		Set("NextYear", from.Year()+1).
		Set("Car", car).
		Set("Stats", stats).
		Set("FavoriteWeekday", getWeekdayLabel(stats.FavoriteWeekday)).
		Set("MostActiveMonth", getMonthName(stats.Records.MostActiveMonth)),
	)
	if err != nil {
		return err
	}

	return t.service.AddNotification(ctx, service.AddNotificationParams{
		Key:    domain.NewNotificationKey("yearly-stats", fmt.Sprint(car.ID), fmt.Sprint(from.Year())),
		UserID: car.CreatedBy,
		Text:   text,
	})
}

func currentYearPeriod(now time.Time) (time.Time, time.Time) {
	from := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	return from, from.AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func getWeekdayLabel(weekday time.Weekday) string {
	weekdays := [...]string{
		"воскресеньям",
		"понедельникам",
		"вторникам",
		"средам",
		"четвергам",
		"пятницам",
		"субботам",
	}

	return weekdays[weekday]
}
