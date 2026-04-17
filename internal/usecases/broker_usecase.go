package usecases

import (
	"context"
	"errors"
	"fmt"
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
	ManagePartitions(ctx context.Context) (*models.PartitionManagementResponse, error)
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

func (u *brokerUsecase) ManagePartitions(ctx context.Context) (*models.PartitionManagementResponse, error) {
	today := time.Now()

	// Helper to find next Monday
	// today.Weekday(): 0 (Sun), 1 (Mon) ... 6 (Sat)
	// If today is Monday (1), we want next Monday (7 days away) or stay (0 days)?
	// The requirement: "hit on 17 Apr (Fri) -> 20 Apr (Mon)".
	// 20 Apr is the next Monday.

	resp := &models.PartitionManagementResponse{
		Details: make([]models.PartitionDetail, 0),
	}

	for i := 0; i < 2; i++ {
		// Calculate target Monday
		// Days to next Monday: 1 - current weekday + (7 if weekday >= 1 else 0)
		daysToMonday := (8 - int(today.Weekday())) % 7
		if daysToMonday == 0 {
			daysToMonday = 7
		}

		// Start of week i (0-based)
		startOfWeek := today.AddDate(0, 0, daysToMonday+(i*7))
		endOfWeek := startOfWeek.AddDate(0, 0, 7) // To Monday next week (exclusive bound)

		year, week := startOfWeek.ISOWeek()
		tableName := fmt.Sprintf("broker_activity_p_%d_w%02d", year, week)
		rangeStr := fmt.Sprintf("%s to %s", startOfWeek.Format("2006-01-02"), endOfWeek.AddDate(0, 0, -1).Format("2006-01-02"))

		exists, err := u.repo.CheckPartitionExists(ctx, tableName)
		if err != nil {
			return nil, err
		}

		detail := models.PartitionDetail{
			Name:  tableName,
			Range: rangeStr,
		}

		if exists {
			detail.Status = "exists"
		} else {
			err := u.repo.CreatePartition(ctx, tableName, startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02"))
			if err != nil {
				return nil, err
			}
			detail.Status = "created"
			resp.PartitionsCreated++
		}

		resp.Details = append(resp.Details, detail)
	}

	return resp, nil
}
