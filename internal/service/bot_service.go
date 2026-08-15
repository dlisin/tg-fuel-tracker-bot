package service

import (
	"context"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
)

type AddCarParams struct {
	RegNumber domain.RegNumber
	FuelType  domain.FuelType
}

type AddRefuelParams struct {
	RegNumber domain.RegNumber
	Odometer  domain.Mileage

	Liters     float64
	PriceTotal float64
}

type DeleteRefuelParams struct {
	RegNumber domain.RegNumber
	Odometer  domain.Mileage
}

type GetRefuelsForPeriodParams struct {
	RegNumber domain.RegNumber

	From time.Time
	To   time.Time
}

type BotService interface {
	AddCar(ctx context.Context, userID domain.TelegramID, params AddCarParams) (*domain.Car, error)

	GetUserCars(ctx context.Context, userID domain.TelegramID) ([]domain.Car, error)

	AddRefuel(ctx context.Context, userID domain.TelegramID, params AddRefuelParams) (*domain.Refuel, error)

	DeleteRefuel(ctx context.Context, userID domain.TelegramID, params DeleteRefuelParams) (*domain.Refuel, error)

	GetRefuelsForPeriod(ctx context.Context, userID domain.TelegramID, params GetRefuelsForPeriodParams) ([]domain.Refuel, error)

	GetRefuelStatsForPeriod(ctx context.Context, userID domain.TelegramID, params GetRefuelsForPeriodParams) (*RefuelStats, error)

	GetLatestRefuelStats(ctx context.Context, userID domain.TelegramID, regNumber domain.RegNumber) (*RefuelStats, error)
}
