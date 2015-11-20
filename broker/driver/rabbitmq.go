package driver

import (
	"github.com/streadway/amqp"
	"github.com/cincout/mssfcore/broker"
)

type rbroker struct {
	conn  *rabbitMQConn
	addrs []string
}

type subscriber struct {
	topic string
	ch    *rabbitMQChannel
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	return s.ch.Close()
}

func (r *rbroker) Publish(topic string, msg *broker.Message) error {
	m := amqp.Publishing{
		Body:    msg.Body,
		Headers: amqp.Table{},
	}

	for k, v := range msg.Header {
		m.Headers[k] = v
	}

	return r.conn.Publish("", topic, m)
}

func (r *rbroker) Subscribe(topic string, handler broker.Handler) (broker.Subscriber, error) {
	ch, sub, err := r.conn.Consume(topic)
	if err != nil {
		return nil, err
	}

	fn := func(msg amqp.Delivery) {
		header := make(map[string]string)
		for k, v := range msg.Headers {
			header[k], _ = v.(string)
		}
		handler(&broker.Message{
			Header: header,
			Body:   msg.Body,
		})
	}

	go func() {
		for d := range sub {
			go fn(d)
		}
	}()

	return &subscriber{ch: ch, topic: topic}, nil
}

func (r *rbroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return ""
}

func (r *rbroker) Init() error {
	return nil
}

func (r *rbroker) Connect() error {
	<-r.conn.Init()
	return nil
}

func (r *rbroker) Disconnect() error {
	r.conn.Close()
	return nil
}

func NewRabbitmqBroker(addrs []string, opt ...broker.Option) broker.Broker {
	return &rbroker{
		conn:  newRabbitMQConn("", addrs),
		addrs: addrs,
	}
}

func init() {
	broker.Register("rabbitmq", NewRabbitmqBroker)
}
