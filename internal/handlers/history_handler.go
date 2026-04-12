package handlers

import (
	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type HistoryHandler struct {
	usecase usecases.HistoryUsecase
}

func NewHistoryHandler(usecase usecases.HistoryUsecase) *HistoryHandler {
	return &HistoryHandler{
		usecase: usecase,
	}
}

func (h *HistoryHandler) SyncStockHistoryHandler(c fiber.Ctx) error {
	source := c.Query("source")

	if source == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Source is required",
		})
	}

	var req models.SyncHistoryRequest

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Year == 0 || req.Month == 0 || req.Day == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year, month, and day are required",
		})
	}

	err := h.usecase.SyncStockHistory(c.Context(), req, source)
	if err != nil {
		logrus.Errorf("Failed to sync stock history for %02d/%02d/%04d from %s: %v", req.Month, req.Day, req.Year, source, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stock history synchronization completed successfully",
		"date":    req,
		"source":  source,
	})
}
