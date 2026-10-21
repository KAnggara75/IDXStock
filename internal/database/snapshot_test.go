package database

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
)

func TestEnsureMasterDatabase(t *testing.T) {
	t.Run("memory_database_skipped", func(t *testing.T) {
		os.Setenv("DATABASE_URL", ":memory:")
		defer os.Unsetenv("DATABASE_URL")
		assert.NotPanics(t, func() {
			EnsureMasterDatabase()
		})
	})

	t.Run("existing_file_skipped", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_db_snapshot_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		filePath := filepath.Join(tempDir, "existing.db")
		assert.NoError(t, os.WriteFile(filePath, []byte("existing"), 0644))

		os.Setenv("DATABASE_URL", filePath)
		defer os.Unsetenv("DATABASE_URL")

		assert.NotPanics(t, func() {
			EnsureMasterDatabase()
		})
	})

	t.Run("successful_download_and_decompress", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_db_snapshot_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		var buf bytes.Buffer
		enc, err := zstd.NewWriter(&buf)
		assert.NoError(t, err)
		_, err = enc.Write([]byte("mock decompressed sqlite database"))
		assert.NoError(t, err)
		assert.NoError(t, enc.Close())

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(buf.Bytes())
		}))
		defer mockServer.Close()

		targetDB := filepath.Join(tempDir, "downloaded.db")
		os.Setenv("DATABASE_URL", targetDB)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", mockServer.URL)
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")
		}()

		EnsureMasterDatabase()

		content, err := os.ReadFile(targetDB)
		assert.NoError(t, err)
		assert.Equal(t, "mock decompressed sqlite database", string(content))
	})

	t.Run("server_not_found_logs_error", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_db_snapshot_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer mockServer.Close()

		targetDB := filepath.Join(tempDir, "notfound.db")
		os.Setenv("DATABASE_URL", targetDB)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", mockServer.URL)
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")
		}()

		assert.NotPanics(t, func() {
			EnsureMasterDatabase()
		})

		_, err = os.Stat(targetDB)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("corrupted_zstd_logs_error", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_db_snapshot_*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("not a valid zstd payload"))
		}))
		defer mockServer.Close()

		targetDB := filepath.Join(tempDir, "corrupted.db")
		os.Setenv("DATABASE_URL", targetDB)
		os.Setenv("MASTER_STOCKS_DOWNLOAD_URL", mockServer.URL)
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("MASTER_STOCKS_DOWNLOAD_URL")
		}()

		assert.NotPanics(t, func() {
			EnsureMasterDatabase()
		})

		_, err = os.Stat(targetDB)
		assert.True(t, os.IsNotExist(err))
	})
}
