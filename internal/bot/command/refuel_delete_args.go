package command

import (
	"fmt"
	"strings"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/util/pointerutils"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/util/stringutils"
)

type refuelDeleteCommandArgs struct {
	RegNumber *domain.RegNumber
	Odometer  domain.Mileage
}

func parseRefuelDeleteCommandArgs(cmdArgs string) (*refuelDeleteCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))

	var regNumber *domain.RegNumber
	if len(args) > 0 {
		if value, err := domain.ParseRegNumber(args[0]); err == nil {
			regNumber = pointerutils.AsPointer(value)
			args = args[1:]
		}
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("укажите [<госномер>] <пробег>")
	}

	odometer, err := stringutils.ParseInt64(args[0])
	if err != nil || odometer <= 0 {
		return nil, fmt.Errorf("<пробег> должен быть целым числом > 0")
	}

	return &refuelDeleteCommandArgs{
		RegNumber: regNumber,
		Odometer:  domain.Mileage(odometer),
	}, nil
}
