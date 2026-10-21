package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockHistoryRepo struct {
	mock.Mock
}

func (m *mockHistoryRepo) BatchUpsertStockHistory(ctx context.Context, records []models.StockHistory) error {
	args := m.Called(ctx, records)
	return args.Error(0)
}

func (m *mockHistoryRepo) GetHistoryByCode(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error) {
	args := m.Called(ctx, code, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockHistory), args.Error(1)
}

func (m *mockHistoryRepo) ApplyStockSplit(ctx context.Context, code string, startDate, endDate time.Time, multiplier float64) error {
	args := m.Called(ctx, code, startDate, endDate, multiplier)
	return args.Error(0)
}

type mockStockRepo struct {
	mock.Mock
}

func (m *mockStockRepo) BatchInsertStocks(ctx context.Context, stocks []models.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *mockStockRepo) UpdateDelistingDate(ctx context.Context, code string, delistingDate string) (*models.StockResponse, error) {
	args := m.Called(ctx, code, delistingDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StockResponse), args.Error(1)
}

func (m *mockStockRepo) UpsertStocksDetail(ctx context.Context, stocks []models.PasardanaStockDetail) ([]models.StockResponse, error) {
	args := m.Called(ctx, stocks)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockResponse), args.Error(1)
}

func (m *mockStockRepo) FindMissingCodes(ctx context.Context, codes []string) ([]string, error) {
	args := m.Called(ctx, codes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockStockRepo) GetAllStocks(ctx context.Context, filter models.GetAllStocksFilter) ([]models.StockDetailResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockDetailResponse), args.Error(1)
}

func (m *mockStockRepo) GetStockByCode(ctx context.Context, code string) (*models.StockDetailResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StockDetailResponse), args.Error(1)
}

func (m *mockStockRepo) GetDistinctBoards(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockStockRepo) GetStocksWithNullSectorOrIndustry(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockStockRepo) UpdateStockSectorAndIndustry(ctx context.Context, code string, sectorID, sectorSubID, industryID, industrySubID *int) (bool, error) {
	args := m.Called(ctx, code, sectorID, sectorSubID, industryID, industrySubID)
	return args.Bool(0), args.Error(1)
}

type mockPasardanaService struct {
	mock.Mock
}

func (m *mockPasardanaService) FetchStockIDs() ([]models.PasardanaStock, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PasardanaStock), args.Error(1)
}

func (m *mockPasardanaService) FetchStockDetailByCode(code string) (*models.PasardanaStockDetail, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PasardanaStockDetail), args.Error(1)
}

func (m *mockPasardanaService) FetchStockHistory(year, month, day int) ([]models.PasardanaHistoryResponse, error) {
	args := m.Called(year, month, day)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PasardanaHistoryResponse), args.Error(1)
}

type mockIdxService struct {
	mock.Mock
}

func (m *mockIdxService) FetchStockSummary(ctx context.Context, year, month, day int, cookie string) ([]models.IdxSummaryData, error) {
	args := m.Called(ctx, year, month, day, cookie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.IdxSummaryData), args.Error(1)
}

func (m *mockIdxService) ParseIdxDate(dateStr string) (string, error) {
	args := m.Called(dateStr)
	return args.String(0), args.Error(1)
}

type mockStockbitService struct {
	mock.Mock
}

func (m *mockStockbitService) FetchDailyPrice(ctx context.Context, code string, from, to string, authorization string) ([]models.StockbitDailyPriceItem, error) {
	args := m.Called(ctx, code, from, to, authorization)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockbitDailyPriceItem), args.Error(1)
}

func TestSyncStockHistory(t *testing.T) {
	ctx := context.Background()

	t.Run("pasardana_after_threshold_error", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2020, Month: 1, Day: 1}, "pasardana", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "restricted for dates after 2019-07-28")
	})

	t.Run("invalid_source", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2018, Month: 1, Day: 1}, "unknown", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid source")
	})

	t.Run("pasardana_fetch_error", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		pSvc.On("FetchStockHistory", 2018, 1, 1).Return(nil, assert.AnError).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2018, Month: 1, Day: 1}, "pasardana", "")
		assert.Error(t, err)
	})

	t.Run("idx_fetch_error", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		idxSvc.On("FetchStockSummary", ctx, 2024, 1, 1, "c").Return(nil, assert.AnError).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2024, Month: 1, Day: 1}, "idx", "c")
		assert.Error(t, err)
	})

	t.Run("empty_records", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		pSvc.On("FetchStockHistory", 2018, 1, 1).Return([]models.PasardanaHistoryResponse{}, nil).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2018, Month: 1, Day: 1}, "pasardana", "")
		assert.NoError(t, err)
	})

	t.Run("pasardana_success_with_missing_codes", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		lastDate := "2018-01-01T00:00:00"
		listDate := "2010-01-01"
		pSvc.On("FetchStockHistory", 2018, 1, 1).Return([]models.PasardanaHistoryResponse{
			{Code: "BBCA", LastDate: &lastDate},
		}, nil).Once()

		sRepo.On("FindMissingCodes", ctx, []string{"BBCA"}).Return([]string{"BBCA"}, nil).Once()
		pSvc.On("FetchStockDetailByCode", "BBCA").Return(&models.PasardanaStockDetail{
			Code:        "BBCA",
			ListingDate: &listDate,
		}, nil).Once()
		sRepo.On("UpsertStocksDetail", ctx, mock.Anything).Return([]models.StockResponse{}, nil).Once()
		hRepo.On("BatchUpsertStockHistory", ctx, mock.Anything).Return(nil).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2018, Month: 1, Day: 1}, "pasardana", "")
		assert.NoError(t, err)
	})

	t.Run("idx_success_flow", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		idxSvc.On("FetchStockSummary", ctx, 2024, 1, 1, "").Return([]models.IdxSummaryData{
			{StockCode: "BBRI", Date: "2024-01-01T00:00:00", DelistingDate: "2025-01-01"},
		}, nil).Once()

		sRepo.On("FindMissingCodes", ctx, []string{"BBRI"}).Return([]string{}, nil).Once()
		hRepo.On("BatchUpsertStockHistory", ctx, mock.Anything).Return(nil).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2024, Month: 1, Day: 1}, "idx", "")
		assert.NoError(t, err)
	})

	t.Run("idx_missing_codes_error_handled", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		idxSvc.On("FetchStockSummary", ctx, 2024, 1, 1, "").Return([]models.IdxSummaryData{
			{StockCode: "BBRI", Date: "2024-01-01T00:00:00"},
		}, nil).Once()

		sRepo.On("FindMissingCodes", ctx, []string{"BBRI"}).Return(nil, assert.AnError).Once()
		hRepo.On("BatchUpsertStockHistory", ctx, mock.Anything).Return(nil).Once()

		err := u.SyncStockHistory(ctx, models.SyncHistoryRequest{Year: 2024, Month: 1, Day: 1}, "idx", "")
		assert.NoError(t, err)
	})
}

func TestGetStockHistory(t *testing.T) {
	hRepo := new(mockHistoryRepo)
	sRepo := new(mockStockRepo)
	pSvc := new(mockPasardanaService)
	idxSvc := new(mockIdxService)
	sbSvc := new(mockStockbitService)
	u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)
	ctx := context.Background()

	hRepo.On("GetHistoryByCode", ctx, "BBCA", (*time.Time)(nil), (*time.Time)(nil)).
		Return([]models.StockHistory{{Code: "BBCA"}}, nil).Once()

	data, err := u.GetStockHistory(ctx, "bbca", nil, nil)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
}

func TestApplyStockSplit(t *testing.T) {
	hRepo := new(mockHistoryRepo)
	sRepo := new(mockStockRepo)
	pSvc := new(mockPasardanaService)
	idxSvc := new(mockIdxService)
	sbSvc := new(mockStockbitService)
	u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)
	ctx := context.Background()

	t.Run("invalid_start_date", func(t *testing.T) {
		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "invalid"})
		assert.Error(t, err)
	})

	t.Run("invalid_end_date", func(t *testing.T) {
		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "2024-01-01", EndDate: "invalid"})
		assert.Error(t, err)
	})

	t.Run("end_before_start", func(t *testing.T) {
		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "2024-01-02", EndDate: "2024-01-01"})
		assert.Error(t, err)
	})

	t.Run("invalid_ratio_format", func(t *testing.T) {
		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "2024-01-01", EndDate: "2024-01-02", Ratio: "1-5"})
		assert.Error(t, err)
	})

	t.Run("invalid_ratio_values", func(t *testing.T) {
		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "2024-01-01", EndDate: "2024-01-02", Ratio: "0:5"})
		assert.Error(t, err)

		err = u.ApplyStockSplit(ctx, models.StockSplitRequest{StartDate: "2024-01-01", EndDate: "2024-01-02", Ratio: "1:0"})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		start, _ := time.Parse("2006-01-02", "2024-01-01")
		end, _ := time.Parse("2006-01-02", "2024-01-02")
		hRepo.On("ApplyStockSplit", ctx, "BBCA", start, end, float64(5)).Return(nil).Once()

		err := u.ApplyStockSplit(ctx, models.StockSplitRequest{
			StockCode: "bbca",
			StartDate: "2024-01-01",
			EndDate:   "2024-01-02",
			Ratio:     "1:5",
		})
		assert.NoError(t, err)
	})
}

func TestSyncStockHistoryFromStockbit(t *testing.T) {
	ctx := context.Background()

	t.Run("empty_code_error", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		count, err := u.SyncStockHistoryFromStockbit(ctx, "", "", "", "")
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "code is required")
	})

	t.Run("fetch_error", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		sRepo.On("FindMissingCodes", ctx, []string{"BBCA"}).Return(nil, nil).Once()
		sbSvc.On("FetchDailyPrice", ctx, "BBCA", mock.Anything, "2000-01-01", "Bearer tok").
			Return(nil, assert.AnError).Once()

		count, err := u.SyncStockHistoryFromStockbit(ctx, "BBCA", "", "", "Bearer tok")
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("success_with_missing_code_auto_fetch", func(t *testing.T) {
		hRepo := new(mockHistoryRepo)
		sRepo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		sbSvc := new(mockStockbitService)
		u := NewHistoryUsecase(hRepo, sRepo, pSvc, idxSvc, sbSvc)

		sRepo.On("FindMissingCodes", ctx, []string{"BBCA"}).Return([]string{"BBCA"}, nil).Once()
		dateStr := "2020-01-01"
		pSvc.On("FetchStockDetailByCode", "BBCA").Return(&models.PasardanaStockDetail{
			Code:        "BBCA",
			ListingDate: &dateStr,
		}, nil).Once()
		sRepo.On("UpsertStocksDetail", ctx, mock.Anything).Return([]models.StockResponse{{Code: "BBCA"}}, nil).Once()

		openVal := 4740.0
		items := []models.StockbitDailyPriceItem{
			{Date: "2026-10-20", Open: &openVal},
			{Date: "invalid-date", Open: &openVal}, // tests skipped branch
		}
		sbSvc.On("FetchDailyPrice", ctx, "BBCA", "2026-10-20", "2000-01-01", "Bearer tok").
			Return(items, nil).Once()
		hRepo.On("BatchUpsertStockHistory", ctx, mock.MatchedBy(func(recs []models.StockHistory) bool {
			return len(recs) == 1 && recs[0].Code == "BBCA"
		})).Return(nil).Once()

		count, err := u.SyncStockHistoryFromStockbit(ctx, "BBCA", "2026-10-20", "2000-01-01", "Bearer tok")
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
