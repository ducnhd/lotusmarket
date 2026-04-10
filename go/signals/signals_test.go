package signals

import (
	"testing"

	"github.com/lotusmarket/lotusmarket-go/types"
)

func TestClassifyVolumeSurge_NoSurge(t *testing.T) {
	s := ClassifyVolumeSurge(1000, 1000, 0, 1.0)
	if s.Signal != "" {
		t.Errorf("Signal = %q, want empty", s.Signal)
	}
}

func TestClassifyVolumeSurge_Accumulation(t *testing.T) {
	s := ClassifyVolumeSurge(3000, 1000, 0, 2.0)
	if s.Signal != "accumulation" {
		t.Errorf("Signal = %q, want accumulation", s.Signal)
	}
	if s.Intensity != "strong_surge" {
		t.Errorf("Intensity = %q, want strong_surge", s.Intensity)
	}
}

func TestClassifyVolumeSurge_Distribution(t *testing.T) {
	s := ClassifyVolumeSurge(2500, 1000, 0, -1.5)
	if s.Signal != "distribution" {
		t.Errorf("Signal = %q, want distribution", s.Signal)
	}
}

func TestClassifyVolumeSurge_ZScore(t *testing.T) {
	s := ClassifyVolumeSurge(5000, 1000, 500, 1.0)
	if s.Intensity != "strong_surge" {
		t.Errorf("Intensity = %q, want strong_surge", s.Intensity)
	}
}

func TestVolumeSurgeLabel(t *testing.T) {
	tests := []struct {
		signal, intensity, want string
	}{
		{"accumulation", "strong_surge", "Tiền đổ vào mạnh"},
		{"accumulation", "surge", "Tiền đang đổ vào"},
		{"distribution", "strong_surge", "Bán tháo mạnh"},
		{"distribution", "surge", "Đang bán tháo"},
		{"high_activity", "surge", "Giao dịch đột biến"},
		{"", "", ""},
	}
	for _, tt := range tests {
		got := VolumeSurgeLabel(tt.signal, tt.intensity)
		if got != tt.want {
			t.Errorf("VolumeSurgeLabel(%q, %q) = %q, want %q", tt.signal, tt.intensity, got, tt.want)
		}
	}
}

func TestComputeDomesticFlow_BuyPressure(t *testing.T) {
	stocks := []types.StockData{
		{Ticker: "ACB", Volume: 1000, BidVol: 700, AskVol: 300, ForeignBuyVol: 50, ForeignSellVol: 50},
		{Ticker: "VNM", Volume: 2000, BidVol: 1400, AskVol: 600, ForeignBuyVol: 100, ForeignSellVol: 100},
	}
	flow := ComputeDomesticFlow(stocks)
	if flow.Signal != "mua mạnh" {
		t.Errorf("Signal = %q, want 'mua mạnh'", flow.Signal)
	}
	if flow.TotalVolume != 3000 {
		t.Errorf("TotalVolume = %d, want 3000", flow.TotalVolume)
	}
}

func TestComputeDomesticFlow_SellPressure(t *testing.T) {
	stocks := []types.StockData{
		{Ticker: "ACB", Volume: 1000, BidVol: 200, AskVol: 800, ForeignBuyVol: 50, ForeignSellVol: 50},
	}
	flow := ComputeDomesticFlow(stocks)
	if flow.Signal != "bán mạnh" {
		t.Errorf("Signal = %q, want 'bán mạnh'", flow.Signal)
	}
}
