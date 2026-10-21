package handlers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/usecases"
)

type StockHandler struct {
	usecase usecases.StockUsecase
}

func NewStockHandler(usecase usecases.StockUsecase) *StockHandler {
	return &StockHandler{
		usecase: usecase,
	}
}

// SyncStockDetailHandler godoc
// @Summary      Sync stock details asynchronously
// @Description  Triggers background sync for stock details from Pasardana
// @Tags         stocks
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Router       /stocks/sync [put]
func (h *StockHandler) SyncStockDetailHandler(c fiber.Ctx) error {
	// Run synchronization in background
	go func(ctx context.Context) {
		_, err := h.usecase.SyncStockDetail(ctx)
		if err != nil {
			logrus.Errorf("Background stock sync failed: %v", err)
		}
	}(context.WithoutCancel(c.Context()))

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stock synchronization process started in the background",
	})
}

// SyncStockSectorAndIndustryHandler godoc
// @Summary      Sync stock sector and industry asynchronously
// @Description  Fetches fkNewSectorId, fkNewSubSectorId, fkNewIndustryId, and fkNewSubIndustryId from Pasardana for stocks that currently have null sector or industry values
// @Tags         stocks
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Router       /stocks/sectors-industries/sync [put]
func (h *StockHandler) SyncStockSectorAndIndustryHandler(c fiber.Ctx) error {
	go func(ctx context.Context) {
		res, err := h.usecase.SyncStockSectorAndIndustry(ctx)
		if err != nil {
			logrus.Errorf("Background sector and industry sync failed: %v", err)
			return
		}
		if res != nil {
			logrus.Infof("Background sector and industry sync finished: TotalChecked=%d, TotalUpdated=%d", res.TotalChecked, res.TotalUpdated)
		} else {
			logrus.Info("Background sector and industry sync finished")
		}
	}(context.WithoutCancel(c.Context()))

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stock sector and industry synchronization process started in the background",
	})
}

// GetAllStocksHandler godoc
// @Summary      Get all stocks / emiten
// @Description  Retrieves list of all stocks/emiten with optional filter by board and search query
// @Tags         stocks
// @Produce      json
// @Param        board   query     string  false  "Board filter (e.g. UTAMA, PENGEMBANGAN, AKSELERASI)"
// @Param        search  query     string  false  "Search by stock code or company name"
// @Success      200     {array}   models.StockDetailResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /stocks [get]
func (h *StockHandler) GetAllStocksHandler(c fiber.Ctx) error {
	var filter models.GetAllStocksFilter
	_ = c.Bind().Query(&filter)
	if filter.Board == "" {
		filter.Board = c.Query("board")
	}
	if filter.Search == "" {
		filter.Search = c.Query("search")
	}

	stocks, err := h.usecase.GetAllStocks(c.Context(), filter)
	if err != nil {
		logrus.Errorf("Failed to retrieve stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

// GetStockByCodeHandler godoc
// @Summary      Get stock / emiten by code
// @Description  Retrieves detailed information for a specific stock/emiten by ticker code
// @Tags         stocks
// @Produce      json
// @Param        code   path      string  true  "Stock Ticker Code (e.g. BBCA)"
// @Success      200    {object}  models.StockDetailResponse
// @Failure      400    {object}  models.ErrorResponse
// @Failure      404    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /stocks/{code} [get]
func (h *StockHandler) GetStockByCodeHandler(c fiber.Ctx) error {
	code := strings.TrimSpace(c.Params("code"))
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Stock code is required",
		})
	}

	stock, err := h.usecase.GetStockByCode(c.Context(), code)
	if err != nil {
		logrus.Errorf("Failed to retrieve stock %s: %v", code, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if stock == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Stock not found",
		})
	}

	return c.JSON(stock)
}

// GetBoardsHandler godoc
// @Summary      Get distinct trading boards
// @Description  Retrieves list of all distinct trading boards available in IDX (e.g. Main, Development, Acceleration, Watchlist, Ekonomi Baru)
// @Tags         stocks
// @Produce      json
// @Success      200     {array}   string
// @Failure      500     {object}  models.ErrorResponse
// @Router       /stocks/boards [get]
func (h *StockHandler) GetBoardsHandler(c fiber.Ctx) error {
	boards, err := h.usecase.GetDistinctBoards(c.Context())
	if err != nil {
		logrus.Errorf("Failed to retrieve boards: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(boards)
}

// GetListingBoardsHandler godoc
// @Summary      Get master listing boards
// @Description  Retrieves list of master listing boards (code, name, name_id, description)
// @Tags         stocks
// @Produce      json
// @Success      200     {array}   models.ListingBoard
// @Failure      500     {object}  models.ErrorResponse
// @Router       /stocks/listing-boards [get]
func (h *StockHandler) GetListingBoardsHandler(c fiber.Ctx) error {
	boards, err := h.usecase.GetListingBoards(c.Context())
	if err != nil {
		logrus.Errorf("Failed to retrieve listing boards: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(boards)
}

// GetStocksByBoardHandler godoc
// @Summary      Get stocks / emiten by board
// @Description  Retrieves list of stocks/emiten filtered by specific trading board
// @Tags         stocks
// @Produce      json
// @Param        board   path      string  true  "Trading Board (e.g. Main, Development, Acceleration, Watchlist, Ekonomi Baru, or Indonesian name UTAMA, PENGEMBANGAN, etc.)"
// @Success      200     {array}   models.StockDetailResponse
// @Failure      400     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /stocks/boards/{board} [get]
func (h *StockHandler) GetStocksByBoardHandler(c fiber.Ctx) error {
	board := strings.TrimSpace(c.Params("board"))
	if board == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Board is required",
		})
	}

	stocks, err := h.usecase.GetStocksByBoard(c.Context(), board)
	if err != nil {
		logrus.Errorf("Failed to retrieve stocks for board %s: %v", board, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}
