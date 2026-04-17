package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestInsertBrokerActivity(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2024, 4, 17, 0, 0, 0, 0, time.UTC)
	record := models.BrokerActivity{
		BrokerCode: "CC",
		StockCode:  "BBCA",
		Date:       date,
		Side:       "BUY",
		Lot:        100,
		Value:      1000000,
		AvgPrice:   10000,
		Freq:       1,
	}

	t.Run("success_insert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()
		repo := NewBrokerActivityRepository(mock)

		mock.ExpectExec("INSERT INTO idxstock.broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date, record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		inserted, err := repo.InsertBrokerActivity(ctx, record)
		assert.NoError(t, err)
		assert.True(t, inserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("already_exists_do_nothing", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()
		repo := NewBrokerActivityRepository(mock)

		mock.ExpectExec("INSERT INTO idxstock.broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date, record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
			WillReturnResult(pgxmock.NewResult("INSERT", 0))

		inserted, err := repo.InsertBrokerActivity(ctx, record)
		assert.NoError(t, err)
		assert.False(t, inserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()
		repo := NewBrokerActivityRepository(mock)

		mock.ExpectExec("INSERT INTO idxstock.broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date, record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
			WillReturnError(assert.AnError)

		inserted, err := repo.InsertBrokerActivity(ctx, record)
		assert.Error(t, err)
		assert.False(t, inserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCheckPartitionExists(t *testing.T) {
	ctx := context.Background()
	tableName := "broker_activity_p_2024_w16"

	t.Run("exists", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()
		repo := NewBrokerActivityRepository(mock)

		rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnRows(rows)

		exists, err := repo.CheckPartitionExists(ctx, tableName)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_exists", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()
		repo := NewBrokerActivityRepository(mock)

		rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnRows(rows)

		exists, err := repo.CheckPartitionExists(ctx, tableName)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
