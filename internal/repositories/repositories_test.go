package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/models"
)

func TestIndustryRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewIndustryRepository(mockExec)
		res1, err := repo.UpsertIndustries(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res1)

		res2, err := repo.UpsertIndustrySub(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res2)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewIndustryRepository(mockExec)
		res, err := repo.UpsertIndustries(ctx, []models.Industry{{Id: 1, Name: "Tech"}})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_batch_upsert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO industry").
			WithArgs(1, "Tech").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Tech"))
		mock.ExpectCommit()

		repo := NewIndustryRepository(mockExec)
		res, err := repo.UpsertIndustries(ctx, []models.Industry{{Id: 1, Name: "Tech"}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_sub_industries_batch_upsert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO industry_sub").
			WithArgs(1, "Software", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Software"))
		mock.ExpectCommit()

		repo := NewIndustryRepository(mockExec)
		res, err := repo.UpsertIndustrySub(ctx, []models.IndustrySub{{Id: 1, Name: "Software", IndustryId: 1}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectorSearchRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewSectorSearchRepository(mockExec)
		res1, err := repo.UpsertNewSectors(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res1)

		res2, err := repo.UpsertNewSectorSub(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res2)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewSectorSearchRepository(mockExec)
		res, err := repo.UpsertNewSectors(ctx, []models.SectorNew{{Id: 1, Name: "Finance"}})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_batch_upsert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO sector").
			WithArgs(1, (*string)(nil), "Finance", (*string)(nil), (*string)(nil)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(1, nil, "Finance"))
		mock.ExpectCommit()

		repo := NewSectorSearchRepository(mockExec)
		res, err := repo.UpsertNewSectors(ctx, []models.SectorNew{{Id: 1, Name: "Finance"}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_subsectors_batch_upsert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO sector_sub").
			WithArgs(1, (*string)(nil), "Bank", (*string)(nil), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(1, nil, "Bank"))
		mock.ExpectCommit()

		repo := NewSectorSearchRepository(mockExec)
		res, err := repo.UpsertNewSectorSub(ctx, []models.SectorSub{{Id: 1, Name: "Bank", SectorId: 1}})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStockRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)
		err = repo.BatchInsertStocks(ctx, nil)
		assert.NoError(t, err)

		res, err := repo.UpsertStocksDetail(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("success_batch_insert_stocks", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		d := "2020-01-01"
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO stocks").
			WithArgs("BBCA", "Bank Central Asia", d, "", int64(100), "Main").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewStockRepository(mockExec)
		err = repo.BatchInsertStocks(ctx, []models.Stock{
			{Code: "BBCA", CompanyName: "Bank Central Asia", ListingDate: &d, DelistingDate: nil, Shares: 100, ListingBoard: "Main"},
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_upsert_stocks_detail", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		d := "2020-01-01"
		s := models.PasardanaStockDetail{
			Id: 1, Code: "BBCA", Name: "Bank Central Asia", ListingDate: &d, FoundingDate: &d,
		}
		mock.ExpectQuery("INSERT INTO stocks").
			WithArgs(s.Code, s.Name, *s.ListingDate, s.FkNewSectorId, s.FkNewSubSectorId,
				s.FkNewIndustryId, s.FkNewSubIndustryId).
			WillReturnRows(sqlmock.NewRows([]string{"code", "name"}).AddRow("BBCA", "Bank Central Asia"))
		mock.ExpectCommit()

		repo := NewStockRepository(mockExec)
		res, err := repo.UpsertStocksDetail(ctx, []models.PasardanaStockDetail{s})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)
		code := "BBCA"
		dateStr := "2024-04-17"

		rows := sqlmock.NewRows([]string{"code", "name"}).
			AddRow("BBCA", "Bank Central Asia")

		mock.ExpectQuery("UPDATE stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnRows(rows)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.NoError(t, err)
		assert.NotNil(t, stock)
		assert.Equal(t, "BBCA", stock.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_not_found_ErrNoRows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)
		code := "UNKNOWN"
		dateStr := "2024-04-17"

		mock.ExpectQuery("UPDATE stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnError(sql.ErrNoRows)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.NoError(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateDelistingDate_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)
		code := "ERR"
		dateStr := "2024-04-17"

		mock.ExpectQuery("UPDATE stocks SET delisting_date").
			WithArgs(dateStr, code).
			WillReturnError(assert.AnError)

		stock, err := repo.UpdateDelistingDate(ctx, code, dateStr)
		assert.Error(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("FindMissingCodes", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		t.Run("empty", func(t *testing.T) {
			res, err := repo.FindMissingCodes(ctx, nil)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})

		t.Run("with_codes", func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"code"}).AddRow("BBRI")
			mock.ExpectQuery("WITH input\\(code\\) AS").
				WithArgs("BBCA", "BBRI").
				WillReturnRows(rows)

			missing, err := repo.FindMissingCodes(ctx, []string{"BBCA", "BBRI"})
			assert.NoError(t, err)
			assert.Equal(t, []string{"BBRI"}, missing)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})

	t.Run("GetAllStocks", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			mockExec := database.NewDBExecutor(db)

			repo := NewStockRepository(mockExec)

			cols := []string{
				"code", "name", "listing_date", "delisting_date", "shares", "board",
				"sector_id", "sector_sub_id", "industry_id", "industry_sub_id",
			}
			rows := sqlmock.NewRows(cols).
				AddRow("BBCA", "Bank Central Asia", "2000-05-31", nil, int64(123456), "UTAMA",
					1, 2, 3, 4).
				AddRow("AADI", "Adaro Andalan", "2024-12-05", "2025-01-01", int64(7890), "PENGEMBANGAN",
					nil, nil, nil, nil)

			mock.ExpectQuery("SELECT (.+) FROM stocks WHERE").
				WithArgs("UTAMA", "BCA").
				WillReturnRows(rows)

			stocks, err := repo.GetAllStocks(ctx, models.GetAllStocksFilter{Board: "UTAMA", Search: "BCA"})
			assert.NoError(t, err)
			assert.Len(t, stocks, 2)
			assert.Equal(t, "BBCA", stocks[0].Code)
			assert.Equal(t, 1, *stocks[0].SectorId)
			assert.Equal(t, "AADI", stocks[1].Code)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("query_error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			mockExec := database.NewDBExecutor(db)

			repo := NewStockRepository(mockExec)

			mock.ExpectQuery("SELECT (.+) FROM stocks WHERE").
				WithArgs("", "").
				WillReturnError(assert.AnError)

			stocks, err := repo.GetAllStocks(ctx, models.GetAllStocksFilter{})
			assert.Error(t, err)
			assert.Nil(t, stocks)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("scan_error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			mockExec := database.NewDBExecutor(db)

			repo := NewStockRepository(mockExec)

			cols := []string{"code"} // missing columns will cause scan error
			rows := sqlmock.NewRows(cols).AddRow("BBCA")

			mock.ExpectQuery("SELECT (.+) FROM stocks WHERE").
				WithArgs("", "").
				WillReturnRows(rows)

			stocks, err := repo.GetAllStocks(ctx, models.GetAllStocksFilter{})
			assert.Error(t, err)
			assert.Nil(t, stocks)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})
}

func TestHistoryRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_upsert", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewHistoryRepository(mockExec)
		err = repo.BatchUpsertStockHistory(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("begin_tx_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin().WillReturnError(assert.AnError)

		repo := NewHistoryRepository(mockExec)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{{Code: "BBCA"}})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_batch_upsert_stock_history", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rec := models.StockHistory{Code: "BBCA", Date: date}
		mock.ExpectExec("INSERT INTO history").
			WithArgs(
				rec.Code, "2024-01-01", rec.Previous, rec.OpenPrice, rec.FirstTrade,
				rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
				rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
				rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
				rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, (*string)(nil),
				rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewHistoryRepository(mockExec)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{rec})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("batch_upsert_stock_history_exec_error_continues", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		mock.ExpectBegin()
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rec := models.StockHistory{Code: "BBCA", Date: date}
		mock.ExpectExec("INSERT INTO history").
			WithArgs(
				rec.Code, "2024-01-01", rec.Previous, rec.OpenPrice, rec.FirstTrade,
				rec.High, rec.Low, rec.Close, rec.Change, rec.Volume, rec.Value,
				rec.Frequency, rec.IndexIndividual, rec.Offer, rec.OfferVolume,
				rec.Bid, rec.BidVolume, rec.ListedShares, rec.TradebleShares,
				rec.WeightForIndex, rec.ForeignSell, rec.ForeignBuy, (*string)(nil),
				rec.NonRegularVolume, rec.NonRegularValue, rec.NonRegularFrequency,
			).
			WillReturnError(assert.AnError)
		mock.ExpectCommit()

		repo := NewHistoryRepository(mockExec)
		err = repo.BatchUpsertStockHistory(ctx, []models.StockHistory{rec})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetHistoryByCode", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewHistoryRepository(mockExec)

		dateStr := "2024-01-01"
		lastModStr := "2024-01-01 10:00:00"
		rows := sqlmock.NewRows([]string{
			"code", "date", "previous", "open_price", "first_trade", "high", "low", "close", "change",
			"volume", "value", "frequency", "index_individual", "offer", "offer_volume", "bid", "bid_volume",
			"listed_shares", "tradeble_shares", "weight_for_index", "foreign_sell", "foreign_buy",
			"delisting_date", "non_regular_volume", "non_regular_value", "non_regular_frequency", "last_modified",
		}).AddRow(
			"BBCA", dateStr, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil,
			nil, nil, lastModStr,
		)

		mock.ExpectQuery("SELECT (.+) FROM history").
			WithArgs("BBCA", nil, nil).
			WillReturnRows(rows)

		hist, err := repo.GetHistoryByCode(ctx, "BBCA", nil, nil)
		assert.NoError(t, err)
		assert.Len(t, hist, 1)
		assert.Equal(t, "BBCA", hist[0].Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ApplyStockSplit", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewHistoryRepository(mockExec)
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		mock.ExpectExec("UPDATE history").
			WithArgs("BBCA", 5.0, "2020-01-01", "2021-01-01").
			WillReturnResult(sqlmock.NewResult(1, 10))

		err = repo.ApplyStockSplit(ctx, "BBCA", start, end, 5.0)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetStockByCode_Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		columns := []string{
			"code", "name", "listing_date", "delisting_date", "shares", "board",
			"sector_id", "sector_sub_id", "industry_id", "industry_sub_id",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			"BBCA", "Bank Central Asia Tbk", "1982-01-01", "2024-01-01", int64(1000000), "UTAMA",
			1, 2, 3, 4,
		)

		mock.ExpectQuery("SELECT (.+) FROM stocks WHERE UPPER\\(code\\) = UPPER\\(\\?1\\) LIMIT 1").
			WithArgs("BBCA").
			WillReturnRows(rows)

		stock, err := repo.GetStockByCode(ctx, "BBCA")
		assert.NoError(t, err)
		assert.NotNil(t, stock)
		assert.Equal(t, "BBCA", stock.Code)
		assert.Equal(t, "Bank Central Asia Tbk", stock.Name)
		assert.NotNil(t, stock.DelistingDate)
		assert.Equal(t, "2024-01-01", *stock.DelistingDate)
		assert.NotNil(t, stock.SectorId)
		assert.Equal(t, 1, *stock.SectorId)
		assert.NotNil(t, stock.SectorSubId)
		assert.Equal(t, 2, *stock.SectorSubId)
		assert.NotNil(t, stock.IndustryId)
		assert.Equal(t, 3, *stock.IndustryId)
		assert.NotNil(t, stock.IndustrySubId)
		assert.Equal(t, 4, *stock.IndustrySubId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetStockByCode_NotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		mock.ExpectQuery("SELECT (.+) FROM stocks WHERE UPPER\\(code\\) = UPPER\\(\\?1\\) LIMIT 1").
			WithArgs("UNKN").
			WillReturnError(sql.ErrNoRows)

		stock, err := repo.GetStockByCode(ctx, "UNKN")
		assert.NoError(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetStockByCode_Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		mock.ExpectQuery("SELECT (.+) FROM stocks WHERE UPPER\\(code\\) = UPPER\\(\\?1\\) LIMIT 1").
			WithArgs("BBCA").
			WillReturnError(assert.AnError)

		stock, err := repo.GetStockByCode(ctx, "BBCA")
		assert.Error(t, err)
		assert.Nil(t, stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetDistinctBoards_Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		rows := sqlmock.NewRows([]string{"board"}).
			AddRow("Acceleration").
			AddRow("Development").
			AddRow("Main")

		mock.ExpectQuery("SELECT DISTINCT board FROM stocks WHERE board IS NOT NULL AND TRIM\\(board\\) != '' ORDER BY board ASC").
			WillReturnRows(rows)

		boards, err := repo.GetDistinctBoards(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []string{"Acceleration", "Development", "Main"}, boards)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetDistinctBoards_Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

		repo := NewStockRepository(mockExec)

		mock.ExpectQuery("SELECT DISTINCT board FROM stocks WHERE board IS NOT NULL AND TRIM\\(board\\) != '' ORDER BY board ASC").
			WillReturnError(assert.AnError)

		boards, err := repo.GetDistinctBoards(ctx)
		assert.Error(t, err)
		assert.Nil(t, boards)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
