package helper

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/go-nats-streaming"
	"testing"
	"time"
)

func TestNATSStreamingProducer_Start(t *testing.T) {
	var (
		URL       = "nats://10.200.204.104:4222,nats://10.200.204.133:4222,nats://10.200.204.48:4222"
		group     = "oceanus"
		topic     = "nautilus-test"
		clusterId = "yonghui"
	)

	logs.Info("hahahahha..........")

	consumer := &Consumer{
		Group:       group,
		Topic:       topic,
		URL:         URL,
		ClusterId:   clusterId,
		ClientId:    "consumer-5",
		DurableName: "mydurableName",
		CallBack: func(m *stan.Msg) {
			logs.Info("consumer-1 receive msg:", (string)(m.Data))
			return
		},
	}

	consumer_b := &Consumer{
		Group:       group,
		Topic:       topic,
		URL:         URL,
		ClusterId:   clusterId,
		ClientId:    "consumer-6",
		DurableName: "mydurableName",
		CallBack: func(m *stan.Msg) {
			logs.Info("consumer-2 receive msg:", (string)(m.Data))
			return
		},
	}

	producer := &Producer{
		URL:       URL,
		ClusterId: clusterId,
		ClientId:  "producer-5",
	}

	if err := producer.Start(); err != nil {
		t.Error(err)
	}

	if err := consumer.RunWithSubscribe(); err != nil {
		t.Error(err)
	}

	if err := consumer_b.RunWithSubscribe(); err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)

		msg1 := fmt.Sprintf("hello world! %d", i)

		logs.Info("begin to send msg1: %d", i)
		if err := producer.Send(topic, ([]byte)(msg1)); err != nil {
			logs.Error(err)
		}
		//
		//msg2 := fmt.Sprintf("goodbye world! %d", i)
		//logs.Info("begin to send msg2: %d", i)
		//if err := producer.SendAsync(topic, ([]byte)(msg2)); err != nil {
		//	logs.Error(err)
		//}
	}

	time.Sleep(1 * time.Second)
	consumer.Close()
	consumer_b.Close()

	logs.Info("begin to send msg: %d", "msg1")
	if err := producer.Send(topic, ([]byte)("msg1")); err != nil {
		logs.Error(err)
	}

	logs.Info("begin to send msg: %d", "msg2")
	if err := producer.Send(topic, ([]byte)("msg2")); err != nil {
		logs.Error(err)
	}

	producer.Shutdown()
}
