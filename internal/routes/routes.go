package routes

import (
	"nova-cdn/internal/config"
	"nova-cdn/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	app.Use(middleware.GlobalLimiter())
	app.Static("/", "./public")

	api := app.Group("/api")

	api.Get("/documentation/*", swagger.HandlerDefault)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	AuthRoutes(api, db)
	GalleryRoutes(api, db)
}
