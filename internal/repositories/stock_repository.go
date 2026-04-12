package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type StockRepository interface {
	BatchInsertStocks(ctx context.Context, stocks []models.Stock) error
	UpdateStockIDs(ctx context.Context, data []models.PasardanaStock) ([]models.StockResponse, error)
	UpsertStocksDetail(ctx context.Context, data []models.PasardanaStockDetail) ([]models.StockResponse, error)
	UpdateDelistingDate(ctx context.Context, code, delistingDate string) (*models.StockResponse, error)
}

type stockRepository struct {
	pool *pgxpool.Pool
}

func NewStockRepository(pool *pgxpool.Pool) StockRepository {
	return &stockRepository{
		pool: pool,
	}
}

func (r *stockRepository) BatchInsertStocks(ctx context.Context, stocks []models.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.stocks (code, name, listing_date, delisting_date, shares, board)
		VALUES ($1, $2, $3::DATE, NULLIF($4, '')::DATE, $5, $6)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			listing_date = EXCLUDED.listing_date,
			delisting_date = EXCLUDED.delisting_date,
			shares = EXCLUDED.shares,
			board = EXCLUDED.board,
			last_modified = now()
		WHERE
			stocks.name IS DISTINCT FROM EXCLUDED.name OR
			stocks.listing_date IS DISTINCT FROM EXCLUDED.listing_date OR
			stocks.delisting_date IS DISTINCT FROM EXCLUDED.delisting_date OR
			stocks.shares IS DISTINCT FROM EXCLUDED.shares OR
			stocks.board IS DISTINCT FROM EXCLUDED.board
	`

	batch := &pgx.Batch{}
	for _, s := range stocks {
		batch.Queue(query, s.Code, s.CompanyName, s.ListingDate, s.DelistingDate, s.Shares, s.ListingBoard)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	var changedCount int64
	for i := range stocks {
		cmdTag, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch statement %d: %w", i, err)
		}

		if cmdTag.RowsAffected() > 0 {
			changedCount++
			logrus.Debugf("Data changed for stock: %s (%s)", stocks[i].Code, stocks[i].CompanyName)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully processed %d stocks, %d records updated/inserted", len(stocks), changedCount)

	return nil
}

func (r *stockRepository) UpdateStockIDs(ctx context.Context, data []models.PasardanaStock) ([]models.StockResponse, error) {
	if len(data) == 0 {
		return make([]models.StockResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE idxstock.stocks
		SET id = $1, last_modified = now()
		WHERE code = $2 AND (id IS NULL OR id != $1)
		RETURNING id, code, name
	`

	batch := &pgx.Batch{}
	for _, s := range data {
		batch.Queue(query, s.Id, s.Code)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updatedStocks := make([]models.StockResponse, 0)
	for i := range data {
		rows, err := br.Query()
		if err != nil {
			return nil, fmt.Errorf("failed to execute batch query %d: %w", i, err)
		}

		for rows.Next() {
			var sr models.StockResponse
			if err := rows.Scan(&sr.Id, &sr.Code, &sr.Name); err == nil {
				updatedStocks = append(updatedStocks, sr)
			}
		}
		rows.Close()
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully updated %d stock IDs", len(updatedStocks))

	return updatedStocks, nil
}

func (r *stockRepository) UpdateDelistingDate(ctx context.Context, code, delistingDate string) (*models.StockResponse, error) {
	query := `
		UPDATE idxstock.stocks
		SET delisting_date = $1::DATE, last_modified = now()
		WHERE code = $2
		RETURNING id, code, name
	`

	var sr models.StockResponse
	err := r.pool.QueryRow(ctx, query, delistingDate, code).Scan(&sr.Id, &sr.Code, &sr.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No stock found with this code
		}
		return nil, fmt.Errorf("failed to update delisting date for %s: %w", code, err)
	}

	return &sr, nil
}

func (r *stockRepository) UpsertStocksDetail(ctx context.Context, data []models.PasardanaStockDetail) ([]models.StockResponse, error) {
	if len(data) == 0 {
		return make([]models.StockResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.stocks (
			id, code, name, listing_date, total_employees, annual_dividend,
			general_information, founding_date, sector_id, sub_sector_id,
			industry_id, sub_industry_id, last_modified
		)
		VALUES ($1, $2, $3, $4::DATE, $5, $6, $7, $8, $9, $10, $11, $12, now())
		ON CONFLICT (code) DO UPDATE SET
			id = EXCLUDED.id,
			name = EXCLUDED.name,
			listing_date = EXCLUDED.listing_date,
			total_employees = EXCLUDED.total_employees,
			annual_dividend = EXCLUDED.annual_dividend,
			general_information = EXCLUDED.general_information,
			founding_date = EXCLUDED.founding_date,
			sector_id = EXCLUDED.sector_id,
			sub_sector_id = EXCLUDED.sub_sector_id,
			industry_id = EXCLUDED.industry_id,
			sub_industry_id = EXCLUDED.sub_industry_id,
			last_modified = now()
		WHERE
			stocks.id IS DISTINCT FROM EXCLUDED.id OR
			stocks.name IS DISTINCT FROM EXCLUDED.name OR
			stocks.listing_date IS DISTINCT FROM EXCLUDED.listing_date OR
			stocks.total_employees IS DISTINCT FROM EXCLUDED.total_employees OR
			stocks.annual_dividend IS DISTINCT FROM EXCLUDED.annual_dividend OR
			stocks.general_information IS DISTINCT FROM EXCLUDED.general_information OR
			stocks.founding_date IS DISTINCT FROM EXCLUDED.founding_date OR
			stocks.sector_id IS DISTINCT FROM EXCLUDED.sector_id OR
			stocks.sub_sector_id IS DISTINCT FROM EXCLUDED.sub_sector_id OR
			stocks.industry_id IS DISTINCT FROM EXCLUDED.industry_id OR
			stocks.sub_industry_id IS DISTINCT FROM EXCLUDED.sub_industry_id
		RETURNING id, code, name
	`

	batch := &pgx.Batch{}
	for _, s := range data {
		batch.Queue(query,
			s.Id, s.Code, s.Name, s.ListingDate, s.TotalEmployees, s.AnnualDividend,
			s.GeneralInformation, s.FoundingDate, s.FkNewSectorId, s.FkNewSubSectorId,
			s.FkNewIndustryId, s.FkNewSubIndustryId,
		)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updatedStocks := make([]models.StockResponse, 0)
	for i, s := range data {
		rows, err := br.Query()
		if err != nil {
			logrus.Errorf("[%s] Sync failed: %v", s.Code, err)
			continue
		}

		hasChange := false
		for rows.Next() {
			var sr models.StockResponse
			if err := rows.Scan(&sr.Id, &sr.Code, &sr.Name); err == nil {
				updatedStocks = append(updatedStocks, sr)
				hasChange = true
			}
		}
		rows.Close()

		if hasChange {
			logrus.Infof("[%d/%d] [%s]: Data Change Detected - Synced", i+1, len(data), s.Code)
		} else {
			logrus.Debugf("[%d/%d] [%s]: No Change Detected", i+1, len(data), s.Code)
		}
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedStocks, nil
}
