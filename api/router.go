package api

import (
	"fmt"
	"html/template"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/knadh/koanf"

	_ "github.com/openrfsense/backend/docs"
	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("api").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

var swaggerConfig = swagger.Config{
	Title:  "Swagger UI",
	Layout: "BaseLayout",
	Plugins: []template.JS{
		template.JS("SwaggerUIBundle.plugins.DownloadUrl"),
	},
	Presets: []template.JS{
		template.JS("SwaggerUIBundle.presets.apis"),
		template.JS("SwaggerUIStandalonePreset"),
	},
	DeepLinking:              true,
	DefaultModelsExpandDepth: 2,
	DefaultModelExpandDepth:  2,
	DefaultModelRendering:    "model",
	DocExpansion:             "list",
	SyntaxHighlight: &swagger.SyntaxHighlightConfig{
		Activate: true,
		Theme:    "agate",
	},
	ShowCommonExtensions:   true,
	ShowMutatedRequest:     true,
	SupportedSubmitMethods: []string{},
	TryItOutEnabled:        false,
}

// Creates a router for the public API. Initializes all REST endpoints under the given prefix
// and servers swagger documentation on /swagger.
func Start(config *koanf.Koanf, prefix string, routerConfig ...fiber.Config) *fiber.App {
	creds := config.MustStringMap("backend.users")

	router := fiber.New(routerConfig...)

	// TODO: rate limiting?
	router.Use(
		helmet.New(),
		logger.New(),
		recover.New(),
		requestid.New(),
	)

	// Backend router for /api/v1
	router.Route(prefix, func(router fiber.Router) {
		// TODO: pass auth to UI or just rate limit these for unauthenticated requests
		router.Get("/campaigns", CampaignsGet)
		router.Get("/campaigns/:campaign_id", CampaignGet)
		router.Get("/campaigns/:campaign_id/samples", CampaignSamplesGet)
		router.Get("/nodes/:sensor_id/campaigns", NodeCampaignsGet)
		router.Get("/nodes/:sensor_id/campaigns/:campaign_id", NodeCampaignSamplesGet)
		router.Get("/nodes/:sensor_id/samples", NodeSamplesGet)

		router.Use(basicauth.New(basicauth.Config{
			Users: creds,
		}))
		router.Get("/nodes", NodesGet)
		router.Get("/nodes/:sensor_id/stats", NodeStatsGet)
		router.Post("/aggregated", AggregatedPost)
		router.Post("/raw", RawPost)
	})

	// Setup documentation routes
	router.Get("/api/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})
	router.Get("/api/docs/*", swagger.New(swaggerConfig))

	// Metrics page and API, but only if enabled in configuration
	if config.Bool("backend.metrics") {
		router.Get("/metrics", monitor.New())
	}

	addr := fmt.Sprintf(":%d", config.MustInt("backend.port"))

	go func() {
		if err := router.Listen(addr); err != nil {
			log.Fatal(err)
		}
	}()

	return router
}
