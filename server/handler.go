package server

import (
	"github.com/cincout/mssfcore/registry"
)

type Handler interface {
	Name() string
	Handler() interface{}
	Endpoints() []*registry.Endpoint
}

type Subscriber interface {
	Topic() string
	Subscriber() interface{}
	Endpoints() []*registry.Endpoint
}
