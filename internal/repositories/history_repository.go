package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type HistoryRepository interface {
	BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error
	GetHistoryByCode(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error)
}

type historyRepository struct {
	pool *pgxpool.Pool
}

func NewHistoryRepository(pool *pgxpool.Pool) HistoryRepository {
	return &historyRepository{
		pool: pool,
	}
}

func (r *historyRepository) BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.history AS h (
			code, date, previous, open_price, first_trade, high, low, close, change,
			volume, value, frequency, index_individual, offer, offer_volume,
			bid, bid_volume, listed_shares, tradeble_shares, weight_for_index,
			foreign_sell, foreign_buy, delisting_date, non_regular_volume,
			non_regular_value, non_regular_frequency, last_modified
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23::DATE, $24, $25, $26, now()
		)
		ON CONFLICT (code, date) DO UPDATE SET
			previous = COALESCE(EXCLUDED.previous, h.previous),
			open_price = COALESCE(EXCLUDED.open_price, h.open_price),
			first_trade = COALESCE(EXCLUDED.first_trade, h.first_trade),
			high = COALESCE(EXCLUDED.high, h.high),
			low = COALESCE(EXCLUDED.low, h.low),
			close = COALESCE(EXCLUDED.close, h.close),
			change = COALESCE(EXCLUDED.change, h.change),
			volume = COALESCE(EXCLUDED.volume, h.volume),
			value = COALESCE(EXCLUDED.value, h.value),
			frequency = COALESCE(EXCLUDED.frequency, h.frequency),
			index_individual = COALESCE(EXCLUDED.index_individual, h.index_individual),
			offer = COALESCE(EXCLUDED.offer, h.offer),
			offer_volume = COALESCE(EXCLUDED.offer_volume, h.offer_volume),
			bid = COALESCE(EXCLUDED.bid, h.bid),
			bid_volume = COALESCE(EXCLUDED.bid_volume, h.bid_volume),
			listed_shares = COALESCE(EXCLUDED.listed_shares, h.listed_shares),
			tradeble_shares = COALESCE(EXCLUDED.tradeble_shares, h.tradeble_shares),
			weight_for_index = COALESCE(EXCLUDED.weight_for_index, h.weight_for_index),
			foreign_sell = COALESCE(EXCLUDED.foreign_sell, h.foreign_sell),
			foreign_buy = COALESCE(EXCLUDED.foreign_buy, h.foreign_buy),
			delisting_date = COALESCE(EXCLUDED.delisting_date, h.delisting_date),
			non_regular_volume = COALESCE(EXCLUDED.non_regular_volume, h.non_regular_volume),
			non_regular_value = COALESCE(EXCLUDED.non_regular_value, h.non_regular_value),
			non_regular_frequency = COALESCE(EXCLUDED.non_regular_frequency, h.non_regular_frequency),
			last_modified = now()
		WHERE
			(EXCLUDED.previous IS NOT NULL AND h.previous IS DISTINCT FROM EXCLUDED.previous) OR
			(EXCLUDED.open_price IS NOT NULL AND h.open_price IS DISTINCT FROM EXCLUDED.open_price) OR
			(EXCLUDED.first_trade IS NOT NULL AND h.first_trade IS DISTINCT FROM EXCLUDED.first_trade) OR
			(EXCLUDED.high IS NOT NULL AND h.high IS DISTINCT FROM EXCLUDED.high) OR
			(EXCLUDED.low IS NOT NULL AND h.low IS DISTINCT FROM EXCLUDED.low) OR
			(EXCLUDED.close IS NOT NULL AND h.close IS DISTINCT FROM EXCLUDED.close) OR
			(EXCLUDED.change IS NOT NULL AND h.change IS DISTINCT FROM EXCLUDED.change) OR
			(EXCLUDED.volume IS NOT NULL AND h.volume IS DISTINCT FROM EXCLUDED.volume) OR
			(EXCLUDED.value IS NOT NULL AND h.value IS DISTINCT FROM EXCLUDED.value) OR
			(EXCLUDED.frequency IS NOT NULL AND h.frequency IS DISTINCT FROM EXCLUDED.frequency) OR
			(EXCLUDED.index_individual IS NOT NULL AND h.index_individual IS DISTINCT FROM EXCLUDED.index_individual) OR
			(EXCLUDED.offer IS NOT NULL AND h.offer IS DISTINCT FROM EXCLUDED.offer) OR
			(EXCLUDED.offer_volume IS NOT NULL AND h.offer_volume IS DISTINCT FROM EXCLUDED.offer_volume) OR
			(EXCLUDED.bid IS NOT NULL AND h.bid IS DISTINCT FROM EXCLUDED.bid) OR
			(EXCLUDED.bid_volume IS NOT NULL AND h.bid_volume IS DISTINCT FROM EXCLUDED.bid_volume) OR
			(EXCLUDED.listed_shares IS NOT NULL AND h.listed_shares IS DISTINCT FROM EXCLUDED.listed_shares) OR
			(EXCLUDED.tradeble_shares IS NOT NULL AND h.tradeble_shares IS DISTINCT FROM EXCLUDED.tradeble_shares) OR
			(EXCLUDED.weight_for_index IS NOT NULL AND h.weight_for_index IS DISTINCT FROM EXCLUDED.weight_for_index) OR
			(EXCLUDED.foreign_sell IS NOT NULL AND h.foreign_sell IS DISTINCT FROM EXCLUDED.foreign_sell) OR
			(EXCLUDED.foreign_buy IS NOT NULL AND h.foreign_buy IS DISTINCT FROM EXCLUDED.foreign_buy) OR
			(EXCLUDED.delisting_date IS NOT NULL AND h.delisting_date IS DISTINCT FROM EXCLUDED.delisting_date) OR
			(EXCLUDED.non_regular_volume IS NOT NULL AND h.non_regular_volume IS DISTINCT FROM EXCLUDED.non_regular_volume) OR
			(EXCLUDED.non_regular_value IS NOT NULL AND h.non_regular_value IS DISTINCT FROM EXCLUDED.non_regular_value) OR
			(EXCLUDED.non_regular_frequency IS NOT NULL AND h.non_regular_frequency IS DISTINCT FROM EXCLUDED.non_regular_frequency)
	`

	batch := &pgx.Batch{}
	for _, rec := range records {
		var dd *string
		if rec.DelistingDate != nil && *rec.DelistingDate != "" {
			dd = rec.DelistingDate
		}

		batch.Queue(query,
			rec.Code, rec.Date, rec.Previous, rec.OpenPrice, rec.FirstTrade,
			rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
			rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
			rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
			rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, dd,
			rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
		)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	var affected int64
	for i := range records {
		ct, err := br.Exec()
		if err != nil {
			// If it fails because of foreign key constraint, we might want to log it and continue
			// but for now let's fail the whole batch to be safe, or just log.
			// Actually history depends on stocks table.
			logrus.Warnf("Failed to upsert history for %s on %s: %v", records[i].Code, records[i].Date.Format("2006-01-02"), err)
			continue
		}
		affected += ct.RowsAffected()
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Batch upsert completed. Affected rows: %d", affected)
	return nil
}
func (r *historyRepository) GetHistoryByCode(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error) {
	query := `
		SELECT
			code, date, previous, open_price, first_trade, high, low, close, change,
			volume, value, frequency, index_individual, offer, offer_volume,
			bid, bid_volume, listed_shares, tradeble_shares, weight_for_index,
			foreign_sell, foreign_buy, delisting_date, non_regular_volume,
			non_regular_value, non_regular_frequency, last_modified
		FROM idxstock.history
		WHERE code = $1
			AND ($2::DATE IS NULL OR date >= $2)
			AND ($3::DATE IS NULL OR date <= $3)
		ORDER BY date ASC
	`

	rows, err := r.pool.Query(ctx, query, code, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %w", code, err)
	}
	defer rows.Close()

	var records []models.StockHistory
	for rows.Next() {
		var h models.StockHistory
		err := rows.Scan(
			&h.Code, &h.Date, &h.Previous, &h.OpenPrice, &h.FirstTrade,
			&h.High, &h.Low, &h.Close, &h.Change, &h.Volume, &h.Value,
			&h.Frequency, &h.IndexIndividual, &h.Offer, &h.OfferVolume,
			&h.Bid, &h.BidVolume, &h.ListedShares, &h.TradebleShares,
			&h.WeightForIndex, &h.ForeignSell, &h.ForeignBuy, &h.DelistingDate,
			&h.NonRegularVolume, &h.NonRegularValue, &h.NonRegularFrequency, &h.LastModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history record: %w", err)
		}
		records = append(records, h)
	}

	return records, nil
}
