package helper

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/go-nats"
	"testing"
	"time"
)

func TestNATSProducer_Start(t *testing.T) {
	var (
		URL   = "nats://10.216.155.34:4222,nats://10.216.155.35:4222,nats://10.216.155.36:4222"
		topic = "nautilus-test"
	)

	logs.Info("hahahahha..........")

	consumer := &NATSConsumer{
		Topic: topic,
		URL:   URL,
		CallBack: func(m *nats.Msg) {
			logs.Info("receive msg:", (string)(m.Data))
			return
		},
	}

	producer := &NATSProducer{
		URL: URL,
	}

	if err := producer.Start(); err != nil {
		t.Error(err)
	}

	if err := consumer.Run(); err != nil {
		t.Error(err)
	}

	for i := 0; i < 2000; i++ {
		time.Sleep(1 * time.Second)

		msg1 := fmt.Sprintf("hello world! %d", i)

		if err := producer.Send(topic, ([]byte)(msg1)); err != nil {
			t.Error(err)
		}

		msg2 := fmt.Sprintf("goodbye world! %d", i)

		if err := producer.Send(topic, ([]byte)(msg2)); err != nil {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
	producer.Shutdown()
	consumer.Shutdown()
}
