package domain

import (
	"strings"
	"time"
)

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSuccess NotificationStatus = "SUCCESS"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

type NotificationID uint64

type NotificationKey string

type NotificationStatus string

type Notification struct {
	ID NotificationID `db:"id"`

	UserID TelegramID      `db:"user_id"`
	Key    NotificationKey `db:"key"`

	Text     string             `db:"text"`
	Status   NotificationStatus `db:"status"`
	Attempts uint32             `db:"attempts"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewNotificationKey(parts ...string) NotificationKey {
	return NotificationKey(strings.Join(parts, ":"))
}

func (k NotificationKey) String() string {
	return string(k)
}
