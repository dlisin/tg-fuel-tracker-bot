package command

import (
	"fmt"
	"strings"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/util/stringutils"
)

type deleteCommandArgs struct {
	Odometer int64
}

func parseDeleteCommandArgs(cmdArgs string) (*deleteCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))
	if len(args) < 1 {
		return nil, fmt.Errorf("недостаточно параметров, укажите <odometer>")
	}

	odometer, err := stringutils.ParseInt64(args[0])
	if err != nil || odometer < 0 {
		return nil, fmt.Errorf("пробег должен быть целым числом ≥ 0")
	}

	return &deleteCommandArgs{
		Odometer: odometer,
	}, nil
}
