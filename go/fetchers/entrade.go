package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lotusmarket/lotusmarket-go/types"
)

const entradeBaseURL = "https://services.entrade.com.vn/chart-api/v2/ohlcs/stock"

type entradeResponse struct {
	T []int64   `json:"t"`
	O []float64 `json:"o"`
	H []float64 `json:"h"`
	L []float64 `json:"l"`
	C []float64 `json:"c"`
	V []float64 `json:"v"`
}

func EntradeHistory(ctx context.Context, ticker string, days int) ([]types.StockData, error) {
	to := time.Now().Unix()
	from := time.Now().AddDate(0, 0, -days).Unix()
	url := fmt.Sprintf("%s?from=%d&to=%d&symbol=%s&resolution=1D", entradeBaseURL, from, to, ticker)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("entrade request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("entrade read: %w", err)
	}
	return parseEntradeResponse(body, ticker)
}

func parseEntradeResponse(data []byte, ticker string) ([]types.StockData, error) {
	var resp entradeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	n := len(resp.T)
	prices := make([]types.StockData, 0, n)
	for i := range n {
		t := time.Unix(resp.T[i], 0)
		sd := types.StockData{
			Ticker: ticker,
			Date:   t.Format("2006-01-02"),
			Open:   resp.O[i],
			High:   resp.H[i],
			Low:    resp.L[i],
			Close:  resp.C[i],
			Volume: int64(resp.V[i]),
		}
		sd.ToVND()
		prices = append(prices, sd)
	}
	return prices, nil
}

func EntradeLatest(ctx context.Context, ticker string) (*types.StockData, error) {
	prices, err := EntradeHistory(ctx, ticker, 5)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no entrade data for %s", ticker)
	}
	latest := prices[len(prices)-1]
	return &latest, nil
}
