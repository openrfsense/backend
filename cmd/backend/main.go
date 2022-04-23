package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	emitter "github.com/emitter-io/go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/backend/api"
	"github.com/openrfsense/backend/docs"
	"github.com/openrfsense/backend/mqtt"
	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
)

type Backend struct {
	Port  int               `koanf:"port"`
	Users map[string]string `koanf:"users"`
}

type MQTT struct {
	Protocol string `koanf:"protocol"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Secret   string `koanf:"secret"`
}

type BackendConfig struct {
	Backend `koanf:"backend"`
	MQTT    `koanf:"mqtt"`
}

// FIXME: move elsewhere
var DefaultConfig = BackendConfig{
	Backend: Backend{
		Port: 8081,
	},
	MQTT: MQTT{
		Protocol: "tcp",
		Port:     8080,
	},
}

var (
	version = ""
	commit  = ""
	date    = ""

	log = logging.New().
		WithPrefix("main").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// @title                      OpenRFSense backend API
// @description                OpenRFSense backend API
// @contact.name               OpenRFSense
// @contact.url                https://github.com/openrfsense/backend/issues
// @license.name               AGPLv3
// @license.url                https://spdx.org/licenses/AGPL-3.0-or-later.html
// @BasePath                   /api/v1
// @securityDefinitions.basic  BasicAuth
func main() {
	configPath := pflag.StringP("config", "c", "config.yml", "path to yaml config file")
	pflag.Parse()

	log.Info("Loading config")
	conf := &BackendConfig{}
	err := config.Load(*configPath, DefaultConfig, &conf)
	if err != nil {
		log.Fatal(err)
	}

	docs.SwaggerInfo.Host = config.GetOrDefault("backend.host", "localhost")
	docs.SwaggerInfo.Version = version

	log.Info("Starting keystore")
	keystore.Init(mqtt.NewBrokerRetriever(), mqtt.DefaultTTL)

	log.Info("Connecting to MQTT")
	mqtt.Init()

	// TODO: remove these
	mqtt.Client().OnPresence(func(_ *emitter.Client, ev emitter.PresenceEvent) {
		log.Debugf("[emitter] -> [B] %d subscriber(s) at topic '%s': %v\n", len(ev.Who), ev.Channel, ev.Who)
	})
	// mqtt.Presence("node/", true, false)
	// mqtt.Subscribe("node/+/output/", func(_ *emitter.Client, msg emitter.Message) {
	// 	log.Debugf("[emitter] -> [B] received on specific handler: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	// })
	mqtt.Subscribe("node/all/", func(_ *emitter.Client, msg emitter.Message) {
		log.Debugf("[emitter] -> [B] received on specific handler: '%s' topic: '%s'", msg.Payload(), msg.Topic())
	})

	log.Info("Starting API")
	router := fiber.New(fiber.Config{
		AppName:               "openrfsense-backend",
		DisableStartupMessage: true,
	})
	api.Use(router, "/api/v1")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{})

	go func() {
		<-c
		router.Shutdown()
		mqtt.Disconnect(time.Second)
		log.Info("Shutting down")
		shutdown <- struct{}{}
	}()

	addr := fmt.Sprintf(":%d", config.Get[int]("backend.port"))
	if err := router.Listen(addr); err != nil {
		log.Fatal(err)
	}

	<-shutdown
}
