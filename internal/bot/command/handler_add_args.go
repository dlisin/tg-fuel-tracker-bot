package command

import (
	"fmt"
	"strings"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/stringutils"
)

type addCommandArgs struct {
	Odometer   int64
	Liters     float64
	TotalPrice float64
}

func parseAddCommandArgs(cmdArgs string, prevRefuel *model.Refuel) (*addCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))
	if len(args) < 3 {
		return nil, fmt.Errorf("недостаточно параметров, укажите <пробег> <литры> <сумма чека>")
	}

	odometer, err := stringutils.ParseInt64(args[0])
	if err != nil || odometer < 0 {
		return nil, fmt.Errorf("пробег должен быть целым числом ≥ 0")
	}

	liters, err := stringutils.ParseFloat64(args[1])
	if err != nil || liters <= 0 {
		return nil, fmt.Errorf("литры должны быть числом > 0")
	}

	totalPrice, err := stringutils.ParseFloat64(args[2])
	if err != nil || totalPrice <= 0 {
		return nil, fmt.Errorf("сумма чека должна быть числом > 0")
	}

	if prevRefuel != nil {
		prevOdometer := prevRefuel.Odometer
		if prevOdometer >= odometer {
			return nil, fmt.Errorf("пробег должен быть больше предыдущего (%d)", prevOdometer)
		}
	}

	return &addCommandArgs{
		Odometer:   odometer,
		Liters:     liters,
		TotalPrice: totalPrice,
	}, nil
}
