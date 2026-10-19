package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/database"
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
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewBrokerRepository(mockExec)
		err = repo.BatchInsertBrokers(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewBrokerRepository(mockExec)
		err = repo.BatchInsertBrokers(ctx, []models.Broker{{Code: "CC"}})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("exec_statement_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		b := models.Broker{Code: "CC", Name: "Mandiri"}
		mock.ExpectExec("INSERT INTO brokers").
			WithArgs(b.Code, b.Name, b.InvestorType, b.TotalValue, b.NetValue, b.BuyValue, b.SellValue, b.TotalVolume, b.TotalFrequency, b.Group).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		repo := NewBrokerRepository(mockExec)
		err = repo.BatchInsertBrokers(ctx, []models.Broker{b})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_batch_insert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		b := models.Broker{Code: "CC", Name: "Mandiri"}
		mock.ExpectExec("INSERT INTO brokers").
			WithArgs(b.Code, b.Name, b.InvestorType, b.TotalValue, b.NetValue, b.BuyValue, b.SellValue, b.TotalVolume, b.TotalFrequency, b.Group).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewBrokerRepository(mockExec)
		err = repo.BatchInsertBrokers(ctx, []models.Broker{b})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
