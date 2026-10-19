package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockStockUsecase struct {
	mock.Mock
}

func (m *mockStockUsecase) PreviewStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	args := m.Called(ctx, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *mockStockUsecase) UploadStocks(ctx context.Context, file io.Reader) ([]models.Stock, error) {
	args := m.Called(ctx, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *mockStockUsecase) SyncStockDetail(ctx context.Context) ([]models.StockResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockResponse), args.Error(1)
}

func (m *mockStockUsecase) SyncDelistingStocks(ctx context.Context, year, month int, cookie string) ([]models.StockResponse, error) {
	args := m.Called(ctx, year, month, cookie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockResponse), args.Error(1)
}

func (m *mockStockUsecase) GetAllStocks(ctx context.Context, filter models.GetAllStocksFilter) ([]models.StockDetailResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StockDetailResponse), args.Error(1)
}

func createMultipartRequest(url, fieldName, fileName string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if fileName != "" {
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(fileContent); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestStockHandler_PreviewAndUpload(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockStockUsecase)
	handler := NewStockHandler(mockUsecase)

	app.Post("/stocks/preview", handler.PreviewHandler)
	app.Post("/stocks/upload", handler.UploadHandler)

	t.Run("preview_no_file", func(t *testing.T) {
		req, _ := createMultipartRequest("/stocks/preview", "wrong_field", "", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("preview_invalid_extension", func(t *testing.T) {
		req, _ := createMultipartRequest("/stocks/preview", "file", "stocks.txt", []byte("hello"))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("preview_success", func(t *testing.T) {
		mockUsecase.On("PreviewStocks", mock.Anything, mock.Anything).
			Return([]models.Stock{{Code: "BBCA"}}, nil).Once()

		req, _ := createMultipartRequest("/stocks/preview", "file", "stocks.xlsx", []byte("dummy excel"))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("preview_usecase_error", func(t *testing.T) {
		mockUsecase.On("PreviewStocks", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("parse error")).Once()

		req, _ := createMultipartRequest("/stocks/preview", "file", "stocks.xls", []byte("dummy excel"))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("upload_no_file", func(t *testing.T) {
		req, _ := createMultipartRequest("/stocks/upload", "wrong_field", "", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("upload_success", func(t *testing.T) {
		mockUsecase.On("UploadStocks", mock.Anything, mock.Anything).
			Return([]models.Stock{{Code: "BBCA"}}, nil).Once()

		req, _ := createMultipartRequest("/stocks/upload", "file", "stocks.xlsx", []byte("dummy excel"))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("upload_usecase_error", func(t *testing.T) {
		mockUsecase.On("UploadStocks", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("upload error")).Once()

		req, _ := createMultipartRequest("/stocks/upload", "file", "stocks.xlsx", []byte("dummy excel"))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestStockHandler_SyncStockDetailHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockStockUsecase)
	handler := NewStockHandler(mockUsecase)

	app.Post("/stocks/sync-detail", handler.SyncStockDetailHandler)

	mockUsecase.On("SyncStockDetail", mock.Anything).Return([]models.StockResponse{}, nil).Maybe()

	req := httptest.NewRequest("POST", "/stocks/sync-detail", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestStockHandler_SyncDelistingStocksHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockStockUsecase)
	handler := NewStockHandler(mockUsecase)

	app.Post("/stocks/sync-delisting", handler.SyncDelistingStocksHandler)

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/stocks/sync-delisting", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing_year_or_month", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/stocks/sync-delisting", bytes.NewBufferString(`{"year":0,"month":5}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("SyncDelistingStocks", mock.Anything, 2024, 4, "test-cookie").
			Return([]models.StockResponse{{Code: "BBCA"}}, nil).Once()

		req := httptest.NewRequest("POST", "/stocks/sync-delisting", bytes.NewBufferString(`{"year":2024,"month":4}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "test-cookie")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.On("SyncDelistingStocks", mock.Anything, 2024, 4, "").
			Return(nil, fmt.Errorf("sync err")).Once()

		req := httptest.NewRequest("POST", "/stocks/sync-delisting", bytes.NewBufferString(`{"year":2024,"month":4}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestStockHandler_GetAllStocksHandler(t *testing.T) {
	mockUsecase := new(mockStockUsecase)
	handler := NewStockHandler(mockUsecase)

	app := fiber.New()
	app.Get("/stocks", handler.GetAllStocksHandler)

	t.Run("success", func(t *testing.T) {
		expected := []models.StockDetailResponse{
			{Code: "BBCA", Name: "Bank Central Asia", Board: "UTAMA"},
		}
		mockUsecase.On("GetAllStocks", mock.Anything, models.GetAllStocksFilter{
			Board:  "UTAMA",
			Search: "BBCA",
		}).Return(expected, nil).Once()

		req := httptest.NewRequest("GET", "/stocks?board=UTAMA&search=BBCA", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.On("GetAllStocks", mock.Anything, models.GetAllStocksFilter{}).
			Return(nil, fmt.Errorf("db error")).Once()

		req := httptest.NewRequest("GET", "/stocks", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
