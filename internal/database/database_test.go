package database

import (
	"context"
	"os"
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
		assert.Equal(t, "idxstock.db", dsn)
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
		pool, err := connectWithConfig(ctx, ":memory:", nil)
		assert.NoError(t, err)
		assert.NotNil(t, pool)
		assert.NoError(t, pool.Ping(ctx))
		assert.NoError(t, pool.Close())
	})

	t.Run("with_migrate_fn", func(t *testing.T) {
		var migrated bool
		pool, err := connectWithConfig(ctx, ":memory:", func(executor ...DBExecutor) {
			migrated = true
		})
		assert.NoError(t, err)
		assert.True(t, migrated)
		assert.NotNil(t, pool)
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
		}, func(ctx context.Context, dsn string, migrateFn func(executor ...DBExecutor)) (DBExecutor, error) {
			return nil, assert.AnError
		})
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	t.Run("connector_success", func(t *testing.T) {
		var migrated bool
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := NewDBExecutor(db)

		pool, err := connectDB(func() (string, error) {
			return "idxstock.db", nil
		}, func(ctx context.Context, dsn string, migrateFn func(executor ...DBExecutor)) (DBExecutor, error) {
			if migrateFn != nil {
				migrated = true
			}
			return mockExec, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, mockExec, pool)
		assert.True(t, migrated)
	})
}

func TestMigrate(t *testing.T) {
	t.Run("nil_executor_and_nil_pool", func(t *testing.T) {
		oldPool := Pool
		Pool = nil
		defer func() { Pool = oldPool }()

		assert.NotPanics(t, func() {
			Migrate()
		})
	})

	t.Run("migration_table_creation_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		exec := NewDBExecutor(db)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
			WillReturnError(assert.AnError)

		Migrate(exec)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("stocks_table_exists_query_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		exec := NewDBExecutor(db)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnError(assert.AnError)

		for i := 0; i < 7; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
				WithArgs(sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(exec)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("stocks_table_exists_seeds_migrations_and_initial_insert_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		exec := NewDBExecutor(db)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectExec("INSERT OR IGNORE INTO schema_migrations").
			WillReturnError(assert.AnError)

		for i := 0; i < 7; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
				WithArgs(sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(exec)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("applies_all_migrations_successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		exec := NewDBExecutor(db)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// 7 migration files with their respective statement counts:
		// 000: 1 stmt
		// 001: 2 stmts
		// 002: 2 stmts
		// 003: 4 stmts
		// 004: 4 stmts
		// 005: 1 stmt
		// 006: 3 stmts
		stmtCounts := []int{1, 2, 2, 4, 4, 1, 3}

		for _, count := range stmtCounts {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
				WithArgs(sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			for s := 0; s < count; s++ {
				mock.ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			mock.ExpectExec("INSERT INTO schema_migrations").
				WithArgs(sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		Migrate(exec)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("migration_exec_error_and_record_error_handled", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		exec := NewDBExecutor(db)

		mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// 1st file (000): check status error
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
			WithArgs(sqlmock.AnyArg()).
			WillReturnError(assert.AnError)

		// 2nd file (001): exec error on 1st statement
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec(".*").
			WillReturnError(assert.AnError)

		// 3rd file (002): 2 statements success, record error
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec(".*").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(".*").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO schema_migrations").
			WithArgs(sqlmock.AnyArg()).
			WillReturnError(assert.AnError)

		// Remaining 4 files (003, 004, 005, 006): already applied
		for i := 0; i < 4; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM schema_migrations WHERE filename = \\?1\\)").
				WithArgs(sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(exec)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
