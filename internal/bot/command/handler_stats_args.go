package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

type statsCommandArgs struct {
	Label  string
	Period *model.Range[time.Time]
}

func parseStatsCommandArgs(cmdArgs string) (*statsCommandArgs, error) {
	cmdArgs = strings.TrimSpace(cmdArgs)

	if cmdArgs == "" { // last 30 days
		now := time.Now()

		return &statsCommandArgs{
			Label: "за последние месяц",
			Period: &model.Range[time.Time]{
				Start: now.AddDate(0, -1, 0),
				End:   now,
			},
		}, nil
	}

	if cmdArgs == "*" {
		return &statsCommandArgs{
			Label: "за всё время",
		}, nil
	}

	return nil, fmt.Errorf("используйте /stats или /stats *")
}
