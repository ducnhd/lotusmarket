package fetchers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// cafefFixture is a minimal CafeF JSON response with two transactions.
// PlanBeginDate is set far in the future so the 90-day cutoff never filters them.
func cafefFixtureJSON(beginMs, endMs int64) string {
	return `{
		"Success": true,
		"Data": {
			"TotalCount": 2,
			"Data": [
				{
					"Stock": "HPG",
					"TransactionMan": "<a href='#'>Nguyễn Văn A</a>",
					"TransactionManPosition": "Chủ tịch HĐQT",
					"RelatedMan": "",
					"RelatedManPosition": "",
					"VolumeBeforeTransaction": 1000000,
					"PlanBuyVolume": 500000,
					"PlanSellVolume": 0,
					"PlanBeginDate": "/Date(` + itoa64(beginMs) + `)/",
					"PlanEndDate": "/Date(` + itoa64(endMs) + `)/",
					"RealBuyVolume": 480000,
					"RealSellVolume": 0,
					"RealEndDate": null,
					"VolumeAfterTransaction": 1480000,
					"TyLeSoHuu": 2.56789
				},
				{
					"Stock": "HPG",
					"TransactionMan": "Trần Thị B",
					"TransactionManPosition": "Phó TGĐ",
					"RelatedMan": "Nguyễn Văn C",
					"RelatedManPosition": "Anh ruột",
					"VolumeBeforeTransaction": 200000,
					"PlanBuyVolume": 0,
					"PlanSellVolume": 100000,
					"PlanBeginDate": "/Date(` + itoa64(beginMs) + `)/",
					"PlanEndDate": "/Date(` + itoa64(endMs) + `)/",
					"RealBuyVolume": 0,
					"RealSellVolume": 95000,
					"RealEndDate": "/Date(` + itoa64(endMs) + `)/",
					"VolumeAfterTransaction": 105000,
					"TyLeSoHuu": 0.12345
				}
			]
		}
	}`
}

func itoa64(n int64) string {
	s := ""
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func TestParseInsiderJSON_Basic(t *testing.T) {
	// Use a timestamp well within the 90-day window.
	beginMs := time.Now().UnixMilli()
	endMs := time.Now().Add(30 * 24 * time.Hour).UnixMilli()

	body := cafefFixtureJSON(beginMs, endMs)
	rows, err := parseInsiderJSON(strings.NewReader(body), "HPG")
	if err != nil {
		t.Fatalf("parseInsiderJSON error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	// First row — HTML name stripped
	r0 := rows[0]
	if r0.InsiderName != "Nguyễn Văn A" {
		t.Errorf("InsiderName = %q, want 'Nguyễn Văn A'", r0.InsiderName)
	}
	if r0.Ticker != "HPG" {
		t.Errorf("Ticker = %q, want HPG", r0.Ticker)
	}
	if r0.BuyRegistered != 500000 {
		t.Errorf("BuyRegistered = %d, want 500000", r0.BuyRegistered)
	}
	if r0.BuyResult != 480000 {
		t.Errorf("BuyResult = %d, want 480000", r0.BuyResult)
	}
	if r0.SharesAfter != 1480000 {
		t.Errorf("SharesAfter = %d, want 1480000", r0.SharesAfter)
	}
	if r0.OwnershipPct != 2.57 {
		t.Errorf("OwnershipPct = %f, want 2.57", r0.OwnershipPct)
	}
	if r0.CompletionDate != "" {
		t.Errorf("CompletionDate should be empty for null RealEndDate, got %q", r0.CompletionDate)
	}

	// Second row — related party
	r1 := rows[1]
	if r1.RelatedParty == "" {
		t.Error("RelatedParty should not be empty")
	}
	if !strings.Contains(r1.RelatedParty, "Nguyễn Văn C") {
		t.Errorf("RelatedParty = %q, want to contain 'Nguyễn Văn C'", r1.RelatedParty)
	}
	if r1.CompletionDate == "" {
		t.Error("CompletionDate should be set when RealEndDate is not null")
	}
}

func TestParseInsiderJSON_Cutoff(t *testing.T) {
	// Begin date 200 days ago — should be filtered out by 90-day cutoff
	beginMs := time.Now().Add(-200 * 24 * time.Hour).UnixMilli()
	endMs := time.Now().Add(-180 * 24 * time.Hour).UnixMilli()

	body := cafefFixtureJSON(beginMs, endMs)
	rows, err := parseInsiderJSON(strings.NewReader(body), "HPG")
	if err != nil {
		t.Fatalf("parseInsiderJSON error: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("got %d rows, want 0 (old transactions filtered by 90-day cutoff)", len(rows))
	}
}

func TestParseInsiderJSON_SuccessFalse(t *testing.T) {
	body := `{"Success": false, "Data": {"TotalCount": 0, "Data": []}}`
	_, err := parseInsiderJSON(strings.NewReader(body), "XYZ")
	if err == nil {
		t.Error("expected error for success=false, got nil")
	}
}

func TestParseInsiderJSON_BadJSON(t *testing.T) {
	_, err := parseInsiderJSON(strings.NewReader("{invalid}"), "XYZ")
	if err == nil {
		t.Error("expected error for bad JSON")
	}
}

func TestParseCafefDate_Valid(t *testing.T) {
	// Known timestamp: 2024-01-15 in VN tz
	ms := int64(1705276800000) // 2024-01-15 00:00:00 UTC → 2024-01-15 07:00 VN
	s := "/Date(" + itoa64(ms) + ")/"
	got := parseCafefDate(s)
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
	if got.Year() != 2024 {
		t.Errorf("year = %d, want 2024", got.Year())
	}
}

func TestParseCafefDate_Empty(t *testing.T) {
	if got := parseCafefDate(""); !got.IsZero() {
		t.Error("expected zero time for empty string")
	}
}

func TestCleanInsiderName_WithHTMLTag(t *testing.T) {
	if got := cleanInsiderName("<a href='#'>Nguyễn Văn A</a>"); got != "Nguyễn Văn A" {
		t.Errorf("got %q, want 'Nguyễn Văn A'", got)
	}
}

func TestCleanInsiderName_Plain(t *testing.T) {
	if got := cleanInsiderName("  Trần Thị B  "); got != "Trần Thị B" {
		t.Errorf("got %q", got)
	}
}

func TestCafefInsider_HTTPTest(t *testing.T) {
	beginMs := time.Now().UnixMilli()
	endMs := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
	body := cafefFixtureJSON(beginMs, endMs)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	// Patch the URL constant by calling parseInsiderJSON directly via HTTP client.
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http: %v", err)
	}
	defer resp.Body.Close()

	rows, err := parseInsiderJSON(resp.Body, "HPG")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("got %d rows via httptest, want 2", len(rows))
	}
}
