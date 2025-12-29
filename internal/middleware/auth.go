package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"nova-cdn/internal/config"
	"nova-cdn/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: No token provided",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		parts := strings.SplitN(tokenString, "|", 2)
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		tokenID := parts[0]
		plainTextToken := parts[1]

		hash := sha256.Sum256([]byte(plainTextToken))
		hashedToken := hex.EncodeToString(hash[:])

		db := config.GetDB()
		var token models.PersonalAccessToken

		result := db.Where("id = ? AND token = ?", tokenID, hashedToken).First(&token)
		if result.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token",
			})
		}

		if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Token expired",
			})
		}

		db.Model(&token).Update("last_used_at", time.Now())

		c.Locals("token", token)
		c.Locals("user_id", token.TokenableID)

		return c.Next()
	}
}
