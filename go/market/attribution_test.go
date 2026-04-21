package market

import (
	"math"
	"testing"
)

func TestAttributeDrivers_FlatMove(t *testing.T) {
	f := DriverFeatures{PriceChangePct: 0.02}
	a := AttributeDrivers(f)
	if a.Dominant != "UNKNOWN" {
		t.Errorf("flat move should yield UNKNOWN dominant, got %q", a.Dominant)
	}
}

func TestAttributeDrivers_NewsDriven(t *testing.T) {
	// Strong positive news, weak other signals → NEWS dominant on up move
	f := DriverFeatures{
		PriceChangePct:   2.5,
		NewsCount:        5,
		NewsSentimentAvg: 0.8,
		MarketDeltaPct:   0.1,
		ForeignNetVND:    1e9,
	}
	a := AttributeDrivers(f)
	if a.Dominant != "NEWS" {
		t.Errorf("expected NEWS dominant, got %q", a.Dominant)
	}
	if a.NewsPct < 30 {
		t.Errorf("NewsPct should be largest (>30), got %.1f", a.NewsPct)
	}
}

func TestAttributeDrivers_ForeignFlow(t *testing.T) {
	// Big foreign net buy on an up day
	f := DriverFeatures{
		PriceChangePct: 1.5,
		ForeignNetVND:  50e9, // 50 tỷ mua ròng
		MarketDeltaPct: 0.2,
	}
	a := AttributeDrivers(f)
	if a.Dominant != "FOREIGN" {
		t.Errorf("expected FOREIGN dominant, got %q", a.Dominant)
	}
}

func TestAttributeDrivers_MisalignedSignal(t *testing.T) {
	// Price UP but news NEGATIVE — misaligned, NEWS gets only 30% credit
	f := DriverFeatures{
		PriceChangePct:   1.0,
		NewsCount:        3,
		NewsSentimentAvg: -0.7,
		MarketDeltaPct:   0.8,
	}
	a := AttributeDrivers(f)
	// Dominant should be a POSITIVE signal matching the up move (MARKET)
	if a.Dominant == "NEWS" {
		t.Errorf("negative news on up day should NOT be dominant, got %q", a.Dominant)
	}
}

func TestAttributeDrivers_PctsSumSane(t *testing.T) {
	f := DriverFeatures{
		PriceChangePct:   1.5,
		ForeignNetVND:    20e9,
		MarketDeltaPct:   0.5,
		SectorDeltaPct:   1.0,
		NewsCount:        2,
		NewsSentimentAvg: 0.5,
		RSI:              55,
		MASignal:         "BUY",
		VolumeRatio:      1.8,
	}
	a := AttributeDrivers(f)
	total := a.NewsPct + a.ForeignPct + a.SectorPct + a.MarketPct + a.TechnicalPct + a.VolumePct + a.ResidualPct
	if math.Abs(total-100) > 1 {
		t.Errorf("percentages should sum to ~100, got %.1f", total)
	}
}
