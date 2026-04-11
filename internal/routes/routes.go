package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/KAnggara75/IDXStocks/internal/database"
)

func Setup(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		err := database.Pool.Ping(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "error",
				"message": "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "pong",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "up",
		})
	})
}
