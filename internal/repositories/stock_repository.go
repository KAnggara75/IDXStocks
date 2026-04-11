package repositories

import (
	"context"
	"github.com/KAnggara75/IDXStocks/internal/models"
)

type StockRepository interface {
	BatchInsertStocks(ctx context.Context, stocks []models.Stock) error
}

type stockRepository struct {
	// db connection if needed later
}

func NewStockRepository() StockRepository {
	return &stockRepository{}
}

func (r *stockRepository) BatchInsertStocks(ctx context.Context, stocks []models.Stock) error {
	// Plan placeholder: Placeholder for database insertion
	return nil
}
