package routes

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/database"
)

func TestSetup_And_Health(t *testing.T) {
	t.Run("health_up", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPing()

		executor := database.NewDBExecutor(db)
		app := fiber.New()
		Setup(app, executor)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("health_error_ping", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPing().WillReturnError(assert.AnError)

		executor := database.NewDBExecutor(db)
		app := fiber.New()
		Setup(app, executor)

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

	t.Run("swagger_route", func(t *testing.T) {
		app := fiber.New()
		Setup(app, nil)

		req := httptest.NewRequest("GET", "/swagger/index.html", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("health_context_timeout", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPing().WillReturnError(context.DeadlineExceeded)

		executor := database.NewDBExecutor(db)
		app := fiber.New()
		Setup(app, executor)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
