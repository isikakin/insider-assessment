package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Recover(app *fiber.App) {

	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Hata loglama
				log.Println("Panic caught:", r)

				// Panik durumunda özelleştirilmiş bir hata yanıtı döndür
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   true,
					"message": fmt.Sprintf("%v", r), // Hata mesajını burada döndür
				})
			}
		}()
		return c.Next()
	})
}
