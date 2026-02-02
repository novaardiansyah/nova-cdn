package routes

import (
	"nova-cdn/internal/controllers"
	"nova-cdn/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GalleryRoutes(api fiber.Router, db *gorm.DB) {
	galleryController := controllers.NewGalleryController(db)

	galleries := api.Group("/galleries", middleware.Auth(db))

	galleries.Get("/", galleryController.Index)
	galleries.Post("/upload", galleryController.Upload)
	galleries.Get("/:id<int>", galleryController.Show)
	galleries.Get("/:group_code<string>", galleryController.ShowByGroupCode)

	galleries.Post("/:id<int>/restore", galleryController.Restore)
	galleries.Post("/:group_code<string>/restore", galleryController.RestoreByGroupCode)

	galleries.Delete("/:id<int>", galleryController.Destroy)
	galleries.Delete("/:group_code<string>", galleryController.DestroyByGroupCode)

	galleries.Delete("/:id<int>/force", galleryController.ForceDelete)
	galleries.Delete("/:group_code<string>/force", galleryController.ForceDeleteByGroupCode)
}
