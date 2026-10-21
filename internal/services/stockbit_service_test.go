package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStockbitService(t *testing.T) {
	ctx := context.Background()

	t.Run("NewStockbitService", func(t *testing.T) {
		svc := NewStockbitService(nil)
		assert.NotNil(t, svc)
	})

	t.Run("FetchDailyPrice_success", func(t *testing.T) {
		jsonResp := `{"message":"success","data":{"chartbit":[{"date":"2026-10-20","open":4740,"high":5200,"low":4740,"close":5200,"volume":1131013100}]}}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(jsonResp))
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
		svc := NewStockbitService(client)

		items, err := svc.FetchDailyPrice(ctx, "BBCA", "2026-10-20", "2000-01-01", "Bearer test-token")
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "2026-10-20", items[0].Date)
		assert.Equal(t, 4740.0, *items[0].Open)
	})

	t.Run("FetchDailyPrice_http_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
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
		svc := NewStockbitService(client)

		items, err := svc.FetchDailyPrice(ctx, "BBCA", "2026-10-20", "2000-01-01", "Bearer invalid")
		assert.Error(t, err)
		assert.Nil(t, items)
	})

	t.Run("FetchDailyPrice_decode_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
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
		svc := NewStockbitService(client)

		items, err := svc.FetchDailyPrice(ctx, "BBCA", "2026-10-20", "2000-01-01", "")
		assert.Error(t, err)
		assert.Nil(t, items)
	})
}
