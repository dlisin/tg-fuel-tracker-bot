package model

import (
	"log"
	"sort"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/sliceutils"
)

type RefuelStats struct {
	Period          Range[time.Time]
	Entries         int
	TotalDistance   int64
	TotalCost       float64
	TotalLiters     float64
	FuelConsumption float64

	PricePerLiterAverage  float64
	PricePerLiterFirst    float64
	PricePerLiterLast     float64
	PricePerLiterDeltaAbs float64
	PricePerLiterDeltaPct float64
}

func CreateRefuelStats(refuels []Refuel) *RefuelStats {
	stats := &RefuelStats{Entries: len(refuels)}

	if len(refuels) > 1 {
		sort.Slice(refuels, func(i, j int) bool {
			return refuels[i].CreatedAt.Before(refuels[j].CreatedAt)
		})

		first := sliceutils.First(refuels)
		last := sliceutils.Last(refuels)

		stats.Period = Range[time.Time]{
			Start: first.CreatedAt,
			End:   last.CreatedAt,
		}

		for _, refuel := range refuels {
			stats.TotalCost += refuel.PriceTotal
			stats.TotalLiters += refuel.Liters
		}

		stats.TotalDistance = last.Odometer - first.Odometer
		if stats.TotalDistance > 0 {
			fuelUsed := stats.TotalLiters - first.Liters
			stats.FuelConsumption = (fuelUsed / float64(stats.TotalDistance)) * 100
		}

		if stats.TotalLiters > 0 {
			stats.PricePerLiterAverage = stats.TotalCost / stats.TotalLiters
		}
		stats.PricePerLiterFirst = first.PricePerLiter
		stats.PricePerLiterLast = last.PricePerLiter
		stats.PricePerLiterDeltaAbs = last.PricePerLiter - first.PricePerLiter
		stats.PricePerLiterDeltaPct = (last.PricePerLiter - first.PricePerLiter) / first.PricePerLiter * 100
	}

	log.Printf("Refuel stats: %+v\n", stats)
	return stats
}
