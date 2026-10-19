package repositories

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/models"
)

type BrokerActivityRepository interface {
	BatchInsertBrokerActivity(ctx context.Context, records []models.BrokerActivity) error
	InsertBrokerActivity(ctx context.Context, record models.BrokerActivity) (bool, error)
	CheckPartitionExists(ctx context.Context, tableName string) (bool, error)
	CreatePartition(ctx context.Context, tableName, startDate, endDate string) error
}

// DBExecutor is an alias to database.DBExecutor
type DBExecutor = database.DBExecutor

type brokerActivityRepository struct {
	pool DBExecutor
}

func NewBrokerActivityRepository(pool DBExecutor) BrokerActivityRepository {
	return &brokerActivityRepository{
		pool: pool,
	}
}

func (r *brokerActivityRepository) BatchInsertBrokerActivity(ctx context.Context, records []models.BrokerActivity) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO broker_activity (
			broker_code, stock_code, date, side, lot, value, avg_price, freq, updated_at
		)
		VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, datetime('now'))
		ON CONFLICT (broker_code, date, stock_code, side) DO UPDATE SET
			lot = excluded.lot,
			value = excluded.value,
			avg_price = excluded.avg_price,
			freq = excluded.freq,
			updated_at = datetime('now')
		WHERE
			broker_activity.lot IS NOT excluded.lot OR
			broker_activity.value IS NOT excluded.value OR
			broker_activity.avg_price IS NOT excluded.avg_price OR
			broker_activity.freq IS NOT excluded.freq
	`

	var affected int64
	for i, rec := range records {
		ct, err := tx.Exec(ctx, query,
			rec.BrokerCode, rec.StockCode, rec.Date.Format("2006-01-02"), rec.Side,
			rec.Lot, rec.Value, rec.AvgPrice, rec.Freq,
		)
		if err != nil {
			return fmt.Errorf("failed to execute batch upsert for broker activity at index %d: %w", i, err)
		}
		affected += ct.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Successfully upserted %d broker activity records (affected: %d)", len(records), affected)
	return nil
}

func (r *brokerActivityRepository) InsertBrokerActivity(ctx context.Context, rec models.BrokerActivity) (bool, error) {
	query := `
		INSERT INTO broker_activity (
			broker_code, stock_code, date, side, lot, value, avg_price, freq, updated_at
		)
		VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, datetime('now'))
		ON CONFLICT (broker_code, date, stock_code, side) DO NOTHING
	`

	ct, err := r.pool.Exec(ctx, query,
		rec.BrokerCode, rec.StockCode, rec.Date.Format("2006-01-02"), rec.Side,
		rec.Lot, rec.Value, rec.AvgPrice, rec.Freq,
	)
	if err != nil {
		return false, fmt.Errorf("failed to insert broker activity (%s, %s, %s, %s): %w",
			rec.BrokerCode, rec.StockCode, rec.Date.Format("2006-01-02"), rec.Side, err)
	}

	return ct.RowsAffected() > 0, nil
}

func (r *brokerActivityRepository) CheckPartitionExists(ctx context.Context, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM sqlite_master
			WHERE type = 'table' AND name = ?1
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check partition existence for %s: %w", tableName, err)
	}
	return exists, nil
}

func (r *brokerActivityRepository) CreatePartition(ctx context.Context, tableName, startDate, endDate string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			broker_code TEXT NOT NULL REFERENCES brokers(code),
			stock_code TEXT NOT NULL REFERENCES stocks(code),
			date TEXT NOT NULL CHECK (date >= '%s' AND date <= '%s'),
			side TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
			lot INTEGER,
			value INTEGER,
			avg_price REAL,
			freq INTEGER,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			PRIMARY KEY (broker_code, date, stock_code, side)
		)
	`, tableName, startDate, endDate)

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create partition %s [%s to %s]: %w", tableName, startDate, endDate, err)
	}

	return nil
}
