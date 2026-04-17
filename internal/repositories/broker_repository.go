package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type BrokerRepository interface {
	IsBrokersTableEmpty(ctx context.Context) (bool, error)
	BatchInsertBrokers(ctx context.Context, brokers []models.Broker) error
}

type brokerRepository struct {
	pool *pgxpool.Pool
}

func NewBrokerRepository(pool *pgxpool.Pool) BrokerRepository {
	return &brokerRepository{
		pool: pool,
	}
}

func (r *brokerRepository) IsBrokersTableEmpty(ctx context.Context) (bool, error) {
	var count int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM brokers").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if brokers table is empty: %w", err)
	}
	return count == 0, nil
}

func (r *brokerRepository) BatchInsertBrokers(ctx context.Context, brokers []models.Broker) error {
	if len(brokers) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO brokers (
			code, name, investor_type, total_value, net_value, buy_value,
			sell_value, total_volume, total_frequency, broker_group
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (code) DO NOTHING
	`

	batch := &pgx.Batch{}
	for _, b := range brokers {
		batch.Queue(query, b.Code, b.Name, b.InvestorType, b.TotalValue, b.NetValue, b.BuyValue, b.SellValue, b.TotalVolume, b.TotalFrequency, b.Group)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(brokers); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert for broker %s: %w", brokers[i].Code, err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Successfully seeded %d brokers", len(brokers))
	return nil
}
