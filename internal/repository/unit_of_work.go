package repository

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error

	CarRepository() CarRepository
	UserCarInviteRepository() UserCarInviteRepository
	UserCarRepository() UserCarRepository
	RefuelRepository() RefuelRepository
}

func WithTransaction(ctx context.Context, uow UnitOfWork, fn func(tx Transaction) error) error {
	tx, err := uow.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
