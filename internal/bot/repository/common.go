package repository

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error

	UserRepository() UserRepository
	RefuelRepository() RefuelRepository
}

func WithTransaction(ctx context.Context, uow UnitOfWork, fn func(ctx context.Context, tx Transaction) error) error {
	tx, err := uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx Transaction) {
		_ = tx.Rollback()
	}(tx)

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
