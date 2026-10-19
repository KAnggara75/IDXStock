package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/usecases"
)

type SectorHandler struct {
	usecase usecases.SectorUsecase
}

func NewSectorHandler(usecase usecases.SectorUsecase) *SectorHandler {
	return &SectorHandler{
		usecase: usecase,
	}
}

// SyncNewSectorsHandler godoc
// @Summary      Sync new sectors and sub-sectors
// @Description  Fetches sectors and subsectors from Pasardana and updates the database
// @Tags         sectors
// @Produce      json
// @Success      200  {object}  models.SectorSyncNewResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /sectors/sync [put]
func (h *SectorHandler) SyncNewSectorsHandler(c fiber.Ctx) error {

	results, err := h.usecase.SyncNewSectors(c.Context())
	if err != nil {
		logrus.Errorf("Failed to sync new sectors: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
