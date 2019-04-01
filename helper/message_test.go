package helper

import (
	"github.com/sevenNt/rocketmq"
	"testing"
	"time"
)

func TestProducer_Start(t *testing.T) {
	var (
		group   = "oceanus"
		nameSrv = "10.1.54.234:9876"
		topic   = "nautilus-test"
	)

	consumer := &RocketMqConsumer{
		Group:      group,
		Topic:      topic,
		NameServer: nameSrv,
		CallBack: func(msgs []*rocketmq.MessageExt) error {
			for _, msg := range msgs {
				t.Log("receive msg:", (string)(msg.Body))
			}
			return nil
		},
	}

	producer := &RocketMqProducer{
		Group:      group,
		NameServer: nameSrv,
	}

	if err := producer.Start(); err != nil {
		t.Error(err)
	}

	if err := consumer.Run(); err != nil {
		t.Error(err)
	}

	time.Sleep(6 * time.Second)

	if err := producer.Send(topic, ([]byte)("hello world!")); err != nil {
		t.Error(err)
	}

	if err := producer.Send(topic, ([]byte)("goodbye world!")); err != nil {
		t.Error(err)
	}

	time.Sleep(10 * time.Second)
	producer.Shutdown()
	consumer.Shutdown()
}
