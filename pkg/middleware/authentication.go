package middleware

import (
	_const "assesment/internal/const"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func BasicAuth(app *fiber.App, key string) {
	app.Use(func(c *fiber.Ctx) error {

		if c.Path() != "/ping" && !strings.Contains(c.Path(), "swagger") {
			authKey := c.Get(_const.AuthorizationKey)

			if authKey == "" || authKey != key {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}
		}

		return c.Next()
	})
}
