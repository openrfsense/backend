package api

import (
	"html/template"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"

	_ "github.com/openrfsense/backend/docs"
	"github.com/openrfsense/common/config"
)

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
	ShowMutatedRequest:     true,
	TryItOutEnabled:        false,
	SupportedSubmitMethods: []string{},
}

// Configure a router for the public API. Initializes all REST endpoints under the given prefix
// and servers swagger documentation on /swagger.
func Use(router *fiber.App, prefix string) {
	creds := config.MustMap[string, string]("backend.users")

	// TODO: rate limiting?
	router.Use(
		helmet.New(),
		logger.New(),
		recover.New(),
		requestid.New(),
	)

	router.Route(prefix, func(router fiber.Router) {
		router.Use(basicauth.New(basicauth.Config{
			Users: creds,
		}))
		router.Post("/key", KeyPost)
		router.Get("/nodes", ListGet)
		router.Get("/nodes/:id/stats", NodeStatsGet)
	})

	// Setup documentation routes
	router.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
	router.Get("/swagger/*", swagger.New(swaggerConfig))
}
