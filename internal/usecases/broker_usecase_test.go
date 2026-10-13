package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockBrokerRepo struct {
	mock.Mock
}

func (m *mockBrokerRepo) BatchInsertBrokerActivity(ctx context.Context, records []models.BrokerActivity) error {
	args := m.Called(ctx, records)
	return args.Error(0)
}

func (m *mockBrokerRepo) InsertBrokerActivity(ctx context.Context, record models.BrokerActivity) (bool, error) {
	args := m.Called(ctx, record)
	return args.Bool(0), args.Error(1)
}

func (m *mockBrokerRepo) CheckPartitionExists(ctx context.Context, tableName string) (bool, error) {
	args := m.Called(ctx, tableName)
	return args.Bool(0), args.Error(1)
}

func (m *mockBrokerRepo) CreatePartition(ctx context.Context, tableName, startDate, endDate string) error {
	args := m.Called(ctx, tableName, startDate, endDate)
	return args.Error(0)
}

type mockBrokerService struct {
	mock.Mock
}

func (m *mockBrokerService) FetchBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) (*models.ExodusBrokerActivityResponse, error) {
	args := m.Called(ctx, token, params)
	return args.Get(0).(*models.ExodusBrokerActivityResponse), args.Error(1)
}

func TestSyncBrokerActivity(t *testing.T) {
	mockRepo := new(mockBrokerRepo)
	mockSvc := new(mockBrokerService)
	usecase := NewBrokerUsecase(mockRepo, mockSvc)
	ctx := context.Background()

	t.Run("success_sync_with_filtering", func(t *testing.T) {
		token := "token"
		params := models.SyncBrokerActivityParams{}
		resp := &models.ExodusBrokerActivityResponse{
			Data: struct {
				BrokerActivityTransaction struct {
					BrokersBuy  []models.ExodusBrokerActivityItem `json:"brokers_buy"`
					BrokersSell []models.ExodusBrokerActivityItem `json:"brokers_sell"`
				} `json:"broker_activity_transaction"`
			}{
				BrokerActivityTransaction: struct {
					BrokersBuy  []models.ExodusBrokerActivityItem `json:"brokers_buy"`
					BrokersSell []models.ExodusBrokerActivityItem `json:"brokers_sell"`
				}{
					BrokersBuy: []models.ExodusBrokerActivityItem{
						{BrokerCode: "CC", StockCode: "BBCA", Date: "2024-04-17", Lot: 10, Value: 100},       // Valid
						{BrokerCode: "CC", StockCode: "PTRODRCM6A", Date: "2024-04-17", Lot: 10, Value: 100}, // > 4 chars
						{BrokerCode: "CC", StockCode: "X-IDXS", Date: "2024-04-17", Lot: 10, Value: 100},     // Starts with X
					},
					BrokersSell: []models.ExodusBrokerActivityItem{},
				},
			},
		}

		mockSvc.On("FetchBrokerActivity", ctx, token, params).Return(resp, nil)
		mockRepo.On("InsertBrokerActivity", ctx, mock.MatchedBy(func(rec models.BrokerActivity) bool {
			return rec.StockCode == "BBCA"
		})).Return(true, nil)

		results, err := usecase.SyncBrokerActivity(ctx, token, params)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "BBCA", results[0].StockCode)

		mockSvc.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fetch_error", func(t *testing.T) {
		mockSvc.On("FetchBrokerActivity", ctx, "err-token", mock.Anything).
			Return((*models.ExodusBrokerActivityResponse)(nil), assert.AnError).Once()

		results, err := usecase.SyncBrokerActivity(ctx, "err-token", models.SyncBrokerActivityParams{})
		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("insert_error_and_sell_side", func(t *testing.T) {
		token := "token2"
		params := models.SyncBrokerActivityParams{}
		resp := &models.ExodusBrokerActivityResponse{
			Data: struct {
				BrokerActivityTransaction struct {
					BrokersBuy  []models.ExodusBrokerActivityItem `json:"brokers_buy"`
					BrokersSell []models.ExodusBrokerActivityItem `json:"brokers_sell"`
				} `json:"broker_activity_transaction"`
			}{
				BrokerActivityTransaction: struct {
					BrokersBuy  []models.ExodusBrokerActivityItem `json:"brokers_buy"`
					BrokersSell []models.ExodusBrokerActivityItem `json:"brokers_sell"`
				}{
					BrokersBuy: []models.ExodusBrokerActivityItem{},
					BrokersSell: []models.ExodusBrokerActivityItem{
						{BrokerCode: "YP", StockCode: "BBRI", Date: "2024-04-17", Lot: 5, Value: 50},
					},
				},
			},
		}

		mockSvc.On("FetchBrokerActivity", ctx, token, params).Return(resp, nil).Once()
		mockRepo.On("InsertBrokerActivity", ctx, mock.Anything).Return(false, assert.AnError).Once()

		results, err := usecase.SyncBrokerActivity(ctx, token, params)
		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestManagePartitions(t *testing.T) {
	mockRepo := new(mockBrokerRepo)
	mockSvc := new(mockBrokerService)
	usecase := NewBrokerUsecase(mockRepo, mockSvc)
	ctx := context.Background()

	t.Run("success_create_and_exists", func(t *testing.T) {
		mockRepo.On("CheckPartitionExists", ctx, mock.Anything).Return(false, nil).Once()
		mockRepo.On("CreatePartition", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockRepo.On("CheckPartitionExists", ctx, mock.Anything).Return(true, nil).Once()

		resp, err := usecase.ManagePartitions(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.PartitionsCreated)
		assert.Len(t, resp.Details, 2)
		assert.Equal(t, "created", resp.Details[0].Status)
		assert.Equal(t, "exists", resp.Details[1].Status)
	})

	t.Run("check_partition_error", func(t *testing.T) {
		mockRepo.On("CheckPartitionExists", ctx, mock.Anything).Return(false, assert.AnError).Once()

		resp, err := usecase.ManagePartitions(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("create_partition_error", func(t *testing.T) {
		mockRepo.On("CheckPartitionExists", ctx, mock.Anything).Return(false, nil).Once()
		mockRepo.On("CreatePartition", ctx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		resp, err := usecase.ManagePartitions(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
