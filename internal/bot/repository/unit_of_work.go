package repository

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error

	UserRepository() UserRepository
	RefuelRepository() RefuelRepository
}
