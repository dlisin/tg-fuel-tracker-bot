package model

import (
	"time"
)

type Refuel struct {
	ID            int64      `db:"id"`
	UserID        TelegramID `db:"user_id"`
	Odometer      int64      `db:"odometer"`
	Liters        float64    `db:"liters"`
	PricePerLiter float64    `db:"price_per_liter"`
	PriceTotal    float64    `db:"price_total"`
	CreatedAt     time.Time  `db:"created_at"`
}
