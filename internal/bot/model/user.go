package model

import (
	"time"
)

type UserID = int64
type TelegramID = int64

type User struct {
	ID         UserID     `json:"id"`
	TelegramID TelegramID `json:"telegram_id"`
	FuelType   string     `json:"fuel_type"`
	Currency   string     `json:"currency"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
