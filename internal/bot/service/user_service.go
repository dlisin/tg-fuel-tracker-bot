package service

import (
	"context"
	"errors"
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type UserService struct {
	cfg *config.Config
	uow repository.UnitOfWork
}

func NewUserService(cfg *config.Config, uow repository.UnitOfWork) *UserService {
	return &UserService{
		cfg: cfg,
		uow: uow,
	}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID model.TelegramID) (*model.User, error) {
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

	var user *model.User

	user, err = s.findUser(ctx, tx, telegramID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.createUser(ctx, tx, telegramID)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) findUser(ctx context.Context, tx repository.Transaction, telegramID model.TelegramID) (*model.User, error) {
	user, err := tx.UserRepository().GetByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, nil
		}

		return nil, err
	}

	log.Printf("Found user: %+v\n", user)
	return user, nil
}

func (s *UserService) createUser(ctx context.Context, tx repository.Transaction, telegramID model.TelegramID) (*model.User, error) {
	user := &model.User{
		TelegramID: telegramID,
		FuelType:   s.cfg.DefaultFuelType,
		Currency:   s.cfg.DefaultCurrency,
	}

	err := tx.UserRepository().Create(ctx, user)
	if err != nil {
		return nil, err
	}

	log.Printf("User created: %+v\n", user)
	return user, nil
}
