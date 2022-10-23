package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/backend/api"
	"github.com/openrfsense/backend/docs"
	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/backend/ui"
	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/logging"
)

var (
	version = ""
	commit  = ""
	date    = ""

	log = logging.New().
		WithPrefix("main").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// @title                     OpenRFSense backend API
// @description               OpenRFSense backend API
// @contact.name              OpenRFSense
// @contact.url               https://github.com/openrfsense/backend/issues
// @license.name              AGPLv3
// @license.url               https://spdx.org/licenses/AGPL-3.0-or-later.html
// @BasePath                  /api/v1
// @securityDefinitions.basic BasicAuth
func main() {
	configPath := pflag.StringP("config", "c", "config.yml", "path to yaml config file")
	showVersion := pflag.BoolP("version", "v", false, "shows program version and build info")
	pflag.Parse()

	if *showVersion {
		fmt.Printf("openrfsense-node v%s (%s) built on %s\n", version, commit, date)
		return
	}

	log.Info("Loading config")
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	docs.SwaggerInfo.Host = config.GetOrDefault("backend.host", "localhost")
	docs.SwaggerInfo.Version = version

	log.Info("Starting NATS server")
	err = nats.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer nats.Disconnect()

	log.Info("Starting API")
	router := api.Start("/api/v1", fiber.Config{
		AppName:               "openrfsense-backend",
		DisableStartupMessage: true,
		Views:                 ui.NewEngine(),
	})
	// Initialize UI (templated web pages)
	ui.Init(router)
	defer func() {
		err = router.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	log.Info("Shutting down")
}
