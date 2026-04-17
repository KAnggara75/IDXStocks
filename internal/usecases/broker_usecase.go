package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/KAnggara75/IDXStocks/internal/repositories"
	"github.com/KAnggara75/IDXStocks/internal/services"
	"github.com/sirupsen/logrus"
)

type BrokerUsecase interface {
	SyncBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) ([]models.BrokerActivity, error)
}

type brokerUsecase struct {
	repo          repositories.BrokerActivityRepository
	brokerService services.BrokerService
}

func NewBrokerUsecase(repo repositories.BrokerActivityRepository, brokerService services.BrokerService) BrokerUsecase {
	return &brokerUsecase{
		repo:          repo,
		brokerService: brokerService,
	}
}

func (u *brokerUsecase) SyncBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) ([]models.BrokerActivity, error) {
	exodusResp, err := u.brokerService.FetchBrokerActivity(ctx, token, params)
	if err != nil {
		return nil, err
	}

	var records []models.BrokerActivity

	// Helper function for mapping
	mapItem := func(item models.ExodusBrokerActivityItem, side string) models.BrokerActivity {
		t, _ := time.Parse("2006-01-02", item.Date)
		return models.BrokerActivity{
			BrokerCode: item.BrokerCode,
			StockCode:  item.StockCode,
			Date:       t,
			Side:       side,
			Lot:        item.Lot,
			Value:      item.Value,
			AvgPrice:   item.AvgPrice,
			Freq:       item.Freq,
		}
	}

	// Map Buy items
	for _, item := range exodusResp.Data.BrokerActivityTransaction.BrokersBuy {
		records = append(records, mapItem(item, "buy"))
	}

	// Map Sell items
	for _, item := range exodusResp.Data.BrokerActivityTransaction.BrokersSell {
		records = append(records, mapItem(item, "sell"))
	}

	if len(records) > 0 {
		if err := u.repo.BatchUpsertBrokerActivity(ctx, records); err != nil {
			logrus.Errorf("Failed to upsert broker activity: %v", err)
			return nil, fmt.Errorf("failed to save broker activity to database: %w", err)
		}
	}

	return records, nil
}
