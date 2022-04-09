package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	emitter "github.com/emitter-io/go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/backend/api"
	"github.com/openrfsense/backend/docs"
	"github.com/openrfsense/backend/mqtt"
	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
)

type Backend struct {
	Port  int
	Users map[string]string
}

type MQTT struct {
	Protocol string
	Host     string
	Port     int
	Secret   string
}

type BackendConfig struct {
	Backend
	MQTT
}

// FIXME: move elsewhere
var DefaultConfig = BackendConfig{
	Backend: Backend{
		Port: 8080,
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
)

// @title                      OpenRFSense backend API
// @version                    0.1
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

	log.Println("Loading config")
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	docs.SwaggerInfo.Host = config.GetOrDefault("backend.host", "localhost")

	log.Println("Starting keystore")
	keystore.Init(mqtt.NewBrokerRetriever(), mqtt.DefaultTTL)

	log.Println("Connecting to MQTT")
	err = mqtt.InitClient()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: remove these
	mqtt.Client.OnPresence(func(_ *emitter.Client, ev emitter.PresenceEvent) {
		log.Printf("[emitter] -> [B] %d subscriber(s) at topic '%s': %v\n", len(ev.Who), ev.Channel, ev.Who)
	})
	// mqtt.Presence("sensors/", true, false)
	mqtt.Subscribe("sensors/+/output/", func(_ *emitter.Client, msg emitter.Message) {
		log.Printf("[emitter] -> [B] received on specific handler: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	log.Println("Starting API")
	router := fiber.New()
	api.Use(router, "/api/v1")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{})

	go func() {
		<-c
		router.Shutdown()
		log.Println("Shutting down")
		shutdown <- struct{}{}
	}()

	addr := fmt.Sprintf(":%d", config.Get[int]("backend.port"))
	if err := router.Listen(addr); err != nil {
		log.Fatal(err)
	}

	<-shutdown
}
