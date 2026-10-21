package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

func TestStockUsecase_SyncStockDetail(t *testing.T) {
	ctx := context.Background()

	t.Run("fetch_stock_ids_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		pSvc.On("FetchStockIDs").Return(nil, assert.AnError).Once()

		res, err := u.SyncStockDetail(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success_sync_stock_detail", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

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

func TestStockUsecase_SyncStockSectorAndIndustry(t *testing.T) {
	ctx := context.Background()

	t.Run("get_null_stocks_error", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		repo.On("GetStocksWithNullSectorOrIndustry", ctx).Return(nil, assert.AnError).Once()

		res, err := u.SyncStockSectorAndIndustry(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("empty_null_stocks", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		repo.On("GetStocksWithNullSectorOrIndustry", ctx).Return([]string{}, nil).Once()

		res, err := u.SyncStockSectorAndIndustry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, res.TotalChecked)
		assert.Equal(t, 0, res.TotalUpdated)
	})

	t.Run("success_sync_sector_and_industry", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		repo.On("GetStocksWithNullSectorOrIndustry", ctx).Return([]string{"PTBA"}, nil).Once()

		secID, subSecID, indID, subIndID := 1, 1, 2, 3
		pSvc.On("FetchStockDetailByCode", "PTBA").Return(&models.PasardanaStockDetail{
			Code:               "PTBA",
			FkNewSectorId:      &secID,
			FkNewSubSectorId:   &subSecID,
			FkNewIndustryId:    &indID,
			FkNewSubIndustryId: &subIndID,
		}, nil).Once()

		repo.On("UpdateStockSectorAndIndustry", mock.Anything, "PTBA", &secID, &subSecID, &indID, &subIndID).
			Return(true, nil).Once()

		res, err := u.SyncStockSectorAndIndustry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, res.TotalChecked)
		assert.Equal(t, 1, res.TotalUpdated)
	})
}

func TestStockUsecase_GetAllStocks(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

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
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

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
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

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
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		repo.On("GetStockByCode", ctx, "UNKN").Return((*models.StockDetailResponse)(nil), nil).Once()

		result, err := u.GetStockByCode(ctx, "UNKN")
		assert.NoError(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

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
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		expected := []string{"Acceleration", "Development", "Main"}
		repo.On("GetDistinctBoards", ctx).Return(expected, nil).Once()

		result, err := u.GetDistinctBoards(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockStockRepo)
		pSvc := new(mockPasardanaService)
		u := NewStockUsecase(repo, pSvc)

		repo.On("GetDistinctBoards", ctx).Return(nil, assert.AnError).Once()

		result, err := u.GetDistinctBoards(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}
