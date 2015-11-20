// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mssfcore

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	
	"github.com/cincout/mssfcore/broker"
	broker_driver "github.com/cincout/mssfcore/broker/driver"
	"github.com/cincout/mssfcore/client"
	"github.com/cincout/mssfcore/registry"
	registry_driver "github.com/cincout/mssfcore/registry/driver"
	"github.com/cincout/mssfcore/server"
	"github.com/cincout/mssfcore/transport"
	transport_driver "github.com/cincout/mssfcore/transport/driver"
	
	"github.com/cincout/mssfcore/config"
	"github.com/cincout/mssfcore/logs"
	"github.com/cincout/mssfcore/utils"
)

var (
	ServerName            string //Name of the server
	ServerVersion         string //Version of the server
	ServerId              string //Id of the server
	ServerAddress         string //Bind address for the server
	ServerAdvertise       string //Used instead of the server_address when registering with discovery
	ServerMetadata 	      string //A list of key-value pairs defining metadata
	Broker                string //Broker for pub/sub. http, nats, rabbitmq
	BrokerAddress		  string
	Registry              string //Registry for discovery. memory, consul, etcd, kubernetes
	RegistryAddress		  string
	Transport             string //Transport mechanism used; http, rabbitmq, nats
	TransportAddress	  string
	
	AppConfig             *mssfAppConfig
	RunMode               string           // run mode, "dev" or "prod"
	AppConfigProvider     string
	AppPath               string
	workPath              string
	AppConfigPath         string
	
	EnableAdmin           bool   // flag of enable admin module to log every request info.
	AdminHttpAddr         string // http server configurations for admin module.
	AdminHttpPort         int
)

type mssfAppConfig struct {
	innerConfig config.ConfigContainer
}

func newAppConfig(AppConfigProvider string, AppConfigPath string) (*mssfAppConfig, error) {
	ac, err := config.NewConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return nil, err
	}
	rac := &mssfAppConfig{ac}
	return rac, nil
}

func (b *mssfAppConfig) Set(key, val string) error {
	err := b.innerConfig.Set(RunMode+"::"+key, val)
	if err == nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *mssfAppConfig) String(key string) string {
	v := b.innerConfig.String(RunMode + "::" + key)
	if v == "" {
		return b.innerConfig.String(key)
	}
	return v
}

func (b *mssfAppConfig) Strings(key string) []string {
	v := b.innerConfig.Strings(RunMode + "::" + key)
	if v[0] == "" {
		return b.innerConfig.Strings(key)
	}
	return v
}

func (b *mssfAppConfig) Int(key string) (int, error) {
	v, err := b.innerConfig.Int(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int(key)
	}
	return v, nil
}

func (b *mssfAppConfig) Int64(key string) (int64, error) {
	v, err := b.innerConfig.Int64(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int64(key)
	}
	return v, nil
}

func (b *mssfAppConfig) Bool(key string) (bool, error) {
	v, err := b.innerConfig.Bool(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Bool(key)
	}
	return v, nil
}

func (b *mssfAppConfig) Float(key string) (float64, error) {
	v, err := b.innerConfig.Float(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Float(key)
	}
	return v, nil
}

func (b *mssfAppConfig) DefaultString(key string, defaultval string) string {
	v := b.String(key)
	if v != "" {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DefaultStrings(key string, defaultval []string) []string {
	v := b.Strings(key)
	if len(v) != 0 {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DefaultInt(key string, defaultval int) int {
	v, err := b.Int(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DefaultInt64(key string, defaultval int64) int64 {
	v, err := b.Int64(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DefaultBool(key string, defaultval bool) bool {
	v, err := b.Bool(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DefaultFloat(key string, defaultval float64) float64 {
	v, err := b.Float(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *mssfAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *mssfAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *mssfAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}

func init() {
	workPath, _ = os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")

	if workPath != AppPath {
		if utils.FileExists(AppConfigPath) {
			os.Chdir(AppPath)
		} else {
			AppConfigPath = filepath.Join(workPath, "conf", "app.conf")
		}
	}

	AppConfigProvider = "ini"
	
	ServerName = "go.micro.srv.example"
	ServerVersion = "0.1"
	ServerId = ""
	ServerAddress = "127.0.0.1:8080"
	ServerAdvertise = "127.0.0.1:8080"
	ServerMetadata = "version=1.0.0"
	Broker = "http"
	BrokerAddress = ""
	Registry = "consul"
	RegistryAddress = ""
	Transport = "http"
	TransportAddress = ""


	EnableAdmin = false
	AdminHttpAddr = "127.0.0.1"
	AdminHttpPort = 8088

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init MssfLogger
	MssfLogger = logs.NewLogger(10000)
	logs.DefaultLog = MssfLogger;
	err := MssfLogger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	SetLogFuncCall(true)

	err = ParseConfig()
	if err != nil && os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		ac := config.NewFakeConfig()
		AppConfig = &mssfAppConfig{ac}
		Warning(err)
	}
	
	broker_driver.InitBrokerDriver()
	broker.DefaultBroker,err = broker.NewBroker(Broker,strings.Split(BrokerAddress, ","))
	if(err != nil) {
		fmt.Println("init broker error:", err)
	}
	
	registry_driver.InitRegistryDriver()
	registry.DefaultRegistry,err = registry.NewRegistry(Registry,strings.Split(RegistryAddress, ","))
	if(err != nil) {
		fmt.Println("init registry error:", err)
	}
	transport_driver.InitTransportDriver()
	transport.DefaultTransport,err = transport.NewTransport(Transport,strings.Split(TransportAddress, ","))
	if(err != nil) {
		fmt.Println("init transport error:", err)
	}
	
	metadata := make(map[string]string)
	parts := strings.Split(ServerMetadata, "=")
	var key,val string
	key = parts[0]
	if len(parts) > 1 {
		val = strings.Join(parts[1:], "=")
	}
	metadata[key] = val
	
	server.DefaultServer = server.NewServer(
		server.Name(ServerName),
		server.Version(ServerVersion),
		server.Id(ServerId),
		server.Address(ServerAddress),
		server.Advertise(ServerAdvertise),
		server.Metadata(metadata),
	)

	client.DefaultClient = client.NewClient()
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {
	AppConfig, err = newAppConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return err
	}
	
	if runmode := AppConfig.String("RunMode"); runmode != "" {
		RunMode = runmode
	}

	if server_name := AppConfig.String("server_name"); server_name != "" {
		ServerName = server_name
	}

	if server_version := AppConfig.String("server_version"); server_version != "" {
		ServerVersion = server_version
	}

	if server_id := AppConfig.String("server_id"); server_id != "" {
		ServerId = server_id
	}

	if server_address := AppConfig.String("server_address"); server_address != "" {
		ServerAddress = server_address
	}

	if server_advertise := AppConfig.String("server_advertise"); server_advertise != "" {
		ServerAdvertise = server_advertise
	}
	
	if tmpBroker := AppConfig.String("broker"); tmpBroker != "" {
		Broker = tmpBroker
		if broker_address := AppConfig.String("broker_address"); broker_address != "" {
			BrokerAddress = broker_address
		}
	}
	
	if tmpRegistry := AppConfig.String("registry"); tmpRegistry != "" {
		Registry = tmpRegistry
		if registry_address := AppConfig.String("registry_address"); registry_address != "" {
			RegistryAddress = registry_address
		}
	}
	
	if tmpTransport := AppConfig.String("transport"); tmpTransport != "" {
		Transport = tmpTransport
		if transport_address := AppConfig.String("transport_address"); transport_address != "" {
			TransportAddress = transport_address
		}
	}

	if enableadmin, err := AppConfig.Bool("EnableAdmin"); err == nil {
		EnableAdmin = enableadmin
	}

	if adminhttpaddr := AppConfig.String("AdminHttpAddr"); adminhttpaddr != "" {
		AdminHttpAddr = adminhttpaddr
	}

	if adminhttpport, err := AppConfig.Int("AdminHttpPort"); err == nil {
		AdminHttpPort = adminhttpport
	}
	return nil
}
