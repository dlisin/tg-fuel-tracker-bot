package model

import (
	"time"
)

type UserID = int64
type TelegramID = int64

type User struct {
	ID         UserID     `db:"id"`
	TelegramID TelegramID `db:"telegram_id"`
	CreatedAt  time.Time  `db:"created_at"`
}
