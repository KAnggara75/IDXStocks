package handlers

import (
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
		return err
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
		return err
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
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only .xlsx and .xls files are allowed",
		})
	}

	return file.Open()
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
