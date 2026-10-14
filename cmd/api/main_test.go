package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	t.Run("default_port", func(t *testing.T) {
		_ = os.Unsetenv("PORT")
		assert.Equal(t, "3000", getPort())
	})

	t.Run("custom_port", func(t *testing.T) {
		_ = os.Setenv("PORT", "8080")
		defer os.Unsetenv("PORT")
		assert.Equal(t, "8080", getPort())
	})
}

func TestBuildApp(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	mock.ExpectPing()

	app := buildApp(mock)
	assert.NotNil(t, app)

	// Test request logging middleware and route handler
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRunSeeder(t *testing.T) {
	t.Run("nil_pool", func(t *testing.T) {
		assert.NotPanics(t, func() {
			runSeeder(nil)
		})
	})

	t.Run("mock_pool", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		assert.NotPanics(t, func() {
			runSeeder(mock)
		})
	})
}

func TestRunMain(t *testing.T) {
	t.Run("successful_execution", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		var connected bool
		var listened bool

		fakeConnect := func() {
			connected = true
		}

		fakeListen := func(app *fiber.App, port string) error {
			listened = true
			assert.Equal(t, "3000", port)
			return nil
		}

		runMain(fakeConnect, fakeListen, mock)
		assert.True(t, connected)
		assert.True(t, listened)
	})

	t.Run("default_pool_branch", func(t *testing.T) {
		var connected bool
		var listened bool

		fakeConnect := func() {
			connected = true
		}

		fakeListen := func(app *fiber.App, port string) error {
			listened = true
			return nil
		}

		runMain(fakeConnect, fakeListen)
		assert.True(t, connected)
		assert.True(t, listened)
	})

	t.Run("startServer_function", func(t *testing.T) {

		app := fiber.New()
		go func() {
			_ = startServer(app, "0")
		}()
		_ = app.Shutdown()
	})
}
