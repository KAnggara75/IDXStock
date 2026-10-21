package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"

	"github.com/KAnggara75/IDXStock/internal/database"
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()
	mockExec := database.NewDBExecutor(db)

	mock.ExpectPing()

	app := buildApp(mockExec)
	assert.NotNil(t, app)

	// Test request logging middleware and route handler
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRunMain(t *testing.T) {
	t.Run("successful_execution", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mockExec := database.NewDBExecutor(db)

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

		runMain(fakeConnect, fakeListen, mockExec)
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

	t.Run("startServer_real_invoked", func(t *testing.T) {
		app := fiber.New()
		err := startServer(app, "invalid_port_!")
		assert.Error(t, err)
	})
}

func TestEnsureMasterDatabase(t *testing.T) {
	t.Run("memory_database_skipped", func(t *testing.T) {
		os.Setenv("DATABASE_URL", ":memory:")
		defer os.Unsetenv("DATABASE_URL")
		assert.NotPanics(t, func() {
			ensureMasterDatabase()
		})
	})

	t.Run("existing_file_skipped", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_main_db_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		filePath := filepath.Join(tempDir, "existing.db")
		assert.NoError(t, os.WriteFile(filePath, []byte("existing"), 0644))

		os.Setenv("DATABASE_URL", filePath)
		defer os.Unsetenv("DATABASE_URL")

		assert.NotPanics(t, func() {
			ensureMasterDatabase()
		})
	})

	t.Run("successful_download_and_decompress", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_main_db_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		var buf bytes.Buffer
		enc, err := zstd.NewWriter(&buf)
		assert.NoError(t, err)
		payload := []byte("master-database-content")
		_, err = enc.Write(payload)
		assert.NoError(t, err)
		assert.NoError(t, enc.Close())

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(buf.Bytes())
		}))
		defer server.Close()

		destPath := filepath.Join(tempDir, "downloaded.db")
		os.Setenv("DATABASE_URL", destPath)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", server.URL)
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")

		ensureMasterDatabase()

		content, err := os.ReadFile(destPath)
		assert.NoError(t, err)
		assert.Equal(t, payload, content)
	})

	t.Run("server_not_found_logs_error", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_main_db_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		destPath := filepath.Join(tempDir, "notfound.db")
		os.Setenv("DATABASE_URL", destPath)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", server.URL)
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")

		assert.NotPanics(t, func() {
			ensureMasterDatabase()
		})
	})

	t.Run("corrupted_zstd_logs_error", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_main_db_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid-zstd-bytes"))
		}))
		defer server.Close()

		destPath := filepath.Join(tempDir, "corrupted.db")
		os.Setenv("DATABASE_URL", destPath)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", server.URL)
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")

		assert.NotPanics(t, func() {
			ensureMasterDatabase()
		})
	})
}
