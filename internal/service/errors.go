package service

import (
	"context"
	"errors"
	"log/slog"
)

var (
	// Car
	ErrCarNotFound      = &ServiceError{err: errors.New("car with the specified registration number not found")}
	ErrCarAlreadyExists = &ServiceError{err: errors.New("car with the specified registration number already exists")}

	// Refuel
	ErrRefuelNotFound       = &ServiceError{err: errors.New("refuel with provided odometer value not found")}
	ErrRefuelOdometerTooLow = &ServiceError{err: errors.New("refuel odometer value must be greater than current car odometer value")}

	// RefuelStats
	ErrStatsNotEnoughRefuels = &ServiceError{err: errors.New("not enough refuels to calculate statistics")}
)

type ServiceError struct {
	err error
}

func (e *ServiceError) Error() string {
	return e.err.Error()
}

func (e *ServiceError) Unwrap() error {
	return e.err
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
