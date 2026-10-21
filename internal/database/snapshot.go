package database

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/sirupsen/logrus"
)

const defaultReleaseDownloadURL = "https://github.com/KAnggara75/IDXStock/releases/latest/download/master_stocks.db.zst"

// EnsureMasterDatabase checks if the master_stocks.db exists locally.
// If missing, it downloads master_stocks.db.zst from GitHub release and decompresses it.
func EnsureMasterDatabase() {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "database/master_stocks.db"
	}

	if strings.HasPrefix(dbPath, ":memory:") || strings.Contains(dbPath, "mode=memory") {
		return
	}

	cleanPath := dbPath
	if idx := strings.Index(cleanPath, "?"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	if _, err := os.Stat(cleanPath); err == nil {
		return
	}

	dir := filepath.Dir(cleanPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logrus.Errorf("Failed to create directory for database: %v", err)
			return
		}
	}

	downloadURL := os.Getenv("MASTER_STOCKS_DOWNLOAD_URL")
	if downloadURL == "" {
		downloadURL = defaultReleaseDownloadURL
	}

	logrus.Infof("Database file not found at '%s'. Downloading from %s...", cleanPath, downloadURL)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, downloadURL, nil)
	if err != nil {
		logrus.Errorf("Failed to create download request: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Failed to download database snapshot from %s: %v", downloadURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Failed to download database snapshot, server returned status: %d", resp.StatusCode)
		return
	}

	tmpFile := cleanPath + ".tmp"
	outFile, err := os.Create(tmpFile)
	if err != nil {
		logrus.Errorf("Failed to create temporary file for database: %v", err)
		return
	}

	zstdReader, err := zstd.NewReader(resp.Body)
	if err != nil {
		outFile.Close()
		os.Remove(tmpFile)
		logrus.Errorf("Failed to initialize zstd decompressor: %v", err)
		return
	}
	defer zstdReader.Close()

	written, err := io.Copy(outFile, zstdReader)
	outFile.Close()
	if err != nil {
		os.Remove(tmpFile)
		logrus.Errorf("Failed during database decompression: %v", err)
		return
	}

	if err := os.Rename(tmpFile, cleanPath); err != nil {
		os.Remove(tmpFile)
		logrus.Errorf("Failed to finalize database file: %v", err)
		return
	}

	logrus.Infof("Successfully downloaded and decompressed database (%d bytes) to '%s'", written, cleanPath)
}
