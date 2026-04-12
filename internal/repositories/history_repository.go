package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type HistoryRepository interface {
	BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error
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
		INSERT INTO idxstock.history (
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
			previous = EXCLUDED.previous,
			open_price = EXCLUDED.open_price,
			first_trade = EXCLUDED.first_trade,
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			close = EXCLUDED.close,
			change = EXCLUDED.change,
			volume = EXCLUDED.volume,
			value = EXCLUDED.value,
			frequency = EXCLUDED.frequency,
			index_individual = EXCLUDED.index_individual,
			offer = EXCLUDED.offer,
			offer_volume = EXCLUDED.offer_volume,
			bid = EXCLUDED.bid,
			bid_volume = EXCLUDED.bid_volume,
			listed_shares = EXCLUDED.listed_shares,
			tradeble_shares = EXCLUDED.tradeble_shares,
			weight_for_index = EXCLUDED.weight_for_index,
			foreign_sell = EXCLUDED.foreign_sell,
			foreign_buy = EXCLUDED.foreign_buy,
			delisting_date = EXCLUDED.delisting_date,
			non_regular_volume = EXCLUDED.non_regular_volume,
			non_regular_value = EXCLUDED.non_regular_value,
			non_regular_frequency = EXCLUDED.non_regular_frequency,
			last_modified = now()
		WHERE
			history.previous IS DISTINCT FROM EXCLUDED.previous OR
			history.open_price IS DISTINCT FROM EXCLUDED.open_price OR
			history.first_trade IS DISTINCT FROM EXCLUDED.first_trade OR
			history.high IS DISTINCT FROM EXCLUDED.high OR
			history.low IS DISTINCT FROM EXCLUDED.low OR
			history.close IS DISTINCT FROM EXCLUDED.close OR
			history.change IS DISTINCT FROM EXCLUDED.change OR
			history.volume IS DISTINCT FROM EXCLUDED.volume OR
			history.value IS DISTINCT FROM EXCLUDED.value OR
			history.frequency IS DISTINCT FROM EXCLUDED.frequency OR
			history.index_individual IS DISTINCT FROM EXCLUDED.index_individual OR
			history.offer IS DISTINCT FROM EXCLUDED.offer OR
			history.offer_volume IS DISTINCT FROM EXCLUDED.offer_volume OR
			history.bid IS DISTINCT FROM EXCLUDED.bid OR
			history.bid_volume IS DISTINCT FROM EXCLUDED.bid_volume OR
			history.listed_shares IS DISTINCT FROM EXCLUDED.listed_shares OR
			history.tradeble_shares IS DISTINCT FROM EXCLUDED.tradeble_shares OR
			history.weight_for_index IS DISTINCT FROM EXCLUDED.weight_for_index OR
			history.foreign_sell IS DISTINCT FROM EXCLUDED.foreign_sell OR
			history.foreign_buy IS DISTINCT FROM EXCLUDED.foreign_buy OR
			history.delisting_date IS DISTINCT FROM EXCLUDED.delisting_date OR
			history.non_regular_volume IS DISTINCT FROM EXCLUDED.non_regular_volume OR
			history.non_regular_value IS DISTINCT FROM EXCLUDED.non_regular_value OR
			history.non_regular_frequency IS DISTINCT FROM EXCLUDED.non_regular_frequency
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
	for i := 0; i < len(records); i++ {
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
