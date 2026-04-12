package usecases

import (
	"context"
	"io"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
	"github.com/sirupsen/logrus"
)

type StockUsecase interface {
	PreviewStocks(ctx context.Context, file io.Reader) ([]models.Stock, error)
	UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error)
	SyncStockIDs(ctx context.Context) ([]models.StockResponse, error)
	SyncStockDetail(ctx context.Context) ([]models.StockResponse, error)
}

type stockUsecase struct {
	repo             repositories.StockRepository
	service          services.StockService
	pasardanaService services.PasardanaService
}

func NewStockUsecase(
	repo repositories.StockRepository,
	service services.StockService,
	pasardanaService services.PasardanaService,
) StockUsecase {
	return &stockUsecase{
		repo:             repo,
		service:          service,
		pasardanaService: pasardanaService,
	}
}

func (u *stockUsecase) PreviewStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	return u.service.ParseExcel(file)
}

func (u *stockUsecase) UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	stocks, err := u.service.ParseExcel(file)
	if err != nil {
		return nil, err
	}

	err = u.repo.BatchInsertStocks(ctx, stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (u *stockUsecase) SyncStockIDs(ctx context.Context) ([]models.StockResponse, error) {
	pasardanaStocks, err := u.pasardanaService.FetchStockIDs()
	if err != nil {
		return nil, err
	}

	return u.repo.UpdateStockIDs(ctx, pasardanaStocks)
}

func (u *stockUsecase) SyncStockDetail(ctx context.Context) ([]models.StockResponse, error) {
	// 1. Get List Stock
	simpleStocks, err := u.pasardanaService.FetchStockIDs()
	if err != nil {
		return nil, err
	}

	allUpdated := make([]models.StockResponse, 0)

	// 2. Loop to get detail and immediately upsert
	for _, s := range simpleStocks {
		logrus.Debugf("Fetching and syncing detail for stock: %s", s.Code)
		detail, err := u.pasardanaService.FetchStockDetailByCode(s.Code)
		if err != nil {
			logrus.Errorf("Failed to fetch detail for %s: %v", s.Code, err)
			continue
		}

		if detail != nil {
			// Immediately upsert to db
			updated, err := u.repo.UpsertStocksDetail(ctx, []models.PasardanaStockDetail{*detail})
			if err != nil {
				logrus.Errorf("Failed to upsert detail for %s: %v", s.Code, err)
				continue
			}

			if len(updated) > 0 {
				allUpdated = append(allUpdated, updated...)
			}
		}
	}

	return allUpdated, nil
}
