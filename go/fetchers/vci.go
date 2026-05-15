// VCI (Vietcap) IQ-Insight fetcher — corporate actions (dividends, bonus
// issues, ESOP) for VN tickers. Source: iq.vietcap.com.vn — the same backend
// used by vnstock python library. Free, no API key.
package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const vciIQEventsURL = "https://iq.vietcap.com.vn/api/iq-insight-service/v1/events"

// DividendEvent is a corporate action on a VN ticker.
type DividendEvent struct {
	Ticker      string
	EventType   int    // 1 = cash dividend, 2 = stock issuance (bonus/rights/ESOP)
	ExDate      string // "YYYY-MM-DD"
	RecordDate  string
	PaymentDate string
	RateCash    float64 // % of par value (10,000 VND/share). 10 = 1,000 VND/share
	RateSplit   float64 // % bonus shares. 50 = 1 new share per 2 existing
	Year        int
	Source      string
	SourceID    string
	TitleVi     string
}

// AmountVND returns VND per share for cash dividends. Returns 0 for non-cash.
func (d DividendEvent) AmountVND() float64 {
	if d.EventType != 1 {
		return 0
	}
	return d.RateCash / 100.0 * 10000.0
}

type vciIQResponse struct {
	Status int `json:"status"`
	Data   struct {
		Content []vciIQEvent `json:"content"`
	} `json:"data"`
}

type vciIQEvent struct {
	ID            string  `json:"id"`
	Ticker        string  `json:"ticker"`
	OrganCode     string  `json:"organCode"`
	EventCode     string  `json:"eventCode"`
	EventTitleVi  string  `json:"eventTitleVi"`
	EventTitleEn  string  `json:"eventTitleEn"`
	PublicDate    string  `json:"publicDate"`
	RecordDate    string  `json:"recordDate"`
	ExrightDate   string  `json:"exrightDate"`
	PayoutDate    string  `json:"payoutDate"`
	ListingDate   string  `json:"listingDate"`
	ValuePerShare float64 `json:"valuePerShare"`
	ExerciseRatio float64 `json:"exerciseRatio"`
	Category      string  `json:"category"`
}

var vciCodeToType = map[string]int{
	"DIV": 1,
	"ISS": 2,
}

// VCIDividends fetches all corporate action events for a ticker.
// Returns events sorted oldest-first by ExDate.
func VCIDividends(ctx context.Context, ticker string) ([]DividendEvent, error) {
	url := fmt.Sprintf("%s?ticker=%s&pageSize=100&page=0", vciIQEventsURL, strings.ToUpper(ticker))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("vci %s build: %w", ticker, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 lotusmarket")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vci %s fetch: %w", ticker, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vci %s status %d", ticker, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vci %s read: %w", ticker, err)
	}
	var out vciIQResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("vci %s decode: %w", ticker, err)
	}
	return parseVCIEvents(out.Data.Content, ticker), nil
}

func parseVCIEvents(events []vciIQEvent, ticker string) []DividendEvent {
	out := make([]DividendEvent, 0, len(events))
	for _, e := range events {
		t, ok := vciCodeToType[e.EventCode]
		if !ok {
			continue
		}
		exDate := normVCIDate(e.ExrightDate)
		if exDate == "" {
			continue
		}
		d := DividendEvent{
			Ticker:      strings.ToUpper(ticker),
			EventType:   t,
			ExDate:      exDate,
			RecordDate:  normVCIDate(e.RecordDate),
			PaymentDate: normVCIDate(e.PayoutDate),
			Source:      "vci_iq",
			SourceID:    e.ID,
			TitleVi:     e.EventTitleVi,
		}
		if t == 1 {
			if e.ValuePerShare > 0 {
				d.RateCash = e.ValuePerShare / 100
			} else {
				d.RateCash = e.ExerciseRatio * 100
			}
		} else {
			d.RateSplit = e.ExerciseRatio * 100
		}
		if parsed, err := time.Parse("2006-01-02", exDate); err == nil {
			d.Year = parsed.Year()
		}
		out = append(out, d)
	}
	return out
}

func normVCIDate(s string) string {
	if len(s) < 10 {
		return ""
	}
	d := s[:10]
	if t, err := time.Parse("2006-01-02", d); err != nil || t.Year() < 2000 {
		return ""
	}
	return d
}
