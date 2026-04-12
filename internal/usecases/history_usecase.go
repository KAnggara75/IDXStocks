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

type HistoryUsecase interface {
	SyncStockHistory(ctx context.Context, req models.SyncHistoryRequest, source string) error
}

type historyUsecase struct {
	repo             repositories.HistoryRepository
	pasardanaService services.PasardanaService
	idxService       services.IdxService
}

func NewHistoryUsecase(
	repo repositories.HistoryRepository,
	pasardanaService services.PasardanaService,
	idxService services.IdxService,
) HistoryUsecase {
	return &historyUsecase{
		repo:             repo,
		pasardanaService: pasardanaService,
		idxService:       idxService,
	}
}

func (u *historyUsecase) SyncStockHistory(ctx context.Context, req models.SyncHistoryRequest, source string) error {
	var records []models.StockHistory
	targetDate := time.Date(req.Year, time.Month(req.Month), req.Day, 0, 0, 0, 0, time.Local)

	switch source {
	case "pasardana":
		data, err := u.pasardanaService.FetchStockHistory(req.Year, req.Month, req.Day)
		if err != nil {
			return err
		}
		for _, d := range data {
			// Parsing LastDate if provided, otherwise use targetDate
			var date time.Time
			if d.LastDate != nil && *d.LastDate != "" {
				// Format expected: "2024-12-10T00:00:00"
				t, err := time.Parse("2006-01-02T15:04:05", *d.LastDate)
				if err == nil {
					date = t
				} else {
					date = targetDate
				}
			} else {
				date = targetDate
			}

			records = append(records, models.StockHistory{
				Code:      d.Code,
				Date:      date,
				Previous:  d.PrevClosingPrice,
				OpenPrice: d.AdjustedOpenPrice,
				High:      d.AdjustedHighPrice,
				Low:       d.AdjustedLowPrice,
				Close:     d.AdjustedClosingPrice,
				Volume:    d.Volume,
				Frequency: d.Frequency,
				Value:     d.Value,
			})
		}

	case "idx":
		data, err := u.idxService.FetchStockSummary(req.Year, req.Month, req.Day)
		if err != nil {
			return err
		}
		for _, d := range data {
			var date time.Time
			if d.Date != "" {
				t, err := time.Parse("2006-01-02T15:04:05", d.Date)
				if err == nil {
					date = t
				} else {
					date = targetDate
				}
			} else {
				date = targetDate
			}

			// DelistingDate mapping
			var dd *string
			if d.DelistingDate != "" {
				dd = &d.DelistingDate
			}

			records = append(records, models.StockHistory{
				Code:                d.StockCode,
				Date:                date,
				Previous:            d.Previous,
				OpenPrice:           d.OpenPrice,
				FirstTrade:          d.FirstTrade,
				High:                d.High,
				Low:                 d.Low,
				Close:               d.Close,
				Change:              d.Change,
				Volume:              d.Volume,
				Value:               d.Value,
				Frequency:           d.Frequency,
				IndexIndividual:     d.IndexIndividual,
				Offer:               d.Offer,
				OfferVolume:         d.OfferVolume,
				Bid:                 d.Bid,
				BidVolume:           d.BidVolume,
				ListedShares:        d.ListedShares,
				TradebleShares:      d.TradebleShares,
				WeightForIndex:      d.WeightForIndex,
				ForeignSell:         d.ForeignSell,
				ForeignBuy:          d.ForeignBuy,
				DelistingDate:       dd,
				NonRegularVolume:    d.NonRegularVolume,
				NonRegularValue:     d.NonRegularValue,
				NonRegularFrequency: d.NonRegularFrequency,
			})
		}
	default:
		return fmt.Errorf("invalid source: %s", source)
	}

	if len(records) == 0 {
		logrus.Warnf("No records found to sync for date %s from %s", targetDate.Format("2006-01-02"), source)
		return nil
	}

	return u.repo.BatchUpsertStockHistory(ctx, records)
}
