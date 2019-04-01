package helper

import (
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	logs "github.com/sirupsen/logrus"
)

///////////////////////////////////
// consumer
//////////////////////////////////

type Consumer struct {
	Group       string
	Topic       string
	URL         string
	ClusterId   string
	ClientId    string
	DurableName string
	CallBack    func(*stan.Msg)
	NConn       *nats.Conn
	SConn       stan.Conn
	sSub        stan.Subscription
}

func (c *Consumer) Run() error {
	nc, err := nats.Connect(c.URL, nats.MaxReconnects(-1))
	if err != nil {
		logs.Error(err)
		return err
	}

	c.NConn = nc
	sc, err := stan.Connect(
		c.ClusterId,
		c.ClientId,
		stan.NatsConn(nc),
		stan.Pings(10, 360), // retain this connection during 1 hour
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			logs.Error("stan connection lost, reason: ", reason)
			panic("NATS-Streaming is not available, after 10 * 1000")
		}))
	if err != nil {
		logs.Error(err)
		return err
	}
	c.SConn = sc
	c.sSub, err = sc.QueueSubscribe(c.Topic, c.Group, c.CallBack, stan.DurableName(c.DurableName))
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c *Consumer) RunWithSubscribe() error {
	nc, err := nats.Connect(c.URL, nats.MaxReconnects(-1))
	if err != nil {
		logs.Error(err)
		return err
	}

	c.NConn = nc
	sc, err := stan.Connect(
		c.ClusterId,
		c.ClientId,
		stan.NatsConn(nc),
		stan.Pings(10, 360), // retain this connection during 1 hour
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			logs.Error("stan connection lost, reason: ", reason)
			panic("NATS-Streaming is not available, after 10 * 1000")
		}))
	if err != nil {
		logs.Error(err)
		return err
	}
	c.SConn = sc
	c.sSub, err = sc.Subscribe(c.Topic, c.CallBack, stan.DurableName(c.DurableName))
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c *Consumer) Shutdown() {
	if err := c.sSub.Unsubscribe(); err != nil {
		logs.Error(err)
	}

	if err := c.SConn.Close(); err != nil {
		logs.Error(err)
	}
}

func (c *Consumer) Close() {
	if err := c.SConn.Close(); err != nil {
		logs.Error(err)
		return
	}
}

///////////////////////////////////
// producer
//////////////////////////////////

type Producer struct {
	Group     string // no used for NATS
	URL       string
	ClusterId string
	ClientId  string
	NConn     *nats.Conn
	SConn     stan.Conn
}

func (p *Producer) Start() error {
	//nc, err := nats.Connect(p.URL, nats.MaxReconnects(-1))
	//if err != nil {
	//	logs.Error(err)
	//	return err
	//}
	//
	//p.NConn = nc
	//sc, err := stan.Connect(
	//	p.ClusterId,
	//	p.ClientId,
	//	stan.NatsConn(nc),
	//	stan.Pings(10, 360), // retain this connection during 1 hour
	//	stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
	//		logs.Error("stan connection lost, reason: ", reason)
	//		panic("NATS-Streaming is not available, after 10 * 1000")
	//	}))
	//
	//p.SConn = sc
	return nil
}

func (p *Producer) Send(topic string, body []byte) error {
	nc, err := nats.Connect(p.URL, nats.MaxReconnects(-1))
	if err != nil {
		logs.Error(err)
		return err
	}

	p.NConn = nc
	sc, err := stan.Connect(p.ClusterId, p.ClientId, stan.NatsConn(nc))
	if err != nil {
		logs.Error(err)
		nc.Close()
		return err
	}

	defer sc.Close()
	p.SConn = sc

	if err := p.SConn.Publish(topic, body); err != nil {
		logs.Error(err)
		return err
	}

	logs.Debug("send message success!", topic)
	return nil
}

func (p *Producer) SendAsync(topic string, body []byte) error {
	nc, err := nats.Connect(p.URL, nats.MaxReconnects(-1))
	if err != nil {
		logs.Error(err)
		return err
	}

	p.NConn = nc
	sc, err := stan.Connect(p.ClusterId, p.ClientId, stan.NatsConn(nc))
	if err != nil {
		logs.Error(err)
		nc.Close()
		return err
	}

	defer sc.Close()
	p.SConn = sc

	ackHandler := func(ackedNuid string, err error) {
		if err != nil {
			logs.Error("Warning: error publishing msg id ", ackedNuid, err.Error())
		} else {
			logs.Error("Received ack for msg id ", ackedNuid)
		}
	}

	if nuid, err := p.SConn.PublishAsync(topic, body, ackHandler); err != nil {
		logs.Error("Error publishing msg ", nuid, err.Error())
		return err
	}

	logs.Debug("send message success!", topic)
	return nil
}

func (p *Producer) Shutdown() {
	p.SConn.Close()
}
