package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/services"
	"github.com/KAnggara75/IDXStock/internal/utils"
	"github.com/sirupsen/logrus"
)

type HistoryUsecase interface {
	SyncStockHistory(ctx context.Context, req models.SyncHistoryRequest, source string, cookie string) error
	GetStockHistory(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error)
	ApplyStockSplit(ctx context.Context, req models.StockSplitRequest) error
}

type historyUsecase struct {
	repo             repositories.HistoryRepository
	stockRepo        repositories.StockRepository
	pasardanaService services.PasardanaService
	idxService       services.IdxService
}

func NewHistoryUsecase(
	repo repositories.HistoryRepository,
	stockRepo repositories.StockRepository,
	pasardanaService services.PasardanaService,
	idxService services.IdxService,
) HistoryUsecase {
	return &historyUsecase{
		repo:             repo,
		stockRepo:        stockRepo,
		pasardanaService: pasardanaService,
		idxService:       idxService,
	}
}

func (u *historyUsecase) SyncStockHistory(ctx context.Context, req models.SyncHistoryRequest, source string, cookie string) error {
	var records []models.StockHistory
	targetDate := time.Date(req.Year, time.Month(req.Month), req.Day, 0, 0, 0, 0, time.UTC)

	if source == "pasardana" {
		threshold := time.Date(2019, 7, 28, 0, 0, 0, 0, time.UTC)
		if targetDate.After(threshold) {
			return fmt.Errorf("sync from pasardana is restricted for dates after 2019-07-28")
		}
	}

	switch source {
	case "pasardana":
		data, err := u.pasardanaService.FetchStockHistory(req.Year, req.Month, req.Day)
		if err != nil {
			return err
		}
		for _, d := range data {
			// Parsing LastDate if provided, otherwise use targetDate
			var date time.Time
			if d.LastDate != nil && *d.LastDate != "" {
				// Format expected: "2024-12-10T00:00:00"
				t, err := time.Parse("2006-01-02T15:04:05", *d.LastDate)
				if err == nil {
					date = t
				} else {
					date = targetDate
				}
			} else {
				date = targetDate
			}

			records = append(records, models.StockHistory{
				Code:      d.Code,
				Date:      date,
				Previous:  d.PrevClosingPrice,
				OpenPrice: d.AdjustedOpenPrice,
				High:      d.AdjustedHighPrice,
				Low:       d.AdjustedLowPrice,
				Close:     d.AdjustedClosingPrice,
				Volume:    d.Volume,
				Frequency: d.Frequency,
				Value:     d.Value,
			})
		}

	case "idx":
		data, err := u.idxService.FetchStockSummary(ctx, req.Year, req.Month, req.Day, cookie)
		if err != nil {
			return err
		}
		for _, d := range data {
			var date time.Time
			if d.Date != "" {
				t, err := time.Parse("2006-01-02T15:04:05", d.Date)
				if err == nil {
					date = t
				} else {
					date = targetDate
				}
			} else {
				date = targetDate
			}

			// DelistingDate mapping
			var dd *string
			if d.DelistingDate != "" {
				dd = &d.DelistingDate
			}

			records = append(records, models.StockHistory{
				Code:                d.StockCode,
				Date:                date,
				Previous:            d.Previous,
				OpenPrice:           d.OpenPrice,
				FirstTrade:          d.FirstTrade,
				High:                d.High,
				Low:                 d.Low,
				Close:               d.Close,
				Change:              d.Change,
				Volume:              d.Volume,
				Value:               d.Value,
				Frequency:           d.Frequency,
				IndexIndividual:     d.IndexIndividual,
				Offer:               d.Offer,
				OfferVolume:         d.OfferVolume,
				Bid:                 d.Bid,
				BidVolume:           d.BidVolume,
				ListedShares:        d.ListedShares,
				TradebleShares:      d.TradebleShares,
				WeightForIndex:      d.WeightForIndex,
				ForeignSell:         d.ForeignSell,
				ForeignBuy:          d.ForeignBuy,
				DelistingDate:       dd,
				NonRegularVolume:    d.NonRegularVolume,
				NonRegularValue:     d.NonRegularValue,
				NonRegularFrequency: d.NonRegularFrequency,
			})
		}
	default:
		return fmt.Errorf("invalid source: %s", source)
	}

	if len(records) == 0 {
		logrus.Warnf("No records found to sync for date %s from %s", targetDate.Format("2006-01-02"), source)
		return nil
	}

	// Enhancement: Ensure all stock codes exist in DB before upserting history
	codeMap := make(map[string]bool)
	var codes []string
	for _, r := range records {
		if !codeMap[r.Code] {
			codeMap[r.Code] = true
			codes = append(codes, r.Code)
		}
	}

	missingCodes, err := u.stockRepo.FindMissingCodes(ctx, codes)
	if err != nil {
		logrus.Errorf("Failed to identify missing codes: %v", err)
	}

	if len(missingCodes) > 0 {
		logrus.Infof("Found %d missing stock codes, fetching from Pasardana...", len(missingCodes))
		for _, code := range missingCodes {
			detail, err := u.pasardanaService.FetchStockDetailByCode(code)
			if err != nil {
				logrus.Warnf("Failed to fetch detail for %s from Pasardana: %v", code, err)
				continue
			}

			if detail != nil {
				// Normalize dates with fallback to Epoch 0 (consistent with StockUsecase)
				epoch0 := "1970-01-01"
				if detail.ListingDate != nil && *detail.ListingDate != "" {
					parsed := utils.NormalizeDate(*detail.ListingDate)
					if parsed == "" {
						parsed = epoch0
					}
					detail.ListingDate = &parsed
				} else {
					detail.ListingDate = &epoch0
				}

				if detail.FoundingDate != nil && *detail.FoundingDate != "" {
					parsed := utils.NormalizeDate(*detail.FoundingDate)
					if parsed == "" {
						parsed = epoch0
					}
					detail.FoundingDate = &parsed
				} else {
					detail.FoundingDate = &epoch0
				}

				_, err = u.stockRepo.UpsertStocksDetail(ctx, []models.PasardanaStockDetail{*detail})
				if err != nil {
					logrus.Errorf("Failed to auto-insert missing stock %s: %v", code, err)
					continue
				}
				logrus.Infof("Auto-inserted missing stock: %s", code)
			}
		}
	}

	return u.repo.BatchUpsertStockHistory(ctx, records)
}
func (u *historyUsecase) GetStockHistory(ctx context.Context, code string, startDate, endDate *time.Time) ([]models.StockHistory, error) {
	code = strings.ToUpper(code)
	return u.repo.GetHistoryByCode(ctx, code, startDate, endDate)
}

func (u *historyUsecase) ApplyStockSplit(ctx context.Context, req models.StockSplitRequest) error {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("end_date cannot be before start_date")
	}

	multiplier, err := u.parseRatio(req.Ratio)
	if err != nil {
		return fmt.Errorf("invalid ratio: %w", err)
	}

	req.StockCode = strings.ToUpper(req.StockCode)

	return u.repo.ApplyStockSplit(ctx, req.StockCode, startDate, endDate, multiplier)
}

func (u *historyUsecase) parseRatio(ratio string) (float64, error) {
	parts := strings.Split(ratio, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("ratio must be in format A:B")
	}

	// Clean up spaces
	aStr := strings.TrimSpace(parts[0])
	bStr := strings.TrimSpace(parts[1])

	var a, b float64
	_, err := fmt.Sscanf(aStr, "%f", &a)
	if err != nil || a == 0 {
		return 0, fmt.Errorf("invalid ratio numerator")
	}

	_, err = fmt.Sscanf(bStr, "%f", &b)
	if err != nil || b == 0 {
		return 0, fmt.Errorf("invalid ratio denominator")
	}

	return b / a, nil
}
