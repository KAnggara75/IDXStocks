package usecases

import (
	"context"
	"errors"
	"math"
	"time"

	"strings"

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
		return make([]models.BrokerActivity, 0), err
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
			Lot:        int64(math.Round(item.Lot)),
			Value:      int64(math.Round(item.Value)),
			AvgPrice:   item.AvgPrice,
			Freq:       item.Freq,
		}
	}

	// Map Buy items
	for _, item := range exodusResp.Data.BrokerActivityTransaction.BrokersBuy {
		if len(item.StockCode) > 4 || strings.HasPrefix(item.StockCode, "X") {
			continue
		}
		records = append(records, mapItem(item, "buy"))
	}

	// Map Sell items
	for _, item := range exodusResp.Data.BrokerActivityTransaction.BrokersSell {
		if len(item.StockCode) > 4 || strings.HasPrefix(item.StockCode, "X") {
			continue
		}
		records = append(records, mapItem(item, "sell"))
	}

	var errs []error
	insertedRecords := make([]models.BrokerActivity, 0)
	for _, rec := range records {
		inserted, err := u.repo.InsertBrokerActivity(ctx, rec)
		if err != nil {
			logrus.Errorf("Error inserting record: %v", err)
			errs = append(errs, err)
			continue
		}
		if inserted {
			insertedRecords = append(insertedRecords, rec)
		}
	}

	if len(errs) > 0 {
		return insertedRecords, errors.Join(errs...)
	}

	return insertedRecords, nil
}
