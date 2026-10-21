package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockHistoryUsecase struct {
	mock.Mock
}

func (m *mockHistoryUsecase) SyncStockHistory(ctx context.Context, req models.SyncHistoryRequest, source, cookie string) error {
	args := m.Called(ctx, req, source, cookie)
	return args.Error(0)
}

func (m *mockHistoryUsecase) GetStockHistory(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error) {
	args := m.Called(ctx, code, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockHistory), args.Error(1)
}

func (m *mockHistoryUsecase) ApplyStockSplit(ctx context.Context, req models.StockSplitRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockHistoryUsecase) SyncStockHistoryFromStockbit(ctx context.Context, code string, from, to string, authorization string) (int, error) {
	args := m.Called(ctx, code, from, to, authorization)
	return args.Int(0), args.Error(1)
}

func TestSyncStockHistoryHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockHistoryUsecase)
	handler := NewHistoryHandler(mockUsecase)

	app.Post("/history/sync", handler.SyncStockHistoryHandler)

	t.Run("missing_source", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/sync", bytes.NewBufferString(`{"year":2024,"month":4,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/sync?source=idx", bytes.NewBufferString(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing_dates", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/sync?source=idx", bytes.NewBufferString(`{"year":0,"month":4,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("pasardana_restricted_after_threshold", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/sync?source=pasardana", bytes.NewBufferString(`{"year":2020,"month":1,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("pasardana_valid_background_sync", func(t *testing.T) {
		mockUsecase.On("SyncStockHistory", mock.Anything, models.SyncHistoryRequest{Year: 2018, Month: 1, Day: 1}, "pasardana", "").
			Return(nil).Maybe()

		req := httptest.NewRequest("POST", "/history/sync?source=pasardana", bytes.NewBufferString(`{"year":2018,"month":1,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("idx_success", func(t *testing.T) {
		mockUsecase.On("SyncStockHistory", mock.Anything, models.SyncHistoryRequest{Year: 2024, Month: 4, Day: 1}, "idx", "cookie").
			Return(nil).Once()

		req := httptest.NewRequest("POST", "/history/sync?source=idx", bytes.NewBufferString(`{"year":2024,"month":4,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "cookie")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("idx_error", func(t *testing.T) {
		mockUsecase.On("SyncStockHistory", mock.Anything, models.SyncHistoryRequest{Year: 2024, Month: 4, Day: 1}, "idx", "").
			Return(fmt.Errorf("idx error")).Once()

		req := httptest.NewRequest("POST", "/history/sync?source=idx", bytes.NewBufferString(`{"year":2024,"month":4,"day":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestGetStockHistoryHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockHistoryUsecase)
	handler := NewHistoryHandler(mockUsecase)

	app.Get("/history/:code", handler.GetStockHistoryHandler)

	val := 100.0
	delist := "2024-01-01"
	sampleHistory := []models.StockHistory{
		{
			Code:          "BBCA",
			Date:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Previous:      &val,
			OpenPrice:     &val,
			FirstTrade:    &val,
			High:          &val,
			Low:           &val,
			Close:         &val,
			Change:        &val,
			Volume:        &val,
			Value:         &val,
			Frequency:     &val,
			DelistingDate: &delist,
			LastModified:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			Code:          "BBCA",
			Date:          time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			DelistingDate: nil,
			Previous:      nil,
		},
	}

	t.Run("invalid_start_date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/history/BBCA?start_date=invalid", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid_end_date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/history/BBCA?end_date=invalid", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("end_before_start", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/history/BBCA?start_date=2024-05-01&end_date=2024-04-01", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.On("GetStockHistory", mock.Anything, "BBCA", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("db error")).Once()

		req := httptest.NewRequest("GET", "/history/BBCA?start_date=2024-01-01&end_date=2024-01-02", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("success_json_all_fields", func(t *testing.T) {
		mockUsecase.On("GetStockHistory", mock.Anything, "BBCA", mock.Anything, mock.Anything).
			Return(sampleHistory, nil).Once()

		req := httptest.NewRequest("GET", "/history/BBCA?fields=code,date,previous,open_price,first_trade,high,low,close,change,volume,value,frequency,index_individual,offer,offer_volume,bid,bid_volume,listed_shares,tradeble_shares,weight_for_index,foreign_sell,foreign_buy,delisting_date,non_regular_volume,non_regular_value,non_regular_frequency,last_modified,unknown", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("success_csv", func(t *testing.T) {
		mockUsecase.On("GetStockHistory", mock.Anything, "BBCA", mock.Anything, mock.Anything).
			Return(sampleHistory, nil).Once()

		req := httptest.NewRequest("GET", "/history/BBCA?output=csv&fields=code,date,previous,open_price,first_trade,high,low,close,change,volume,value,frequency,index_individual,offer,offer_volume,bid,bid_volume,listed_shares,tradeble_shares,weight_for_index,foreign_sell,foreign_buy,delisting_date,non_regular_volume,non_regular_value,non_regular_frequency,last_modified,unknown", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
		mockUsecase.AssertExpectations(t)
	})
}

func TestApplyStockSplitHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockHistoryUsecase)
	handler := NewHistoryHandler(mockUsecase)

	app.Post("/history/split", handler.ApplyStockSplitHandler)

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/split", bytes.NewBufferString(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing_fields", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/history/split", bytes.NewBufferString(`{"stock_code":""}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("usecase_error", func(t *testing.T) {
		splitReq := models.StockSplitRequest{
			StockCode: "BBCA",
			StartDate: "2020-01-01",
			EndDate:   "2021-01-01",
			Ratio:     "1:5",
		}
		mockUsecase.On("ApplyStockSplit", mock.Anything, splitReq).Return(fmt.Errorf("split err")).Once()

		req := httptest.NewRequest("POST", "/history/split", bytes.NewBufferString(`{"stock_code":"BBCA","start_date":"2020-01-01","end_date":"2021-01-01","ratio":"1:5"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		splitReq := models.StockSplitRequest{
			StockCode: "BBCA",
			StartDate: "2020-01-01",
			EndDate:   "2021-01-01",
			Ratio:     "1:5",
		}
		mockUsecase.On("ApplyStockSplit", mock.Anything, splitReq).Return(nil).Once()

		req := httptest.NewRequest("POST", "/history/split", bytes.NewBufferString(`{"stock_code":"BBCA","start_date":"2020-01-01","end_date":"2021-01-01","ratio":"1:5"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestSyncStockHistoryByCodeHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockHistoryUsecase)
	handler := NewHistoryHandler(mockUsecase)

	app.Put("/stocks/:code/history/sync", handler.SyncStockHistoryByCodeHandler)

	t.Run("missing_authorization", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/stocks/BBCA/history/sync", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.On("SyncStockHistoryFromStockbit", mock.Anything, "BBCA", "", "", "Bearer tok").
			Return(0, fmt.Errorf("sync err")).Once()

		req := httptest.NewRequest("PUT", "/stocks/BBCA/history/sync", nil)
		req.Header.Set("Authorization", "Bearer tok")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("SyncStockHistoryFromStockbit", mock.Anything, "BBCA", "2026-10-20", "2000-01-01", "Bearer tok").
			Return(100, nil).Once()

		req := httptest.NewRequest("PUT", "/stocks/BBCA/history/sync?from=2026-10-20&to=2000-01-01", nil)
		req.Header.Set("Authorization", "Bearer tok")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
