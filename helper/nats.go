package helper

import (
	"github.com/nats-io/go-nats"
	logs "github.com/sirupsen/logrus"
)

///////////////////////////////////
// consumer
//////////////////////////////////

type NATSConsumer struct {
	Group    string // no used for NATS
	Topic    string
	URL      string
	CallBack func(*nats.Msg)
	Conn     *nats.Conn
}

func (c *NATSConsumer) Run() error {
	nc, err := nats.Connect(c.URL)
	if err != nil {
		logs.Error(err)
		return err
	}
	c.Conn = nc

	_, err = nc.Subscribe(c.Topic, c.CallBack)
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c *NATSConsumer) Shutdown() {
	if c.Conn.IsClosed() {
		return
	}

	c.Conn.Close()
}

///////////////////////////////////
// producer
//////////////////////////////////

type NATSProducer struct {
	Group string // no used for NATS
	URL   string
	Conn  *nats.Conn
}

func (p *NATSProducer) Start() error {
	nc, err := nats.Connect(p.URL)
	if err != nil {
		logs.Error(err)
		return err
	}
	p.Conn = nc

	return nil
}

func (p *NATSProducer) Send(topic string, body []byte) error {
	if err := p.Conn.Publish(topic, body); err != nil {
		logs.Error(err)
		return err
	}

	logs.Debug("send message success!", topic)
	return nil
}

func (p *NATSProducer) Shutdown() {
	if p.Conn.IsClosed() {
		return
	}

	p.Conn.Close()
}
