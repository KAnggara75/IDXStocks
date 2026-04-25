package handlers

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/usecases"
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
	code := c.Params("code")
	output := c.Query("output", "json")
	fieldsRaw := c.Query("fields")
	startDateRaw := c.Query("start_date")
	endDateRaw := c.Query("end_date")

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is required",
		})
	}

	var startDate, endDate *time.Time
	if startDateRaw != "" {
		t, err := time.Parse("2006-01-02", startDateRaw)
		if err == nil {
			startDate = &t
		} else {
			logrus.Warnf("Invalid start_date format: %s, expected YYYY-MM-DD", startDateRaw)
		}
	}
	if endDateRaw != "" {
		t, err := time.Parse("2006-01-02", endDateRaw)
		if err == nil {
			endDate = &t
		} else {
			logrus.Warnf("Invalid end_date format: %s, expected YYYY-MM-DD", endDateRaw)
		}
	}

	data, err := h.usecase.GetStockHistory(c.Context(), code, startDate, endDate)
	if err != nil {
		logrus.Errorf("Failed to get stock history for %s: %v", code, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Define all available fields and their getter/formatter
	allFields := []string{
		"code", "date", "previous", "open_price", "first_trade", "high", "low", "close", "change",
		"volume", "value", "frequency", "index_individual", "offer", "offer_volume", "bid", "bid_volume",
		"listed_shares", "tradeble_shares", "weight_for_index", "foreign_sell", "foreign_buy",
		"delisting_date", "non_regular_volume", "non_regular_value", "non_regular_frequency", "last_modified",
	}

	var requestedFields []string
	if fieldsRaw != "" {
		requestedFields = strings.Split(fieldsRaw, ",")
		// Trim spaces and lowercase
		for i, f := range requestedFields {
			requestedFields[i] = strings.ToLower(strings.TrimSpace(f))
		}
	} else {
		// Default fields (exclude code and last_modified)
		for _, f := range allFields {
			if f != "code" && f != "last_modified" {
				requestedFields = append(requestedFields, f)
			}
		}
	}

	if output == "csv" {
		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.csv\"", strings.ToUpper(code)))

		writer := csv.NewWriter(c)
		// Write header
		if err := writer.Write(requestedFields); err != nil {
			return err
		}

		for _, row := range data {
			line := make([]string, len(requestedFields))
			for i, field := range requestedFields {
				line[i] = getFieldAsString(row, field)
			}
			if err := writer.Write(line); err != nil {
				return err
			}
		}

		writer.Flush()
		return nil
	}

	// JSON Output
	results := make([]map[string]any, len(data))
	for i, row := range data {
		results[i] = make(map[string]any)
		for _, field := range requestedFields {
			results[i][field] = getFieldValue(row, field)
		}
	}

	return c.JSON(results)
}

func getFieldAsString(row models.StockHistory, field string) string {
	switch field {
	case "code":
		return row.Code
	case "date":
		return row.Date.Format("2006-01-02")
	case "previous":
		return formatFloat(row.Previous)
	case "open_price":
		return formatFloat(row.OpenPrice)
	case "first_trade":
		return formatFloat(row.FirstTrade)
	case "high":
		return formatFloat(row.High)
	case "low":
		return formatFloat(row.Low)
	case "close":
		return formatFloat(row.Close)
	case "change":
		return formatFloat(row.Change)
	case "volume":
		return formatFloat(row.Volume)
	case "value":
		return formatFloat(row.Value)
	case "frequency":
		return formatFloat(row.Frequency)
	case "index_individual":
		return formatFloat(row.IndexIndividual)
	case "offer":
		return formatFloat(row.Offer)
	case "offer_volume":
		return formatFloat(row.OfferVolume)
	case "bid":
		return formatFloat(row.Bid)
	case "bid_volume":
		return formatFloat(row.BidVolume)
	case "listed_shares":
		return formatFloat(row.ListedShares)
	case "tradeble_shares":
		return formatFloat(row.TradebleShares)
	case "weight_for_index":
		return formatFloat(row.WeightForIndex)
	case "foreign_sell":
		return formatFloat(row.ForeignSell)
	case "foreign_buy":
		return formatFloat(row.ForeignBuy)
	case "delisting_date":
		return formatString(row.DelistingDate)
	case "non_regular_volume":
		return formatFloat(row.NonRegularVolume)
	case "non_regular_value":
		return formatFloat(row.NonRegularValue)
	case "non_regular_frequency":
		return formatFloat(row.NonRegularFrequency)
	case "last_modified":
		return row.LastModified.Format(time.RFC3339)
	default:
		return ""
	}
}

func getFieldValue(row models.StockHistory, field string) any {
	switch field {
	case "code":
		return row.Code
	case "date":
		return row.Date
	case "previous":
		return row.Previous
	case "open_price":
		return row.OpenPrice
	case "first_trade":
		return row.FirstTrade
	case "high":
		return row.High
	case "low":
		return row.Low
	case "close":
		return row.Close
	case "change":
		return row.Change
	case "volume":
		return row.Volume
	case "value":
		return row.Value
	case "frequency":
		return row.Frequency
	case "index_individual":
		return row.IndexIndividual
	case "offer":
		return row.Offer
	case "offer_volume":
		return row.OfferVolume
	case "bid":
		return row.Bid
	case "bid_volume":
		return row.BidVolume
	case "listed_shares":
		return row.ListedShares
	case "tradeble_shares":
		return row.TradebleShares
	case "weight_for_index":
		return row.WeightForIndex
	case "foreign_sell":
		return row.ForeignSell
	case "foreign_buy":
		return row.ForeignBuy
	case "delisting_date":
		return row.DelistingDate
	case "non_regular_volume":
		return row.NonRegularVolume
	case "non_regular_value":
		return row.NonRegularValue
	case "non_regular_frequency":
		return row.NonRegularFrequency
	case "last_modified":
		return row.LastModified
	default:
		return nil
	}
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
