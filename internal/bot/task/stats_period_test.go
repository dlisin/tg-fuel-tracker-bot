package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldSkipMonthlyStats(t *testing.T) {
	assert.True(t, shouldSkipMonthlyStats(
		time.Date(2027, time.January, 1, 0, 1, 0, 0, time.UTC),
		true,
	))

	assert.False(t, shouldSkipMonthlyStats(
		time.Date(2027, time.January, 1, 0, 1, 0, 0, time.UTC),
		false,
	))

	assert.False(t, shouldSkipMonthlyStats(
		time.Date(2027, time.February, 1, 0, 1, 0, 0, time.UTC),
		true,
	))
}

func TestCurrentYearPeriod(t *testing.T) {
	location := time.FixedZone("MSK", 3*60*60)

	t.Run("December 31", func(t *testing.T) {
		now := time.Date(2026, time.December, 31, 23, 0, 0, 0, location)

		from, to := currentYearPeriod(now)

		assert.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, location), from)
		assert.Equal(t, time.Date(2026, time.December, 31, 23, 59, 59, 999999999, location), to)
	})

	t.Run("January retry", func(t *testing.T) {
		now := time.Date(2027, time.January, 1, 0, 5, 0, 0, location)

		from, to := currentYearPeriod(now)

		assert.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, location), from)
		assert.Equal(t, time.Date(2026, time.December, 31, 23, 59, 59, 999999999, location), to)
	})
}
