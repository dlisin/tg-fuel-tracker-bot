package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/util/pointerutils"
)

type listCommandArgs struct {
	RegNumber *domain.RegNumber
	Label     string
	From      time.Time
	To        time.Time
}

func parseListCommandArgs(cmdArgs string) (*listCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))

	var regNumber *domain.RegNumber
	if len(args) > 0 {
		if value, err := domain.ParseRegNumber(args[0]); err == nil {
			regNumber = pointerutils.AsPointer(value)
			args = args[1:]
		}
	}

	now := time.Now()
	switch len(args) {
	case 0:
		return &listCommandArgs{
			RegNumber: regNumber,
			Label:     "за последний месяц",
			From:      now.AddDate(0, -1, 0),
			To:        now,
		}, nil

	case 1:
		if args[0] != "*" {
			return nil, fmt.Errorf("укажите [<гос номер>] [<начало периода> <конец периода> | *]")
		}

		return &listCommandArgs{
			RegNumber: regNumber,
			Label:     "за всё время",
		}, nil

	case 2:
		from, err := time.Parse(time.DateOnly, args[0])
		if err != nil {
			return nil, fmt.Errorf("<начало периода> должно быть в формате yyyy-mm-dd")
		}

		to, err := time.Parse(time.DateOnly, args[1])
		if err != nil {
			return nil, fmt.Errorf("<конец периода> должен быть в формате yyyy-mm-dd")
		}

		return &listCommandArgs{
			RegNumber: regNumber,
			Label:     fmt.Sprintf("за период с %s по %s", from.Format(time.DateOnly), to.Format(time.DateOnly)),
			From:      from,
			To:        to,
		}, nil

	default:
		return nil, fmt.Errorf("укажите [<гос номер>] [<начало периода> <конец периода> | *]")
	}
}
