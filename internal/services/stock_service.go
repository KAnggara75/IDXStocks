package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/xuri/excelize/v2"
)

type StockService interface {
	ParseExcel(file io.Reader) ([]models.Stock, error)
	FetchPasardanaStockIDs() ([]models.PasardanaStock, error)
	FetchPasardanaSectors() ([]models.PasardanaSector, error)
	FetchPasardanaStockSearchResult() ([]models.PasardanaSearchResult, error)
}

type stockService struct{}

func NewStockService() StockService {
	return &stockService{}
}

func (s *stockService) ParseExcel(file io.Reader) ([]models.Stock, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel reader: %w", err)
	}
	defer f.Close()

	// Assuming data is in the first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet %s: %w", sheetName, err)
	}

	var stocks []models.Stock
	// Skip header row (index 0)
	for i, row := range rows {
		if i == 0 || len(row) < 6 {
			continue
		}

		// Handle comma separators in shares (e.g., 11,766,313,488)
		sharesStr := strings.ReplaceAll(row[4], ",", "")
		shares, _ := strconv.ParseInt(sharesStr, 10, 64)

		stock := models.Stock{
			Code:         row[1],
			CompanyName:  row[2],
			ListingDate:  parseIndoDate(row[3]),
			ListingBoard: row[5],
			Shares:       shares,
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func parseIndoDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}

	// Check if already in YYYY-MM-DD format
	if len(dateStr) == 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		return dateStr
	}

	// Mapping Indonesian/English month abbreviations to numerical months
	months := map[string]string{
		"Jan": "01",
		"Feb": "02",
		"Mar": "03",
		"Apr": "04",
		"Mei": "05",
		"May": "05",
		"Jun": "06",
		"Jul": "07",
		"Agu": "08",
		"Agt": "08",
		"Aug": "08",
		"Sep": "09",
		"Okt": "10",
		"Oct": "10",
		"Nov": "11",
		"Des": "12",
		"Dec": "12",
	}

	parts := strings.Split(dateStr, " ")
	if len(parts) != 3 {
		return dateStr // Return original if format is unexpected
	}

	day := parts[0]
	if len(day) == 1 {
		day = "0" + day
	}

	month, ok := months[parts[1]]
	if !ok {
		return dateStr
	}

	year := parts[2]

	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

func (s *stockService) FetchPasardanaStockIDs() ([]models.PasardanaStock, error) {
	url := "https://www.pasardana.id/api/Stock/GetAllSimpleStocks?username=anonymous"
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from pasardana API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pasardana API returned status: %d", resp.StatusCode)
	}

	var pasardanaStocks []models.PasardanaStock
	if err := json.NewDecoder(resp.Body).Decode(&pasardanaStocks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return pasardanaStocks, nil
}

func (s *stockService) FetchPasardanaSectors() ([]models.PasardanaSector, error) {
	url := "https://www.pasardana.id/api/StockNewSector/GetAll"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from pasardana API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pasardana API returned status: %d", resp.StatusCode)
	}

	var pasardanaSectors []models.PasardanaSector
	if err := json.NewDecoder(resp.Body).Decode(&pasardanaSectors); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return pasardanaSectors, nil
}

func (s *stockService) FetchPasardanaStockSearchResult() ([]models.PasardanaSearchResult, error) {
	url := "https://www.pasardana.id/api/StockSearchResult/GetAll"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from pasardana API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pasardana API returned status: %d", resp.StatusCode)
	}

	var pasardanaResults []models.PasardanaSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&pasardanaResults); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return pasardanaResults, nil
}
