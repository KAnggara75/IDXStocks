package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			listing_date = EXCLUDED.listing_date,
			delisting_date = EXCLUDED.delisting_date,
			shares = EXCLUDED.shares,
			board = EXCLUDED.board,
			last_modified = now()
	`

	batch := &pgx.Batch{}
	for _, s := range stocks {
		batch.Queue(query, s.Code, s.CompanyName, s.ListingDate, s.DelistingDate, s.Shares, s.ListingBoard)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
