package usecases

import (
	"context"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
)

type SectorUsecase interface {
	SyncSectors(ctx context.Context) ([]models.SectorResponse, error)
}

type sectorUsecase struct {
	repo             repositories.SectorRepository
	pasardanaService services.PasardanaService
}

func NewSectorUsecase(repo repositories.SectorRepository, pasardanaService services.PasardanaService) SectorUsecase {
	return &sectorUsecase{
		repo:             repo,
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
