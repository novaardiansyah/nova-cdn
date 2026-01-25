package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"nova-cdn/internal/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) fiber.Handler {
	PersonalAccessTokenRepo := repositories.NewPersonalAccessTokenRepository(db)

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

		tokenIDStr := parts[0]
		plainTextToken := parts[1]

		tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		hash := sha256.Sum256([]byte(plainTextToken))
		hashedToken := hex.EncodeToString(hash[:])

		token, err := PersonalAccessTokenRepo.FindByIDAndHashedToken(tokenID, hashedToken)

		if err != nil {
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

		fields := map[string]interface{}{"last_used_at": time.Now()}
		PersonalAccessTokenRepo.UpdateFields(token, fields)

		UserId := token.TokenableID

		c.Locals("token", *token)
		c.Locals("user_id", UserId)

		return c.Next()
	}
}
