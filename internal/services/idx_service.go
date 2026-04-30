package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/sirupsen/logrus"
)

type IdxService interface {
	FetchDelistedStocks(year, month int) ([]models.IdxDelistedStock, error)
	FetchStockSummary(year, month, day int) ([]models.IdxSummaryData, error)
	ParseIdxDate(dateStr string) (string, error)
}

type idxService struct {
	client *http.Client
}

func NewIdxService(client *http.Client) IdxService {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &idxService{
		client: client,
	}
}

func (s *idxService) FetchDelistedStocks(year, month int) ([]models.IdxDelistedStock, error) {
	url := fmt.Sprintf("https://idx.co.id/primary/DigitalStatistic/GetApiDataPaginated?urlName=LINK_DELISTING&periodYear=%d&periodMonth=%d", year, month)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	// Log request URL only
	logrus.Infof("Requesting IDX: %s", url)

	// #nosec G107
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from IDX: %w", err)
	}
	defer resp.Body.Close()

	logrus.Infof("IDX Response Status: %d", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	logrus.Debugf("IDX Response received, length: %d bytes", len(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDX API returned status: %d", resp.StatusCode)
	}

	var idxResp models.IdxDelistingResponse
	if err := json.Unmarshal(bodyBytes, &idxResp); err != nil {
		logrus.Errorf("Failed to decode IDX response: %v", err)
		return nil, fmt.Errorf("failed to decode IDX response: %w", err)
	}

	return idxResp.Data, nil
}

func (s *idxService) FetchStockSummary(year, month, day int) ([]models.IdxSummaryData, error) {
	url := fmt.Sprintf("https://idx.co.id/primary/TradingSummary/GetStockSummary?date=%04d%02d%02d", year, month, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	logrus.Infof("Requesting IDX Summary: %s", url)

	// #nosec G107
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from IDX: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDX API returned status: %d", resp.StatusCode)
	}

	var idxResp models.IdxSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&idxResp); err != nil {
		return nil, fmt.Errorf("failed to decode IDX summary response: %w", err)
	}

	return idxResp.Data, nil
}

func (s *idxService) ParseIdxDate(dateStr string) (string, error) {
	// Format example: "18 July 2025"
	// time.Parse layout for this: "02 January 2006"
	t, err := time.Parse("02 January 2006", strings.TrimSpace(dateStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}
	return t.Format("2006-01-02"), nil
}

func (s *idxService) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("Origin", "https://idx.co.id")
	req.Header.Set("Referer", "https://idx.co.id/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)

	// Use cookie from environment variable if available, otherwise use default
	cookie := os.Getenv("IDX_COOKIE")
	if cookie == "" {
		cookie = "EGRUM_BTM=756eff19-92e1-4e3e-9701-4f7ddc93fe4e#~#1.B.11||0; __cf_bm=1vj162Z37dpo3zNyPnpz8CrDvs4BcdN7_WlqjpYvWQ8-1777555066.174328-1.0.1.1-dVrYohKm.aRGa.HmUqjoxYlQ4A.YGp9LGfZ72YZTd4YxNm5D39cgPysiOdveXNbN7DQcFgMBbRKQPsnXc97frTbYoUHoOFUN7YDosFxjTLlq5Y8IUuX2Pz14jGEvtFqz; _cfuvid=BN3x.yGYjO.L0DeEfqo68wDJ09E0gSS0.paXQ2oQCL8-1777555066.174328-1.0.1.1-4XH2ZsWQxqvI9APck5pm5saXAXL1RYJq2XZcJePq89k"
	}
	req.Header.Set("Cookie", cookie)
}
