package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/models"
)

func TestBrokerService_New(t *testing.T) {
	svcDefault := NewBrokerService(nil)
	assert.NotNil(t, svcDefault)

	svcCustom := NewBrokerService(&http.Client{})
	assert.NotNil(t, svcCustom)
}

func TestBrokerService_FetchBrokerActivity(t *testing.T) {

	ctx := context.Background()

	t.Run("http_error", func(t *testing.T) {
		client := &http.Client{
			Transport: &roundTripperFunc{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					return nil, assert.AnError
				},
			},
		}
		svc := NewBrokerService(client)

		resp, err := svc.FetchBrokerActivity(ctx, "Bearer token", models.SyncBrokerActivityParams{BrokerCode: "CC"})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("non_200_status", func(t *testing.T) {
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
		svc := NewBrokerService(client)

		resp, err := svc.FetchBrokerActivity(ctx, "Bearer token", models.SyncBrokerActivityParams{From: "2024-01-01", To: "2024-01-02"})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("success", func(t *testing.T) {
		jsonResp := `{"data":{"broker_activity_transaction":{"brokers_buy":[{"broker_code":"CC","stock_code":"BBCA","date":"2024-01-01","lot":10,"value":100}]}}}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
			assert.Equal(t, "CC", r.URL.Query().Get("broker_code"))
			assert.Equal(t, "RG", r.URL.Query().Get("market_board"))
			assert.Equal(t, "NET", r.URL.Query().Get("transaction_type"))
			assert.Equal(t, "D", r.URL.Query().Get("investor_type"))
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
		svc := NewBrokerService(client)

		resp, err := svc.FetchBrokerActivity(ctx, "Bearer token", models.SyncBrokerActivityParams{
			BrokerCode:      "CC",
			MarketBoard:     "RG",
			TransactionType: "NET",
			InvestorType:    "D",
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Data.BrokerActivityTransaction.BrokersBuy, 1)
	})
}
