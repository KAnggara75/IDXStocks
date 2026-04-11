package handlers

import (
	"path/filepath"
	"strings"

	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type StockHandler struct {
	usecase usecases.StockUsecase
}

func NewStockHandler(usecase usecases.StockUsecase) *StockHandler {
	return &StockHandler{
		usecase: usecase,
	}
}

func (h *StockHandler) UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("Failed to get file from form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only .xlsx and .xls files are allowed",
		})
	}

	f, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer f.Close()

	stocks, err := h.usecase.UploadStocks(c.Context(), f)
	if err != nil {
		logrus.Errorf("Failed to upload stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}
