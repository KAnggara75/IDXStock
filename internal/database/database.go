package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DBExecutor is an interface that matches the methods we use from pgxpool.Pool
type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
}

var Pool *pgxpool.Pool

func getDsn() (string, error) {
	dsn := viper.GetString("db.idxstock.dsn")
	if dsn == "" {
		return "", fmt.Errorf("DSN environment variable is not set")
	}
	logrus.Debug("DSN environment variable is set")
	return dsn, nil
}

// connectWithConfig handles creating pool, pinging, and running migrations.
func connectWithConfig(ctx context.Context, config *pgxpool.Config, migrateFn func(executor ...DBExecutor)) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	if migrateFn != nil {
		migrateFn(pool)
	}

	return pool, nil
}

func connectDB(dsnGetter func() (string, error), connector func(ctx context.Context, config *pgxpool.Config, migrateFn func(executor ...DBExecutor)) (*pgxpool.Pool, error)) (*pgxpool.Pool, error) {
	logrus.Debug("Connecting to the database")
	dsn, err := dsnGetter()
	if err != nil {
		return nil, err
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}
	logrus.Debugf("Connecting to database at %s:%d", config.ConnConfig.Host, config.ConnConfig.Port)

	pool, err := connector(context.Background(), config, Migrate)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Connection pool created and pinged successfully")
	fmt.Println("Successfully connected to the database")
	return pool, nil
}

func Connect() {
	pool, err := connectDB(getDsn, connectWithConfig)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	Pool = pool
}

func Migrate(executor ...DBExecutor) {
	var exec DBExecutor
	if len(executor) > 0 && executor[0] != nil {
		exec = executor[0]
	} else if Pool != nil {
		exec = Pool
	}

	if exec == nil {
		logrus.Warn("Database pool is nil, skipping migrations")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Ensure schema exists
	_, err := exec.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS idxstock")
	if err != nil {
		logrus.Errorf("Failed to create schema idxstock: %v", err)
		return
	}

	// 2. Create migrations tracking table
	_, err = exec.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS idxstock.schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT now()
		)
	`)
	if err != nil {
		logrus.Errorf("Failed to create migration tracking table: %v", err)
		return
	}

	// 3. Handle transition from old skip logic:
	// If stocks table exists but migrations haven't been tracked, mark 000 and 001 as applied
	var stocksExists bool
	err = exec.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'idxstock' AND table_name = 'stocks')").Scan(&stocksExists)
	if err != nil {
		logrus.Debugf("Note: could not check if stocks table exists: %v", err)
	}

	if stocksExists {
		_, err = exec.Exec(ctx, "INSERT INTO idxstock.schema_migrations (filename) VALUES ('migrations/000_board_type.sql'), ('migrations/001_create_stocks_table.sql') ON CONFLICT DO NOTHING")
		if err != nil {
			logrus.Debugf("Note: could not seed initial migrations: %v", err)
		}
	}

	// 4. Get migration files
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil || len(files) == 0 {
		files, err = filepath.Glob("../../migrations/*.sql")
	}
	if err != nil {
		logrus.Errorf("Failed to read migrations directory: %v", err)
		return
	}

	sort.Strings(files)

	for _, file := range files {
		// Check if already applied
		var applied bool
		err = exec.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM idxstock.schema_migrations WHERE filename = $1)", file).Scan(&applied)
		if err != nil {
			logrus.Errorf("Failed to check migration status for %s: %v", file, err)
			continue
		}

		if applied {
			continue
		}

		logrus.Infof("Running migration: %s", file)
		cleanFile := filepath.Clean(file)
		// #nosec G304
		content, err := os.ReadFile(cleanFile)
		if err != nil {
			logrus.Errorf("Failed to read migration file %s: %v", file, err)
			continue
		}

		_, err = exec.Exec(ctx, string(content))
		if err != nil {
			logrus.Errorf("Failed to execute migration %s: %v", file, err)
			continue
		}

		// Mark as applied
		_, err = exec.Exec(ctx, "INSERT INTO idxstock.schema_migrations (filename) VALUES ($1)", file)
		if err != nil {
			logrus.Errorf("Failed to record migration %s: %v", file, err)
			continue
		}
	}

	logrus.Info("Migrations synchronization completed")
}
