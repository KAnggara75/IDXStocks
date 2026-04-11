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
	for i := 0; i < len(stocks); i++ {
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
