package service

import (
	"slices"
	"sort"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
)

type RefuelStats struct {
	From time.Time
	To   time.Time

	Entries       int
	TotalDistance domain.Mileage
	TotalCost     float64
	TotalLiters   float64

	FuelConsumption float64

	PricePerLiterAverage  float64
	PricePerLiterFirst    float64
	PricePerLiterLast     float64
	PricePerLiterDeltaAbs float64
	PricePerLiterDeltaPct float64
}

func CalculateRefuelStats(refuels []domain.Refuel) (*RefuelStats, error) {
	if len(refuels) < 2 {
		return nil, ErrStatsNotEnoughRefuels
	}

	refuels = slices.Clone(refuels)
	sort.Slice(refuels, func(i, j int) bool {
		return refuels[i].Odometer < refuels[j].Odometer
	})

	first := refuels[0]
	last := refuels[len(refuels)-1]

	stats := &RefuelStats{
		From:          first.CreatedAt,
		To:            last.CreatedAt,
		Entries:       len(refuels),
		TotalDistance: last.Odometer - first.Odometer,

		PricePerLiterFirst: first.PricePerLiter,
		PricePerLiterLast:  last.PricePerLiter,
	}

	var fuelUsed float64
	for i, refuel := range refuels {
		stats.TotalCost += refuel.PriceTotal
		stats.TotalLiters += refuel.Liters

		if i > 0 {
			fuelUsed += refuel.Liters
		}
	}

	if stats.TotalDistance > 0 {
		stats.FuelConsumption = fuelUsed / float64(stats.TotalDistance) * 100
	}

	if stats.TotalLiters > 0 {
		stats.PricePerLiterAverage = stats.TotalCost / stats.TotalLiters
	}

	stats.PricePerLiterDeltaAbs = stats.PricePerLiterLast - stats.PricePerLiterFirst

	if stats.PricePerLiterFirst > 0 {
		stats.PricePerLiterDeltaPct = stats.PricePerLiterDeltaAbs / stats.PricePerLiterFirst * 100
	}

	return stats, nil
}
