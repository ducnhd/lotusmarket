package earnings

import (
	"strings"
	"testing"
)

func TestExtractTickers_Basic(t *testing.T) {
	tickers := ExtractTickers("HPG báo lãi quý 3 đạt 4.500 tỷ")
	if len(tickers) != 1 || tickers[0] != "HPG" {
		t.Errorf("got %v, want [HPG]", tickers)
	}
}

func TestExtractTickers_SkipCEO(t *testing.T) {
	tickers := ExtractTickers("CEO của VCB phát biểu")
	// CEO skipped; VCB kept
	found := map[string]bool{}
	for _, t := range tickers {
		found[t] = true
	}
	if found["CEO"] {
		t.Error("CEO should be filtered out")
	}
	if !found["VCB"] {
		t.Error("VCB should be extracted")
	}
}

func TestExtractTickers_SkipGDP(t *testing.T) {
	tickers := ExtractTickers("GDP tăng 7%, FPT tăng trưởng tốt")
	found := map[string]bool{}
	for _, tk := range tickers {
		found[tk] = true
	}
	if found["GDP"] {
		t.Error("GDP should be filtered")
	}
	if !found["FPT"] {
		t.Error("FPT should be extracted")
	}
}

func TestExtractTickers_NoDuplicates(t *testing.T) {
	tickers := ExtractTickers("HPG tăng, HPG tiếp tục tăng")
	if len(tickers) != 1 {
		t.Errorf("expected 1 unique ticker, got %d: %v", len(tickers), tickers)
	}
}

func TestParseVietnameseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"35.000", 35000},
		{"4.500", 4500},
		{"12.000", 12000},
		{"6.200", 6200},
		{"100", 100},
		{"", 0},
		{"abc", 0},
	}
	for _, tt := range tests {
		if got := ParseVietnameseNumber(tt.input); got != tt.want {
			t.Errorf("ParseVietnameseNumber(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestDetectPeriod_Quarter(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"HPG lãi ròng 6.200 tỷ Q3/2025", "Q3/2025"},
		{"Lợi nhuận quý 1/2026 đạt 2.000 tỷ", "Q1/2026"},
		{"quý 4 năm 2025 đạt kỷ lục", "Q4/2025"},
		{"9 tháng đầu năm 2025 lãi 20.000 tỷ", "9M/2025"},
		{"6 tháng năm 2026 doanh thu 50.000 tỷ", "6M/2026"},
		{"Kế hoạch lãi năm 2026 đạt 15.000 tỷ", "FY/2026"},
		{"không có kỳ nào", ""},
	}
	for _, tt := range tests {
		if got := DetectPeriod(tt.title); got != tt.want {
			t.Errorf("DetectPeriod(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}

func TestExtractProfit(t *testing.T) {
	profit, period := ExtractProfit("HPG lãi ròng 6.200 tỷ Q3/2025")
	if profit != 6200 {
		t.Errorf("profit = %d, want 6200", profit)
	}
	if period != "Q3/2025" {
		t.Errorf("period = %q, want Q3/2025", period)
	}
}

func TestExtractProfit_NotFound(t *testing.T) {
	profit, period := ExtractProfit("HPG thông báo kế hoạch sản xuất")
	if profit != 0 {
		t.Errorf("profit = %d, want 0", profit)
	}
	if period != "" {
		t.Errorf("period = %q, want empty", period)
	}
}

func TestExtractRevenue(t *testing.T) {
	rev := ExtractRevenue("doanh thu đạt 35.000 tỷ trong quý 2")
	if rev != 35000 {
		t.Errorf("revenue = %d, want 35000", rev)
	}
}

func TestExtractRevenue_NotFound(t *testing.T) {
	if rev := ExtractRevenue("HPG tăng trưởng tốt"); rev != 0 {
		t.Errorf("revenue = %d, want 0", rev)
	}
}

func TestExtractAnnualTarget(t *testing.T) {
	target := ExtractAnnualTarget("Kế hoạch lãi 12.000 tỷ năm 2026")
	if target != 12000 {
		t.Errorf("target = %d, want 12000", target)
	}
}

func TestExtractAnnualTarget_NotFound(t *testing.T) {
	if target := ExtractAnnualTarget("HPG lãi ròng 6.200 tỷ Q3/2025"); target != 0 {
		t.Errorf("target = %d, want 0", target)
	}
}

func TestBCTCDeadlines_Count(t *testing.T) {
	deadlines := BCTCDeadlines(2025)
	if len(deadlines) != 6 {
		t.Errorf("got %d deadlines, want 6", len(deadlines))
	}
}

func TestBCTCDeadlines_Dates(t *testing.T) {
	dl := BCTCDeadlines(2025)
	expected := []string{
		"2025-04-20",
		"2025-07-20",
		"2025-08-14",
		"2025-10-20",
		"2026-01-20", // Q4 deadline is in next year
		"2026-03-31", // FY audited is in next year
	}
	for i, e := range expected {
		if dl[i].Date != e {
			t.Errorf("deadline[%d].Date = %q, want %q", i, dl[i].Date, e)
		}
	}
}

func TestBCTCDeadlines_Descriptions(t *testing.T) {
	dl := BCTCDeadlines(2025)
	for _, d := range dl {
		if d.Description == "" {
			t.Errorf("empty description for deadline %q", d.Date)
		}
		if !strings.Contains(d.Description, "BCTC") {
			t.Errorf("description %q does not contain 'BCTC'", d.Description)
		}
	}
}

func TestBCTCDeadlines_Year2026(t *testing.T) {
	dl := BCTCDeadlines(2026)
	if dl[0].Date != "2026-04-20" {
		t.Errorf("Q1 2026 deadline = %q, want 2026-04-20", dl[0].Date)
	}
	if dl[4].Date != "2027-01-20" {
		t.Errorf("Q4 2026 deadline = %q, want 2027-01-20", dl[4].Date)
	}
}
