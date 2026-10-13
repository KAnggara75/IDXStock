package services

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasardanaService(t *testing.T) {
	t.Run("NewPasardanaService", func(t *testing.T) {
		svcDefault := NewPasardanaService(nil)
		assert.NotNil(t, svcDefault)

		svcCustom := NewPasardanaService(&http.Client{})
		assert.NotNil(t, svcCustom)
	})

	t.Run("FetchStockIDs_success", func(t *testing.T) {

		jsonResp := `[{"Code":"BBCA","Name":"Bank Central Asia"}]`
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
		svc := NewPasardanaService(client)

		stocks, err := svc.FetchStockIDs()
		assert.NoError(t, err)
		assert.Len(t, stocks, 1)
		assert.Equal(t, "BBCA", stocks[0].Code)
	})

	t.Run("FetchStockSearchResult_gzip_success", func(t *testing.T) {
		jsonResp := `[{"Id":1,"Code":"BBCA","NewIndustryId":2,"NewSectorId":3}]`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusOK)

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, _ = gz.Write([]byte(jsonResp))
			_ = gz.Close()

			_, _ = w.Write(buf.Bytes())
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
		svc := NewPasardanaService(client)

		res, err := svc.FetchStockSearchResult()
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "BBCA", res[0].Code)
	})

	t.Run("FetchNewSectors_and_SubSectors_success", func(t *testing.T) {
		jsonResp := `[{"Id":1,"Name":"Finance"}]`
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
		svc := NewPasardanaService(client)

		sec, err := svc.FetchNewSectors()
		assert.NoError(t, err)
		assert.Len(t, sec, 1)

		subsec, err := svc.FetchNewSubSectors()
		assert.NoError(t, err)
		assert.Len(t, subsec, 1)
	})

	t.Run("FetchStockDetailByCode_and_History_success", func(t *testing.T) {
		jsonResp := `{"Code":"BBCA","Name":"Bank Central Asia"}`
		histResp := `[{"Code":"BBCA","AdjustedClosingPrice":9000}]`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if r.URL.Path == "/api/Stock/GetByCode" {
				_, _ = w.Write([]byte(jsonResp))
			} else {
				_, _ = w.Write([]byte(histResp))
			}
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
		svc := NewPasardanaService(client)

		detail, err := svc.FetchStockDetailByCode("BBCA")
		assert.NoError(t, err)
		assert.NotNil(t, detail)
		assert.Equal(t, "BBCA", detail.Code)

		hist, err := svc.FetchStockHistory(2018, 1, 1)
		assert.NoError(t, err)
		assert.Len(t, hist, 1)
	})

	t.Run("http_status_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
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
		svc := NewPasardanaService(client)

		_, err := svc.FetchStockIDs()
		assert.Error(t, err)
	})
}
