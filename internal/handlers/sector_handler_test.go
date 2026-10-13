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

type mockSectorUsecase struct {
	mock.Mock
}

func (m *mockSectorUsecase) SyncNewSectors(ctx context.Context) (*models.SectorSyncNewResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SectorSyncNewResponse), args.Error(1)
}

func TestSyncNewSectorsHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockSectorUsecase)
	handler := NewSectorHandler(mockUsecase)

	app.Get("/sectors/sync", handler.SyncNewSectorsHandler)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("SyncNewSectors", mock.Anything).Return(&models.SectorSyncNewResponse{
			Sectors: []models.BasicResponseWithCode{
				{Id: 1, Name: "Sector 1"},
			},
		}, nil).Once()

		req := httptest.NewRequest("GET", "/sectors/sync", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUsecase.On("SyncNewSectors", mock.Anything).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest("GET", "/sectors/sync", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
