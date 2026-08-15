package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/repository"
)

type BotServiceImpl struct {
	logger *slog.Logger
	uow    repository.UnitOfWork
}

func NewBotService(logger *slog.Logger, uow repository.UnitOfWork) *BotServiceImpl {
	return &BotServiceImpl{
		logger: logger.With(
			slog.String("component", "BotService"),
		),
		uow: uow,
	}
}

func (s *BotServiceImpl) AddCar(ctx context.Context, userID domain.TelegramID, params AddCarParams) (*domain.Car, error) {
	logger := s.logger.With(
		slog.String("operation", "AddCar"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(params.RegNumber)),
	)

	logger.InfoContext(ctx, "operation started")

	now := time.Now()
	car := &domain.Car{
		RegNumber: params.RegNumber,
		FuelType:  params.FuelType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := repository.WithTransaction(ctx, s.uow, func(tx repository.Transaction) error {
		if err := tx.CarRepository().Create(ctx, car); err != nil {
			if errors.Is(err, repository.ErrEntityAlreadyExists) {
				return ErrCarAlreadyExists
			}

			return err
		}

		logger.DebugContext(ctx, "car created",
			slog.Uint64("carId", uint64(car.ID)),
			slog.String("fuelType", string(params.FuelType)),
		)

		userCar := &domain.UserCar{
			UserID:    userID,
			CarID:     car.ID,
			IsOwner:   true,
			CreatedAt: now,
		}
		if err := tx.UserCarRepository().Create(ctx, userCar); err != nil {
			return err
		}

		logger.DebugContext(ctx, "car owner created", slog.Uint64("userCarId", userCar.ID))
		return nil
	})
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Uint64("carId", uint64(car.ID)))
	return car, nil
}

func (s *BotServiceImpl) GetUserCars(ctx context.Context, userID domain.TelegramID) ([]domain.Car, error) {
	logger := s.logger.With(
		slog.String("operation", "GetUserCars"),
		slog.Uint64("userId", uint64(userID)),
	)

	logger.InfoContext(ctx, "operation started")

	var cars []domain.Car
	err := repository.WithTransaction(ctx, s.uow, func(tx repository.Transaction) error {
		var err error
		cars, err = tx.UserCarRepository().List(ctx, userID)
		return err
	})
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("carsCount", len(cars)))

	return cars, nil
}

func (s *BotServiceImpl) AddRefuel(ctx context.Context, userID domain.TelegramID, params AddRefuelParams) (*domain.Refuel, error) {
	logger := s.logger.With(
		slog.String("operation", "AddRefuel"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(params.RegNumber)),
		slog.Uint64("odometer", uint64(params.Odometer)),
		slog.Float64("liters", params.Liters),
		slog.Float64("priceTotal", params.PriceTotal),
	)

	logger.InfoContext(ctx, "operation started")

	now := time.Now()
	var refuel *domain.Refuel
	err := repository.WithTransaction(ctx, s.uow, func(tx repository.Transaction) error {
		car, err := tx.CarRepository().Get(ctx, params.RegNumber)
		if err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrCarNotFound
			}

			return err
		}
		logger.DebugContext(ctx, "car found", slog.Uint64("carId", uint64(car.ID)))

		if _, err := tx.UserCarRepository().Get(ctx, userID, car.ID); err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrUserHasNoAccessToCar
			}

			return err
		}
		logger.DebugContext(ctx, "car access verified")

		if params.Odometer <= car.Odometer {
			return ErrRefuelOdometerTooLow
		}

		refuel = &domain.Refuel{
			CarID:         car.ID,
			Odometer:      params.Odometer,
			Liters:        params.Liters,
			PricePerLiter: params.PriceTotal / params.Liters,
			PriceTotal:    params.PriceTotal,
			CreatedAt:     now,
			CreatedBy:     userID,
		}

		if err := tx.RefuelRepository().Create(ctx, refuel); err != nil {
			return err
		}
		logger.DebugContext(ctx, "refuel created", slog.Int64("refuelId", refuel.ID))

		car.Odometer = params.Odometer
		car.UpdatedAt = now
		if err := tx.CarRepository().Update(ctx, car); err != nil {
			return err
		}
		logger.DebugContext(ctx, "car odometer updated", slog.Uint64("odometer", uint64(car.Odometer)))

		return nil
	})
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Int64("refuelId", refuel.ID))
	return refuel, nil
}

func (s *BotServiceImpl) DeleteRefuel(ctx context.Context, userID domain.TelegramID, params DeleteRefuelParams) (*domain.Refuel, error) {
	logger := s.logger.With(
		slog.String("operation", "DeleteRefuel"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(params.RegNumber)),
		slog.Uint64("odometer", uint64(params.Odometer)),
	)

	logger.InfoContext(ctx, "operation started")

	var refuel *domain.Refuel
	err := repository.WithTransaction(ctx, s.uow, func(tx repository.Transaction) error {
		car, err := tx.CarRepository().Get(ctx, params.RegNumber)
		if err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrCarNotFound
			}

			return err
		}
		logger.DebugContext(ctx, "car found", slog.Uint64("carId", uint64(car.ID)))

		if _, err := tx.UserCarRepository().Get(ctx, userID, car.ID); err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrUserHasNoAccessToCar
			}

			return err
		}
		logger.DebugContext(ctx, "car access verified")

		refuel, err = tx.RefuelRepository().Get(ctx, car.ID, params.Odometer)
		if err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrRefuelNotFound
			}

			return err
		}
		logger.DebugContext(ctx, "refuel found", slog.Int64("refuelId", refuel.ID))

		if err := tx.RefuelRepository().Delete(ctx, refuel); err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrRefuelNotFound
			}

			return err
		}
		logger.DebugContext(ctx, "refuel deleted", slog.Int64("refuelId", refuel.ID))

		if refuel.Odometer != car.Odometer {
			return nil
		}

		refuels, err := tx.RefuelRepository().List(ctx, car.ID, repository.RefuelListParams{
			Limit: 1,
			Order: repository.SortOrderDesc,
		})
		if err != nil {
			return err
		}

		if len(refuels) > 0 {
			car.Odometer = refuels[0].Odometer
		} else {
			car.Odometer = 0
		}
		car.UpdatedAt = time.Now()

		if err := tx.CarRepository().Update(ctx, car); err != nil {
			return err
		}

		logger.DebugContext(ctx, "car odometer updated", slog.Uint64("odometer", uint64(car.Odometer)))
		return nil
	})
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed")
	return refuel, nil
}

func (s *BotServiceImpl) GetRefuelsForPeriod(ctx context.Context, userID domain.TelegramID, params GetRefuelsForPeriodParams) ([]domain.Refuel, error) {
	logger := s.logger.With(
		slog.String("operation", "GetRefuelsForPeriod"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(params.RegNumber)),
		slog.Time("from", params.From),
		slog.Time("to", params.To),
	)

	logger.InfoContext(ctx, "operation started")

	refuels, err := s.listRefuels(logger, ctx, userID, params.RegNumber, repository.RefuelListParams{
		From: params.From,
		To:   params.To,
	})
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("count", len(refuels)))
	return refuels, nil
}

func (s *BotServiceImpl) GetRefuelStatsForPeriod(ctx context.Context, userID domain.TelegramID, params GetRefuelsForPeriodParams) (*RefuelStats, error) {
	logger := s.logger.With(
		slog.String("operation", "GetRefuelStatsForPeriod"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(params.RegNumber)),
		slog.Time("from", params.From),
		slog.Time("to", params.To),
	)

	logger.InfoContext(ctx, "operation started")

	var stats *RefuelStats
	err := func() error {
		refuels, err := s.listRefuels(logger, ctx, userID, params.RegNumber, repository.RefuelListParams{
			From: params.From,
			To:   params.To,
		})
		if err != nil {
			return err
		}

		stats, err = CalculateRefuelStats(refuels)
		if err != nil {
			return err
		}
		logger.DebugContext(ctx, "refuel stats calculated",
			slog.Int("entries", stats.Entries),
			slog.Uint64("totalDistance", uint64(stats.TotalDistance)),
			slog.Float64("totalLiters", stats.TotalLiters),
			slog.Float64("totalCost", stats.TotalCost),
			slog.Float64("fuelConsumption", stats.FuelConsumption),
		)

		return nil
	}()
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("entries", stats.Entries))
	return stats, nil
}

func (s *BotServiceImpl) GetLatestRefuelStats(ctx context.Context, userID domain.TelegramID, regNumber domain.RegNumber) (*RefuelStats, error) {
	logger := s.logger.With(
		slog.String("operation", "GetLatestRefuelStats"),
		slog.Uint64("userId", uint64(userID)),
		slog.String("regNumber", string(regNumber)),
	)

	logger.InfoContext(ctx, "operation started")

	var stats *RefuelStats
	err := func() error {
		refuels, err := s.listRefuels(logger, ctx, userID, regNumber, repository.RefuelListParams{
			Limit: 2,
			Order: repository.SortOrderDesc,
		})
		if err != nil {
			return err
		}

		stats, err = CalculateRefuelStats(refuels)
		if err != nil {
			return err
		}

		logger.DebugContext(ctx, "refuel stats calculated",
			slog.Int("entries", stats.Entries),
			slog.Uint64("totalDistance", uint64(stats.TotalDistance)),
			slog.Float64("totalLiters", stats.TotalLiters),
			slog.Float64("totalCost", stats.TotalCost),
			slog.Float64("fuelConsumption", stats.FuelConsumption),
		)

		return nil
	}()
	if err != nil {
		return nil, handleServiceError(logger, ctx, err)
	}

	logger.InfoContext(ctx, "operation completed", slog.Int("entries", stats.Entries))
	return stats, nil
}

func (s *BotServiceImpl) listRefuels(logger *slog.Logger, ctx context.Context, userID domain.TelegramID, regNumber domain.RegNumber, params repository.RefuelListParams) ([]domain.Refuel, error) {
	var refuels []domain.Refuel
	err := repository.WithTransaction(ctx, s.uow, func(tx repository.Transaction) error {
		car, err := tx.CarRepository().Get(ctx, regNumber)
		if err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrCarNotFound
			}
			return err
		}
		logger.DebugContext(ctx, "car found", slog.Uint64("carId", uint64(car.ID)))

		if _, err := tx.UserCarRepository().Get(ctx, userID, car.ID); err != nil {
			if errors.Is(err, repository.ErrEntityNotFound) {
				return ErrUserHasNoAccessToCar
			}

			return err
		}
		logger.DebugContext(ctx, "car access verified")

		refuels, err = tx.RefuelRepository().List(ctx, car.ID, params)
		if err != nil {
			return err
		}
		logger.DebugContext(ctx, "refuels found", slog.Int("count", len(refuels)))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return refuels, nil
}

func handleServiceError(logger *slog.Logger, ctx context.Context, err error) error {
	var serviceErr *ServiceError
	if errors.As(err, &serviceErr) {
		logger.WarnContext(ctx, "operation aborted", slog.Any("error", err))
		return err
	}

	logger.ErrorContext(ctx, "operation failed", slog.Any("error", err))
	return err
}
