package domain

import (
	"time"
)

type CarID uint64

type FuelType string

type Mileage uint64

type TelegramID uint64

type Token string

type Car struct {
	ID CarID `db:"id"`

	RegNumber RegNumber `db:"reg_number"`
	FuelType  FuelType  `db:"fuel_type"`
	Odometer  Mileage   `db:"odometer"`

	CreatedBy TelegramID `db:"created_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

type UserCarInvite struct {
	ID uint64

	CarID     CarID
	Token     Token
	ExpiresAt time.Time

	CreatedBy TelegramID
	CreatedAt time.Time
}

type UserCar struct {
	ID uint64 `db:"id"`

	UserID TelegramID `db:"user_id"`
	CarID  CarID      `db:"car_id"`

	CreatedAt time.Time `db:"created_at"`
}

type Refuel struct {
	ID int64 `db:"id"`

	CarID    CarID   `db:"car_id"`
	Odometer Mileage `db:"odometer"`

	Liters        float64 `db:"liters"`
	PricePerLiter float64 `db:"price_per_liter"`
	PriceTotal    float64 `db:"price_total"`

	CreatedBy TelegramID `db:"created_by"`
	CreatedAt time.Time  `db:"created_at"`
}
