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

type mockBrokerUsecase struct {
	mock.Mock
}

func (m *mockBrokerUsecase) SyncBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) ([]models.BrokerActivity, error) {
	args := m.Called(ctx, token, params)
	return args.Get(0).([]models.BrokerActivity), args.Error(1)
}

func (m *mockBrokerUsecase) ManagePartitions(ctx context.Context) (*models.PartitionManagementResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.PartitionManagementResponse), args.Error(1)
}

func TestSyncBrokerActivityHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockBrokerUsecase)
	handler := NewBrokerHandler(mockUsecase)

	app.Get("/broker/sync", handler.SyncBrokerActivityHandler)

	t.Run("unauthorized_missing_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/broker/sync", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("success_sync", func(t *testing.T) {
		mockUsecase.On("SyncBrokerActivity", mock.Anything, "Bearer test-token", mock.Anything).
			Return([]models.BrokerActivity{}, nil)

		req := httptest.NewRequest("GET", "/broker/sync", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.On("SyncBrokerActivity", mock.Anything, "Bearer test-token-err", mock.Anything).
			Return([]models.BrokerActivity{}, assert.AnError)

		req := httptest.NewRequest("GET", "/broker/sync", nil)
		req.Header.Set("Authorization", "Bearer test-token-err")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestManagePartitionsHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockBrokerUsecase)
	handler := NewBrokerHandler(mockUsecase)

	app.Post("/broker/partitions", handler.ManagePartitionsHandler)

	t.Run("success_created", func(t *testing.T) {
		respData := &models.PartitionManagementResponse{
			PartitionsCreated: 2,
			Details: []models.PartitionDetail{
				{Name: "p1", Status: "created"},
				{Name: "p2", Status: "created"},
			},
		}
		mockUsecase.On("ManagePartitions", mock.Anything).Return(respData, nil).Once()

		req := httptest.NewRequest("POST", "/broker/partitions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("success_ok_no_partitions_created", func(t *testing.T) {
		respData := &models.PartitionManagementResponse{
			PartitionsCreated: 0,
			Details: []models.PartitionDetail{
				{Name: "p1", Status: "exists"},
			},
		}
		mockUsecase.On("ManagePartitions", mock.Anything).Return(respData, nil).Once()

		req := httptest.NewRequest("POST", "/broker/partitions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("error_usecase", func(t *testing.T) {
		mockUsecase.On("ManagePartitions", mock.Anything).Return((*models.PartitionManagementResponse)(nil), assert.AnError).Once()

		req := httptest.NewRequest("POST", "/broker/partitions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
