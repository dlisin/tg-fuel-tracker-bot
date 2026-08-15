package service

import "errors"

var (
	// User
	ErrUserHasNoAccessToCar = &ServiceError{err: errors.New("user has no access to the specified car")}

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
