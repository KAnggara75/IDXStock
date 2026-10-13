package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockIndustryUsecase struct {
	mock.Mock
}

func (m *mockIndustryUsecase) SyncIndustry(ctx context.Context) (*models.IndustrySyncResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IndustrySyncResponse), args.Error(1)
}

func TestIndustrySyncHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockIndustryUsecase)
	handler := NewIndustryHandler(mockUsecase)

	app.Get("/industries/sync", handler.IndustrySyncHandler)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("SyncIndustry", mock.Anything).Return(&models.IndustrySyncResponse{
			Industries: []models.BasicResponse{
				{Id: 1, Name: "Industry 1"},
			},
		}, nil).Once()

		req := httptest.NewRequest("GET", "/industries/sync", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUsecase.On("SyncIndustry", mock.Anything).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest("GET", "/industries/sync", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
