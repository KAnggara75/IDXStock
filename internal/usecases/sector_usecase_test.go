package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockSectorSearchRepo struct {
	mock.Mock
}

func (m *mockSectorSearchRepo) UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponseWithCode, error) {
	args := m.Called(ctx, sectors)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BasicResponseWithCode), args.Error(1)
}

func (m *mockSectorSearchRepo) UpsertNewSectorSub(ctx context.Context, sectorSubs []models.SectorSub) ([]models.BasicResponseWithCode, error) {
	args := m.Called(ctx, sectorSubs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BasicResponseWithCode), args.Error(1)
}

func TestSectorUsecase_SyncNewSectors(t *testing.T) {
	ctx := context.Background()

	t.Run("fetch_sectors_error", func(t *testing.T) {
		repo := new(mockSectorSearchRepo)
		pSvc := new(mockPasardanaService)
		u := NewSectorUsecase(repo, pSvc)

		pSvc.On("FetchNewSectors").Return(nil, assert.AnError).Once()

		resp, err := u.SyncNewSectors(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("upsert_sectors_error", func(t *testing.T) {
		repo := new(mockSectorSearchRepo)
		pSvc := new(mockPasardanaService)
		u := NewSectorUsecase(repo, pSvc)

		pSvc.On("FetchNewSectors").Return([]models.PasardanaNewSector{{Id: 1, Name: "Fin"}}, nil).Once()
		repo.On("UpsertNewSectors", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		resp, err := u.SyncNewSectors(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("fetch_subsectors_error", func(t *testing.T) {
		repo := new(mockSectorSearchRepo)
		pSvc := new(mockPasardanaService)
		u := NewSectorUsecase(repo, pSvc)

		pSvc.On("FetchNewSectors").Return([]models.PasardanaNewSector{{Id: 1, Name: "Fin"}}, nil).Once()
		repo.On("UpsertNewSectors", ctx, mock.Anything).Return([]models.BasicResponseWithCode{{Id: 1, Name: "Fin"}}, nil).Once()
		pSvc.On("FetchNewSubSectors").Return(nil, assert.AnError).Once()

		resp, err := u.SyncNewSectors(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("upsert_subsectors_error", func(t *testing.T) {
		repo := new(mockSectorSearchRepo)
		pSvc := new(mockPasardanaService)
		u := NewSectorUsecase(repo, pSvc)

		pSvc.On("FetchNewSectors").Return([]models.PasardanaNewSector{{Id: 1, Name: "Fin"}}, nil).Once()
		repo.On("UpsertNewSectors", ctx, mock.Anything).Return([]models.BasicResponseWithCode{{Id: 1, Name: "Fin"}}, nil).Once()
		pSvc.On("FetchNewSubSectors").Return([]models.PasardanaNewSubSector{{Id: 2, Name: "Bank"}}, nil).Once()
		repo.On("UpsertNewSectorSub", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		resp, err := u.SyncNewSectors(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(mockSectorSearchRepo)
		pSvc := new(mockPasardanaService)
		u := NewSectorUsecase(repo, pSvc)

		pSvc.On("FetchNewSectors").Return([]models.PasardanaNewSector{{Id: 1, Name: "Fin"}}, nil).Once()
		repo.On("UpsertNewSectors", ctx, mock.Anything).Return([]models.BasicResponseWithCode{{Id: 1, Name: "Fin"}}, nil).Once()
		pSvc.On("FetchNewSubSectors").Return([]models.PasardanaNewSubSector{{Id: 2, Name: "Bank"}}, nil).Once()
		repo.On("UpsertNewSectorSub", ctx, mock.Anything).Return([]models.BasicResponseWithCode{{Id: 2, Name: "Bank"}}, nil).Once()

		resp, err := u.SyncNewSectors(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Sectors, 1)
		assert.Len(t, resp.SectorSub, 1)
	})
}
