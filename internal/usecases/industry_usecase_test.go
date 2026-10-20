package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockIndustryRepo struct {
	mock.Mock
}

func (m *mockIndustryRepo) UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error) {
	args := m.Called(ctx, industries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BasicResponse), args.Error(1)
}

func (m *mockIndustryRepo) UpsertIndustrySub(ctx context.Context, industrySubs []models.IndustrySub) ([]models.BasicResponse, error) {
	args := m.Called(ctx, industrySubs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BasicResponse), args.Error(1)
}

func TestIndustryUsecase_SyncIndustry(t *testing.T) {
	ctx := context.Background()

	t.Run("fetch_error", func(t *testing.T) {
		repo := new(mockIndustryRepo)
		pSvc := new(mockPasardanaService)
		u := NewIndustryUsecase(repo, pSvc)

		pSvc.On("FetchStockSearchResult").Return(nil, assert.AnError).Once()

		resp, err := u.SyncIndustry(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("upsert_industries_error", func(t *testing.T) {
		repo := new(mockIndustryRepo)
		pSvc := new(mockPasardanaService)
		u := NewIndustryUsecase(repo, pSvc)

		pSvc.On("FetchStockSearchResult").Return([]models.PasardanaSearchResult{
			{NewIndustryId: 1, NewIndustryName: "Tech", NewSubIndustryId: 2, NewSubIndustryName: "Software"},
		}, nil).Once()

		repo.On("UpsertIndustries", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		resp, err := u.SyncIndustry(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("upsert_sub_industries_error", func(t *testing.T) {
		repo := new(mockIndustryRepo)
		pSvc := new(mockPasardanaService)
		u := NewIndustryUsecase(repo, pSvc)

		pSvc.On("FetchStockSearchResult").Return([]models.PasardanaSearchResult{
			{NewIndustryId: 1, NewIndustryName: "Tech", NewSubIndustryId: 2, NewSubIndustryName: "Software"},
		}, nil).Once()

		repo.On("UpsertIndustries", ctx, mock.Anything).Return([]models.BasicResponse{{Id: 1, Name: "Tech"}}, nil).Once()
		repo.On("UpsertIndustrySub", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		resp, err := u.SyncIndustry(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(mockIndustryRepo)
		pSvc := new(mockPasardanaService)
		u := NewIndustryUsecase(repo, pSvc)

		pSvc.On("FetchStockSearchResult").Return([]models.PasardanaSearchResult{
			{NewIndustryId: 1, NewIndustryName: "Tech", NewSubIndustryId: 2, NewSubIndustryName: "Software"},
		}, nil).Once()

		repo.On("UpsertIndustries", ctx, mock.Anything).Return([]models.BasicResponse{{Id: 1, Name: "Tech"}}, nil).Once()
		repo.On("UpsertIndustrySub", ctx, mock.Anything).Return([]models.BasicResponse{{Id: 2, Name: "Software"}}, nil).Once()

		resp, err := u.SyncIndustry(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Industries, 1)
		assert.Len(t, resp.IndustrySub, 1)
	})
}
