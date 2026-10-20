package repositories

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type IndustryRepository interface {
	UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error)
	UpsertIndustrySub(ctx context.Context, industrySubs []models.IndustrySub) ([]models.BasicResponse, error)
}

type industryRepository struct {
	pool DBExecutor
}

func NewIndustryRepository(pool DBExecutor) IndustryRepository {
	return &industryRepository{
		pool: pool,
	}
}

func (r *industryRepository) UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error) {
	if len(industries) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO industry (id, name, last_modified)
		VALUES (?1, ?2, datetime('now'))
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			last_modified = datetime('now')
		WHERE (industry.name IS NOT excluded.name)
		RETURNING id, name
	`

	updated := make([]models.BasicResponse, 0)
	for _, ind := range industries {
		var res models.BasicResponse
		err := tx.QueryRow(ctx, query, ind.Id, ind.Name).Scan(&res.Id, &res.Name)
		if err != nil {
			continue
		}
		updated = append(updated, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d industries", len(updated))
	return updated, nil
}

func (r *industryRepository) UpsertIndustrySub(ctx context.Context, industrySubs []models.IndustrySub) ([]models.BasicResponse, error) {
	if len(industrySubs) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO industry_sub (id, name, industry_id, last_modified)
		VALUES (?1, ?2, ?3, datetime('now'))
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			industry_id = excluded.industry_id,
			last_modified = datetime('now')
		WHERE (industry_sub.name IS NOT excluded.name OR
		       industry_sub.industry_id IS NOT excluded.industry_id)
		RETURNING id, name
	`

	updated := make([]models.BasicResponse, 0)
	for _, sub := range industrySubs {
		var res models.BasicResponse
		err := tx.QueryRow(ctx, query, sub.Id, sub.Name, sub.IndustryId).Scan(&res.Id, &res.Name)
		if err != nil {
			continue
		}
		updated = append(updated, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d sub-industries", len(updated))
	return updated, nil
}
