package routes

import (
	"nova-cdn/internal/config"
	"nova-cdn/internal/controllers"
	"nova-cdn/internal/middleware"
	"nova-cdn/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	app.Static("/", "./public")

	userRepo := repositories.NewUserRepository(db)
	authController := controllers.NewAuthController(*userRepo)

	api := app.Group("/api")

	api.Get("/documentation/*", swagger.HandlerDefault)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)

	galleryRepo := repositories.NewGalleryRepository(db)
	genRepo := repositories.NewGenerateRepository(db)

	galleryController := controllers.NewGalleryController(galleryRepo, genRepo)

	galleries := api.Group("/galleries", middleware.Auth())
	galleries.Get("/", galleryController.Index)
	galleries.Post("/upload", galleryController.Upload)
	galleries.Get("/:id<int>", galleryController.Show)
	galleries.Get("/:group_code<string>", galleryController.ShowByGroupCode)
	galleries.Delete("/:id<int>", galleryController.Destroy)
	galleries.Delete("/:id<int>/force", galleryController.ForceDelete)
	galleries.Post("/:id<int>/restore", galleryController.Restore)
}
