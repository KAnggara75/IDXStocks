package handlers

import (
	"fmt"
	"mime/multipart"
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

func (h *StockHandler) PreviewHandler(c *fiber.Ctx) error {
	f, err := h.getFileAndValidate(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer f.Close()

	stocks, err := h.usecase.PreviewStocks(c.Context(), f)
	if err != nil {
		logrus.Errorf("Failed to preview stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

func (h *StockHandler) UploadHandler(c *fiber.Ctx) error {
	f, err := h.getFileAndValidate(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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

func (h *StockHandler) getFileAndValidate(c *fiber.Ctx) (multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("Failed to get file from form: %v", err)
		return nil, fmt.Errorf("file is required")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return nil, fmt.Errorf("only .xlsx and .xls files are allowed")
	}

	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}

func (h *StockHandler) SyncIDHandler(c *fiber.Ctx) error {
	stocks, err := h.usecase.SyncStockIDs(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync stock IDs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

func (h *StockHandler) SyncSectorHandler(c *fiber.Ctx) error {
	sectors, err := h.usecase.SyncSectors(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync sectors: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sectors)
}
