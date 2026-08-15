package repository

import (
	"context"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type TimeRange struct {
	From time.Time
	To   time.Time
}

type RefuelListParams struct {
	From time.Time
	To   time.Time

	Limit int
	Order SortOrder
}

type CarRepository interface {
	Get(ctx context.Context, regNum domain.RegNumber) (*domain.Car, error)

	Create(ctx context.Context, car *domain.Car) error

	Update(ctx context.Context, car *domain.Car) error

	Delete(ctx context.Context, car *domain.Car) error
}

type UserCarRepository interface {
	Get(ctx context.Context, userID domain.TelegramID, carID domain.CarID) (*domain.UserCar, error)

	List(ctx context.Context, userID domain.TelegramID) ([]domain.Car, error)

	Create(ctx context.Context, userCar *domain.UserCar) error

	Delete(ctx context.Context, userCar *domain.UserCar) error
}

type UserCarInviteRepository interface {
}

type RefuelRepository interface {
	Get(ctx context.Context, carID domain.CarID, odometer domain.Mileage) (*domain.Refuel, error)

	List(ctx context.Context, carID domain.CarID, params RefuelListParams) ([]domain.Refuel, error)

	Create(ctx context.Context, refuel *domain.Refuel) error

	Delete(ctx context.Context, refuel *domain.Refuel) error
}
