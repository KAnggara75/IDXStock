package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
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

// PreviewHandler godoc
// @Summary      Preview stocks from Excel file
// @Description  Parses uploaded stock Excel (.xlsx or .xls) file without persisting to database
// @Tags         stocks
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Excel file (.xlsx or .xls)"
// @Success      200   {array}   models.Stock
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /stocks/upload [post]
func (h *StockHandler) PreviewHandler(c fiber.Ctx) error {
	f, err := h.getFileAndValidate(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer f.Close()

	stocks, err := h.usecase.PreviewStocks(c.Context(), f)
	if err != nil {
		logrus.Errorf("Failed to preview stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

// UploadHandler godoc
// @Summary      Upload stocks from Excel file
// @Description  Parses and inserts uploaded stock Excel file into the database
// @Tags         stocks
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Excel file (.xlsx or .xls)"
// @Success      200   {array}   models.StockResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /stocks/upload [patch]
func (h *StockHandler) UploadHandler(c fiber.Ctx) error {
	f, err := h.getFileAndValidate(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer f.Close()

	stocks, err := h.usecase.UploadStocks(c.Context(), f)
	if err != nil {
		logrus.Errorf("Failed to upload stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
}

func (h *StockHandler) getFileAndValidate(c fiber.Ctx) (multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("Failed to get file from form: %v", err)
		return nil, fmt.Errorf("file is required")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return nil, fmt.Errorf("only .xlsx and .xls files are allowed")
	}

	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
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

// SyncDelistingStocksHandler godoc
// @Summary      Sync delisted stocks
// @Description  Fetches delisted stocks for specific year/month from IDX and updates the database
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        request  body      models.SyncDelistingRequest  true  "Delisting sync request"
// @Param        Cookie   header    string                       false "IDX Session Cookie"
// @Success      200      {array}   models.IdxDelistedStock
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /stocks/delisting/sync [put]
func (h *StockHandler) SyncDelistingStocksHandler(c fiber.Ctx) error {
	var req models.SyncDelistingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Year == 0 || req.Month == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year and month are required",
		})
	}

	stocks, err := h.usecase.SyncDelistingStocks(c.Context(), req.Year, req.Month, c.Get("Cookie"))
	if err != nil {
		logrus.Errorf("Failed to sync delisting stocks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stocks)
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
