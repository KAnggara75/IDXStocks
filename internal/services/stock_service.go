package services

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/xuri/excelize/v2"
)

type StockService interface {
	ParseExcel(file io.Reader) ([]models.Stock, error)
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
			ListingDate:  row[3],
			ListingBoard: row[5],
			Shares:       shares,
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}
