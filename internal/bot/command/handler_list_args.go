package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

type listCommandArgs struct {
	Label  string
	Period model.Range[time.Time]
}

func parseListCommandArgs(cmdArgs string) (*listCommandArgs, error) {
	args := strings.Fields(strings.TrimSpace(cmdArgs))

	now := time.Now()
	if len(args) == 0 {
		return &listCommandArgs{
			Label: "за последний месяц",
			Period: model.Range[time.Time]{
				Start: now.AddDate(0, -1, 0),
				End:   now,
			},
		}, nil
	} else if len(args) == 1 && args[0] == "*" {
		return &listCommandArgs{
			Label: "за всё время",
		}, nil
	} else if len(args) == 2 {
		startDate, err := time.Parse(time.DateOnly, args[0])
		if err != nil {
			return nil, fmt.Errorf("дата должна быть в формате yyyy-mm-dd")
		}

		endDate, err := time.Parse(time.DateOnly, args[1])
		if err != nil {
			return nil, fmt.Errorf("дата должна быть в формате yyyy-mm-dd")
		}

		return &listCommandArgs{
			Label: fmt.Sprintf("за период с %s по %s", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)),
			Period: model.Range[time.Time]{
				Start: startDate,
				End:   endDate,
			},
		}, nil
	}

	return nil, fmt.Errorf("используйте /stats, /stats <from> <to>, /stats *")
}
