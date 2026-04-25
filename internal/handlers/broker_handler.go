package handlers

import (
	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type BrokerHandler struct {
	usecase usecases.BrokerUsecase
}

func NewBrokerHandler(usecase usecases.BrokerUsecase) *BrokerHandler {
	return &BrokerHandler{
		usecase: usecase,
	}
}

func (h *BrokerHandler) SyncBrokerActivityHandler(c fiber.Ctx) error {
	var params models.SyncBrokerActivityParams
	if err := c.Bind().Query(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse query parameters",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	activities, err := h.usecase.SyncBrokerActivity(c.Context(), token, params)
	if err != nil {
		logrus.Errorf("Broker activity sync error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(activities)
}

func (h *BrokerHandler) ManagePartitionsHandler(c fiber.Ctx) error {
	resp, err := h.usecase.ManagePartitions(c.Context())
	if err != nil {
		logrus.Errorf("Partition management error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status := fiber.StatusOK
	if resp.PartitionsCreated > 0 {
		status = fiber.StatusCreated
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}
