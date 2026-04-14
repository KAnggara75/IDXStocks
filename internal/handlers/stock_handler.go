package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/usecases"
	"github.com/gofiber/fiber/v3"
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

func (h *StockHandler) PreviewHandler(c fiber.Ctx) error {
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

func (h *StockHandler) UploadHandler(c fiber.Ctx) error {
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

func (h *StockHandler) getFileAndValidate(c fiber.Ctx) (multipart.File, error) {
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

func (h *StockHandler) SyncIDHandler(c fiber.Ctx) error {
	stocks, err := h.usecase.SyncStockIDs(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync stock IDs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

func (h *StockHandler) SyncStockDetailHandler(c fiber.Ctx) error {
	// Run synchronization in background
	go func(ctx context.Context) {
		_, err := h.usecase.SyncStockDetail(ctx)
		if err != nil {
			logrus.Errorf("Background stock sync failed: %v", err)
		}
	}(context.WithoutCancel(c.Context()))

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stock synchronization process started in the background",
	})
}

func (h *StockHandler) SyncDelistingStocksHandler(c fiber.Ctx) error {
	var req models.SyncDelistingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Year == 0 || req.Month == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year and month are required",
		})
	}

	stocks, err := h.usecase.SyncDelistingStocks(c.Context(), req.Year, req.Month)
	if err != nil {
		logrus.Errorf("Failed to sync delisting stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}
