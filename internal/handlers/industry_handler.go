package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/usecases"
)

type IndustryHandler struct {
	usecase usecases.IndustryUsecase
}

func NewIndustryHandler(usecase usecases.IndustryUsecase) *IndustryHandler {
	return &IndustryHandler{
		usecase: usecase,
	}
}

// IndustrySyncHandler godoc
// @Summary      Sync industries and sub-industries
// @Description  Fetches industries taxonomy from Pasardana and updates the database
// @Tags         industries
// @Produce      json
// @Success      200  {object}  models.IndustrySyncResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /industries/sync [put]
func (h *IndustryHandler) IndustrySyncHandler(c fiber.Ctx) error {

	results, err := h.usecase.SyncIndustry(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync industries: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
