package transport

import(
	"fmt"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
}

type Client interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

type transportType func(addrs []string, opt ...Option) Transport

type Transport interface {
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string) (Listener, error)
}


type Option func(*options)

type options struct{}

type DialOptions struct {
	Stream bool
}

type DialOption func(*DialOptions)

var (
	DefaultTransport Transport
	adapters = make(map[string]transportType)
)

func WithStream() DialOption {
	return func(o *DialOptions) {
		o.Stream = true
	}
}

func Register(name string, transport transportType) {
	if transport == nil {
		panic("transport: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("transport: Register called twice for provider " + name)
	}
	adapters[name] = transport
}

func NewTransport(adaptername string,addrs []string, opt ...Option) (Transport,error) {
	if transport, ok := adapters[adaptername]; ok {
		newTransport := transport(addrs, opt...)
		return newTransport,nil
	} else {
		return nil,fmt.Errorf("transport: unknown adaptername %q (forgotten Register?)", adaptername)
	}
	
}


func Dial(addr string, opts ...DialOption) (Client, error) {
	return DefaultTransport.Dial(addr, opts...)
}

func Listen(addr string) (Listener, error) {
	return DefaultTransport.Listen(addr)
}
