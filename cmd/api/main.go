package main

//go:generate swag init -g main.go -d ./,../../internal -o ../../docs

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/scc2go"

	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/routes"
	"github.com/KAnggara75/IDXStock/internal/services"
)

// @title                      IDXStock API
// @version                    1.0
// @description                Automated stock market data integration and analytics API for the Indonesia Stock Exchange (IDX).
// @contact.name               IDXStock Support
// @contact.url                https://github.com/KAnggara75/IDXStock
// @license.name               MIT
// @BasePath                   /api/v1

func init() {
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("SCC_AUTH"))
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

// runSeeder executes initial data seeders.
func runSeeder(pool repositories.DBExecutor) {
	if pool == nil {
		return
	}
	brokerRepo := repositories.NewBrokerRepository(pool)
	seederService := services.NewSeederService(brokerRepo, nil)
	if err := seederService.SeedBrokersData(context.Background()); err != nil {
		logrus.Errorf("Failed to run brokers seeder: %v", err)
	}
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

	runSeeder(dbPool)
	app := buildApp(dbPool)
	port := getPort()
	if err := listen(app, port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	runMain(database.Connect, startServer)
}
