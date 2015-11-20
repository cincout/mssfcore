package driver

//
// All credit to Mondo
//

import (
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	DefaultExchange  = "micro"
	DefaultRabbitURL = "amqp://guest:guest@127.0.0.1:5672"
)

type rabbitMQConn struct {
	Connection      *amqp.Connection
	Channel         *rabbitMQChannel
	ExchangeChannel *rabbitMQChannel
	notify          chan bool
	exchange        string
	url             string

	connected bool

	mtx    sync.Mutex
	close  chan bool
	closed bool
}

func newRabbitMQConn(exchange string, urls []string) *rabbitMQConn {
	var url string

	if len(urls) > 0 && strings.HasPrefix(urls[0], "amqp://") {
		url = urls[0]
	} else {
		url = DefaultRabbitURL
	}

	if len(exchange) == 0 {
		exchange = DefaultExchange
	}

	return &rabbitMQConn{
		exchange: exchange,
		url:      url,
		notify:   make(chan bool, 1),
		close:    make(chan bool),
	}
}

func (r *rabbitMQConn) Init() chan bool {
	go r.Connect(r.notify)
	return r.notify
}

func (r *rabbitMQConn) Connect(connected chan bool) {
	for {
		if err := r.tryToConnect(); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		connected <- true
		r.connected = true
		notifyClose := make(chan *amqp.Error)
		r.Connection.NotifyClose(notifyClose)

		// Block until we get disconnected, or shut down
		select {
		case <-notifyClose:
			// Spin around and reconnect
			r.connected = false
		case <-r.close:
			// Shut down connection
			if err := r.Connection.Close(); err != nil {
			}
			r.connected = false
			return
		}
	}
}

func (r *rabbitMQConn) IsConnected() bool {
	return r.connected
}

func (r *rabbitMQConn) Close() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.closed {
		return
	}

	close(r.close)
	r.closed = true
}

func (r *rabbitMQConn) tryToConnect() error {
	var err error
	r.Connection, err = amqp.Dial(r.url)
	if err != nil {
		return err
	}
	r.Channel, err = newRabbitChannel(r.Connection)
	if err != nil {
		return err
	}
	r.Channel.DeclareExchange(r.exchange)
	r.ExchangeChannel, err = newRabbitChannel(r.Connection)
	if err != nil {
		return err
	}
	return nil
}

func (r *rabbitMQConn) Consume(queue string) (*rabbitMQChannel, <-chan amqp.Delivery, error) {
	consumerChannel, err := newRabbitChannel(r.Connection)
	if err != nil {
		return nil, nil, err
	}

	err = consumerChannel.DeclareQueue(queue)
	if err != nil {
		return nil, nil, err
	}

	deliveries, err := consumerChannel.ConsumeQueue(queue)
	if err != nil {
		return nil, nil, err
	}

	err = consumerChannel.BindQueue(queue, r.exchange)
	if err != nil {
		return nil, nil, err
	}

	return consumerChannel, deliveries, nil
}

func (r *rabbitMQConn) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.ExchangeChannel.Publish(exchange, key, msg)
}
