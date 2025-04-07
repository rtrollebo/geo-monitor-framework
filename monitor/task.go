package monitor

import (
	"context"
	"time"
)

type Runner interface {
	Run(comm chan TaskResult, ctx context.Context)
}

type TaskResult struct {
	TimeTaken     int64
	Completed     bool
	TimeStarted   time.Time
	TimeCompleted time.Time
	Cause         string
}

type TaskNotifyDefault struct {
	Name     string
	Notifier Notifier
}

func (t TaskNotifyDefault) Run(ch chan TaskResult, ctx context.Context) {
	err := Run(ctx, t.Notifier)
	if err != nil {
		ch <- TaskResult{Cause: err.Error()}
		return
	}
	ch <- TaskResult{TimeTaken: 1, TimeStarted: time.Now(), TimeCompleted: time.Now(), Completed: false}
}
