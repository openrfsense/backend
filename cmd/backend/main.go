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
	"github.com/openrfsense/backend/samples"
	"github.com/openrfsense/backend/ui"
	"github.com/openrfsense/common/logging"

	_ "net/http/pprof"
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

func main() {
	configPath := pflag.StringP("config", "c", "config.yml", "path to yaml config file")
	showVersion := pflag.BoolP("version", "v", false, "shows program version and build info")
	pflag.Parse()

	if *showVersion {
		fmt.Printf("openrfsense-node v%s (%s) built on %s\n", version, commit, date)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

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

	router := fiber.New(fiber.Config{
		AppName:               "openrfsense-backend",
		DisableStartupMessage: true,
		ErrorHandler:          ui.ErrorHandler,
		Views:                 ui.NewEngine(),
	})

	log.Info("Starting API")
	docs.SwaggerInfo.Host = konfig.String("backend.host")
	docs.SwaggerInfo.Version = version
	api.Init(konfig, router, "/api/v1")

	// Initialize UI (templated web pages)
	log.Info("Starting UI")
	ui.Init(router)

	log.Info("Starting measurement collector")
	err = samples.StartCollector(ctx, konfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting WebSocket streaming handler")
	samples.StartWebsocket(ctx, konfig, router)

	addr := fmt.Sprintf(":%d", konfig.MustInt("backend.port"))
	go func() {
		if err := router.Listen(addr); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Info("Shutting down")
	_ = router.Shutdown()
	nats.Disconnect()
	database.Close()
}
