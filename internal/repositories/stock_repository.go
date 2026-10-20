package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type StockRepository interface {
	BatchInsertStocks(ctx context.Context, stocks []models.Stock) error
	UpsertStocksDetail(ctx context.Context, data []models.PasardanaStockDetail) ([]models.StockResponse, error)
	UpdateDelistingDate(ctx context.Context, code, delistingDate string) (*models.StockResponse, error)
	FindMissingCodes(ctx context.Context, codes []string) ([]string, error)
	GetAllStocks(ctx context.Context, filter models.GetAllStocksFilter) ([]models.StockDetailResponse, error)
	GetStockByCode(ctx context.Context, code string) (*models.StockDetailResponse, error)
	GetDistinctBoards(ctx context.Context) ([]string, error)
}

type stockRepository struct {
	pool DBExecutor
}

func NewStockRepository(pool DBExecutor) StockRepository {
	return &stockRepository{
		pool: pool,
	}
}

func (r *stockRepository) BatchInsertStocks(ctx context.Context, stocks []models.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO stocks (code, name, listing_date, delisting_date, shares, board)
		VALUES (?1, ?2, ?3, NULLIF(?4, ''), ?5, ?6)
		ON CONFLICT (code) DO UPDATE SET
			name = excluded.name,
			listing_date = excluded.listing_date,
			delisting_date = excluded.delisting_date,
			shares = excluded.shares,
			board = excluded.board,
			last_modified = datetime('now')
		WHERE
			stocks.name IS NOT excluded.name OR
			stocks.listing_date IS NOT excluded.listing_date OR
			stocks.delisting_date IS NOT excluded.delisting_date OR
			stocks.shares IS NOT excluded.shares OR
			stocks.board IS NOT excluded.board
	`

	var changedCount int64
	for i, s := range stocks {
		var listDate, delistDate string
		if s.ListingDate != nil {
			listDate = *s.ListingDate
		}
		if s.DelistingDate != nil {
			delistDate = *s.DelistingDate
		}

		res, err := tx.Exec(ctx, query, s.Code, s.CompanyName, listDate, delistDate, s.Shares, s.ListingBoard)
		if err != nil {
			return fmt.Errorf("failed to execute batch statement %d: %w", i, err)
		}

		if res.RowsAffected() > 0 {
			changedCount++
			logrus.Debugf("Data changed for stock: %s (%s)", stocks[i].Code, stocks[i].CompanyName)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully processed %d stocks, %d records updated/inserted", len(stocks), changedCount)
	return nil
}

func (r *stockRepository) UpdateDelistingDate(ctx context.Context, code, delistingDate string) (*models.StockResponse, error) {
	query := `
		UPDATE stocks
		SET delisting_date = ?1, last_modified = datetime('now')
		WHERE code = ?2
		RETURNING id, code, name
	`

	var sr models.StockResponse
	var id sql.NullInt64
	err := r.pool.QueryRow(ctx, query, delistingDate, code).Scan(&id, &sr.Code, &sr.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No stock found with this code
		}
		return nil, fmt.Errorf("failed to update delisting date for %s: %w", code, err)
	}
	if id.Valid {
		sr.Id = int(id.Int64)
	}

	return &sr, nil
}

func (r *stockRepository) UpsertStocksDetail(ctx context.Context, data []models.PasardanaStockDetail) ([]models.StockResponse, error) {
	if len(data) == 0 {
		return make([]models.StockResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO stocks (
			id, code, name, listing_date, total_employees, annual_dividend,
			general_information, founding_date, sector_id, sector_sub_id,
			industry_id, industry_sub_id, last_modified
		)
		VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, datetime('now'))
		ON CONFLICT (code) DO UPDATE SET
			id = excluded.id,
			name = excluded.name,
			listing_date = excluded.listing_date,
			total_employees = excluded.total_employees,
			annual_dividend = excluded.annual_dividend,
			general_information = excluded.general_information,
			founding_date = excluded.founding_date,
			sector_id = excluded.sector_id,
			sector_sub_id = excluded.sector_sub_id,
			industry_id = excluded.industry_id,
			industry_sub_id = excluded.industry_sub_id,
			last_modified = datetime('now')
		WHERE
			stocks.id IS NOT excluded.id OR
			stocks.name IS NOT excluded.name OR
			stocks.listing_date IS NOT excluded.listing_date OR
			stocks.total_employees IS NOT excluded.total_employees OR
			stocks.annual_dividend IS NOT excluded.annual_dividend OR
			stocks.general_information IS NOT excluded.general_information OR
			stocks.founding_date IS NOT excluded.founding_date OR
			stocks.sector_id IS NOT excluded.sector_id OR
			stocks.sector_sub_id IS NOT excluded.sector_sub_id OR
			stocks.industry_id IS NOT excluded.industry_id OR
			stocks.industry_sub_id IS NOT excluded.industry_sub_id
		RETURNING id, code, name
	`

	updatedStocks := make([]models.StockResponse, 0)
	for i, s := range data {
		var listDate string
		if s.ListingDate != nil {
			listDate = *s.ListingDate
		}

		var sr models.StockResponse
		var id sql.NullInt64
		err := tx.QueryRow(ctx, query,
			s.Id, s.Code, s.Name, listDate, s.TotalEmployees, s.AnnualDividend,
			s.GeneralInformation, s.FoundingDate, s.FkNewSectorId, s.FkNewSubSectorId,
			s.FkNewIndustryId, s.FkNewSubIndustryId,
		).Scan(&id, &sr.Code, &sr.Name)

		if err != nil {
			logrus.Debugf("[%d/%d] [%s]: No Change Detected or skipped", i+1, len(data), s.Code)
			continue
		}

		if id.Valid {
			sr.Id = int(id.Int64)
		}
		updatedStocks = append(updatedStocks, sr)
		logrus.Infof("[%d/%d] [%s]: Data Change Detected - Synced", i+1, len(data), s.Code)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedStocks, nil
}

func (r *stockRepository) FindMissingCodes(ctx context.Context, codes []string) ([]string, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(codes))
	args := make([]any, len(codes))
	for i, c := range codes {
		placeholders[i] = fmt.Sprintf("(?%d)", i+1)
		args[i] = c
	}

	query := fmt.Sprintf(`
		WITH input(code) AS (VALUES %s)
		SELECT input.code FROM input
		LEFT JOIN stocks s ON input.code = s.code
		WHERE s.code IS NULL
	`, strings.Join(placeholders, ", "))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find missing codes: %w", err)
	}
	defer rows.Close()

	var missingCodes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		missingCodes = append(missingCodes, code)
	}

	return missingCodes, nil
}

func (r *stockRepository) GetAllStocks(ctx context.Context, filter models.GetAllStocksFilter) ([]models.StockDetailResponse, error) {
	query := `
		SELECT
			code, name, listing_date, delisting_date, shares, board,
			total_employees, annual_dividend, general_information, founding_date,
			sector_id, sector_sub_id, industry_id, industry_sub_id, last_modified
		FROM stocks
		WHERE (?1 = '' OR UPPER(board) = UPPER(?1))
		  AND (?2 = '' OR UPPER(code) LIKE '%' || UPPER(?2) || '%' OR UPPER(name) LIKE '%' || UPPER(?2) || '%')
		ORDER BY code ASC
	`

	rows, err := r.pool.Query(ctx, query, strings.TrimSpace(filter.Board), strings.TrimSpace(filter.Search))
	if err != nil {
		return nil, fmt.Errorf("failed to get all stocks: %w", err)
	}
	defer rows.Close()

	stocks := make([]models.StockDetailResponse, 0)
	for rows.Next() {
		var s models.StockDetailResponse
		var delistingDate, totalEmployees, generalInfo, foundingDate sql.NullString
		var annualDividend sql.NullFloat64
		var sectorId, sectorSubId, industryId, industrySubId sql.NullInt64

		err := rows.Scan(
			&s.Code, &s.Name, &s.ListingDate, &delistingDate, &s.Shares, &s.Board,
			&totalEmployees, &annualDividend, &generalInfo, &foundingDate,
			&sectorId, &sectorSubId, &industryId, &industrySubId, &s.LastModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock row: %w", err)
		}

		if delistingDate.Valid {
			s.DelistingDate = &delistingDate.String
		}
		if totalEmployees.Valid {
			s.TotalEmployees = &totalEmployees.String
		}
		if annualDividend.Valid {
			s.AnnualDividend = &annualDividend.Float64
		}
		if generalInfo.Valid {
			s.GeneralInformation = &generalInfo.String
		}
		if foundingDate.Valid {
			s.FoundingDate = &foundingDate.String
		}
		if sectorId.Valid {
			v := int(sectorId.Int64)
			s.SectorId = &v
		}
		if sectorSubId.Valid {
			v := int(sectorSubId.Int64)
			s.SectorSubId = &v
		}
		if industryId.Valid {
			v := int(industryId.Int64)
			s.IndustryId = &v
		}
		if industrySubId.Valid {
			v := int(industrySubId.Int64)
			s.IndustrySubId = &v
		}

		stocks = append(stocks, s)
	}

	return stocks, nil
}

func (r *stockRepository) GetStockByCode(ctx context.Context, code string) (*models.StockDetailResponse, error) {
	query := `
		SELECT
			code, name, listing_date, delisting_date, shares, board,
			total_employees, annual_dividend, general_information, founding_date,
			sector_id, sector_sub_id, industry_id, industry_sub_id, last_modified
		FROM stocks
		WHERE UPPER(code) = UPPER(?1)
		LIMIT 1
	`

	var s models.StockDetailResponse
	var delistingDate, totalEmployees, generalInfo, foundingDate sql.NullString
	var annualDividend sql.NullFloat64
	var sectorId, sectorSubId, industryId, industrySubId sql.NullInt64

	err := r.pool.QueryRow(ctx, query, strings.TrimSpace(code)).Scan(
		&s.Code, &s.Name, &s.ListingDate, &delistingDate, &s.Shares, &s.Board,
		&totalEmployees, &annualDividend, &generalInfo, &foundingDate,
		&sectorId, &sectorSubId, &industryId, &industrySubId, &s.LastModified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock by code: %w", err)
	}

	if delistingDate.Valid {
		s.DelistingDate = &delistingDate.String
	}
	if totalEmployees.Valid {
		s.TotalEmployees = &totalEmployees.String
	}
	if annualDividend.Valid {
		s.AnnualDividend = &annualDividend.Float64
	}
	if generalInfo.Valid {
		s.GeneralInformation = &generalInfo.String
	}
	if foundingDate.Valid {
		s.FoundingDate = &foundingDate.String
	}
	if sectorId.Valid {
		v := int(sectorId.Int64)
		s.SectorId = &v
	}
	if sectorSubId.Valid {
		v := int(sectorSubId.Int64)
		s.SectorSubId = &v
	}
	if industryId.Valid {
		v := int(industryId.Int64)
		s.IndustryId = &v
	}
	if industrySubId.Valid {
		v := int(industrySubId.Int64)
		s.IndustrySubId = &v
	}

	return &s, nil
}

func (r *stockRepository) GetDistinctBoards(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT board
		FROM stocks
		WHERE board IS NOT NULL AND TRIM(board) != ''
		ORDER BY board ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct boards: %w", err)
	}
	defer rows.Close()

	boards := make([]string, 0)
	for rows.Next() {
		var board string
		if err := rows.Scan(&board); err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}
		boards = append(boards, board)
	}

	return boards, nil
}
