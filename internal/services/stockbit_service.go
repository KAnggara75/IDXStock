package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/models"
)

type StockbitService interface {
	FetchDailyPrice(ctx context.Context, code string, from, to string, authorization string) ([]models.StockbitDailyPriceItem, error)
}

type stockbitService struct {
	client *http.Client
}

func NewStockbitService(client *http.Client) StockbitService {
	if client == nil {
		client = &http.Client{
			Timeout: 60 * time.Second,
		}
	}
	return &stockbitService{
		client: client,
	}
}

func (s *stockbitService) FetchDailyPrice(ctx context.Context, code string, from, to string, authorization string) ([]models.StockbitDailyPriceItem, error) {
	url := fmt.Sprintf("https://exodus.stockbit.com/chartbit/%s/price/daily?from=%s&to=%s", strings.ToUpper(code), from, to)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("origin", "https://stockbit.com")
	req.Header.Set("referer", "https://stockbit.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	logrus.Infof("Requesting Stockbit daily price for %s: %s", code, url)

	// #nosec G107
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Stockbit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stockbit API returned status: %d", resp.StatusCode)
	}

	var response models.StockbitDailyPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode Stockbit response: %w", err)
	}

	return response.Data.Chartbit, nil
}
