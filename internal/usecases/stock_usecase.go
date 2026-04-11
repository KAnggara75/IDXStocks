package usecases

import (
	"context"
	"io"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
)

type StockUsecase interface {
	UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error)
}

type stockUsecase struct {
	repo    repositories.StockRepository
	service services.StockService
}

func NewStockUsecase(repo repositories.StockRepository, service services.StockService) StockUsecase {
	return &stockUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *stockUsecase) UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	stocks, err := u.service.ParseExcel(file)
	if err != nil {
		return nil, err
	}

	// Optional: Save to repository
	err = u.repo.BatchInsertStocks(ctx, stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
