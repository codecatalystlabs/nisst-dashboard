package middleware

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if fe, ok := err.(*fiber.Error); ok { code = fe.Code }
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
