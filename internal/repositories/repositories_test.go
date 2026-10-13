package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/models"
)

func TestIndustryRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewIndustryRepository(mock)
		res1, err := repo.UpsertIndustries(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res1)

		res2, err := repo.UpsertSubIndustries(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res2)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewIndustryRepository(mock)
		res, err := repo.UpsertIndustries(ctx, []models.Industry{{Id: 1, Name: "Tech"}})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_batch_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		batch.ExpectQuery("INSERT INTO idxstock.industry").
			WithArgs(1, "Tech").
			WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "Tech"))
		mock.ExpectCommit()

		repo := NewIndustryRepository(mock)
		res, err := repo.UpsertIndustries(ctx, []models.Industry{{Id: 1, Name: "Tech"}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_sub_industries_batch_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		batch.ExpectQuery("INSERT INTO idxstock.sub_industry").
			WithArgs(1, "Software", 1).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "Software"))
		mock.ExpectCommit()

		repo := NewIndustryRepository(mock)
		res, err := repo.UpsertSubIndustries(ctx, []models.SubIndustry{{Id: 1, Name: "Software", IndustryId: 1}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectorSearchRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSectorSearchRepository(mock)
		res1, err := repo.UpsertNewSectors(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res1)

		res2, err := repo.UpsertNewSubSectors(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res2)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewSectorSearchRepository(mock)
		res, err := repo.UpsertNewSectors(ctx, []models.SectorNew{{Id: 1, Name: "Finance"}})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("success_batch_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		batch.ExpectQuery("INSERT INTO idxstock.sector").
			WithArgs(1, (*string)(nil), "Finance", (*string)(nil), (*string)(nil)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "code", "name"}).AddRow(1, nil, "Finance"))
		mock.ExpectCommit()

		repo := NewSectorSearchRepository(mock)
		res, err := repo.UpsertNewSectors(ctx, []models.SectorNew{{Id: 1, Name: "Finance"}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_subsectors_batch_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		batch.ExpectQuery("INSERT INTO idxstock.sub_sector").
			WithArgs(1, (*string)(nil), "Bank", (*string)(nil), (*string)(nil), 1).
			WillReturnRows(pgxmock.NewRows([]string{"id", "code", "name"}).AddRow(1, nil, "Bank"))
		mock.ExpectCommit()

		repo := NewSectorSearchRepository(mock)
		res, err := repo.UpsertNewSubSectors(ctx, []models.SubSector{{Id: 1, Name: "Bank", SectorId: 1}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStockRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewStockRepository(mock)
		err = repo.BatchInsertStocks(ctx, nil)
		assert.NoError(t, err)

		res, err := repo.UpsertStocksDetail(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("success_batch_insert_stocks", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		d := "2020-01-01"
		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		batch.ExpectExec("INSERT INTO idxstock.stocks").
			WithArgs("BBCA", "Bank Central Asia", &d, (*string)(nil), int64(100), "Main").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectCommit()

		repo := NewStockRepository(mock)
		err = repo.BatchInsertStocks(ctx, []models.Stock{
			{Code: "BBCA", CompanyName: "Bank Central Asia", ListingDate: &d, DelistingDate: nil, Shares: 100, ListingBoard: "Main"},
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_upsert_stocks_detail", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		d := "2020-01-01"
		s := models.PasardanaStockDetail{
			Id: 1, Code: "BBCA", Name: "Bank Central Asia", ListingDate: &d, FoundingDate: &d,
		}
		batch.ExpectQuery("INSERT INTO idxstock.stocks").
			WithArgs(s.Id, s.Code, s.Name, s.ListingDate, s.TotalEmployees, s.AnnualDividend,
				s.GeneralInformation, s.FoundingDate, s.FkNewSectorId, s.FkNewSubSectorId,
				s.FkNewIndustryId, s.FkNewSubIndustryId).
			WillReturnRows(pgxmock.NewRows([]string{"id", "code", "name"}).AddRow(1, "BBCA", "Bank Central Asia"))
		mock.ExpectCommit()

		repo := NewStockRepository(mock)
		res, err := repo.UpsertStocksDetail(ctx, []models.PasardanaStockDetail{s})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewStockRepository(mock)
		code := "BBCA"
		dateStr := "2024-04-17"

		rows := pgxmock.NewRows([]string{"id", "code", "name"}).
			AddRow(1, "BBCA", "Bank Central Asia")

		mock.ExpectQuery("UPDATE idxstock.stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnRows(rows)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.NoError(t, err)
		assert.NotNil(t, stock)
		assert.Equal(t, "BBCA", stock.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_not_found_ErrNoRows", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewStockRepository(mock)
		code := "UNKNOWN"
		dateStr := "2024-04-17"

		mock.ExpectQuery("UPDATE idxstock.stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnError(pgx.ErrNoRows)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.NoError(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewStockRepository(mock)
		code := "ERR"
		dateStr := "2024-04-17"

		mock.ExpectQuery("UPDATE idxstock.stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnError(assert.AnError)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.Error(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("FindMissingCodes", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewStockRepository(mock)

		t.Run("empty", func(t *testing.T) {
			res, err := repo.FindMissingCodes(ctx, nil)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})

		t.Run("with_codes", func(t *testing.T) {
			rows := pgxmock.NewRows([]string{"c"}).AddRow("BBRI")
			mock.ExpectQuery("SELECT c FROM").
				WithArgs([]string{"BBCA", "BBRI"}).
				WillReturnRows(rows)

			missing, err := repo.FindMissingCodes(ctx, []string{"BBCA", "BBRI"})
			assert.NoError(t, err)
			assert.Equal(t, []string{"BBRI"}, missing)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})
}

func TestHistoryRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewHistoryRepository(mock)
		err = repo.BatchUpsertStockHistory(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewHistoryRepository(mock)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{{Code: "BBCA"}})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_batch_upsert_stock_history", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rec := models.StockHistory{Code: "BBCA", Date: date}
		batch.ExpectExec("INSERT INTO idxstock.history").
			WithArgs(
				rec.Code, rec.Date, rec.Previous, rec.OpenPrice, rec.FirstTrade,
				rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
				rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
				rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
				rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, (*string)(nil),
				rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectCommit()

		repo := NewHistoryRepository(mock)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{rec})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("batch_upsert_stock_history_exec_error_continues", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		batch := mock.ExpectBatch()
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rec := models.StockHistory{Code: "BBCA", Date: date}
		batch.ExpectExec("INSERT INTO idxstock.history").
			WithArgs(
				rec.Code, rec.Date, rec.Previous, rec.OpenPrice, rec.FirstTrade,
				rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
				rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
				rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
				rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, (*string)(nil),
				rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
			).
			WillReturnError(assert.AnError)
		mock.ExpectCommit()

		repo := NewHistoryRepository(mock)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{rec})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetHistoryByCode", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewHistoryRepository(mock)

		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rows := pgxmock.NewRows([]string{
			"code", "date", "previous", "open_price", "first_trade", "high", "low", "close", "change",
			"volume", "value", "frequency", "index_individual", "offer", "offer_volume", "bid", "bid_volume",
			"listed_shares", "tradeble_shares", "weight_for_index", "foreign_sell", "foreign_buy",
			"delisting_date", "non_regular_volume", "non_regular_value", "non_regular_frequency", "last_modified",
		}).AddRow(
			"BBCA", date, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil,
			nil, nil, time.Now(),
		)

		mock.ExpectQuery("SELECT (.+) FROM idxstock.history").
			WithArgs("BBCA", (*time.Time)(nil), (*time.Time)(nil)).
			WillReturnRows(rows)

		hist, err := repo.GetHistoryByCode(ctx, "BBCA", nil, nil)
		assert.NoError(t, err)
		assert.Len(t, hist, 1)
		assert.Equal(t, "BBCA", hist[0].Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ApplyStockSplit", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewHistoryRepository(mock)
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		mock.ExpectExec("UPDATE idxstock.history").
			WithArgs("BBCA", 5.0, start, end).
			WillReturnResult(pgxmock.NewResult("UPDATE", 10))

		err = repo.ApplyStockSplit(ctx, "BBCA", start, end, 5.0)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
