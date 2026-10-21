package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type HistoryRepository interface {
	BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error
	GetHistoryByCode(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error)
	ApplyStockSplit(ctx context.Context, code string, startDate, endDate time.Time, multiplier float64) error
}

type historyRepository struct {
	pool DBExecutor
}

func NewHistoryRepository(pool DBExecutor) HistoryRepository {
	return &historyRepository{
		pool: pool,
	}
}

func (r *historyRepository) BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO stocks_history (
			code, date, previous, open_price, first_trade, high, low, close, change,
			volume, value, frequency, index_individual, offer, offer_volume,
			bid, bid_volume, listed_shares, tradeble_shares, weight_for_index,
			foreign_sell, foreign_buy, delisting_date, non_regular_volume,
			non_regular_value, non_regular_frequency, last_modified
		)
		VALUES (
			?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14, ?15,
			?16, ?17, ?18, ?19, ?20, ?21, ?22, ?23, ?24, ?25, ?26, datetime('now')
		)
		ON CONFLICT (code, date) DO UPDATE SET
			previous = COALESCE(excluded.previous, stocks_history.previous),
			open_price = COALESCE(excluded.open_price, stocks_history.open_price),
			first_trade = COALESCE(excluded.first_trade, stocks_history.first_trade),
			high = COALESCE(excluded.high, stocks_history.high),
			low = COALESCE(excluded.low, stocks_history.low),
			close = COALESCE(excluded.close, stocks_history.close),
			change = COALESCE(excluded.change, stocks_history.change),
			volume = COALESCE(excluded.volume, stocks_history.volume),
			value = COALESCE(excluded.value, stocks_history.value),
			frequency = COALESCE(excluded.frequency, stocks_history.frequency),
			index_individual = COALESCE(excluded.index_individual, stocks_history.index_individual),
			offer = COALESCE(excluded.offer, stocks_history.offer),
			offer_volume = COALESCE(excluded.offer_volume, stocks_history.offer_volume),
			bid = COALESCE(excluded.bid, stocks_history.bid),
			bid_volume = COALESCE(excluded.bid_volume, stocks_history.bid_volume),
			listed_shares = COALESCE(excluded.listed_shares, stocks_history.listed_shares),
			tradeble_shares = COALESCE(excluded.tradeble_shares, stocks_history.tradeble_shares),
			weight_for_index = COALESCE(excluded.weight_for_index, stocks_history.weight_for_index),
			foreign_sell = COALESCE(excluded.foreign_sell, stocks_history.foreign_sell),
			foreign_buy = COALESCE(excluded.foreign_buy, stocks_history.foreign_buy),
			delisting_date = COALESCE(excluded.delisting_date, stocks_history.delisting_date),
			non_regular_volume = COALESCE(excluded.non_regular_volume, stocks_history.non_regular_volume),
			non_regular_value = COALESCE(excluded.non_regular_value, stocks_history.non_regular_value),
			non_regular_frequency = COALESCE(excluded.non_regular_frequency, stocks_history.non_regular_frequency),
			last_modified = datetime('now')
		WHERE
			(excluded.previous IS NOT NULL AND stocks_history.previous IS NOT excluded.previous) OR
			(excluded.open_price IS NOT NULL AND stocks_history.open_price IS NOT excluded.open_price) OR
			(excluded.first_trade IS NOT NULL AND stocks_history.first_trade IS NOT excluded.first_trade) OR
			(excluded.high IS NOT NULL AND stocks_history.high IS NOT excluded.high) OR
			(excluded.low IS NOT NULL AND stocks_history.low IS NOT excluded.low) OR
			(excluded.close IS NOT NULL AND stocks_history.close IS NOT excluded.close) OR
			(excluded.change IS NOT NULL AND stocks_history.change IS NOT excluded.change) OR
			(excluded.volume IS NOT NULL AND stocks_history.volume IS NOT excluded.volume) OR
			(excluded.value IS NOT NULL AND stocks_history.value IS NOT excluded.value) OR
			(excluded.frequency IS NOT NULL AND stocks_history.frequency IS NOT excluded.frequency) OR
			(excluded.index_individual IS NOT NULL AND stocks_history.index_individual IS NOT excluded.index_individual) OR
			(excluded.offer IS NOT NULL AND stocks_history.offer IS NOT excluded.offer) OR
			(excluded.offer_volume IS NOT NULL AND stocks_history.offer_volume IS NOT excluded.offer_volume) OR
			(excluded.bid IS NOT NULL AND stocks_history.bid IS NOT excluded.bid) OR
			(excluded.bid_volume IS NOT NULL AND stocks_history.bid_volume IS NOT excluded.bid_volume) OR
			(excluded.listed_shares IS NOT NULL AND stocks_history.listed_shares IS NOT excluded.listed_shares) OR
			(excluded.tradeble_shares IS NOT NULL AND stocks_history.tradeble_shares IS NOT excluded.tradeble_shares) OR
			(excluded.weight_for_index IS NOT NULL AND stocks_history.weight_for_index IS NOT excluded.weight_for_index) OR
			(excluded.foreign_sell IS NOT NULL AND stocks_history.foreign_sell IS NOT excluded.foreign_sell) OR
			(excluded.foreign_buy IS NOT NULL AND stocks_history.foreign_buy IS NOT excluded.foreign_buy) OR
			(excluded.delisting_date IS NOT NULL AND stocks_history.delisting_date IS NOT excluded.delisting_date) OR
			(excluded.non_regular_volume IS NOT NULL AND stocks_history.non_regular_volume IS NOT excluded.non_regular_volume) OR
			(excluded.non_regular_value IS NOT NULL AND stocks_history.non_regular_value IS NOT excluded.non_regular_value) OR
			(excluded.non_regular_frequency IS NOT NULL AND stocks_history.non_regular_frequency IS NOT excluded.non_regular_frequency)
	`

	var affected int64
	for _, rec := range records {
		var dd *string
		if rec.DelistingDate != nil && *rec.DelistingDate != "" {
			dd = rec.DelistingDate
		}

		dateStr := rec.Date.Format("2006-01-02")
		ct, err := tx.Exec(ctx, query,
			rec.Code, dateStr, rec.Previous, rec.OpenPrice, rec.FirstTrade,
			rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
			rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
			rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
			rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, dd,
			rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
		)
		if err != nil {
			logrus.Warnf("Failed to upsert history for %s on %s: %v", rec.Code, dateStr, err)
			continue
		}
		affected += ct.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Batch upsert completed. Affected rows: %d", affected)
	return nil
}

func (r *historyRepository) GetHistoryByCode(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error) {
	var startStr, endStr any
	if startDate != nil {
		startStr = startDate.Format("2006-01-02")
	}
	if endDate != nil {
		endStr = endDate.Format("2006-01-02")
	}

	query := `
		SELECT
			code, date, previous, open_price, first_trade, high, low, close, change,
			volume, value, frequency, index_individual, offer, offer_volume,
			bid, bid_volume, listed_shares, tradeble_shares, weight_for_index,
			foreign_sell, foreign_buy, delisting_date, non_regular_volume,
			non_regular_value, non_regular_frequency, last_modified
		FROM stocks_history
		WHERE code = ?1
			AND (?2 IS NULL OR date >= ?2)
			AND (?3 IS NULL OR date <= ?3)
		ORDER BY date ASC
	`

	rows, err := r.pool.Query(ctx, query, code, startStr, endStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %w", code, err)
	}
	defer rows.Close()

	var records []models.StockHistory
	for rows.Next() {
		var h models.StockHistory
		var dateStr string
		var lastModifiedStr string
		err := rows.Scan(
			&h.Code, &dateStr, &h.Previous, &h.OpenPrice, &h.FirstTrade,
			&h.High, &h.Low, &h.Close, &h.Change, &h.Volume, &h.Value,
			&h.Frequency, &h.IndexIndividual, &h.Offer, &h.OfferVolume,
			&h.Bid, &h.BidVolume, &h.ListedShares, &h.TradebleShares,
			&h.WeightForIndex, &h.ForeignSell, &h.ForeignBuy, &h.DelistingDate,
			&h.NonRegularVolume, &h.NonRegularValue, &h.NonRegularFrequency, &lastModifiedStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history record: %w", err)
		}
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			h.Date = t
		}
		if t, err := time.Parse("2006-01-02 15:04:05", lastModifiedStr); err == nil {
			h.LastModified = t
		}
		records = append(records, h)
	}

	return records, nil
}

func (r *historyRepository) ApplyStockSplit(ctx context.Context, code string, startDate, endDate time.Time, multiplier float64) error {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	query := `
		UPDATE stocks_history
		SET
			previous = previous / ?2,
			open_price = open_price / ?2,
			first_trade = first_trade / ?2,
			high = high / ?2,
			low = low / ?2,
			close = close / ?2,
			offer = offer / ?2,
			bid = bid / ?2,
			volume = volume * ?2,
			offer_volume = offer_volume * ?2,
			bid_volume = bid_volume * ?2,
			listed_shares = listed_shares * ?2,
			tradeble_shares = tradeble_shares * ?2,
			foreign_sell = foreign_sell * ?2,
			foreign_buy = foreign_buy * ?2,
			non_regular_volume = non_regular_volume * ?2,
			change = (close / ?2) - (previous / ?2),
			last_modified = datetime('now')
		WHERE code = ?1
		  AND date >= ?3
		  AND date <= ?4
	`

	_, err := r.pool.Exec(ctx, query, code, multiplier, startStr, endStr)
	if err != nil {
		return fmt.Errorf("failed to apply stock split for %s: %w", code, err)
	}

	return nil
}
