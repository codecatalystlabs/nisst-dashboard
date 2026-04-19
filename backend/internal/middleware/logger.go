package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Printf("method=%s path=%s status=%d elapsed_ms=%d", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start).Milliseconds())
		return err
	}
}
