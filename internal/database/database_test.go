package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetDsn(t *testing.T) {
	t.Run("dsn_not_set", func(t *testing.T) {
		viper.Set("db.idxstock.dsn", "")
		dsn, err := getDsn()
		assert.Error(t, err)
		assert.Empty(t, dsn)
	})

	t.Run("dsn_set", func(t *testing.T) {
		viper.Set("db.idxstock.dsn", "postgres://user:pass@localhost:5432/dbname")
		defer viper.Set("db.idxstock.dsn", "")
		dsn, err := getDsn()
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", dsn)
	})
}

func TestConnectWithConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("ping_failure_closes_pool", func(t *testing.T) {
		config, err := pgxpool.ParseConfig("postgres://invalid:invalid@127.0.0.1:54321/nonexistent?connect_timeout=1")
		assert.NoError(t, err)

		pool, err := connectWithConfig(ctx, config, nil)
		assert.Error(t, err)
		assert.Nil(t, pool)
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

	t.Run("invalid_dsn_format", func(t *testing.T) {
		pool, err := connectDB(func() (string, error) {
			return "invalid://", nil
		}, nil)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	t.Run("connector_error", func(t *testing.T) {
		pool, err := connectDB(func() (string, error) {
			return "postgres://user:pass@localhost:5432/dbname", nil
		}, func(ctx context.Context, config *pgxpool.Config, migrateFn func(executor ...DBExecutor)) (*pgxpool.Pool, error) {
			return nil, assert.AnError
		})
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	t.Run("connector_success", func(t *testing.T) {
		var migrated bool
		fakePool := &pgxpool.Pool{}
		pool, err := connectDB(func() (string, error) {
			return "postgres://user:pass@localhost:5432/dbname", nil
		}, func(ctx context.Context, config *pgxpool.Config, migrateFn func(executor ...DBExecutor)) (*pgxpool.Pool, error) {
			if migrateFn != nil {
				migrated = true
			}
			return fakePool, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, fakePool, pool)
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

	t.Run("schema_creation_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnError(assert.AnError)

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("migration_table_creation_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS idxstock.schema_migrations").
			WillReturnError(assert.AnError)

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("stocks_table_exists_query_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS idxstock.schema_migrations").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnError(assert.AnError)

		for i := 0; i < 7; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
				WithArgs(pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("stocks_table_exists_seeds_migrations_and_initial_insert_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS idxstock.schema_migrations").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))

		// stocks table exists = true
		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectExec("INSERT INTO idxstock.schema_migrations").
			WillReturnError(assert.AnError)

		for i := 0; i < 7; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
				WithArgs(pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("applies_all_migrations_successfully", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS idxstock.schema_migrations").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		// 7 migration files in migrations/*.sql
		for i := 0; i < 7; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
				WithArgs(pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(".*").
				WillReturnResult(pgxmock.NewResult("CREATE", 1))

			mock.ExpectExec("INSERT INTO idxstock.schema_migrations").
				WithArgs(pgxmock.AnyArg()).
				WillReturnResult(pgxmock.NewResult("INSERT", 1))
		}

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("migration_exec_error_and_record_error_handled", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectExec("CREATE SCHEMA IF NOT EXISTS idxstock").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS idxstock.schema_migrations").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))

		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		// 1st file: check status error
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
			WithArgs(pgxmock.AnyArg()).
			WillReturnError(assert.AnError)

		// 2nd file: exec error
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec(".*").
			WillReturnError(assert.AnError)

		// 3rd file: record error
		mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec(".*").
			WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec("INSERT INTO idxstock.schema_migrations").
			WithArgs(pgxmock.AnyArg()).
			WillReturnError(assert.AnError)

		// Remaining 4 files: already applied
		for i := 0; i < 4; i++ {
			mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM idxstock.schema_migrations WHERE filename = \\$1\\)").
				WithArgs(pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
		}

		Migrate(mock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
