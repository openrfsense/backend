package ui

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"

	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("ui").
	WithFlags(logging.FlagsDevelopment)

//go:embed views/*
var viewsFs embed.FS

//go:embed static/*
var staticFs embed.FS

// Initializes Fiber view engine with embedded HTML templates from views/
func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFs), ".html")
	engine.Reload(true)

	return engine
}

// Configure a router and use a logger for the UI. Initializes routes and view models.
func Init(router *fiber.App) {
	router.Use(
		"/static",
		filesystem.New(filesystem.Config{
			Root:       http.FS(staticFs),
			PathPrefix: "static",
			Browse:     true,
		}),
		compress.New(),
	)

	router.Get("/", renderIndex)
	router.Get("/sensor/:id", renderSensorPage)
}

func renderIndex(c *fiber.Ctx) error {
	sensorStats, err := fetchAllSensorStats()
	if err != nil {
		return err
	}

	return c.Render("views/index", fiber.Map{
		"sensors": sensorStats,
	})
}

func renderSensorPage(c *fiber.Ctx) error {
	stat, err := fetchSensorStats(c.Params("id"))
	if err != nil {
		return err
	}

	return c.Render("views/sensor", stat)
}
