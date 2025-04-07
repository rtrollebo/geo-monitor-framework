package monitor

import (
	"context"
	"log"
	"time"

	"github.com/rtrollebo/geo-monitor-framework/geo"
	"github.com/rtrollebo/geo-monitor-framework/internal"
)

type Notifications struct {
	Time      time.Time
	Recipient string
}

type Notifier interface {
	send() error
	getRecipients() []string
}

func Run(ctx context.Context, notifier Notifier) error {
	logError := ctx.Value("logerror").(*log.Logger)
	logWarning := ctx.Value("logwarning").(*log.Logger)
	logInfo := ctx.Value("loginfo").(*log.Logger)

	notifications, readFileErrorNot := internal.ReadFile[Notifications]("notifications.json")
	if readFileErrorNot != nil {
		logError.Println("Failed to read notifications file")
		return readFileErrorNot
	}

	for _, not := range notifications {
		if not.Time.After(time.Now().Add(time.Duration(-1) * time.Hour)) {
			logInfo.Println("Notification already sent")
			return nil
		}
	}

	events, readFileErr := internal.ReadFile[geo.GeoEvent]("events.json")
	if readFileErr != nil {
		logError.Println("Failed to read events file")
	}
	if events == nil || len(events) == 0 {
		logInfo.Println("No events found")
		return nil
	}

	var recentEvent geo.GeoEvent
	newNotification := false
	for _, event := range events {
		if event.Time.After(recentEvent.Time) && event.Processed && event.Time.After(time.Now().Add(time.Duration(-1)*time.Hour)) {
			recentEvent = event
			newNotification = true
		}
	}
	if !newNotification {
		logInfo.Println("No new events found")
		return nil
	}

	// Write notfications

	recipient := notifier.getRecipients()
	if recipient == nil || len(recipient) == 0 {
		logWarning.Println("No recipients found for notification")
		return nil
	}
	notifications = append(notifications, Notifications{Time: time.Now(), Recipient: notifier.getRecipients()[0]})
	writeErrorNot := internal.WriteFile[Notifications](notifications, "notifications.json")
	if writeErrorNot != nil {
		logError.Println("Failed to write notifications file")
		return writeErrorNot
	}

	errNotify := notifier.send()
	if errNotify != nil {
		logError.Println("Failed to send notification: " + errNotify.Error())
		return errNotify
	}
	return nil
}
