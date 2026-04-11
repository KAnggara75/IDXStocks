package services

import (
	"fmt"
	"io"
	"strconv"

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

		shares, _ := strconv.ParseInt(row[5], 10, 64)

		stock := models.Stock{
			Code:          row[0],
			CompanyName:   row[1],
			ListingDate:   row[2],
			DelistingDate: row[3],
			ListingBoard:  row[4],
			Shares:        shares,
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}
