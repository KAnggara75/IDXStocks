package routes

import (
	"github.com/KAnggara75/IDXStocks/internal/database"
	"github.com/KAnggara75/IDXStocks/internal/handlers"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	// Dependency Injection
	stockRepo := repositories.NewStockRepository(database.Pool)
	sectorSearchRepo := repositories.NewSectorSearchRepository(database.Pool)
	industryRepo := repositories.NewIndustryRepository(database.Pool)
	historyRepo := repositories.NewHistoryRepository(database.Pool)

	stockService := services.NewStockService()
	pasardanaService := services.NewPasardanaService()
	idxService := services.NewIdxService()

	stockUsecase := usecases.NewStockUsecase(stockRepo, stockService, pasardanaService, idxService)
	industryUsecase := usecases.NewIndustryUsecase(industryRepo, pasardanaService)
	sectorUsecase := usecases.NewSectorUsecase(sectorSearchRepo, pasardanaService)
	historyUsecase := usecases.NewHistoryUsecase(historyRepo, stockRepo, pasardanaService, idxService)

	stockHandler := handlers.NewStockHandler(stockUsecase)
	industryHandler := handlers.NewIndustryHandler(industryUsecase)
	sectorHandler := handlers.NewSectorHandler(sectorUsecase)
	historyHandler := handlers.NewHistoryHandler(historyUsecase)

	app.Get("/health", func(c fiber.Ctx) error {
		err := database.Pool.Ping(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status": "up",
		})
	})

	// API V1 Routes
	v1 := app.Group("/api/v1")
	v1.Post("/stocks/upload", stockHandler.PreviewHandler)
	v1.Patch("/stocks/upload", stockHandler.UploadHandler)
	v1.Put("/stocks/id", stockHandler.SyncIDHandler)
	v1.Put("/stocks/sync", stockHandler.SyncStockDetailHandler)
	v1.Put("/stocks/delisting/sync", stockHandler.SyncDelistingStocksHandler)
	v1.Put("/stocks/history/sync", historyHandler.SyncStockHistoryHandler)
	v1.Get("/stocks/:code/history", historyHandler.GetStockHistoryHandler)
	v1.Put("/sectors/sync", sectorHandler.SyncNewSectorsHandler)
	v1.Put("/industries/sync", industryHandler.IndustrySyncHandler)
}
