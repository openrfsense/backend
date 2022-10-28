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
	"github.com/openrfsense/backend/config"
	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/backend/docs"
	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/backend/ui"
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
	konfig, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connecting to database")
	err = database.Init(konfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting NATS server")
	err = nats.Start(konfig)
	if err != nil {
		log.Fatal(err)
	}
	defer nats.Disconnect()

	log.Info("Starting API")
	docs.SwaggerInfo.Host = konfig.String("backend.host")
	docs.SwaggerInfo.Version = version
	router := api.Start(konfig, "/api/v1", fiber.Config{
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
