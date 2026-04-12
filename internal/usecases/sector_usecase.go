package usecases

import (
	"context"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
)

type SectorUsecase interface {
	SyncSectors(ctx context.Context) ([]models.SectorResponse, error)
	SyncNewSectors(ctx context.Context) (*models.SectorSyncNewResponse, error)
}

type sectorUsecase struct {
	repo             repositories.SectorRepository
	searchRepo       repositories.SectorSearchRepository
	pasardanaService services.PasardanaService
}

func NewSectorUsecase(
	repo repositories.SectorRepository,
	searchRepo repositories.SectorSearchRepository,
	pasardanaService services.PasardanaService,
) SectorUsecase {
	return &sectorUsecase{
		repo:             repo,
		searchRepo:       searchRepo,
		pasardanaService: pasardanaService,
	}
}

func (u *sectorUsecase) SyncSectors(ctx context.Context) ([]models.SectorResponse, error) {
	pasardanaSectors, err := u.pasardanaService.FetchSectors()
	if err != nil {
		return nil, err
	}

	return u.repo.UpsertSectors(ctx, pasardanaSectors)
}

func (u *sectorUsecase) SyncNewSectors(ctx context.Context) (*models.SectorSyncNewResponse, error) {
	results, err := u.pasardanaService.FetchStockSearchResult()
	if err != nil {
		return nil, err
	}

	// Deduplicate sectors and sub-sectors
	sectorMap := make(map[int]models.SectorNew)
	subSectorMap := make(map[int]models.SubSector)

	for _, res := range results {
		if res.NewSectorId > 0 {
			sectorMap[res.NewSectorId] = models.SectorNew{
				Id:   res.NewSectorId,
				Name: res.NewSectorName,
			}
		}
		if res.NewSubSectorId > 0 {
			subSectorMap[res.NewSubSectorId] = models.SubSector{
				Id:       res.NewSubSectorId,
				Name:     res.NewSubSectorName,
				SectorId: res.NewSectorId,
			}
		}
	}

	sectors := make([]models.SectorNew, 0, len(sectorMap))
	for _, s := range sectorMap {
		sectors = append(sectors, s)
	}

	subSectors := make([]models.SubSector, 0, len(subSectorMap))
	for _, sub := range subSectorMap {
		subSectors = append(subSectors, sub)
	}

	// Upsert sectors first due to FK constraint
	updatedSectors, err := u.searchRepo.UpsertNewSectors(ctx, sectors)
	if err != nil {
		return nil, err
	}

	updatedSubSectors, err := u.searchRepo.UpsertNewSubSectors(ctx, subSectors)
	if err != nil {
		return nil, err
	}

	return &models.SectorSyncNewResponse{
		Sectors:    updatedSectors,
		SubSectors: updatedSubSectors,
	}, nil
}
