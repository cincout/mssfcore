package broker

import(
	
	"fmt"
)

type brokerType func(addrs []string, opt ...Option) Broker

type Broker interface {
	Address() string
	Connect() error
	Disconnect() error
	Init() error
	Publish(string, *Message) error
	Subscribe(string, Handler) (Subscriber, error)
}

type Handler func(*Message)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Subscriber interface {
	Topic() string
	Unsubscribe() error
}

type options struct{}

type Option func(*options)

var (
	DefaultBroker Broker
	adapters = make(map[string]brokerType)
)

func Register(name string, broker brokerType) {
	if broker == nil {
		panic("broker: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("broker: Register called twice for provider " + name)
	}
	adapters[name] = broker
}

func NewBroker(adaptername string,addrs []string, opt ...Option) (Broker,error) {
	if broker, ok := adapters[adaptername]; ok {
		newBroker := broker(addrs, opt...)
		return newBroker,nil
	} else {
		return nil,fmt.Errorf("broker: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func Init() error {
	return DefaultBroker.Init()
}

func Connect() error {
	return DefaultBroker.Connect()
}

func Disconnect() error {
	return DefaultBroker.Disconnect()
}

func Publish(topic string, msg *Message) error {
	return DefaultBroker.Publish(topic, msg)
}

func Subscribe(topic string, handler Handler) (Subscriber, error) {
	return DefaultBroker.Subscribe(topic, handler)
}