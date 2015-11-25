package driver

import (
	"sync"
	"github.com/streadway/amqp"
	"github.com/cincout/mssfcore/broker"
)

type rbroker struct {
	conn  []*rabbitMQConn
	addrs []string
	mpublicer chan *publicer
	once sync.Once
}

type publicer struct {
	topic string
	msg *broker.Message
}

type subscriber struct {
	topic string
	ch    []*rabbitMQChannel
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	//return s.ch.Close()
	for i := 0; i < len(s.ch); i++ {
		s.ch[i].Close()
	}
	return nil
}

func (r *rbroker) Publish(topic string, msg *broker.Message) error {

	r.once.Do(func() {
			go func() {
				for _,v := range  r.conn {
					go r.publish(v)
				}
			}()
	})
	
	mp := &publicer{
		topic: topic,
		msg: msg,
	}
	r.mpublicer <- mp
	return nil
}

func (r *rbroker) publish(conn  *rabbitMQConn) error {
	for mp := range r.mpublicer {
		
		m := amqp.Publishing{
			Body:    mp.msg.Body,
			Headers: amqp.Table{},
		}
	
		for k, v := range mp.msg.Header {
			m.Headers[k] = v
		}
		
		conn.Publish("", mp.topic, m)
	}
	return nil
}

func (r *rbroker) Subscribe(topic string, handler broker.Handler) (broker.Subscriber, error) {
	
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
	
	chs := make([]*rabbitMQChannel,2)
	
	for i := 0; i < len(r.conn); i++ {
		ch, sub, err := r.conn[i].Consume(topic)
		if err != nil {
			return nil, err
		}
		chs[i] = ch
		go func() {
			for d := range sub {
				go fn(d)
			}
		}()
	}

	return &subscriber{ch: chs, topic: topic}, nil
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
	//<-r.conn.Init()
	for i := 0; i < len(r.conn); i++ {
		<-r.conn[i].Init()
	}
	return nil
}

func (r *rbroker) Disconnect() error {
	//r.conn.Close()
	for i := 0; i < len(r.conn); i++ {
		r.conn[i].Close()
	}
	close(r.mpublicer)
	return nil
}

func NewRabbitmqBroker(addrs []string, opt ...broker.Option) broker.Broker {
	conns := make([]*rabbitMQConn,2)
	for i := 0; i < len(conns); i++ {
		conns[i] = newRabbitMQConn("", addrs)
	}
	
	var once sync.Once
	
	return &rbroker{
		conn:  conns,
		addrs: addrs,
		mpublicer: make(chan *publicer, 3072),
		once: once,
	}
}

func init() {
	broker.Register("rabbitmq", NewRabbitmqBroker)
}
