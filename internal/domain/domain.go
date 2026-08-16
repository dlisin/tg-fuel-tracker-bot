package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var regNumberRegexp = regexp.MustCompile(`(?i)^[АВЕКМНОРСТУХABEKMHOPCTYX]\d{3}[АВЕКМНОРСТУХABEKMHOPCTYX]{2}\d{2,3}$`)

type CarID uint64

type RegNumber string

type FuelType string

type Mileage uint64

type TelegramID uint64

type Token string

type Car struct {
	ID CarID `db:"id"`

	RegNumber RegNumber `db:"reg_number"`

	FuelType  FuelType  `db:"fuel_type"`
	Odometer  Mileage   `db:"odometer"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserCarInvite struct {
	ID uint64

	CarID     CarID
	Token     Token
	ExpiresAt time.Time

	CreatedAt time.Time
	CreatedBy TelegramID
}

type UserCar struct {
	ID uint64 `db:"id"`

	UserID TelegramID `db:"user_id"`
	CarID  CarID      `db:"car_id"`

	IsOwner   bool      `db:"is_owner"`
	CreatedAt time.Time `db:"created_at"`
}

type Refuel struct {
	ID int64 `db:"id"`

	CarID    CarID   `db:"car_id"`
	Odometer Mileage `db:"odometer"`

	Liters        float64    `db:"liters"`
	PricePerLiter float64    `db:"price_per_liter"`
	PriceTotal    float64    `db:"price_total"`
	CreatedBy     TelegramID `db:"created_by"`
	CreatedAt     time.Time  `db:"created_at"`
}

func ParseRegNumber(value string) (RegNumber, error) {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)
	value = strings.NewReplacer(
		" ", "",
		"А", "A",
		"В", "B",
		"Е", "E",
		"К", "K",
		"М", "M",
		"Н", "H",
		"О", "O",
		"Р", "P",
		"С", "C",
		"Т", "T",
		"У", "Y",
		"Х", "X",
	).Replace(value)

	if !regNumberRegexp.MatchString(value) {
		return "", fmt.Errorf("invalid registration number: %s", value)
	}

	return RegNumber(value), nil
}
