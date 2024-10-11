package middleware

import (
	"assesment/pkg/log"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func BindFiberContext(c *fiber.Ctx, logger log.Logger) context.Context {
	var correlationId = c.Locals("correlationId")
	var ctx = context.WithValue(context.Background(), "correlationId", correlationId)
	ctx = context.WithValue(ctx, "fiberContext", c)
	return logger.WithCorrelationId(ctx, correlationId.(string))
}

func AddCorrelationId(app *fiber.App) {

	app.Use(func(c *fiber.Ctx) error {
		var correlationId = c.Get("correlationId")

		if correlationId == "" {
			correlationId = uuid.NewString()
		}
		c.Locals("correlationId", correlationId)

		return c.Next()
	})
}
