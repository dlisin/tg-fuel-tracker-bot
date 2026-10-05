package template

import (
	"testing"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderYearlyStats(t *testing.T) {
	stats := &service.YearStats{
		Year:               2026,
		Entries:            42,
		TotalCost:          104320,
		TotalDistance:      18742,
		CostPerKm:          5.57,
		RefuelIntervalDays: 9,
		FavoriteWeekday:    time.Friday,
		PriceChange: service.ValueStats{
			From:         58.40,
			To:           64.10,
			DeltaPercent: 9.8,
		},
		Records: service.YearRecords{
			MaxDistanceBetweenRefuels: 782,
			MaxRefuelLiters:            78.4,
			MaxRefuelCost:              5860,
			MostActiveMonth:            time.August,
			MostActiveMonthDistance:    2431,
		},
	}

	text, err := Render("task/yearly_stats.jet", jet.VarMap{}.
		Set("Year", 2026).
		Set("NextYear", 2027).
		Set("Car", domain.Car{RegNumber: domain.RegNumber("M484MX150")}).
		Set("Stats", stats).
		Set("FavoriteWeekday", "пятницам").
		Set("MostActiveMonth", "август").
		Set("PriceChanged", true).
		Set("PriceChangeDirection", "выросла").
		Set("PriceChangePercent", 9.8),
	)
	require.NoError(t, err)

	assert.Contains(t, text, "До 2027 года остались считанные часы")
	assert.Contains(t, text, "18 742 км")
	assert.Contains(t, text, "104 320 ₽")
	assert.Contains(t, text, "5,57 ₽")
	assert.Contains(t, text, "58,40 ₽")
	assert.Contains(t, text, "64,10 ₽")
	assert.Contains(t, text, "782 км")
	assert.Contains(t, text, "78,4 л")
	assert.Contains(t, text, "5 860 ₽")
	assert.Contains(t, text, "август — 2 431 км")
}
