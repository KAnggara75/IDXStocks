package handlers

import (
	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type SectorHandler struct {
	usecase usecases.SectorUsecase
}

func NewSectorHandler(usecase usecases.SectorUsecase) *SectorHandler {
	return &SectorHandler{
		usecase: usecase,
	}
}

func (h *SectorHandler) SyncNewSectorsHandler(c fiber.Ctx) error {
	results, err := h.usecase.SyncNewSectors(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync new sectors: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
