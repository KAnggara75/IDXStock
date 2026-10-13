package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestAssetHandler_GetCompanyLogo(t *testing.T) {
	t.Run("success_200", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("fake image data")),
					Header: http.Header{
						"Content-Type":   []string{"image/png"},
						"Content-Length": []string{"15"},
					},
				}, nil
			},
		}

		handler := &AssetHandler{
			client:  mockClient,
			baseURL: "https://test-logos",
		}
		app := fiber.New()
		app.Get("/asset/logo/:code", handler.GetCompanyLogo)

		req := httptest.NewRequest("GET", "/asset/logo/BBCA.png", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	})

	t.Run("not_found_404", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString("not found")),
					Header:     make(http.Header),
				}, nil
			},
		}

		handler := &AssetHandler{
			client:  mockClient,
			baseURL: "https://test-logos",
		}
		app := fiber.New()
		app.Get("/asset/logo/:code", handler.GetCompanyLogo)

		req := httptest.NewRequest("GET", "/asset/logo/UNKNOWN", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("other_status_upstream_error", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(bytes.NewBufferString("bad gateway")),
					Header:     make(http.Header),
				}, nil
			},
		}

		handler := &AssetHandler{
			client:  mockClient,
			baseURL: "https://test-logos",
		}
		app := fiber.New()
		app.Get("/asset/logo/:code", handler.GetCompanyLogo)

		req := httptest.NewRequest("GET", "/asset/logo/BBCA", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadGateway, resp.StatusCode)
	})

	t.Run("client_do_network_error", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("connection refused")
			},
		}

		handler := &AssetHandler{
			client:  mockClient,
			baseURL: "https://test-logos",
		}
		app := fiber.New()
		app.Get("/asset/logo/:code", handler.GetCompanyLogo)

		req := httptest.NewRequest("GET", "/asset/logo/BBCA", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("missing_code", func(t *testing.T) {
		handler := NewAssetHandler()
		app := fiber.New()
		app.Get("/asset/logo/", handler.GetCompanyLogo)

		req := httptest.NewRequest("GET", "/asset/logo/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
