package command

import (
	"fmt"
	"strings"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/util/pointerutils"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/util/stringutils"
)

type refuelAddCommandArgs struct {
	RegNumber  *domain.RegNumber
	Odometer   domain.Mileage
	Liters     float64
	TotalPrice float64
}

func parseRefuelAddCommandArgs(cmdArgs string) (*refuelAddCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))

	var regNumber *domain.RegNumber
	if len(args) > 0 {
		if value, err := domain.ParseRegNumber(args[0]); err == nil {
			regNumber = pointerutils.AsPointer(value)
			args = args[1:]
		}
	}

	if len(args) != 3 {
		return nil, fmt.Errorf("укажите [<госномер>] <пробег> <литры> <сумма чека>")
	}

	odometer, err := stringutils.ParseInt64(args[0])
	if err != nil || odometer < 0 {
		return nil, fmt.Errorf("<пробег> должен быть целым числом ≥ 0")
	}

	liters, err := stringutils.ParseFloat64(args[1])
	if err != nil || liters <= 0 {
		return nil, fmt.Errorf("<литры> должны быть числом > 0")
	}

	totalPrice, err := stringutils.ParseFloat64(args[2])
	if err != nil || totalPrice <= 0 {
		return nil, fmt.Errorf("<сумма чека> должна быть числом > 0")
	}

	return &refuelAddCommandArgs{
		RegNumber:  regNumber,
		Odometer:   domain.Mileage(odometer),
		Liters:     liters,
		TotalPrice: totalPrice,
	}, nil
}
