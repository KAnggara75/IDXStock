package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
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

func connectWithConfig(ctx context.Context, dsn string) (DBExecutor, error) {
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

	return NewDBExecutor(db), nil
}

func connectDB(dsnGetter func() (string, error), connector func(ctx context.Context, dsn string) (DBExecutor, error)) (DBExecutor, error) {
	logrus.Debug("Connecting to the database")
	dsn, err := dsnGetter()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Connecting to SQLite database at %s", dsn)

	pool, err := connector(context.Background(), dsn)
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
