// FRED (Federal Reserve Economic Data) fetcher — US macro indicators that
// matter for VN equity decisions: Fed funds rate, 10Y Treasury yield, VIX,
// WTI oil, yield-curve spread. Requires a free API key from
// https://fred.stlouisfed.org/docs/api/api_key.html — pass it explicitly so
// library users keep control of their own credentials.
package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultFREDSeries — 5 core series tracked alongside VN markets.
// Keys are the FRED series ID; values are Vietnamese display names.
var DefaultFREDSeries = map[string]string{
	"DFF":        "Lãi suất Fed",
	"DGS10":      "Lợi suất trái phiếu 10 năm Mỹ",
	"VIXCLS":     "Chỉ số biến động VIX",
	"DCOILWTICO": "Giá dầu WTI",
	"T10Y2Y":     "Chênh lệch lợi suất 10Y-2Y",
}

// FREDObservation is one (date, value) point.
type FREDObservation struct {
	Date  string // YYYY-MM-DD
	Value float64
}

type fredAPIResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

// FREDSeries fetches the latest `days` of observations for one series.
// Returns observations newest-first. Missing values (FRED's ".") are skipped.
func FREDSeries(ctx context.Context, apiKey, seriesID string, days int) ([]FREDObservation, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("fred: api key required")
	}
	if days <= 0 {
		days = 30
	}
	start := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	url := fmt.Sprintf(
		"https://api.stlouisfed.org/fred/series/observations"+
			"?series_id=%s&observation_start=%s&api_key=%s"+
			"&file_type=json&sort_order=desc&limit=30",
		seriesID, start, apiKey,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fred %s build: %w", seriesID, err)
	}
	req.Header.Set("User-Agent", "lotusmarket")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fred %s fetch: %w", seriesID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fred %s status %d: %s", seriesID, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var data fredAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("fred %s decode: %w", seriesID, err)
	}
	out := make([]FREDObservation, 0, len(data.Observations))
	for _, o := range data.Observations {
		if o.Value == "." {
			continue
		}
		v, err := strconv.ParseFloat(o.Value, 64)
		if err != nil {
			continue
		}
		out = append(out, FREDObservation{Date: o.Date, Value: v})
	}
	return out, nil
}

// FREDLatest returns the most recent observation for a series. Useful for
// "current Fed rate" / "current VIX" style queries.
func FREDLatest(ctx context.Context, apiKey, seriesID string) (*FREDObservation, error) {
	obs, err := FREDSeries(ctx, apiKey, seriesID, 30)
	if err != nil {
		return nil, err
	}
	if len(obs) == 0 {
		return nil, fmt.Errorf("fred %s: no observations", seriesID)
	}
	return &obs[0], nil
}

// FREDSnapshot is one row of the default-series snapshot table.
type FREDSnapshot struct {
	SeriesID string
	NameVi   string
	Date     string
	Value    float64
}

// FREDAllLatest fetches the latest value for each series in DefaultFREDSeries.
// Failures on individual series are silently skipped — caller gets partial
// results. Useful for a one-shot dashboard query.
func FREDAllLatest(ctx context.Context, apiKey string) []FREDSnapshot {
	out := []FREDSnapshot{}
	for id, nameVi := range DefaultFREDSeries {
		o, err := FREDLatest(ctx, apiKey, id)
		if err != nil {
			continue
		}
		out = append(out, FREDSnapshot{
			SeriesID: id, NameVi: nameVi,
			Date: o.Date, Value: o.Value,
		})
	}
	return out
}
