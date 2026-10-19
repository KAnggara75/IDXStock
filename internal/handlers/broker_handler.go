package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/usecases"
)

type BrokerHandler struct {
	usecase usecases.BrokerUsecase
}

func NewBrokerHandler(usecase usecases.BrokerUsecase) *BrokerHandler {
	return &BrokerHandler{
		usecase: usecase,
	}
}

// SyncBrokerActivityHandler godoc
// @Summary      Sync broker activity
// @Description  Fetches broker transaction activities from Exodus API using access token
// @Tags         brokers
// @Accept       json
// @Produce      json
// @Param        Authorization    header    string  true  "Bearer token"
// @Param        broker_code      query     string  false "Broker code (e.g. CC)"
// @Param        from             query     string  false "Start date (YYYY-MM-DD)"
// @Param        to               query     string  false "End date (YYYY-MM-DD)"
// @Param        transaction_type query     string  false "Transaction type (e.g. NET)"
// @Param        market_board     query     string  false "Market board (e.g. RG)"
// @Param        investor_type    query     string  false "Investor type (e.g. ALL, D, F)"
// @Success      200              {array}   models.BrokerActivity
// @Failure      400              {object}  models.ErrorResponse
// @Failure      401              {object}  models.ErrorResponse
// @Failure      500              {object}  models.ErrorResponse
// @Router       /broker/sync [get]
func (h *BrokerHandler) SyncBrokerActivityHandler(c fiber.Ctx) error {
	var params models.SyncBrokerActivityParams
	if err := c.Bind().Query(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse query parameters",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	activities, err := h.usecase.SyncBrokerActivity(c.Context(), token, params)
	if err != nil {
		logrus.Errorf("Broker activity sync error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(activities)
}

// ManagePartitionsHandler godoc
// @Summary      Manage broker activity table partitions
// @Description  Automatically detects and provisions monthly table partitions for broker activity
// @Tags         partitions
// @Produce      json
// @Success      200  {object}  map[string]any
// @Success      201  {object}  map[string]any
// @Failure      500  {object}  models.ErrorResponse
// @Router       /partition/broker-activity [put]
func (h *BrokerHandler) ManagePartitionsHandler(c fiber.Ctx) error {

	resp, err := h.usecase.ManagePartitions(c.Context())
	if err != nil {
		logrus.Errorf("Partition management error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status := fiber.StatusOK
	if resp.PartitionsCreated > 0 {
		status = fiber.StatusCreated
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}
