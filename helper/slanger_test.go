package helper

import (
	"testing"
)

func TestPushEvent(t *testing.T) {
	wsServer := WsServer{
		AppId:  "oceanus",
		Key:    "yhcloud",
		Secret: "yh2016",
		Host:   "10.0.91.169:4567",
	}
	wsServer.GetClient()
	wsServer.PushOperationStatusToUI(append([]string{}, "stack"), "create stack success!")
}
