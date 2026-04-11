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
	SyncIndustry(ctx context.Context) (*models.IndustrySyncResponse, error)
}

type stockUsecase struct {
	repo         repositories.StockRepository
	industryRepo repositories.IndustryRepository
	service      services.StockService
}

func NewStockUsecase(repo repositories.StockRepository, industryRepo repositories.IndustryRepository, service services.StockService) StockUsecase {
	return &stockUsecase{
		repo:         repo,
		industryRepo: industryRepo,
		service:      service,
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
	pasardanaStocks, err := u.service.FetchPasardanaStockIDs()
	if err != nil {
		return nil, err
	}

	return u.repo.UpdateStockIDs(ctx, pasardanaStocks)
}

func (u *stockUsecase) SyncIndustry(ctx context.Context) (*models.IndustrySyncResponse, error) {
	results, err := u.service.FetchPasardanaStockSearchResult()
	if err != nil {
		return nil, err
	}

	// Deduplicate industries and sub-industries
	industryMap := make(map[int]models.Industry)
	subIndustryMap := make(map[int]models.SubIndustry)

	for _, res := range results {
		if res.NewIndustryId > 0 {
			industryMap[res.NewIndustryId] = models.Industry{
				Id:   res.NewIndustryId,
				Name: res.NewIndustryName,
			}
		}
		if res.NewSubIndustryId > 0 {
			subIndustryMap[res.NewSubIndustryId] = models.SubIndustry{
				Id:         res.NewSubIndustryId,
				Name:       res.NewSubIndustryName,
				IndustryId: res.NewIndustryId,
			}
		}
	}

	industries := make([]models.Industry, 0, len(industryMap))
	for _, ind := range industryMap {
		industries = append(industries, ind)
	}

	subIndustries := make([]models.SubIndustry, 0, len(subIndustryMap))
	for _, sub := range subIndustryMap {
		subIndustries = append(subIndustries, sub)
	}

	// Upsert industries first due to FK constraint
	updatedIndustries, err := u.industryRepo.UpsertIndustries(ctx, industries)
	if err != nil {
		return nil, err
	}

	updatedSubIndustries, err := u.industryRepo.UpsertSubIndustries(ctx, subIndustries)
	if err != nil {
		return nil, err
	}

	return &models.IndustrySyncResponse{
		Industries:    updatedIndustries,
		SubIndustries: updatedSubIndustries,
	}, nil
}
