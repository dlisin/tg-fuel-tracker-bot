package scheduler

import (
	"context"
)

type ScheduledTask interface {
	Run(ctx context.Context) error
}

type Scheduler interface {
	Schedule(name string, schedule string, task ScheduledTask) error

	Run(ctx context.Context) error
}
