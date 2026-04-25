package usecases

import (
	"context"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/services"
)

type IndustryUsecase interface {
	SyncIndustry(ctx context.Context) (*models.IndustrySyncResponse, error)
}

type industryUsecase struct {
	repo             repositories.IndustryRepository
	pasardanaService services.PasardanaService
}

func NewIndustryUsecase(repo repositories.IndustryRepository, pasardanaService services.PasardanaService) IndustryUsecase {
	return &industryUsecase{
		repo:             repo,
		pasardanaService: pasardanaService,
	}
}

func (u *industryUsecase) SyncIndustry(ctx context.Context) (*models.IndustrySyncResponse, error) {
	results, err := u.pasardanaService.FetchStockSearchResult()
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
	updatedIndustries, err := u.repo.UpsertIndustries(ctx, industries)
	if err != nil {
		return nil, err
	}

	updatedSubIndustries, err := u.repo.UpsertSubIndustries(ctx, subIndustries)
	if err != nil {
		return nil, err
	}

	return &models.IndustrySyncResponse{
		Industries:    updatedIndustries,
		SubIndustries: updatedSubIndustries,
	}, nil
}
