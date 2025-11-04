package stats

import (
	"fmt"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/db"
)

type Stats struct {
	DistanceKm     float64
	AvgConsumption float64 // L/100km
	PriceDeltaAbs  float64 // last - first
	PriceDeltaPct  float64 // percent
	FirstPPL       float64
	LastPPL        float64
	Entries        int
	Label          string
	StartEnd       *[2]time.Time
}

func MonthRange(now time.Time) [2]time.Time {
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)
	return [2]time.Time{start, end}
}

func MonthRangeYM(y int, m time.Month, loc *time.Location) [2]time.Time {
	start := time.Date(y, m, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	return [2]time.Time{start, end}
}

func YearRange(y int, loc *time.Location) [2]time.Time {
	start := time.Date(y, 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(y+1, 1, 1, 0, 0, 0, 0, loc)
	return [2]time.Time{start, end}
}

func ParseStatsRange(arg string, now time.Time) (*[2]time.Time, string, error) {
	arg = trim(arg)
	if arg == "" {
		return nil, "за всё время", nil
	}
	if equalFold(arg, "month") || equalFold(arg, "месяц") {
		r := MonthRange(now)
		return &r, "за текущий месяц", nil
	}
	// YYYY-MM
	if len(arg) == 7 && count(arg, "-") == 1 {
		var y, m int
		_, err := fmt.Sscanf(arg, "%d-%d", &y, &m)
		if err == nil && y >= 1970 && m >= 1 && m <= 12 {
			r := MonthRangeYM(y, time.Month(m), now.Location())
			return &r, fmt.Sprintf("за %04d-%02d", y, m), nil
		}
	}
	// Year
	var y int
	if _, err := fmt.Sscanf(arg, "%d", &y); err == nil && y >= 1970 && y <= 3000 {
		r := YearRange(y, now.Location())
		return &r, fmt.Sprintf("за %d год", y), nil
	}
	return nil, "", fmt.Errorf("используйте /stats, /stats month, /stats <год> или /stats <год-месяц>")
}

func Calc(fs []db.Fillup) *Stats {
	if len(fs) < 2 {
		return &Stats{Entries: len(fs)}
	}
	minOdo := fs[0].Odometer
	maxOdo := fs[0].Odometer
	var sumLiters float64
	for _, f := range fs {
		if f.Odometer < minOdo {
			minOdo = f.Odometer
		}
		if f.Odometer > maxOdo {
			maxOdo = f.Odometer
		}
		sumLiters += f.Liters
	}
	dist := float64(maxOdo - minOdo)
	if dist <= 0 {
		return &Stats{Entries: len(fs)}
	}
	avg := (sumLiters / dist) * 100
	firstPPL := fs[0].PricePerLiter
	lastPPL := fs[len(fs)-1].PricePerLiter
	delta := lastPPL - firstPPL
	pct := 0.0
	if firstPPL != 0 {
		pct = (delta / firstPPL) * 100
	}
	return &Stats{
		DistanceKm:     dist,
		AvgConsumption: avg,
		PriceDeltaAbs:  delta,
		PriceDeltaPct:  pct,
		FirstPPL:       firstPPL,
		LastPPL:        lastPPL,
		Entries:        len(fs),
	}
}

// tiny helpers (no extra deps)
func trim(s string) string       { return stringsTrimSpace(s) }
func equalFold(a, b string) bool { return stringsEqualFold(a, b) }
func count(s, sep string) int    { return stringsCount(s, sep) }

// minimal inline replacements to avoid importing strings in this file
func stringsTrimSpace(s string) string {
	i, j := 0, len(s)-1
	for i <= j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n' || s[j] == '\r') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}
func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ai, bi := a[i], b[i]
		if ai >= 'A' && ai <= 'Z' {
			ai += 32
		}
		if bi >= 'A' && bi <= 'Z' {
			bi += 32
		}
		if ai != bi {
			return false
		}
	}
	return true
}
func stringsCount(s, sub string) int {
	if sub == "" {
		return 0
	}
	c := 0
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			c++
		}
	}
	return c
}
