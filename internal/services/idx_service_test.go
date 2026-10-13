package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdxService(t *testing.T) {
	ctx := context.Background()

	t.Run("ParseIdxDate", func(t *testing.T) {
		svc := NewIdxService(nil)
		date, err := svc.ParseIdxDate("18 July 2025")
		assert.NoError(t, err)
		assert.Equal(t, "2025-07-18", date)

		_, err = svc.ParseIdxDate("invalid-date")
		assert.Error(t, err)
	})

	t.Run("FetchDelistedStocks_success", func(t *testing.T) {
		jsonResp := `{"data":[{"Code":"BBCA","DeListingDate":"2024-04-17"}]}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "my-cookie", r.Header.Get("Cookie"))
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
		svc := NewIdxService(client)

		stocks, err := svc.FetchDelistedStocks(ctx, 2024, 4, "my-cookie")
		assert.NoError(t, err)
		assert.Len(t, stocks, 1)
		assert.Equal(t, "BBCA", stocks[0].Code)
	})

	t.Run("FetchDelistedStocks_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
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
		svc := NewIdxService(client)

		stocks, err := svc.FetchDelistedStocks(ctx, 2024, 4, "")
		assert.Error(t, err)
		assert.Nil(t, stocks)
	})

	t.Run("FetchDelistedStocks_json_error", func(t *testing.T) {
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
		svc := NewIdxService(client)

		stocks, err := svc.FetchDelistedStocks(ctx, 2024, 4, "")
		assert.Error(t, err)
		assert.Nil(t, stocks)
	})

	t.Run("FetchStockSummary_success", func(t *testing.T) {
		jsonResp := `{"data":[{"StockCode":"BBRI","Close":5000}]}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		svc := NewIdxService(client)

		data, err := svc.FetchStockSummary(ctx, 2024, 4, 1, "")
		assert.NoError(t, err)
		assert.Len(t, data, 1)
		assert.Equal(t, "BBRI", data[0].StockCode)
	})

	t.Run("FetchStockSummary_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
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
		svc := NewIdxService(client)

		data, err := svc.FetchStockSummary(ctx, 2024, 4, 1, "")
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}
