package main

import (
	"log"
	"nova-cdn/docs"
	_ "nova-cdn/docs"
	"nova-cdn/internal/config"
	"nova-cdn/internal/middleware"
	"nova-cdn/internal/routes"
	"os"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title Nova CDN API
// @version 1.0
// @description This is an official Nova CDN API documentation.
// @termsOfService https://novaardiansyah.my.id/live/nova-app/terms-of-service

// @contact.name API Support
// @contact.url https://novaardiansyah.my.id
// @contact.email support@novadev.my.id

// @license.name MIT License
// @license.url https://github.com/novaardiansyah/nova-cdn/blob/main/LICENSE

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"

func main() {
	config.LoadEnv()

	config.ConnectDatabase()

	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	routes.SetupRoutes(app)

	if config.AppURL != "" {
		host := config.AppURL
		host = strings.Replace(host, "http://", "", 1)
		host = strings.Replace(host, "https://", "", 1)
		docs.SwaggerInfo.Host = host
	} else {
		docs.SwaggerInfo.Host = "localhost:" + config.AppPort
	}

	addr := config.AppIP + ":" + config.AppPort

	log.Printf("Server starting on %s...\n", addr)
	log.Fatal(app.Listen(addr))
}
