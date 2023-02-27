package api

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/swagger"
	"github.com/knadh/koanf"

	_ "github.com/openrfsense/backend/docs"
)

var swaggerConfig = swagger.Config{
	Title:  "OpenRFSense API",
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
	DocExpansion:             "list",
	SyntaxHighlight: &swagger.SyntaxHighlightConfig{
		Activate: true,
		Theme:    "nord",
	},
	ShowCommonExtensions:   true,
	SupportedSubmitMethods: []string{},
	TryItOutEnabled:        false,
}

// @title                     OpenRFSense backend API
// @description               OpenRFSense backend API
// @contact.name              OpenRFSense
// @contact.url               https://github.com/openrfsense/backend/issues
// @license.name              AGPLv3
// @license.url               https://spdx.org/licenses/AGPL-3.0-or-later.html
// @BasePath                  /api/v1
// @securityDefinitions.basic BasicAuth

// Creates a router for the public API. Initializes all REST endpoints under the given prefix
// and servers swagger documentation on /swagger.
func Init(config *koanf.Koanf, router *fiber.App, prefix string) {
	creds := config.MustStringMap("backend.users")

	// TODO: rate limiting?
	router.Use(
		helmet.New(),
		recover.New(),
		requestid.New(),
	)

	// Backend router for /api/v1
	router.Route(prefix, func(router fiber.Router) {
		router.Use(
			logger.New(),
			basicauth.New(basicauth.Config{
				Users: creds,
			}),
		)
		router.Get("/campaigns", CampaignsGet)
		router.Get("/samples", SamplesGet)
		router.Get("/nodes", NodesGet)
		router.Get("/nodes/:sensor_id", NodeGet)
		router.Post("/aggregated", AggregatedPost)
		router.Post("/raw", RawPost)
	})

	// Setup documentation routes
	router.Get("/api/docs/*", swagger.New(swaggerConfig))

	// Metrics page and API, but only if enabled in configuration
	if config.Bool("backend.metrics") {
		router.Get("/metrics", monitor.New())
	}
}
