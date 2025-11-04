package service

import (
	"context"
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type RefuelService struct {
	cfg *config.Config
	uow repository.UnitOfWork
}

func NewRefuelService(cfg *config.Config, uow repository.UnitOfWork) *RefuelService {
	return &RefuelService{
		cfg: cfg,
		uow: uow,
	}
}

func (s *RefuelService) AddRefuel(ctx context.Context, userID model.UserID, odometer int64, liters float64, totalPrice float64) (*model.Refuel, error) {
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	refuel := &model.Refuel{
		UserID:        userID,
		Odometer:      odometer,
		Liters:        liters,
		PriceTotal:    totalPrice,
		PricePerLiter: totalPrice / liters,
	}

	err = tx.RefuelRepository().Create(ctx, refuel)
	if err != nil {
		return nil, err
	}

	log.Printf("Refuel added: %+v\n", refuel)
	return refuel, nil
}
