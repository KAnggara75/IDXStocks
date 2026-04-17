package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/KAnggara75/IDXStocks/internal/models"
)

type BrokerService interface {
	FetchBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) (*models.ExodusBrokerActivityResponse, error)
}

type brokerService struct{}

func NewBrokerService() BrokerService {
	return &brokerService{}
}

func (s *brokerService) FetchBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) (*models.ExodusBrokerActivityResponse, error) {
	baseURL := "https://exodus.stockbit.com/order-trade/broker/activity"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	if params.BrokerCode != "" {
		q.Set("broker_code", params.BrokerCode)
	}
	if params.From != "" {
		q.Set("from", params.From)
	}
	if params.To != "" {
		q.Set("to", params.To)
	}
	if params.TransactionType != "" {
		q.Set("transaction_type", params.TransactionType)
	}
	if params.MarketBoard != "" {
		q.Set("market_board", params.MarketBoard)
	}
	if params.InvestorType != "" {
		q.Set("investor_type", params.InvestorType)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Forward Bearer token
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Exodus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exodus API returned status: %d", resp.StatusCode)
	}

	var exodusResp models.ExodusBrokerActivityResponse
	if err := json.NewDecoder(resp.Body).Decode(&exodusResp); err != nil {
		return nil, fmt.Errorf("failed to decode Exodus response: %w", err)
	}

	return &exodusResp, nil
}
