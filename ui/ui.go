package ui

import (
	"embed"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("ui").
	WithFlags(logging.FlagsDevelopment)

//go:embed views/*
var viewsFs embed.FS

//go:embed static/*
var staticFs embed.FS

// Initializes Fiber view engine with embedded HTML templates from views/.
func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFs), ".html")
	engine.AddFunc("title", cases.Title(language.AmericanEnglish).String)
	engine.AddFunc("percent", func(this float64, total float64) int {
		part := this / total
		return int(part * 100)
	})
	engine.AddFunc("humanizeDuration", func(span time.Duration) string {
		now := time.Now()
		return humanize.RelTime(now.Add(-span), now, "", "")
	})
	engine.AddFunc("humanizeSizeKB", func(kb float64) string {
		return humanize.Bytes(uint64(kb) * 1024)
	})
	engine.AddFunc("humanizeSizeB", func(b float64) string {
		return humanize.Bytes(uint64(b))
	})

	return engine
}

// Configure a router and use a logger for the UI. Initializes routes and view models.
func Init(router *fiber.App) {
	router.Use(
		"/static",
		compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}),
		func(c *fiber.Ctx) error {
			c.Set("Cache-Control", "public, max-age=31536000")
			return c.Next()
		},
		filesystem.New(filesystem.Config{
			Root:       http.FS(staticFs),
			PathPrefix: "static",
			Browse:     true,
		}),
	)

	router.Get("/", renderIndex)
	router.Get("/sensor/:sensor_id", renderSensorPage)
}

// A custom Fiber error handler which renders a simple web page.
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	values := fiber.Map{
		"code":    code,
		"title":   "Oops… You just found an error page",
		"message": "Try adjusting your search or filter to find what you're looking for.",
	}

	if code == http.StatusNotFound {
		values["message"] = "Couldn't find what you were looking for."
	}

	err = ctx.Render("views/error", values)
	if err != nil {
		// In case the render fails
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}

func renderIndex(ctx *fiber.Ctx) error {
	sensorStats, err := fetchAllSensorStats()
	if err != nil {
		return err
	}

	return ctx.Render("views/index", fiber.Map{
		"sensors": sensorStats,
	})
}

func renderSensorPage(ctx *fiber.Ctx) error {
	id := ctx.Params("sensor_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	stat, err := fetchSensorStats(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	campaigns := []database.Campaign{}
	err = database.Instance().
		Model(&database.Campaign{}).
		Where("? = any (sensors)", id).
		Find(&campaigns).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.Render("views/sensor", fiber.Map{
		"campaigns": campaigns,
		"stats":     stat,
	})
}
