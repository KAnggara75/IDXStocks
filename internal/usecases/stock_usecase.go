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
	SyncSectors(ctx context.Context) ([]models.SectorResponse, error)
}

type stockUsecase struct {
	repo             repositories.StockRepository
	sectorRepo       repositories.SectorRepository
	service          services.StockService
	pasardanaService services.PasardanaService
}

func NewStockUsecase(
	repo repositories.StockRepository,
	sectorRepo repositories.SectorRepository,
	service services.StockService,
	pasardanaService services.PasardanaService,
) StockUsecase {
	return &stockUsecase{
		repo:             repo,
		sectorRepo:       sectorRepo,
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

func (u *stockUsecase) SyncSectors(ctx context.Context) ([]models.SectorResponse, error) {
	pasardanaSectors, err := u.pasardanaService.FetchSectors()
	if err != nil {
		return nil, err
	}

	return u.sectorRepo.UpsertSectors(ctx, pasardanaSectors)
}
