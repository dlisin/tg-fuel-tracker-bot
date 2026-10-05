package service

import (
	"slices"
	"sort"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
)

type YearStats struct {
	Year    int
	Entries int

	TotalCost     float64
	TotalLiters   float64
	TotalDistance domain.Mileage
	CostPerKm     float64

	RefuelIntervalDays float64
	FavoriteWeekday    time.Weekday

	PriceChange ValueStats
	Records     YearRecords
}

type YearRecords struct {
	MaxDistanceBetweenRefuels domain.Mileage
	MaxRefuelLiters            float64
	MaxRefuelCost              float64

	MostActiveMonth         time.Month
	MostActiveMonthDistance domain.Mileage
}

func CalculateYearStats(refuels []domain.Refuel) (*YearStats, error) {
	if len(refuels) < 2 {
		return nil, ErrStatsNotEnoughRefuels
	}

	refuels = slices.Clone(refuels)
	sort.Slice(refuels, func(i, j int) bool {
		return refuels[i].Odometer < refuels[j].Odometer
	})

	first := refuels[0]
	last := refuels[len(refuels)-1]

	stats := &YearStats{
		Year:          first.CreatedAt.Year(),
		Entries:       len(refuels),
		TotalDistance: last.Odometer - first.Odometer,
		PriceChange: calculateValueStats(
			first.PricePerLiter,
			last.PricePerLiter,
		),
	}

	weekdayCounts := make(map[time.Weekday]int)
	monthDistances := make(map[time.Month]domain.Mileage)

	for i, refuel := range refuels {
		stats.TotalCost += refuel.PriceTotal
		stats.TotalLiters += refuel.Liters

		if refuel.Liters > stats.Records.MaxRefuelLiters {
			stats.Records.MaxRefuelLiters = refuel.Liters
		}

		if refuel.PriceTotal > stats.Records.MaxRefuelCost {
			stats.Records.MaxRefuelCost = refuel.PriceTotal
		}

		weekdayCounts[refuel.CreatedAt.Weekday()]++

		if i == 0 {
			continue
		}

		distance := refuel.Odometer - refuels[i-1].Odometer
		if distance > stats.Records.MaxDistanceBetweenRefuels {
			stats.Records.MaxDistanceBetweenRefuels = distance
		}

		monthDistances[refuel.CreatedAt.Month()] += distance
	}

	if stats.TotalDistance > 0 {
		stats.CostPerKm = stats.TotalCost / float64(stats.TotalDistance)
	}

	periodDays := last.CreatedAt.Sub(first.CreatedAt).Hours() / 24
	if periodDays > 0 {
		stats.RefuelIntervalDays = periodDays / float64(len(refuels)-1)
	}

	maxWeekdayCount := 0
	for _, refuel := range refuels {
		weekday := refuel.CreatedAt.Weekday()
		if weekdayCounts[weekday] > maxWeekdayCount {
			stats.FavoriteWeekday = weekday
			maxWeekdayCount = weekdayCounts[weekday]
		}
	}

	for month := time.January; month <= time.December; month++ {
		distance := monthDistances[month]
		if distance > stats.Records.MostActiveMonthDistance {
			stats.Records.MostActiveMonth = month
			stats.Records.MostActiveMonthDistance = distance
		}
	}

	return stats, nil
}
