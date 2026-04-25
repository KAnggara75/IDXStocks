package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/sirupsen/logrus"
)

type SeederService interface {
	SeedBrokersData(ctx context.Context) error
}

type seederService struct {
	brokerRepo repositories.BrokerRepository
	client     *http.Client
}

func NewSeederService(brokerRepo repositories.BrokerRepository, client *http.Client) SeederService {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &seederService{
		brokerRepo: brokerRepo,
		client:     client,
	}
}

func (s *seederService) SeedBrokersData(ctx context.Context) error {
	logrus.Info("Starting brokers data seeding...")

	url := "https://raw.githubusercontent.com/KAnggara75/IDXStock/refs/heads/main/testData/broker.json"

	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch broker data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch broker data: status %d", resp.StatusCode)
	}

	var seederResp models.BrokerSeederResponse
	if err := json.NewDecoder(resp.Body).Decode(&seederResp); err != nil {
		return fmt.Errorf("failed to decode broker data: %w", err)
	}

	if len(seederResp.Data.List) == 0 {
		logrus.Warn("No broker data found in JSON response")
		return nil
	}

	if err := s.brokerRepo.BatchInsertBrokers(ctx, seederResp.Data.List); err != nil {
		return fmt.Errorf("failed to insert broker data: %w", err)
	}

	logrus.Infof("Finished brokers data seeding. Total: %d", len(seederResp.Data.List))
	return nil
}
