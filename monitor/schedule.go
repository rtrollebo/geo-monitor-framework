package monitor

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type TaskSchedule struct {
	Task     Runner
	Delay    time.Duration
	Interval time.Duration
}

func RunSchedule(ctx context.Context, tasks []TaskSchedule, messages chan TaskResult) {
	logError := ctx.Value("logerror").(*log.Logger)
	var wg sync.WaitGroup

	for _, schedule := range tasks {
		wg.Add(1)
		go func(s TaskSchedule) {
			defer wg.Done()

			select {
			case <-time.After(s.Delay):
			case <-ctx.Done():
				return
			}
			for {
				select {
				case <-ctx.Done():
					return
				default:
					s.Task.Run(messages, ctx)
				}

				select {
				case <-time.After(s.Interval):
				case <-ctx.Done():
					return
				}
			}
		}(schedule)
	}

	go func() {
		err := readChannel(ctx, messages, len(tasks))
		if err != nil {
			logError.Println("Failed to finalize tasks: " + err.Error())
		}
	}()

	wg.Wait()
}

func readChannel(ctx context.Context, messages chan TaskResult, totalNumberOfTasks int) error {
	logInfo := ctx.Value("loginfo").(*log.Logger)
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg := <-messages:
			logInfo.Println("Message from task: " + msg.Cause)

			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)

		case <-timer.C:
			if runtime.NumGoroutine() > totalNumberOfTasks+5 {
				return fmt.Errorf("tasks taking too long time")
			}
			return nil

		case <-ctx.Done():
			logInfo.Println("Application shutdown received")
			return nil
		}
	}
}
