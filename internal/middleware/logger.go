package middleware

import (
	"fmt"
	"nova-cdn/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start)
		fmt.Printf("[%s] [%s] [%d] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), utils.FormatDuration(elapsed), c.Response().StatusCode(), c.Method(), c.Path())
		return err
	}
}
