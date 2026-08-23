package service

import (
	"slices"
	"sort"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
)

type RefuelStats struct {
	From    time.Time
	To      time.Time
	Entries int

	TotalCost   float64
	TotalLiters float64

	Price       PriceStats
	Consumption *ConsumptionStats
}

type PriceStats struct {
	Average   float64
	FirstLast ValueStats
	MinMax    ValueStats
}

type ValueStats struct {
	From         float64
	To           float64
	Delta        float64
	DeltaPercent float64
}

type ConsumptionStats struct {
	TotalDistance   domain.Mileage
	FuelConsumption float64
}

func CalculateRefuelStats(refuels []domain.Refuel) (*RefuelStats, error) {
	if len(refuels) == 0 {
		return nil, ErrStatsNotEnoughRefuels
	}

	refuels = slices.Clone(refuels)
	sort.Slice(refuels, func(i, j int) bool {
		return refuels[i].Odometer < refuels[j].Odometer
	})

	first := refuels[0]
	last := refuels[len(refuels)-1]

	stats := &RefuelStats{
		From:    first.CreatedAt,
		To:      last.CreatedAt,
		Entries: len(refuels),
		Price: PriceStats{
			FirstLast: calculateValueStats(
				first.PricePerLiter,
				last.PricePerLiter,
			),
		},
	}

	minPrice := first.PricePerLiter
	maxPrice := first.PricePerLiter

	var fuelUsed float64
	var priceSum float64

	for i, refuel := range refuels {
		stats.TotalCost += refuel.PriceTotal
		stats.TotalLiters += refuel.Liters

		priceSum += refuel.PricePerLiter

		if refuel.PricePerLiter < minPrice {
			minPrice = refuel.PricePerLiter
		}

		if refuel.PricePerLiter > maxPrice {
			maxPrice = refuel.PricePerLiter
		}

		if i > 0 {
			fuelUsed += refuel.Liters
		}
	}

	stats.Price.Average = priceSum / float64(len(refuels))

	stats.Price.MinMax = calculateValueStats(
		minPrice,
		maxPrice,
	)

	if len(refuels) >= 2 {
		distance := last.Odometer - first.Odometer

		stats.Consumption = &ConsumptionStats{
			TotalDistance: distance,
		}

		if distance > 0 {
			stats.Consumption.FuelConsumption =
				fuelUsed / float64(distance) * 100
		}
	}

	return stats, nil
}

func calculateValueStats(from, to float64) ValueStats {
	stats := ValueStats{
		From:  from,
		To:    to,
		Delta: to - from,
	}

	if from > 0 {
		stats.DeltaPercent = stats.Delta / from * 100
	}

	return stats
}
