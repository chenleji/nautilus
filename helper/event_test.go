package helper

import (
	logs "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestEvent_Broadcast(t *testing.T) {
	// registry event
	GetEventMap().Registry(
		&Event{
			Name: "event-1",
			// define callback function
			Callback: func(s []byte) {
				logs.Debug("receive event:", string(s))
			},
		},
	)

	// send event
	GetEventMap().Get("event-1").Broadcast("hello kitty!")

	time.Sleep(2 * time.Second)

	// send event
	GetEventMap().Get("event-1").Broadcast("hello tom!")

	time.Sleep(2 * time.Second)
}
