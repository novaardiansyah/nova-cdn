package routes

import (
	"nova-cdn/internal/controllers"
	"nova-cdn/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoutes(api fiber.Router, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	auth := api.Group("/auth")

	auth.Use(middleware.AuthLimiter())
	auth.Post("/login", authController.Login)
}
