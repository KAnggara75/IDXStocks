package handlers

import (
	"encoding/csv"
	"fmt"
	"strings"

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
func (h *HistoryHandler) GetStockHistoryHandler(c fiber.Ctx) error {
	code := c.Query("code")
	output := c.Query("output", "json")

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is required",
		})
	}

	data, err := h.usecase.GetStockHistory(c.Context(), code)
	if err != nil {
		logrus.Errorf("Failed to get stock history for %s: %v", code, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if output == "csv" {
		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"history_%s.csv\"", strings.ToLower(code)))

		writer := csv.NewWriter(c)
		// Write header
		header := []string{
			"code", "date", "previous", "open_price", "first_trade", "high", "low", "close", "change",
			"volume", "value", "frequency", "index_individual", "offer", "offer_volume", "bid", "bid_volume",
			"listed_shares", "tradeble_shares", "weight_for_index", "foreign_sell", "foreign_buy",
			"delisting_date", "non_regular_volume", "non_regular_value", "non_regular_frequency",
		}
		if err := writer.Write(header); err != nil {
			return err
		}

		for _, row := range data {
			line := []string{
				row.Code,
				row.Date.Format("2006-01-02"),
				formatFloat(row.Previous),
				formatFloat(row.OpenPrice),
				formatFloat(row.FirstTrade),
				formatFloat(row.High),
				formatFloat(row.Low),
				formatFloat(row.Close),
				formatFloat(row.Change),
				formatFloat(row.Volume),
				formatFloat(row.Value),
				formatFloat(row.Frequency),
				formatFloat(row.IndexIndividual),
				formatFloat(row.Offer),
				formatFloat(row.OfferVolume),
				formatFloat(row.Bid),
				formatFloat(row.BidVolume),
				formatFloat(row.ListedShares),
				formatFloat(row.TradebleShares),
				formatFloat(row.WeightForIndex),
				formatFloat(row.ForeignSell),
				formatFloat(row.ForeignBuy),
				formatString(row.DelistingDate),
				formatFloat(row.NonRegularVolume),
				formatFloat(row.NonRegularValue),
				formatFloat(row.NonRegularFrequency),
			}
			if err := writer.Write(line); err != nil {
				return err
			}
		}

		writer.Flush()
		return nil
	}

	return c.JSON(data)
}

func formatFloat(f *float64) string {
	if f == nil {
		return "0.0"
	}
	return fmt.Sprintf("%.1f", *f)
}

func formatString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
