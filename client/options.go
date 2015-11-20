package client

import (
	"github.com/cincout/mssfcore/broker"
	"github.com/cincout/mssfcore/registry"
	"github.com/cincout/mssfcore/transport"
)

type options struct {
	broker    broker.Broker
	registry  registry.Registry
	transport transport.Transport
}

func Broker(b broker.Broker) Option {
	return func(o *options) {
		o.broker = b
	}
}

func Registry(r registry.Registry) Option {
	return func(o *options) {
		o.registry = r
	}
}

func Transport(t transport.Transport) Option {
	return func(o *options) {
		o.transport = t
	}
}
