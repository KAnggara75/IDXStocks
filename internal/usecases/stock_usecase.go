package usecases

import (
	"context"
	"io"

	"sync"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/services"
	"github.com/KAnggara75/IDXStock/internal/utils"
	"github.com/sirupsen/logrus"
)

type StockUsecase interface {
	PreviewStocks(ctx context.Context, file io.Reader) ([]models.Stock, error)
	UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error)
	SyncStockDetail(ctx context.Context) ([]models.StockResponse, error)
	SyncDelistingStocks(ctx context.Context, year, month int) ([]models.StockResponse, error)
}

type stockUsecase struct {
	repo             repositories.StockRepository
	service          services.StockService
	pasardanaService services.PasardanaService
	idxService       services.IdxService
}

func NewStockUsecase(
	repo repositories.StockRepository,
	service services.StockService,
	pasardanaService services.PasardanaService,
	idxService services.IdxService,
) StockUsecase {
	return &stockUsecase{
		repo:             repo,
		service:          service,
		pasardanaService: pasardanaService,
		idxService:       idxService,
	}
}

func (u *stockUsecase) PreviewStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	return u.service.ParseExcel(file)
}

func (u *stockUsecase) UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	stocks, err := u.service.ParseExcel(file)
	if err != nil {
		return nil, err
	}

	err = u.repo.BatchInsertStocks(ctx, stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (u *stockUsecase) SyncStockDetail(ctx context.Context) ([]models.StockResponse, error) {
	// 1. Get List Stock
	simpleStocks, err := u.pasardanaService.FetchStockIDs()
	if err != nil {
		return nil, err
	}

	totalStocks := len(simpleStocks)
	logrus.Infof("Starting parallel sync for %d stocks", totalStocks)

	// 2. Setup Worker Pool
	const numWorkers = 20
	jobs := make(chan models.PasardanaStock, totalStocks)
	results := make(chan []models.StockResponse, totalStocks)

	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		// #nosec G204
		// #nosec G118
		go func(workerID int) {
			defer wg.Done()
			for s := range jobs {
				detail, err := u.pasardanaService.FetchStockDetailByCode(s.Code)
				if err != nil {
					logrus.Errorf("[Worker %d] Failed to fetch %s: %v", workerID, s.Code, err)
					continue
				}

				if detail != nil {
					// Parse dates to YYYY-MM-DD format, default to Unix epoch (1970-01-01) if null or empty
					epoch0 := "1970-01-01"
					if detail.ListingDate != nil && *detail.ListingDate != "" {
						parsed := utils.NormalizeDate(*detail.ListingDate)
						if parsed == "" {
							parsed = epoch0
						}
						detail.ListingDate = &parsed
					} else {
						detail.ListingDate = &epoch0
					}

					if detail.FoundingDate != nil && *detail.FoundingDate != "" {
						parsed := utils.NormalizeDate(*detail.FoundingDate)
						if parsed == "" {
							parsed = epoch0
						}
						detail.FoundingDate = &parsed
					} else {
						detail.FoundingDate = &epoch0
					}

					updated, err := u.repo.UpsertStocksDetail(context.Background(), []models.PasardanaStockDetail{*detail})
					if err != nil {
						logrus.Errorf("[Worker %d] Failed to upsert %s: %v", workerID, s.Code, err)
						continue
					}
					results <- updated
				}
			}
		}(w)
	}

	// Send jobs
	for _, s := range simpleStocks {
		jobs <- s
	}
	close(jobs)

	// Wait for workers in a separate goroutine to close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	allUpdated := make([]models.StockResponse, 0)
	for res := range results {
		allUpdated = append(allUpdated, res...)
	}

	logrus.Infof("Parallel sync completed. Total updated: %d", len(allUpdated))
	return allUpdated, nil
}

func (u *stockUsecase) SyncDelistingStocks(ctx context.Context, year, month int) ([]models.StockResponse, error) {
	// 1. Fetch from IDX
	idxStocks, err := u.idxService.FetchDelistedStocks(year, month)
	if err != nil {
		return nil, err
	}

	updatedStocks := make([]models.StockResponse, 0)

	// 2. Loop and Update
	for _, s := range idxStocks {
		// Parse date
		formattedDate := utils.NormalizeDate(s.DeListingDate)
		if formattedDate == "" {
			logrus.Warnf("Skipping stock %s: invalid delisting date format %s", s.Code, s.DeListingDate)
			continue
		}

		// Update DB
		updated, err := u.repo.UpdateDelistingDate(ctx, s.Code, formattedDate)
		if err != nil {
			logrus.Errorf("Failed to update delisting date for %s: %v", s.Code, err)
			continue
		}

		// Enhancement: If stock doesn't exist in DB, fetch from Pasardana and insert
		if updated == nil {
			logrus.Infof("Stock %s not found in DB, fetching from Pasardana...", s.Code)

			detail, err := u.pasardanaService.FetchStockDetailByCode(s.Code)
			if err != nil {
				logrus.Errorf("Failed to fetch detail for %s from Pasardana: %v", s.Code, err)
				continue
			}

			if detail != nil {
				// Normalize dates with fallback to Epoch 0
				epoch0 := "1970-01-01"
				if detail.ListingDate != nil && *detail.ListingDate != "" {
					parsed := utils.NormalizeDate(*detail.ListingDate)
					if parsed == "" {
						parsed = epoch0
					}
					detail.ListingDate = &parsed
				} else {
					detail.ListingDate = &epoch0
				}

				if detail.FoundingDate != nil && *detail.FoundingDate != "" {
					parsed := utils.NormalizeDate(*detail.FoundingDate)
					if parsed == "" {
						parsed = epoch0
					}
					detail.FoundingDate = &parsed
				} else {
					detail.FoundingDate = &epoch0
				}

				// Insert new stock
				_, err = u.repo.UpsertStocksDetail(ctx, []models.PasardanaStockDetail{*detail})
				if err != nil {
					logrus.Errorf("Failed to insert new stock %s: %v", s.Code, err)
					continue
				}

				logrus.Infof("Inserted new delisted stock from Pasardana: %s", s.Code)

				// Retry update delisting date
				updated, err = u.repo.UpdateDelistingDate(ctx, s.Code, formattedDate)
				if err != nil {
					logrus.Errorf("Failed to update delisting date for %s after insert: %v", s.Code, err)
					continue
				}
			}
		}

		if updated != nil {
			updatedStocks = append(updatedStocks, *updated)
		}
	}

	return updatedStocks, nil
}
