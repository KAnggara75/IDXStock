package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/models"
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
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectExec("INSERT INTO broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date.Format("2006-01-02"), record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
			WillReturnResult(sqlmock.NewResult(1, 1))

		inserted, err := repo.InsertBrokerActivity(ctx, record)
		assert.NoError(t, err)
		assert.True(t, inserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("already_exists_do_nothing", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectExec("INSERT INTO broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date.Format("2006-01-02"), record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
			WillReturnResult(sqlmock.NewResult(0, 0))

		inserted, err := repo.InsertBrokerActivity(ctx, record)
		assert.NoError(t, err)
		assert.False(t, inserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectExec("INSERT INTO broker_activity").
			WithArgs(record.BrokerCode, record.StockCode, record.Date.Format("2006-01-02"), record.Side, record.Lot, record.Value, record.AvgPrice, record.Freq).
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
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnRows(rows)

		exists, err := repo.CheckPartitionExists(ctx, tableName)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnRows(rows)

		exists, err := repo.CheckPartitionExists(ctx, tableName)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnError(assert.AnError)

		exists, err := repo.CheckPartitionExists(ctx, tableName)
		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreatePartition(t *testing.T) {
	ctx := context.Background()
	tableName := "broker_activity_p_2024_w16"
	startDate := "2024-04-15"
	endDate := "2024-04-22"

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableName).WillReturnResult(sqlmock.NewResult(1, 0))

		err = repo.CreatePartition(ctx, tableName, startDate, endDate)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableName).WillReturnError(assert.AnError)

		err = repo.CreatePartition(ctx, tableName, startDate, endDate)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBatchInsertBrokerActivity(t *testing.T) {
	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		err = repo.BatchInsertBrokerActivity(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		mock.ExpectBegin().WillReturnError(assert.AnError)

		err = repo.BatchInsertBrokerActivity(ctx, []models.BrokerActivity{{BrokerCode: "CC"}})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("exec_statement_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		date := time.Date(2024, 4, 17, 0, 0, 0, 0, time.UTC)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO broker_activity").
			WithArgs("CC", "BBCA", date.Format("2006-01-02"), "BUY", int64(100), int64(1000), 10.0, int64(1)).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = repo.BatchInsertBrokerActivity(ctx, []models.BrokerActivity{
			{BrokerCode: "CC", StockCode: "BBCA", Date: date, Side: "BUY", Lot: 100, Value: 1000, AvgPrice: 10.0, Freq: 1},
		})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)
		repo := NewBrokerActivityRepository(mockExec)

		date := time.Date(2024, 4, 17, 0, 0, 0, 0, time.UTC)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO broker_activity").
			WithArgs("CC", "BBCA", date.Format("2006-01-02"), "BUY", int64(100), int64(1000), 10.0, int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.BatchInsertBrokerActivity(ctx, []models.BrokerActivity{
			{BrokerCode: "CC", StockCode: "BBCA", Date: date, Side: "BUY", Lot: 100, Value: 1000, AvgPrice: 10.0, Freq: 1},
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
