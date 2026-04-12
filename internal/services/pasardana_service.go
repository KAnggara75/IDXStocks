package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KAnggara75/IDXStocks/internal/models"
)

type PasardanaService interface {
	FetchStockIDs() ([]models.PasardanaStock, error)
	FetchSectors() ([]models.PasardanaSector, error)
	FetchStockSearchResult() ([]models.PasardanaSearchResult, error)
	FetchNewSectors() ([]models.PasardanaNewSector, error)
	FetchNewSubSectors() ([]models.PasardanaNewSubSector, error)
}

type pasardanaService struct{}

func NewPasardanaService() PasardanaService {
	return &pasardanaService{}
}

func (s *pasardanaService) FetchStockIDs() ([]models.PasardanaStock, error) {
	url := "https://www.pasardana.id/api/Stock/GetAllSimpleStocks?username=anonymous"
	var results []models.PasardanaStock
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) FetchSectors() ([]models.PasardanaSector, error) {
	url := "https://www.pasardana.id/api/StockNewSector/GetAll"
	var results []models.PasardanaSector
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) FetchStockSearchResult() ([]models.PasardanaSearchResult, error) {
	url := "https://www.pasardana.id/api/StockSearchResult/GetAll"
	var results []models.PasardanaSearchResult
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) FetchNewSectors() ([]models.PasardanaNewSector, error) {
	url := "https://www.pasardana.id/api/StockNewSector/GetAll"
	var results []models.PasardanaNewSector
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) FetchNewSubSectors() ([]models.PasardanaNewSubSector, error) {
	url := "https://www.pasardana.id/api/StockNewSubSector/GetAll"
	var results []models.PasardanaNewSubSector
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) fetch(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch from pasardana API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pasardana API returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
