package routes

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetup_And_Health(t *testing.T) {
	t.Run("health_up", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectPing()

		app := fiber.New()
		Setup(app, mock)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("health_error_ping", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		mock.ExpectPing().WillReturnError(assert.AnError)

		app := fiber.New()
		Setup(app, mock)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("health_nil_pool", func(t *testing.T) {
		app := fiber.New()
		Setup(app, nil)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
