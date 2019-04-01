package helper

import (
	"github.com/docker/distribution/uuid"
	"github.com/sevenNt/rocketmq"
	logs "github.com/sirupsen/logrus"
)

///////////////////////////////////
// consumer
//////////////////////////////////

type RocketMqConsumer struct {
	Group      string
	Topic      string
	NameServer string
	CallBack   func([]*rocketmq.MessageExt) error
	worker     rocketmq.Consumer
}

func (c *RocketMqConsumer) Run() error {
	id := uuid.Generate().String()

	conf := &rocketmq.Config{
		Namesrv:      c.NameServer,
		InstanceName: "consumer-" + id,
	}

	consumer, err := rocketmq.NewDefaultConsumer(c.Group, conf)
	if err != nil {
		logs.Error(err)
		return err
	}
	c.worker = consumer
	c.worker.Subscribe(c.Topic, "*")
	c.worker.RegisterMessageListener(c.CallBack)
	c.worker.Start()

	return nil
}

func (c *RocketMqConsumer) Shutdown() {
	c.worker.UnSubscribe(c.Topic)
	c.worker.Shutdown()
}

///////////////////////////////////
// producer
//////////////////////////////////

type RocketMqProducer struct {
	Group      string
	NameServer string
	worker     rocketmq.Producer
}

func (p *RocketMqProducer) Start() error {
	id := uuid.Generate().String()

	conf := &rocketmq.Config{
		Namesrv:      p.NameServer,
		InstanceName: "procedure-" + id,
	}

	producer, err := rocketmq.NewDefaultProducer(p.Group, conf)
	producer.Start()
	if err != nil {
		logs.Error(err)
		return err
	}
	p.worker = producer

	return nil
}

func (p *RocketMqProducer) Send(topic string, body []byte) error {
	msg := rocketmq.NewMessage(topic, body)

	_, err := p.worker.Send(msg)
	if err != nil {
		logs.Error(err)
		return err
	}

	//utils.Display("result:", result)

	//logs.Info("send message success!")
	return nil
}

func (p *RocketMqProducer) Shutdown() {
	p.worker.Shutdown()
}
