package nlu

import "testing"

func TestParse_StockInfo(t *testing.T) {
	parser := New(nil)
	r := parser.Parse("giá ACB hiện tại")
	if r.Intent != "GET_STOCK_INFO" {
		t.Errorf("Intent = %q, want GET_STOCK_INFO", r.Intent)
	}
	if len(r.Tickers) != 1 || r.Tickers[0] != "ACB" {
		t.Errorf("Tickers = %v, want [ACB]", r.Tickers)
	}
}

func TestParse_AnalyzeTrend(t *testing.T) {
	parser := New(nil)
	r := parser.Parse("xu hướng FPT tăng hay giảm")
	if r.Intent != "ANALYZE_TREND" {
		t.Errorf("Intent = %q, want ANALYZE_TREND", r.Intent)
	}
}

func TestParse_NoTicker(t *testing.T) {
	parser := New(nil)
	r := parser.Parse("thị trường hôm nay thế nào")
	if len(r.Tickers) != 0 {
		t.Errorf("Tickers = %v, want empty", r.Tickers)
	}
}

func TestParse_StopWords(t *testing.T) {
	parser := New(nil)
	r := parser.Parse("giá cổ phiếu")
	for _, tk := range r.Tickers {
		if tk == "GIA" {
			t.Error("should not extract 'GIA' as ticker")
		}
	}
}
