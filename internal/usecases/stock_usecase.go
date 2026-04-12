package usecases

import (
	"context"
	"io"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
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

	// 2. Loop to get detail for each stock
	details := make([]models.PasardanaStockDetail, 0, len(simpleStocks))

	// For simplicity and to avoid hitting Pasardana too hard, we can use a serial loop
	// or a limited concurrency loop. Let's start with a simple loop for now as per "junior/AI" request.
	for _, s := range simpleStocks {
		detail, err := u.pasardanaService.FetchStockDetailByCode(s.Code)
		if err != nil {
			// Skip if failed to fetch detail for one stock
			continue
		}
		if detail != nil {
			details = append(details, *detail)
		}
	}

	// 3. Upsert to Repository
	return u.repo.UpsertStocksDetail(ctx, details)
}
