package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/handlers"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/services"
	"github.com/KAnggara75/IDXStock/internal/usecases"
)

func Setup(app *fiber.App, db ...repositories.DBExecutor) {
	var pool repositories.DBExecutor
	if len(db) > 0 {
		pool = db[0]
	} else if database.Pool != nil {
		pool = database.Pool
	}

	// Dependency Injection
	stockRepo := repositories.NewStockRepository(pool)
	sectorSearchRepo := repositories.NewSectorSearchRepository(pool)
	industryRepo := repositories.NewIndustryRepository(pool)
	historyRepo := repositories.NewHistoryRepository(pool)
	brokerActRepo := repositories.NewBrokerActivityRepository(pool)

	stockService := services.NewStockService()
	pasardanaService := services.NewPasardanaService(nil)
	idxService := services.NewIdxService(nil)
	brokerService := services.NewBrokerService(nil)

	stockUsecase := usecases.NewStockUsecase(stockRepo, stockService, pasardanaService, idxService)
	industryUsecase := usecases.NewIndustryUsecase(industryRepo, pasardanaService)
	sectorUsecase := usecases.NewSectorUsecase(sectorSearchRepo, pasardanaService)
	historyUsecase := usecases.NewHistoryUsecase(historyRepo, stockRepo, pasardanaService, idxService)
	brokerUsecase := usecases.NewBrokerUsecase(brokerActRepo, brokerService)

	stockHandler := handlers.NewStockHandler(stockUsecase)
	industryHandler := handlers.NewIndustryHandler(industryUsecase)
	sectorHandler := handlers.NewSectorHandler(sectorUsecase)
	historyHandler := handlers.NewHistoryHandler(historyUsecase)
	brokerHandler := handlers.NewBrokerHandler(brokerUsecase)
	assetHandler := handlers.NewAssetHandler()

	app.Get("/health", func(c fiber.Ctx) error {
		if pool == nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "database connection failed",
			})
		}
		err := pool.Ping(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status": "up",
		})
	})

	// API V1 Routes
	v1 := app.Group("/api/v1")
	v1.Post("/stocks/upload", stockHandler.PreviewHandler)
	v1.Patch("/stocks/upload", stockHandler.UploadHandler)
	v1.Put("/stocks/sync", stockHandler.SyncStockDetailHandler)
	v1.Put("/stocks/delisting/sync", stockHandler.SyncDelistingStocksHandler)
	v1.Put("/stocks/history/sync", historyHandler.SyncStockHistoryHandler)
	v1.Get("/stocks/:code/history", historyHandler.GetStockHistoryHandler)
	v1.Put("/stocks/split", historyHandler.ApplyStockSplitHandler)
	v1.Put("/sectors/sync", sectorHandler.SyncNewSectorsHandler)
	v1.Put("/industries/sync", industryHandler.IndustrySyncHandler)
	v1.Get("/broker/sync", brokerHandler.SyncBrokerActivityHandler)
	v1.Put("/partition/broker-activity", brokerHandler.ManagePartitionsHandler)

	// Asset Routes
	v1.Get("/assets/logos/companies/:code", assetHandler.GetCompanyLogo)
}
