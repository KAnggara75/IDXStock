package main

//go:generate swag init -d ./,../../internal/handlers,../../internal/models -g main.go -o ../../docs

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/klauspost/compress/zstd"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/routes"
)

// @title                      IDXStock API
// @version                    1.0
// @description                Automated stock market data integration and analytics API for the Indonesia Stock Exchange (IDX).
// @contact.name               IDXStock Support
// @contact.url                https://github.com/KAnggara75/IDXStock
// @license.name               MIT
// @BasePath                   /api/v1

const defaultReleaseDownloadURL = "https://github.com/KAnggara75/IDXStock/releases/latest/download/master_stocks.db.zst"

func init() {
	config.InitLogger()
	ensureMasterDatabase()
}

// ensureMasterDatabase checks if the master_stocks.db exists locally.
// If missing, it downloads master_stocks.db.zst from GitHub release and decompresses it.
func ensureMasterDatabase() {
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

// buildApp creates and configures the Fiber application and its middlewares/routes.
func buildApp(pool ...repositories.DBExecutor) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "IDXStock API",
	})

	// Middlewares
	app.Use(recover.New())

	// Custom Logger Middleware using logrus
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		logrus.WithFields(logrus.Fields{
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  c.Response().StatusCode(),
			"latency": stop.Sub(start).String(),
			"ip":      c.IP(),
		}).Info("HTTP Request")

		return err
	})

	// Setup Routes
	routes.Setup(app, pool...)

	return app
}

// getPort retrieves the port from environment variable or defaults to 3000.
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

// startServer starts listening on the specified port.
var startServer = func(app *fiber.App, port string) error {
	logrus.Infof("Server starting on port %s", port)
	return app.Listen(":" + port)
}

// runMain encapsulates main logic with customizable connect and listen functions.
func runMain(connectDB func(), listen func(app *fiber.App, port string) error, pool ...repositories.DBExecutor) {
	config.InitLogger()
	connectDB()

	var dbPool repositories.DBExecutor
	if len(pool) > 0 {
		dbPool = pool[0]
	} else if database.Pool != nil {
		dbPool = database.Pool
	}

	app := buildApp(dbPool)
	port := getPort()
	if err := listen(app, port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	runMain(database.Connect, startServer)
}
