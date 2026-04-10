package market

import (
	"testing"

	"github.com/ducnhd/lotusmarket/go/types"
)

func TestScoreBreadth(t *testing.T) {
	tests := []struct {
		ratio float64
		want  int
	}{{70, 80}, {50, 50}, {30, 20}}
	for _, tt := range tests {
		if got := ScoreBreadth(tt.ratio); got != tt.want {
			t.Errorf("ScoreBreadth(%v) = %d, want %d", tt.ratio, got, tt.want)
		}
	}
}

func TestScoreForeignFlow(t *testing.T) {
	if s := ScoreForeignFlow(200e9); s != 90 {
		t.Errorf("200B = %d, want 90", s)
	}
	if s := ScoreForeignFlow(-200e9); s != 20 {
		t.Errorf("-200B = %d, want 20", s)
	}
	if s := ScoreForeignFlow(0); s != 50 {
		t.Errorf("0 = %d, want 50", s)
	}
}

func TestScoreVIX(t *testing.T) {
	if s := ScoreVIX(12); s != 10 {
		t.Errorf("VIX 12 = %d, want 10", s)
	}
	if s := ScoreVIX(35); s != 90 {
		t.Errorf("VIX 35 = %d, want 90", s)
	}
}

func TestScoreYieldCurve(t *testing.T) {
	if s := ScoreYieldCurve(1.5); s != 10 {
		t.Errorf("1.5 = %d, want 10", s)
	}
	if s := ScoreYieldCurve(-0.5); s != 90 {
		t.Errorf("-0.5 = %d, want 90", s)
	}
}

func TestRankSectorsByFlow_Empty(t *testing.T) {
	if r := RankSectorsByFlow(nil); r != nil {
		t.Errorf("got %v, want nil", r)
	}
}

func TestRankSectorsByFlow_Basic(t *testing.T) {
	stocks := []types.StockData{
		{Ticker: "ACB", Close: 25000, Volume: 1000000, ChangePercent: 1.5, ForeignNetVol: 500},
		{Ticker: "VNM", Close: 70000, Volume: 500000, ChangePercent: -0.5, ForeignNetVol: -200},
	}
	result := RankSectorsByFlow(stocks)
	if len(result) == 0 {
		t.Fatal("expected sectors")
	}
	if result[0].Sector != "Ngân hàng" {
		t.Errorf("top = %q, want Ngân hàng", result[0].Sector)
	}
	if result[0].Rank != 1 {
		t.Errorf("Rank = %d, want 1", result[0].Rank)
	}
}
