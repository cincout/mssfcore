package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cincout/mssfcore/logs"
	"github.com/pborman/uuid"
)

type Server interface {
	Config() options
	Init(...Option)
	Handle(Handler) error
	NewHandler(interface{}) Handler
	NewSubscriber(string, interface{}) Subscriber
	Subscribe(Subscriber) error
	Register() error
	Deregister() error
	Start() error
	Stop() error
}

type Option func(*options)

var (
	DefaultAddress        = ":0"
	DefaultName           = "go-server"
	DefaultVersion        = "1.0.0"
	DefaultId             = uuid.NewUUID().String()
	DefaultServer  Server = newRpcServer()
)

func Config() options {
	return DefaultServer.Config()
}

func Init(opt ...Option) {
	if DefaultServer == nil {
		DefaultServer = newRpcServer(opt...)
	}
	DefaultServer.Init(opt...)
}

func NewServer(opt ...Option) Server {
	return newRpcServer(opt...)
}

func NewSubscriber(topic string, h interface{}) Subscriber {
	return DefaultServer.NewSubscriber(topic, h)
}

func NewHandler(h interface{}) Handler {
	return DefaultServer.NewHandler(h)
}

func Handle(h Handler) error {
	return DefaultServer.Handle(h)
}

func Subscribe(s Subscriber) error {
	return DefaultServer.Subscribe(s)
}

func Register() error {
	return DefaultServer.Register()
}

func Deregister() error {
	return DefaultServer.Deregister()
}

func Run() error {
	if err := Start(); err != nil {
		return err
	}

	if err := DefaultServer.Register(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	logs.DefaultLog.Informational("Received signal %s", <-ch)

	if err := DefaultServer.Deregister(); err != nil {
		return err
	}

	return Stop()
}

func Start() error {
	config := DefaultServer.Config()
	logs.DefaultLog.Informational("Starting server %s id %s", config.Name(), config.Id())
	return DefaultServer.Start()
}

func Stop() error {
	logs.DefaultLog.Informational("Stopping server")
	return DefaultServer.Stop()
}
