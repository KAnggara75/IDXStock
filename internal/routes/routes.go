package routes

import (
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "github.com/KAnggara75/IDXStock/docs"
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
	historyRepo := repositories.NewHistoryRepository(pool)

	pasardanaService := services.NewPasardanaService(nil)
	idxService := services.NewIdxService(nil)

	stockUsecase := usecases.NewStockUsecase(stockRepo, pasardanaService)
	historyUsecase := usecases.NewHistoryUsecase(historyRepo, stockRepo, pasardanaService, idxService)

	stockHandler := handlers.NewStockHandler(stockUsecase)
	historyHandler := handlers.NewHistoryHandler(historyUsecase)
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

	// Swagger OpenAPI documentation
	app.Get("/swagger/*", swaggo.HandlerDefault)

	// API V1 Routes
	v1 := app.Group("/api/v1")

	v1.Get("/stocks", stockHandler.GetAllStocksHandler)
	v1.Get("/stocks/boards", stockHandler.GetBoardsHandler)
	v1.Get("/stocks/listing-boards", stockHandler.GetListingBoardsHandler)
	v1.Get("/stocks/boards/:board", stockHandler.GetStocksByBoardHandler)
	v1.Get("/stocks/:code", stockHandler.GetStockByCodeHandler)
	v1.Put("/stocks/sync", stockHandler.SyncStockDetailHandler)
	v1.Put("/stocks/sectors-industries/sync", stockHandler.SyncStockSectorAndIndustryHandler)
	v1.Put("/stocks/history/sync", historyHandler.SyncStockHistoryHandler)
	v1.Get("/stocks/:code/history", historyHandler.GetStockHistoryHandler)
	v1.Put("/stocks/split", historyHandler.ApplyStockSplitHandler)

	// Asset Routes
	v1.Get("/assets/logos/companies/:code", assetHandler.GetCompanyLogo)
}
