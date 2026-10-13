package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type mockBrokerRepoForSeeder struct {
	mock.Mock
}

func (m *mockBrokerRepoForSeeder) BatchInsertBrokers(ctx context.Context, brokers []models.Broker) error {
	args := m.Called(ctx, brokers)
	return args.Error(0)
}

func TestSeederService_SeedBrokersData(t *testing.T) {
	ctx := context.Background()

	t.Run("NewSeederService", func(t *testing.T) {
		repo := new(mockBrokerRepoForSeeder)
		sDefault := NewSeederService(repo, nil)
		assert.NotNil(t, sDefault)

		sCustom := NewSeederService(repo, &http.Client{})
		assert.NotNil(t, sCustom)
	})

	t.Run("http_error", func(t *testing.T) {

		repo := new(mockBrokerRepoForSeeder)
		client := &http.Client{
			Transport: &roundTripperFunc{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					return nil, assert.AnError
				},
			},
		}
		s := NewSeederService(repo, client)

		err := s.SeedBrokersData(ctx)
		assert.Error(t, err)
	})

	t.Run("non_200_status", func(t *testing.T) {
		repo := new(mockBrokerRepoForSeeder)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := &http.Client{
			Transport: &roundTripperFunc{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					req.URL.Scheme = "http"
					req.URL.Host = server.Listener.Addr().String()
					return http.DefaultTransport.RoundTrip(req)
				},
			},
		}
		s := NewSeederService(repo, client)

		err := s.SeedBrokersData(ctx)
		assert.Error(t, err)
	})

	t.Run("success_seeding", func(t *testing.T) {
		repo := new(mockBrokerRepoForSeeder)
		jsonResponse := `{"data":{"list":[{"code":"CC","name":"Mandiri Sekuritas"}]}}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(jsonResponse))
		}))
		defer server.Close()

		client := &http.Client{
			Transport: &roundTripperFunc{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					req.URL.Scheme = "http"
					req.URL.Host = server.Listener.Addr().String()
					return http.DefaultTransport.RoundTrip(req)
				},
			},
		}
		s := NewSeederService(repo, client)

		repo.On("BatchInsertBrokers", ctx, mock.MatchedBy(func(b []models.Broker) bool {
			return len(b) == 1 && b[0].Code == "CC"
		})).Return(nil).Once()

		err := s.SeedBrokersData(ctx)
		assert.NoError(t, err)
	})

	t.Run("empty_list", func(t *testing.T) {
		repo := new(mockBrokerRepoForSeeder)
		jsonResponse := `{"data":{"list":[]}}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(jsonResponse))
		}))
		defer server.Close()

		client := &http.Client{
			Transport: &roundTripperFunc{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					req.URL.Scheme = "http"
					req.URL.Host = server.Listener.Addr().String()
					return http.DefaultTransport.RoundTrip(req)
				},
			},
		}
		s := NewSeederService(repo, client)

		err := s.SeedBrokersData(ctx)
		assert.NoError(t, err)
	})

}

type roundTripperFunc struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (f *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.roundTrip(req)
}
