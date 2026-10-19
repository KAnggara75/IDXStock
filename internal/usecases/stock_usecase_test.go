package usecases

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockStockService struct {
	mock.Mock
}

func (m *mockStockService) ParseExcel(file io.Reader) ([]models.Stock, error) {
	args := m.Called(file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Stock), args.Error(1)
}

func TestStockUsecase_PreviewAndUpload(t *testing.T) {
	ctx := context.Background()
	buf := bytes.NewBufferString("dummy")

	t.Run("preview_success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		svc.On("ParseExcel", buf).Return([]models.Stock{{Code: "BBCA"}}, nil).Once()

		stocks, err := u.PreviewStocks(ctx, buf)
		assert.NoError(t, err)
		assert.Len(t, stocks, 1)
	})

	t.Run("upload_parse_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		svc.On("ParseExcel", buf).Return(nil, assert.AnError).Once()

		stocks, err := u.UploadStocks(ctx, buf)
		assert.Error(t, err)
		assert.Nil(t, stocks)
	})

	t.Run("upload_batch_insert_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		stocks := []models.Stock{{Code: "BBCA"}}
		svc.On("ParseExcel", buf).Return(stocks, nil).Once()
		repo.On("BatchInsertStocks", ctx, stocks).Return(assert.AnError).Once()

		res, err := u.UploadStocks(ctx, buf)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("upload_success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		stocks := []models.Stock{{Code: "BBCA"}}
		svc.On("ParseExcel", buf).Return(stocks, nil).Once()
		repo.On("BatchInsertStocks", ctx, stocks).Return(nil).Once()

		res, err := u.UploadStocks(ctx, buf)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestStockUsecase_SyncStockDetail(t *testing.T) {
	ctx := context.Background()

	t.Run("fetch_stock_ids_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		pSvc.On("FetchStockIDs").Return(nil, assert.AnError).Once()

		res, err := u.SyncStockDetail(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_sync_stock_detail", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		pSvc.On("FetchStockIDs").Return([]models.PasardanaStock{
			{Code: "BBCA"},
		}, nil).Once()

		dateStr := "2020-01-01"
		pSvc.On("FetchStockDetailByCode", "BBCA").Return(&models.PasardanaStockDetail{
			Code:         "BBCA",
			ListingDate:  &dateStr,
			FoundingDate: &dateStr,
		}, nil).Once()

		repo.On("UpsertStocksDetail", mock.Anything, mock.Anything).
			Return([]models.StockResponse{{Code: "BBCA"}}, nil).Once()

		res, err := u.SyncStockDetail(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestStockUsecase_SyncDelistingStocks(t *testing.T) {
	ctx := context.Background()

	t.Run("fetch_delisted_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		idxSvc.On("FetchDelistedStocks", ctx, 2024, 4, "c").Return(nil, assert.AnError).Once()

		res, err := u.SyncDelistingStocks(ctx, 2024, 4, "c")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_existing_stock", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		idxSvc.On("FetchDelistedStocks", ctx, 2024, 4, "c").Return([]models.IdxDelistedStock{
			{Code: "BBCA", DeListingDate: "2024-04-17"},
		}, nil).Once()

		repo.On("UpdateDelistingDate", ctx, "BBCA", "2024-04-17").
			Return(&models.StockResponse{Code: "BBCA"}, nil).Once()

		res, err := u.SyncDelistingStocks(ctx, 2024, 4, "c")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("stock_not_found_in_db_fetches_pasardana_and_inserts", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		idxSvc.On("FetchDelistedStocks", ctx, 2024, 4, "c").Return([]models.IdxDelistedStock{
			{Code: "NEWC", DeListingDate: "2024-04-17"},
		}, nil).Once()

		// First update returns nil (stock not found in DB)
		repo.On("UpdateDelistingDate", ctx, "NEWC", "2024-04-17").Return(nil, nil).Once()

		d := "2020-01-01"
		pSvc.On("FetchStockDetailByCode", "NEWC").Return(&models.PasardanaStockDetail{
			Code:         "NEWC",
			ListingDate:  &d,
			FoundingDate: &d,
		}, nil).Once()

		repo.On("UpsertStocksDetail", ctx, mock.Anything).Return([]models.StockResponse{{Code: "NEWC"}}, nil).Once()
		// Retry update
		repo.On("UpdateDelistingDate", ctx, "NEWC", "2024-04-17").Return(&models.StockResponse{Code: "NEWC"}, nil).Once()

		res, err := u.SyncDelistingStocks(ctx, 2024, 4, "c")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("skip_invalid_date_and_update_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		idxSvc.On("FetchDelistedStocks", ctx, 2024, 4, "c").Return([]models.IdxDelistedStock{
			{Code: "INV", DeListingDate: ""},
			{Code: "ERR", DeListingDate: "2024-04-17"},
		}, nil).Once()

		repo.On("UpdateDelistingDate", ctx, "ERR", "2024-04-17").
			Return(nil, assert.AnError).Once()

		res, err := u.SyncDelistingStocks(ctx, 2024, 4, "c")
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("stock_not_found_in_db_fetches_pasardana_errors", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		idxSvc.On("FetchDelistedStocks", ctx, 2024, 4, "c").Return([]models.IdxDelistedStock{
			{Code: "FAIL", DeListingDate: "2024-04-17"},
		}, nil).Once()

		repo.On("UpdateDelistingDate", ctx, "FAIL", "2024-04-17").Return(nil, nil).Once()
		pSvc.On("FetchStockDetailByCode", "FAIL").Return(nil, assert.AnError).Once()

		res, err := u.SyncDelistingStocks(ctx, 2024, 4, "c")
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestStockUsecase_GetAllStocks(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		filter := models.GetAllStocksFilter{Board: "UTAMA", Search: "BBCA"}
		expected := []models.StockDetailResponse{
			{Code: "BBCA", Name: "Bank Central Asia", Board: "UTAMA"},
		}

		repo.On("GetAllStocks", ctx, filter).Return(expected, nil).Once()

		result, err := u.GetAllStocks(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		filter := models.GetAllStocksFilter{}
		repo.On("GetAllStocks", ctx, filter).Return(nil, assert.AnError).Once()

		result, err := u.GetAllStocks(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestStockUsecase_GetStockByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		expected := &models.StockDetailResponse{
			Code:  "BBCA",
			Name:  "Bank Central Asia Tbk",
			Board: "UTAMA",
		}

		repo.On("GetStockByCode", ctx, "BBCA").Return(expected, nil).Once()

		result, err := u.GetStockByCode(ctx, "BBCA")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		repo.On("GetStockByCode", ctx, "UNKN").Return((*models.StockDetailResponse)(nil), nil).Once()

		result, err := u.GetStockByCode(ctx, "UNKN")
		assert.NoError(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		repo.On("GetStockByCode", ctx, "BBCA").Return((*models.StockDetailResponse)(nil), assert.AnError).Once()

		result, err := u.GetStockByCode(ctx, "BBCA")
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestStockUsecase_GetDistinctBoards(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		expected := []string{"Acceleration", "Development", "Main"}
		repo.On("GetDistinctBoards", ctx).Return(expected, nil).Once()

		result, err := u.GetDistinctBoards(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		repo.On("GetDistinctBoards", ctx).Return(nil, assert.AnError).Once()

		result, err := u.GetDistinctBoards(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestStockUsecase_GetStocksByBoard(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		expected := []models.StockDetailResponse{
			{Code: "BBCA", Name: "Bank Central Asia Tbk", Board: "Main"},
		}
		repo.On("GetAllStocks", ctx, models.GetAllStocksFilter{Board: "Main"}).Return(expected, nil).Once()

		result, err := u.GetStocksByBoard(ctx, "Main")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		svc := new(mockStockService)
		pSvc := new(mockPasardanaService)
		idxSvc := new(mockIdxService)
		u := NewStockUsecase(repo, svc, pSvc, idxSvc)

		repo.On("GetAllStocks", ctx, models.GetAllStocksFilter{Board: "Main"}).Return(nil, assert.AnError).Once()

		result, err := u.GetStocksByBoard(ctx, "Main")
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}
