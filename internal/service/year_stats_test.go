package service

import (
	"testing"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateYearStats(t *testing.T) {
	refuels := []domain.Refuel{
		{
			Odometer:      10000,
			Liters:        50,
			PricePerLiter: 60,
			PriceTotal:    3000,
			CreatedAt:     time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			Odometer:      10500,
			Liters:        55,
			PricePerLiter: 61,
			PriceTotal:    3355,
			CreatedAt:     time.Date(2026, time.February, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			Odometer:      11200,
			Liters:        70,
			PricePerLiter: 63,
			PriceTotal:    4410,
			CreatedAt:     time.Date(2026, time.August, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			Odometer:      11850,
			Liters:        60,
			PricePerLiter: 65,
			PriceTotal:    3900,
			CreatedAt:     time.Date(2026, time.December, 4, 12, 0, 0, 0, time.UTC),
		},
	}

	stats, err := CalculateYearStats(refuels)
	require.NoError(t, err)

	assert.Equal(t, 2026, stats.Year)
	assert.Equal(t, 4, stats.Entries)
	assert.Equal(t, domain.Mileage(1850), stats.TotalDistance)
	assert.Equal(t, 235.0, stats.TotalLiters)
	assert.Equal(t, 14665.0, stats.TotalCost)
	assert.InDelta(t, 7.927, stats.CostPerKm, 0.001)

	assert.Equal(t, 112.0, stats.RefuelIntervalDays)
	assert.Equal(t, time.Friday, stats.FavoriteWeekday)

	assert.Equal(t, 60.0, stats.PriceChange.From)
	assert.Equal(t, 65.0, stats.PriceChange.To)
	assert.Equal(t, 5.0, stats.PriceChange.Delta)
	assert.InDelta(t, 8.333, stats.PriceChange.DeltaPercent, 0.001)

	assert.Equal(t, domain.Mileage(700), stats.Records.MaxDistanceBetweenRefuels)
	assert.Equal(t, 70.0, stats.Records.MaxRefuelLiters)
	assert.Equal(t, 4410.0, stats.Records.MaxRefuelCost)
	assert.Equal(t, time.August, stats.Records.MostActiveMonth)
	assert.Equal(t, domain.Mileage(700), stats.Records.MostActiveMonthDistance)
}

func TestCalculateYearStatsNotEnoughRefuels(t *testing.T) {
	_, err := CalculateYearStats([]domain.Refuel{{Odometer: 10000}})
	require.ErrorIs(t, err, ErrStatsNotEnoughRefuels)
}
