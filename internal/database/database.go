package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

// CommandResult abstracts sql.Result for rows affected checks
type CommandResult interface {
	RowsAffected() int64
}

type commandResult struct {
	rowsAffected int64
}

func (c commandResult) RowsAffected() int64 {
	return c.rowsAffected
}

// DBExecutor is an interface that matches methods for SQL execution and transactions
type DBExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (CommandResult, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
	Begin(ctx context.Context) (TxExecutor, error)
	Ping(ctx context.Context) error
	Close() error
}

// TxExecutor is an interface that abstracts SQL transaction operations
type TxExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (CommandResult, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type sqlDB struct {
	db *sql.DB
}

func NewDBExecutor(db *sql.DB) DBExecutor {
	return &sqlDB{db: db}
}

func (s *sqlDB) Exec(ctx context.Context, query string, args ...any) (CommandResult, error) {
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		ra = 0
	}
	return commandResult{rowsAffected: ra}, nil
}

func (s *sqlDB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *sqlDB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *sqlDB) Begin(ctx context.Context) (TxExecutor, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx: tx}, nil
}

func (s *sqlDB) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *sqlDB) Close() error {
	return s.db.Close()
}

type sqlTx struct {
	tx *sql.Tx
}

func (s *sqlTx) Exec(ctx context.Context, query string, args ...any) (CommandResult, error) {
	res, err := s.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		ra = 0
	}
	return commandResult{rowsAffected: ra}, nil
}

func (s *sqlTx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.tx.QueryContext(ctx, query, args...)
}

func (s *sqlTx) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.tx.QueryRowContext(ctx, query, args...)
}

func (s *sqlTx) Commit(ctx context.Context) error {
	return s.tx.Commit()
}

func (s *sqlTx) Rollback(ctx context.Context) error {
	err := s.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

var Pool DBExecutor

func getDsn() (string, error) {
	dsn := viper.GetString("db.idxstock.dsn")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = "database/master_stocks.db"
	}
	logrus.Debug("SQLite DSN is set")
	return dsn, nil
}

func connectWithConfig(ctx context.Context, dsn string, migrateFn func(executor ...DBExecutor)) (DBExecutor, error) {
	if dir := filepath.Dir(dsn); dir != "" && dir != "." && !strings.HasPrefix(dsn, ":memory:") {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return nil, fmt.Errorf("unable to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Recommended connection pooling for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	executor := NewDBExecutor(db)
	if migrateFn != nil {
		migrateFn(executor)
	}

	return executor, nil
}

func connectDB(dsnGetter func() (string, error), connector func(ctx context.Context, dsn string, migrateFn func(executor ...DBExecutor)) (DBExecutor, error)) (DBExecutor, error) {
	logrus.Debug("Connecting to the database")
	dsn, err := dsnGetter()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Connecting to SQLite database at %s", dsn)

	pool, err := connector(context.Background(), dsn, Migrate)
	if err != nil {
		return nil, err
	}
	logrus.Debug("SQLite connection created and pinged successfully")
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

	// 1. Ensure migrations tracking table exists
	_, err := exec.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TEXT DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		logrus.Errorf("Failed to create migration tracking table: %v", err)
		return
	}

	// 2. Handle transition from old skip logic:
	// If stocks table exists but migrations haven't been tracked, mark 000 and 001 as applied
	var stocksExists bool
	err = exec.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name = 'stocks')").Scan(&stocksExists)
	if err != nil {
		logrus.Debugf("Note: could not check if stocks table exists: %v", err)
	}

	if stocksExists {
		_, err = exec.Exec(ctx, "INSERT OR IGNORE INTO schema_migrations (filename) VALUES ('migrations/000_board_type.sql'), ('migrations/001_create_stocks_table.sql')")
		if err != nil {
			logrus.Debugf("Note: could not seed initial migrations: %v", err)
		}
	}

	// 3. Get migration files
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
		err = exec.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE filename = ?1)", file).Scan(&applied)
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

		rawSQL := string(content)
		// Split multiple statements if any
		statements := strings.Split(rawSQL, ";")
		var execErr error
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = exec.Exec(ctx, stmt)
			if err != nil {
				execErr = err
				break
			}
		}

		if execErr != nil {
			logrus.Errorf("Failed to execute migration %s: %v", file, execErr)
			continue
		}

		// Mark as applied
		_, err = exec.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES (?1)", file)
		if err != nil {
			logrus.Errorf("Failed to record migration %s: %v", file, err)
			continue
		}
	}

	logrus.Info("Migrations synchronization completed")
}
