package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type BrokerActivityRepository interface {
	BatchInsertBrokerActivity(ctx context.Context, records []models.BrokerActivity) error
	InsertBrokerActivity(ctx context.Context, record models.BrokerActivity) (bool, error)
	CheckPartitionExists(ctx context.Context, tableName string) (bool, error)
	CreatePartition(ctx context.Context, tableName, startDate, endDate string) error
}

type brokerActivityRepository struct {
	pool *pgxpool.Pool
}

func NewBrokerActivityRepository(pool *pgxpool.Pool) BrokerActivityRepository {
	return &brokerActivityRepository{
		pool: pool,
	}
}

func (r *brokerActivityRepository) BatchInsertBrokerActivity(ctx context.Context, records []models.BrokerActivity) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.broker_activity (
			broker_code, stock_code, date, side, lot, value, avg_price, freq, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
		ON CONFLICT (broker_code, stock_code, date, side) DO UPDATE SET
			lot = EXCLUDED.lot,
			value = EXCLUDED.value,
			avg_price = EXCLUDED.avg_price,
			freq = EXCLUDED.freq,
			updated_at = now()
		WHERE
			idxstock.broker_activity.lot IS DISTINCT FROM EXCLUDED.lot OR
			idxstock.broker_activity.value IS DISTINCT FROM EXCLUDED.value OR
			idxstock.broker_activity.avg_price IS DISTINCT FROM EXCLUDED.avg_price OR
			idxstock.broker_activity.freq IS DISTINCT FROM EXCLUDED.freq
	`

	batch := &pgx.Batch{}
	for _, rec := range records {
		batch.Queue(query,
			rec.BrokerCode, rec.StockCode, rec.Date, rec.Side,
			rec.Lot, rec.Value, rec.AvgPrice, rec.Freq,
		)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	var affected int64
	for i := 0; i < len(records); i++ {
		ct, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch upsert for broker activity at index %d: %w", i, err)
		}
		affected += ct.RowsAffected()
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Successfully upserted %d broker activity records (affected: %d)", len(records), affected)
	return nil
}

func (r *brokerActivityRepository) InsertBrokerActivity(ctx context.Context, rec models.BrokerActivity) (bool, error) {
	query := `
		INSERT INTO idxstock.broker_activity (
			broker_code, stock_code, date, side, lot, value, avg_price, freq, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
		ON CONFLICT (broker_code, stock_code, date, side) DO NOTHING
	`

	ct, err := r.pool.Exec(ctx, query,
		rec.BrokerCode, rec.StockCode, rec.Date, rec.Side,
		rec.Lot, rec.Value, rec.AvgPrice, rec.Freq,
	)
	if err != nil {
		return false, fmt.Errorf("failed to insert broker activity (%s, %s, %s, %s): %w",
			rec.BrokerCode, rec.StockCode, rec.Date.Format("2006-01-02"), rec.Side, err)
	}

	return ct.RowsAffected() > 0, nil
}

func (r *brokerActivityRepository) CheckPartitionExists(ctx context.Context, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'idxstock' AND table_name = $1
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check partition existence for %s: %w", tableName, err)
	}
	return exists, nil
}

func (r *brokerActivityRepository) CreatePartition(ctx context.Context, tableName, startDate, endDate string) error {
	// DDL statements cannot use parameters for table names, so we use fmt.Sprintf
	// We assume tableName is internally generated and safe.
	query := fmt.Sprintf(`
		CREATE TABLE idxstock.%s PARTITION OF idxstock.broker_activity
		FOR VALUES FROM ('%s') TO ('%s')
	`, tableName, startDate, endDate)

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create partition %s [%s to %s]: %w", tableName, startDate, endDate, err)
	}

	return nil
}
