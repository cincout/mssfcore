package registry

import(
	"fmt"
)

type registryType func(addrs []string, opt ...Option) Registry

type Registry interface {
	Register(*Service) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch() (Watcher, error)
}

type Watcher interface {
	Stop()
}

type options struct{}

type Option func(*options)

var (
	DefaultRegistry Registry
	adapters = make(map[string]registryType)
)

func RegisterRegistry(name string, registry registryType) {
	if registry == nil {
		panic("registry: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("registry: Register called twice for provider " + name)
	}
	adapters[name] = registry
}


func NewRegistry(adaptername string,addrs []string, opt ...Option) (Registry,error) {
	if registry, ok := adapters[adaptername]; ok {
		newObjRegistry := registry(addrs, opt...)
		return newObjRegistry,nil
	} else {
		return nil,fmt.Errorf("registry: unknown adaptername %q (forgotten Register?)", adaptername)
	}
	
}

func Register(s *Service) error {
	return DefaultRegistry.Register(s)
}

func Deregister(s *Service) error {
	return DefaultRegistry.Deregister(s)
}

func GetService(name string) ([]*Service, error) {
	return DefaultRegistry.GetService(name)
}

func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}
