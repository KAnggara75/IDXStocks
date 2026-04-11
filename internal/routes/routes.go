package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/KAnggara75/IDXStocks/internal/database"
	"github.com/KAnggara75/IDXStocks/internal/handlers"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
	"github.com/KAnggara75/IDXStocks/internal/usecases"
)

func Setup(app *fiber.App) {
	// Dependency Injection
	stockRepo := repositories.NewStockRepository(database.Pool)
	sectorRepo := repositories.NewSectorRepository(database.Pool)
	stockService := services.NewStockService()
	stockUsecase := usecases.NewStockUsecase(stockRepo, sectorRepo, stockService)
	stockHandler := handlers.NewStockHandler(stockUsecase)

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

	// API V1 Routes
	v1 := app.Group("/api/v1")
	v1.Post("/stocks/upload", stockHandler.PreviewHandler)
	v1.Patch("/stocks/upload", stockHandler.UploadHandler)
	v1.Put("/stocks/id", stockHandler.SyncIDHandler)
	v1.Put("/stocks/sync-sector", stockHandler.SyncSectorHandler)
}
