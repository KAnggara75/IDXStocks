package services

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type PasardanaService interface {
	FetchStockIDs() ([]models.PasardanaStock, error)
	FetchStockSearchResult() ([]models.PasardanaSearchResult, error)
	FetchNewSectors() ([]models.PasardanaNewSector, error)
	FetchNewSubSectors() ([]models.PasardanaNewSubSector, error)
	FetchStockDetailByCode(code string) (*models.PasardanaStockDetail, error)
	FetchStockHistory(year, month, day int) ([]models.PasardanaHistoryResponse, error)
}

type pasardanaService struct {
	client *http.Client
}

func NewPasardanaService(client *http.Client) PasardanaService {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &pasardanaService{
		client: client,
	}
}

func (s *pasardanaService) FetchStockIDs() ([]models.PasardanaStock, error) {
	url := "https://www.pasardana.id/api/Stock/GetAllSimpleStocks?username=anonymous"
	var results []models.PasardanaStock
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

func (s *pasardanaService) FetchStockDetailByCode(code string) (*models.PasardanaStockDetail, error) {
	url := fmt.Sprintf("https://www.pasardana.id/api/Stock/GetByCode?code=%s", code)
	var result models.PasardanaStockDetail
	if err := s.fetch(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *pasardanaService) FetchStockHistory(year, month, day int) ([]models.PasardanaHistoryResponse, error) {
	url := fmt.Sprintf("https://www.pasardana.id/api/StockSearchResult/GetAll?date=%02d/%02d/%04d", month, day, year)
	var results []models.PasardanaHistoryResponse
	if err := s.fetch(url, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *pasardanaService) fetch(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers per user request
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Host", "www.pasardana.id")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://pasardana.id/stock/search")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")

	// #nosec G107
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch from pasardana API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pasardana API returned status: %d", resp.StatusCode)
	}

	var body io.ReadCloser = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gz.Close()
		body = gz
	}

	if err := json.NewDecoder(body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
