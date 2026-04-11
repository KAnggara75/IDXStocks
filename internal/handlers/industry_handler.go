package handlers

import (
	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IndustryHandler struct {
	usecase usecases.IndustryUsecase
}

func NewIndustryHandler(usecase usecases.IndustryUsecase) *IndustryHandler {
	return &IndustryHandler{
		usecase: usecase,
	}
}

func (h *IndustryHandler) IndustrySyncHandler(c *fiber.Ctx) error {
	results, err := h.usecase.SyncIndustry(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync industries: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
