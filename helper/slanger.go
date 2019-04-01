package helper

import (
	"github.com/pusher/pusher-http-go"
)

const (
	SlangerEventTypeOperationStatus = "operationStatus"
)

type WsServer struct {
	AppId  string
	Key    string
	Secret string
	Host   string
	client *pusher.Client
}

func (ws *WsServer) GetClient() {
	ws.client = &pusher.Client{
		AppId:  ws.AppId,
		Key:    ws.Key,
		Secret: ws.Secret,
		Host:   ws.Host,
		Secure: false,
	}
}

func (ws *WsServer) PushOperationStatusToUI(channels []string, data interface{}) {
	for _, channel := range channels {
		ws.client.Trigger(channel, SlangerEventTypeOperationStatus, data)
	}
}

func (ws *WsServer) PushMsgToUI(channel string, eventName string, data interface{}) {
	ws.client.Trigger(channel, eventName, data)
}
