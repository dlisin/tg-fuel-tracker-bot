package model

import (
	"time"
)

type Refuel struct {
	ID            int64     `json:"id"`
	UserID        UserID    `json:"user_id"`
	Odometer      int64     `json:"odometer"`
	Liters        float64   `json:"liters"`
	PricePerLiter float64   `json:"price_per_liter"`
	PriceTotal    float64   `json:"price_total"`
	CreatedAt     time.Time `json:"created_at"`
}
