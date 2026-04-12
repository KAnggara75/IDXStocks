package usecases

import (
	"context"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
)

type SectorUsecase interface {
	SyncNewSectors(ctx context.Context) (*models.SectorSyncNewResponse, error)
}

type sectorUsecase struct {
	searchRepo       repositories.SectorSearchRepository
	pasardanaService services.PasardanaService
}

func NewSectorUsecase(
	searchRepo repositories.SectorSearchRepository,
	pasardanaService services.PasardanaService,
) SectorUsecase {
	return &sectorUsecase{
		searchRepo:       searchRepo,
		pasardanaService: pasardanaService,
	}
}

func (u *sectorUsecase) SyncNewSectors(ctx context.Context) (*models.SectorSyncNewResponse, error) {
	// Sync Sectors
	pasardanaSectors, err := u.pasardanaService.FetchNewSectors()
	if err != nil {
		return nil, err
	}

	sectors := make([]models.SectorNew, 0, len(pasardanaSectors))
	for _, ps := range pasardanaSectors {
		sectors = append(sectors, models.SectorNew{
			Id:          ps.Id,
			Code:        ps.Code,
			Name:        ps.Name,
			NameEn:      ps.NameEn,
			Description: ps.Description,
		})
	}

	updatedSectors, err := u.searchRepo.UpsertNewSectors(ctx, sectors)
	if err != nil {
		return nil, err
	}

	// Sync SubSectors
	pasardanaSubSectors, err := u.pasardanaService.FetchNewSubSectors()
	if err != nil {
		return nil, err
	}

	subSectors := make([]models.SubSector, 0, len(pasardanaSubSectors))
	for _, ps := range pasardanaSubSectors {
		subSectors = append(subSectors, models.SubSector{
			Id:          ps.Id,
			SectorId:    ps.FkNewSectorId,
			Code:        ps.Code,
			Name:        ps.Name,
			NameEn:      ps.NameEn,
			Description: ps.Description,
		})
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
