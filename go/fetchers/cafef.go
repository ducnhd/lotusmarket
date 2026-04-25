package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/types"
)

const cafefInsiderURL = "https://cafef.vn/du-lieu/Ajax/PageNew/DataHistory/GDCoDong.ashx?Symbol=%s&PageIndex=1&PageSize=50"

// cafefInsiderResponse is the top-level JSON from the CafeF AJAX API.
type cafefInsiderResponse struct {
	Data struct {
		TotalCount int                `json:"TotalCount"`
		Data       []cafefInsiderItem `json:"Data"`
	} `json:"Data"`
	Success bool `json:"Success"`
}

type cafefInsiderItem struct {
	Stock                   string  `json:"Stock"`
	TransactionMan          string  `json:"TransactionMan"`
	TransactionManPosition  string  `json:"TransactionManPosition"`
	RelatedMan              string  `json:"RelatedMan"`
	RelatedManPosition      string  `json:"RelatedManPosition"`
	VolumeBeforeTransaction int64   `json:"VolumeBeforeTransaction"`
	PlanBuyVolume           int     `json:"PlanBuyVolume"`
	PlanSellVolume          int     `json:"PlanSellVolume"`
	PlanBeginDate           string  `json:"PlanBeginDate"` // "/Date(ms)/"
	PlanEndDate             string  `json:"PlanEndDate"`
	RealBuyVolume           int     `json:"RealBuyVolume"`
	RealSellVolume          int     `json:"RealSellVolume"`
	RealEndDate             *string `json:"RealEndDate"` // null if incomplete
	VolumeAfterTransaction  int64   `json:"VolumeAfterTransaction"`
	TyLeSoHuu               float64 `json:"TyLeSoHuu"`
}

// CafefInsider fetches insider transactions for ticker from CafeF AJAX API.
// Only returns transactions with StartDate within the past 90 days.
func CafefInsider(ctx context.Context, ticker string) ([]types.InsiderTransaction, error) {
	ticker = strings.ToUpper(ticker)
	url := fmt.Sprintf(cafefInsiderURL, ticker)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; lotusmarket/1.0)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch CafeF insider %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CafeF API returned status %d for %s", resp.StatusCode, ticker)
	}

	return parseInsiderJSON(resp.Body, ticker)
}

// parseInsiderJSON decodes the CafeF JSON response body and returns typed rows.
// Exposed for testing with fixture data.
func parseInsiderJSON(r io.Reader, ticker string) ([]types.InsiderTransaction, error) {
	var apiResp cafefInsiderResponse
	if err := json.NewDecoder(r).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	if !apiResp.Success {
		return nil, fmt.Errorf("CafeF API returned success=false for %s", ticker)
	}

	cutoff := time.Now().AddDate(0, 0, -90)
	var rows []types.InsiderTransaction

	for _, item := range apiResp.Data.Data {
		startDate := parseCafefDate(item.PlanBeginDate)
		if startDate.IsZero() {
			continue
		}
		if startDate.Before(cutoff) {
			continue
		}

		relatedParty := ""
		if item.RelatedMan != "" {
			relatedParty = fmt.Sprintf("%s (%s)", item.RelatedMan, item.RelatedManPosition)
		}

		endDate := parseCafefDate(item.PlanEndDate)
		completionDate := time.Time{}
		if item.RealEndDate != nil {
			completionDate = parseCafefDate(*item.RealEndDate)
		}

		rows = append(rows, types.InsiderTransaction{
			Ticker:         ticker,
			InsiderName:    cleanInsiderName(item.TransactionMan),
			Position:       item.TransactionManPosition,
			RelatedParty:   relatedParty,
			BuyRegistered:  item.PlanBuyVolume,
			SellRegistered: item.PlanSellVolume,
			BuyResult:      item.RealBuyVolume,
			SellResult:     item.RealSellVolume,
			StartDate:      formatCafefDate(startDate),
			EndDate:        formatCafefDate(endDate),
			CompletionDate: formatCafefDate(completionDate),
			SharesAfter:    item.VolumeAfterTransaction,
			OwnershipPct:   math.Round(item.TyLeSoHuu*100) / 100,
		})
	}

	return rows, nil
}

// parseCafefDate parses CafeF's "/Date(ms)/" format to time.Time in VN timezone (UTC+7).
func parseCafefDate(s string) time.Time {
	s = strings.TrimPrefix(s, "/Date(")
	s = strings.TrimSuffix(s, ")/")
	if s == "" {
		return time.Time{}
	}
	var ms int64
	if _, err := fmt.Sscanf(s, "%d", &ms); err != nil {
		return time.Time{}
	}
	vnTZ := time.FixedZone("VN", 7*3600)
	return time.Unix(ms/1000, 0).In(vnTZ)
}

// formatCafefDate formats a time.Time to "YYYY-MM-DD", or "" if zero.
func formatCafefDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// cleanInsiderName strips HTML anchor tags sometimes wrapped around names.
func cleanInsiderName(s string) string {
	if idx := strings.Index(s, ">"); idx >= 0 {
		s = s[idx+1:]
	}
	if idx := strings.Index(s, "<"); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}
