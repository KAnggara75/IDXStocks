package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/sirupsen/logrus"
)

type IdxService interface {
	FetchDelistedStocks(year, month int) ([]models.IdxDelistedStock, error)
	ParseIdxDate(dateStr string) (string, error)
}

type idxService struct{}

func NewIdxService() IdxService {
	return &idxService{}
}

func (s *idxService) FetchDelistedStocks(year, month int) ([]models.IdxDelistedStock, error) {
	url := fmt.Sprintf("https://idx.co.id/primary/DigitalStatistic/GetApiDataPaginated?urlName=LINK_DELISTING&periodYear=%d&periodMonth=%d", year, month)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent to avoid WAF blocking as per SOP
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Log request details as per user request
	logrus.Infof("Requesting IDX: %s", url)
	logrus.Debugf("Request Headers: %v", req.Header)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// #nosec G107
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from IDX: %w", err)
	}
	defer resp.Body.Close()

	logrus.Infof("IDX Response Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		// Log response body if possible to diagnose errors (like 403 WAF)
		var bodyMsg interface{}
		_ = json.NewDecoder(resp.Body).Decode(&bodyMsg)
		logrus.Errorf("IDX Error Response: %v", bodyMsg)
		return nil, fmt.Errorf("IDX API returned status: %d", resp.StatusCode)
	}

	var idxResp models.IdxDelistingResponse
	if err := json.NewDecoder(resp.Body).Decode(&idxResp); err != nil {
		logrus.Errorf("Failed to decode IDX response: %v", err)
		return nil, fmt.Errorf("failed to decode IDX response: %w", err)
	}

	return idxResp.Data, nil
}

func (s *idxService) ParseIdxDate(dateStr string) (string, error) {
	// Format example: "18 July 2025"
	// time.Parse layout for this: "02 January 2006"
	t, err := time.Parse("02 January 2006", strings.TrimSpace(dateStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}
	return t.Format("2006-01-02"), nil
}
