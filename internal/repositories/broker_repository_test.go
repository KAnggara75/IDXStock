package repositories

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/models"
)

func TestBrokerRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("nil_pool", func(t *testing.T) {
		repo := NewBrokerRepository(nil)
		err := repo.BatchInsertBrokers(ctx, []models.Broker{{Code: "CC"}})
		assert.Error(t, err)
		assert.Equal(t, "database pool is not initialized", err.Error())
	})

	t.Run("empty_brokers", func(t *testing.T) {

		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewBrokerRepository(mock)
		err = repo.BatchInsertBrokers(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("success_batch_insert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		b := models.Broker{Code: "CC", Name: "Mandiri"}
		batch.ExpectExec("INSERT INTO idxstock.brokers").
			WithArgs(b.Code, b.Name, b.InvestorType, b.TotalValue, b.NetValue, b.BuyValue, b.SellValue, b.TotalVolume, b.TotalFrequency, b.Group).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectCommit()

		repo := NewBrokerRepository(mock)
		err = repo.BatchInsertBrokers(ctx, []models.Broker{b})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
