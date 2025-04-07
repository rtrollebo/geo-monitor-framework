package monitor

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func RunSchedule(ctx context.Context, interval int, sign chan os.Signal, messages chan TaskResult, tasks []Runner) {
	logInfo := ctx.Value("loginfo").(*log.Logger)
	for {
		select {
		case <-sign:
			logInfo.Println("Received signal to stop")
			return
		default:
			for _, v := range tasks {
				go v.Run(messages, ctx)
			}
			readChannelErr := readChannel(messages, 3)
			if readChannelErr != nil {
				logInfo.Println("Error reading channel:", readChannelErr)
				os.Exit(1)
			}
			time.Sleep(time.Duration(interval) * time.Second)

		}
	}
}

func readChannel(messages chan TaskResult, totalNumberOfTasks int) error {
	timeOut := 10
	for {
		select {
		case <-messages:
		case <-time.After(time.Second * time.Duration(timeOut)):
			if runtime.NumGoroutine() > totalNumberOfTasks {
				return fmt.Errorf("Tasks taking too long time")
			}
			return nil
		}
	}
}
