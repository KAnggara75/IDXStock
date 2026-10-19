package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetDsn(t *testing.T) {
	t.Run("dsn_default", func(t *testing.T) {
		viper.Set("db.idxstock.dsn", "")
		os.Unsetenv("DATABASE_URL")
		dsn, err := getDsn()
		assert.NoError(t, err)
		assert.Equal(t, "database/master_stocks.db", dsn)
	})

	t.Run("dsn_env_set", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "test.db")
		defer os.Unsetenv("DATABASE_URL")
		dsn, err := getDsn()
		assert.NoError(t, err)
		assert.Equal(t, "test.db", dsn)
	})

	t.Run("dsn_viper_set", func(t *testing.T) {
		viper.Set("db.idxstock.dsn", "custom.db")
		defer viper.Set("db.idxstock.dsn", "")
		dsn, err := getDsn()
		assert.NoError(t, err)
		assert.Equal(t, "custom.db", dsn)
	})
}

func TestConnectWithConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("success_in_memory", func(t *testing.T) {
		pool, err := connectWithConfig(ctx, ":memory:")
		assert.NoError(t, err)
		assert.NotNil(t, pool)
		assert.NoError(t, pool.Ping(ctx))
		assert.NoError(t, pool.Close())
	})

	t.Run("with_nested_directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "sqlite_test_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "subdir", "test.db")
		pool, err := connectWithConfig(ctx, dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, pool)
		assert.NoError(t, pool.Ping(ctx))
		assert.NoError(t, pool.Close())
	})
}

func TestConnectDB(t *testing.T) {
	t.Run("dsn_getter_error", func(t *testing.T) {
		pool, err := connectDB(func() (string, error) {
			return "", assert.AnError
		}, nil)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	t.Run("connector_error", func(t *testing.T) {
		pool, err := connectDB(func() (string, error) {
			return "idxstock.db", nil
		}, func(ctx context.Context, dsn string) (DBExecutor, error) {
			return nil, assert.AnError
		})
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	t.Run("connector_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := NewDBExecutor(db)

		pool, err := connectDB(func() (string, error) {
			return "idxstock.db", nil
		}, func(ctx context.Context, dsn string) (DBExecutor, error) {
			return mockExec, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, mockExec, pool)
	})
}

func TestDBExecutorAndTx(t *testing.T) {
	ctx := context.Background()

	t.Run("sqlDB_methods", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer db.Close()

		exec := NewDBExecutor(db)

		// Exec
		mock.ExpectExec("UPDATE test").WillReturnResult(sqlmock.NewResult(1, 2))
		res, err := exec.Exec(ctx, "UPDATE test")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.RowsAffected())

		// Query
		mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"val"}).AddRow(1))
		rows, err := exec.Query(ctx, "SELECT 1")
		assert.NoError(t, err)
		assert.True(t, rows.Next())
		rows.Close()

		// QueryRow
		mock.ExpectQuery("SELECT 2").WillReturnRows(sqlmock.NewRows([]string{"val"}).AddRow(2))
		var val int
		err = exec.QueryRow(ctx, "SELECT 2").Scan(&val)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)

		// Ping
		mock.ExpectPing()
		assert.NoError(t, exec.Ping(ctx))

		// Begin error
		mock.ExpectBegin().WillReturnError(assert.AnError)
		_, err = exec.Begin(ctx)
		assert.Error(t, err)

		// Begin success
		mock.ExpectBegin()
		tx, err := exec.Begin(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, tx)

		// Tx Exec
		mock.ExpectExec("INSERT INTO test").WillReturnResult(sqlmock.NewResult(1, 1))
		txRes, err := tx.Exec(ctx, "INSERT INTO test")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), txRes.RowsAffected())

		// Tx Query
		mock.ExpectQuery("SELECT 3").WillReturnRows(sqlmock.NewRows([]string{"val"}).AddRow(3))
		txRows, err := tx.Query(ctx, "SELECT 3")
		assert.NoError(t, err)
		assert.True(t, txRows.Next())
		txRows.Close()

		// Tx QueryRow
		mock.ExpectQuery("SELECT 4").WillReturnRows(sqlmock.NewRows([]string{"val"}).AddRow(4))
		var txVal int
		err = tx.QueryRow(ctx, "SELECT 4").Scan(&txVal)
		assert.NoError(t, err)
		assert.Equal(t, 4, txVal)

		// Tx Commit
		mock.ExpectCommit()
		assert.NoError(t, tx.Commit(ctx))

		// Tx Rollback when already committed
		assert.NoError(t, tx.Rollback(ctx))

		// Close
		mock.ExpectClose()
		assert.NoError(t, exec.Close())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
