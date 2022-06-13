package api

import (
	"fmt"
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

// Create a router for the public API. Initializes all REST endpoints under the given prefix
// and servers swagger documentation on /swagger.
func Start(prefix string, routerConfig ...fiber.Config) *fiber.App {
	creds := config.MustMap[string, string]("backend.users")

	router := fiber.New(routerConfig...)

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
	router.Get("/api/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})
	router.Get("/api/docs/*", swagger.New(swaggerConfig))

	addr := fmt.Sprintf(":%d", config.GetWeakInt("backend.port"))

	go func() {
		if err := router.Listen(addr); err != nil {
			log.Fatal(err)
		}
	}()

	return router
}
