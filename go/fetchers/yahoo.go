// Yahoo Finance fetcher — public chart API v8. No API key required.
// Supports:
//   - VN tickers: use .VN suffix (HOSE) or .HN suffix (HNX), e.g. "ACB.VN", "SHB.HN"
//   - Global indices: ^GSPC, ^DJI, ^HSI, ^N225, ^VIX, ^TNX, etc.
//   - Commodities: GC=F (gold), CL=F (oil), DX-Y.NYB (USD index)
//   - Futures: ES=F, NQ=F, YM=F
//
// All prices returned by Yahoo are split- and dividend-adjusted.
package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ducnhd/lotusmarket/go/types"
)

// YahooQuote is a snapshot from the Yahoo chart endpoint.
type YahooQuote struct {
	Symbol     string
	Price      float64
	PrevClose  float64
	Change     float64
	ChangePct  float64
	DayHigh    float64
	DayLow     float64
	Volume     int64
	Currency   string
	MarketTime time.Time
}

// YahooBar is a single OHLCV bar from historical fetch.
type YahooBar struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type yahooChartResp struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				ChartPreviousClose float64 `json:"chartPreviousClose"`
				RegularMarketTime  int64   `json:"regularMarketTime"`
				Currency           string  `json:"currency"`
				Symbol             string  `json:"symbol"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
				AdjClose []struct {
					AdjClose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

// YahooLatest fetches the latest snapshot quote for a symbol.
// Uses range=5d so we can derive previous-session close from the actual bar,
// not the chart-meta ChartPreviousClose (which is the close before the range starts).
func YahooLatest(ctx context.Context, symbol string) (*YahooQuote, error) {
	resp, err := yahooChart(ctx, symbol, "1d", "5d")
	if err != nil {
		return nil, err
	}
	if len(resp.Chart.Result) == 0 {
		return nil, fmt.Errorf("yahoo %s: empty result", symbol)
	}
	r := resp.Chart.Result[0]
	price := r.Meta.RegularMarketPrice

	prevClose := r.Meta.ChartPreviousClose
	if len(r.Indicators.Quote) > 0 {
		q := r.Indicators.Quote[0]
		n := len(q.Close)
		if n >= 2 && q.Close[n-2] > 0 {
			prevClose = q.Close[n-2]
		}
	}

	change := price - prevClose
	changePct := 0.0
	if prevClose != 0 {
		changePct = (change / prevClose) * 100
	}

	out := &YahooQuote{
		Symbol:     symbol,
		Price:      price,
		PrevClose:  prevClose,
		Change:     change,
		ChangePct:  changePct,
		Currency:   r.Meta.Currency,
		MarketTime: time.Unix(r.Meta.RegularMarketTime, 0),
	}
	if len(r.Indicators.Quote) > 0 {
		q := r.Indicators.Quote[0]
		n := len(q.Close)
		if n > 0 {
			if len(q.High) >= n {
				out.DayHigh = q.High[n-1]
			}
			if len(q.Low) >= n {
				out.DayLow = q.Low[n-1]
			}
			if len(q.Volume) >= n {
				out.Volume = q.Volume[n-1]
			}
		}
	}
	return out, nil
}

// YahooHistory fetches daily OHLCV bars for a symbol over the given range.
// `rangeStr` is Yahoo notation: 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, max.
// Returns bars oldest-first.
func YahooHistory(ctx context.Context, symbol, rangeStr string) ([]YahooBar, error) {
	resp, err := yahooChart(ctx, symbol, "1d", rangeStr)
	if err != nil {
		return nil, err
	}
	if len(resp.Chart.Result) == 0 {
		return nil, fmt.Errorf("yahoo %s: empty result", symbol)
	}
	r := resp.Chart.Result[0]
	if len(r.Indicators.Quote) == 0 {
		return nil, fmt.Errorf("yahoo %s: no quote bars", symbol)
	}
	q := r.Indicators.Quote[0]
	n := len(r.Timestamp)
	bars := make([]YahooBar, 0, n)
	for i := 0; i < n; i++ {
		// Skip rows with missing close
		if i >= len(q.Close) || q.Close[i] == 0 {
			continue
		}
		bar := YahooBar{
			Time:  time.Unix(r.Timestamp[i], 0),
			Close: q.Close[i],
		}
		if i < len(q.Open) {
			bar.Open = q.Open[i]
		}
		if i < len(q.High) {
			bar.High = q.High[i]
		}
		if i < len(q.Low) {
			bar.Low = q.Low[i]
		}
		if i < len(q.Volume) {
			bar.Volume = q.Volume[i]
		}
		bars = append(bars, bar)
	}
	return bars, nil
}

// YahooMultiple fetches latest quotes for many symbols concurrently.
// Failures on individual symbols are silently skipped — caller gets partial results.
func YahooMultiple(ctx context.Context, symbols []string) []YahooQuote {
	type result struct {
		idx int
		q   *YahooQuote
	}
	ch := make(chan result, len(symbols))
	for i, sym := range symbols {
		go func(i int, sym string) {
			q, err := YahooLatest(ctx, sym)
			if err != nil {
				ch <- result{i, nil}
				return
			}
			ch <- result{i, q}
		}(i, sym)
	}
	out := make([]YahooQuote, 0, len(symbols))
	collected := make([]*YahooQuote, len(symbols))
	for range symbols {
		r := <-ch
		collected[r.idx] = r.q
	}
	for _, q := range collected {
		if q != nil {
			out = append(out, *q)
		}
	}
	return out
}

// YahooToStockData converts a YahooQuote to a lotusmarket types.StockData for
// callers that want a unified interface across fetchers.
func YahooToStockData(q YahooQuote) types.StockData {
	return types.StockData{
		Ticker:        q.Symbol,
		Close:         q.Price,
		High:          q.DayHigh,
		Low:           q.DayLow,
		Volume:        q.Volume,
		ChangePercent: q.ChangePct,
		ChangeValue:   q.Change,
		RefPrice:      q.PrevClose,
	}
}

func yahooChart(ctx context.Context, symbol, interval, rangeStr string) (*yahooChartResp, error) {
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=%s&range=%s",
		symbol, interval, rangeStr,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("yahoo %s build request: %w", symbol, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 lotusmarket")
	client := &http.Client{Timeout: 15 * time.Second}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yahoo %s fetch: %w", symbol, err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo %s status %d", symbol, httpResp.StatusCode)
	}
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("yahoo %s read: %w", symbol, err)
	}
	var out yahooChartResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("yahoo %s decode: %w", symbol, err)
	}
	return &out, nil
}

// GlobalIndexRegistry lists default international indices commonly tracked
// alongside VN markets. Use with YahooMultiple for one-shot fetch.
var GlobalIndexRegistry = []struct {
	Symbol string
	Name   string
	NameVi string
	Region string
}{
	{"^DJI", "Dow Jones", "Dow Jones", "us"},
	{"^GSPC", "S&P 500", "S&P 500", "us"},
	{"^IXIC", "NASDAQ", "NASDAQ", "us"},
	{"^HSI", "Hang Seng", "Hang Seng (HK)", "asia"},
	{"000001.SS", "Shanghai", "Thượng Hải", "asia"},
	{"^N225", "Nikkei 225", "Nikkei (Nhật)", "asia"},
	{"^KS11", "KOSPI", "KOSPI (Hàn)", "asia"},
	{"^STI", "STI Singapore", "STI (Singapore)", "asia"},
	{"^FTSE", "FTSE 100", "FTSE 100 (Anh)", "europe"},
	{"^GDAXI", "DAX", "DAX (Đức)", "europe"},
	{"GC=F", "Gold", "Vàng", "commodity"},
	{"CL=F", "WTI Crude", "Dầu WTI", "commodity"},
	{"DX-Y.NYB", "USD Index", "Chỉ số USD", "commodity"},
	{"^VIX", "VIX", "VIX (sợ hãi)", "macro"},
}
