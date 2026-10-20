package repositories

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type SectorSearchRepository interface {
	UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponseWithCode, error)
	UpsertNewSectorSub(ctx context.Context, sectorSubs []models.SectorSub) ([]models.BasicResponseWithCode, error)
}

type sectorSearchRepository struct {
	pool DBExecutor
}

func NewSectorSearchRepository(pool DBExecutor) SectorSearchRepository {
	return &sectorSearchRepository{
		pool: pool,
	}
}

func (r *sectorSearchRepository) UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponseWithCode, error) {
	if len(sectors) == 0 {
		return make([]models.BasicResponseWithCode, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO sector (id, code, name, name_en, description, last_modified)
		VALUES (?1, ?2, ?3, ?4, ?5, datetime('now'))
		ON CONFLICT (id) DO UPDATE SET
			code = excluded.code,
			name = excluded.name,
			name_en = excluded.name_en,
			description = excluded.description,
			last_modified = datetime('now')
		WHERE (sector.name IS NOT excluded.name OR
		       sector.code IS NOT excluded.code OR
		       sector.description IS NOT excluded.description OR
		       sector.name_en IS NOT excluded.name_en)
		RETURNING id, code, name
	`

	updated := make([]models.BasicResponseWithCode, 0)
	for _, s := range sectors {
		var res models.BasicResponseWithCode
		err := tx.QueryRow(ctx, query, s.Id, s.Code, s.Name, s.NameEn, s.Description).Scan(&res.Id, &res.Code, &res.Name)
		if err != nil {
			continue
		}
		updated = append(updated, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d new sectors", len(updated))
	return updated, nil
}

func (r *sectorSearchRepository) UpsertNewSectorSub(ctx context.Context, sectorSubs []models.SectorSub) ([]models.BasicResponseWithCode, error) {
	if len(sectorSubs) == 0 {
		return make([]models.BasicResponseWithCode, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO sector_sub (id, code, name, name_en, sector_id, last_modified)
		VALUES (?1, ?2, ?3, ?4, ?5, datetime('now'))
		ON CONFLICT (id) DO UPDATE SET
			code = excluded.code,
			name = excluded.name,
			name_en = excluded.name_en,
			sector_id = excluded.sector_id,
			last_modified = datetime('now')
		WHERE (sector_sub.name IS NOT excluded.name OR
		       sector_sub.code IS NOT excluded.code OR
		       sector_sub.name_en IS NOT excluded.name_en OR
		       sector_sub.sector_id IS NOT excluded.sector_id)
		RETURNING id, code, name
	`

	updated := make([]models.BasicResponseWithCode, 0)
	for _, sub := range sectorSubs {
		var res models.BasicResponseWithCode
		err := tx.QueryRow(ctx, query, sub.Id, sub.Code, sub.Name, sub.NameEn, sub.SectorId).Scan(&res.Id, &res.Code, &res.Name)
		if err != nil {
			continue
		}
		updated = append(updated, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d new sub-sectors", len(updated))
	return updated, nil
}
