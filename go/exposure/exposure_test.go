package exposure

import "testing"

// stubProvider returns a fixed sequence of closes for any symbol.
func stubProvider(closes []float64) HistoryProvider {
	return func(symbol string, days int) []DriverClose {
		out := make([]DriverClose, len(closes))
		for i, c := range closes {
			out[i] = DriverClose{Date: "2025-01-01", Close: c}
		}
		return out
	}
}

func TestGetMappingsForTicker_Mapped(t *testing.T) {
	m := GetMappingsForTicker("HPG")
	if len(m) == 0 {
		t.Fatal("expected mappings for HPG, got none")
	}
	// First HPG driver should be CSI 300
	if m[0].Driver != "000300.SS" {
		t.Errorf("first HPG driver = %q, want 000300.SS", m[0].Driver)
	}
}

func TestGetMappingsForTicker_Unmapped(t *testing.T) {
	m := GetMappingsForTicker("UNKNOWN")
	if m != nil {
		t.Errorf("expected nil for unknown ticker, got %v", m)
	}
}

func TestGetMappingsForTicker_Banks(t *testing.T) {
	banks := []string{"VCB", "TCB", "BID", "CTG", "ACB", "MBB", "HDB", "LPB", "SHB", "SSB", "STB", "TPB", "VIB", "VPB"}
	for _, b := range banks {
		if m := GetMappingsForTicker(b); len(m) == 0 {
			t.Errorf("bank %s has no mappings", b)
		}
	}
}

func TestUniqueDriverSymbols_Count(t *testing.T) {
	syms := UniqueDriverSymbols()
	if len(syms) < 18 {
		t.Errorf("expected >=18 unique driver symbols, got %d: %v", len(syms), syms)
	}
	// Should be sorted
	for i := 1; i < len(syms); i++ {
		if syms[i] < syms[i-1] {
			t.Errorf("not sorted at index %d: %q > %q", i, syms[i-1], syms[i])
		}
	}
}

func TestUniqueDriverSymbols_ContainsExpected(t *testing.T) {
	syms := UniqueDriverSymbols()
	symSet := make(map[string]bool, len(syms))
	for _, s := range syms {
		symSet[s] = true
	}
	expected := []string{"^VIX", "VNM", "EEM", "BHP", "TIO=F", "QQQ", "CL=F"}
	for _, e := range expected {
		if !symSet[e] {
			t.Errorf("expected symbol %q not found in UniqueDriverSymbols", e)
		}
	}
}

func TestAnalyze_Unmapped(t *testing.T) {
	hp := stubProvider([]float64{100, 110})
	a := Analyze("UNKNOWN", hp)
	if a.HasMapping {
		t.Error("expected HasMapping=false for unknown ticker")
	}
	if a.Ticker != "UNKNOWN" {
		t.Errorf("Ticker = %q, want UNKNOWN", a.Ticker)
	}
}

func TestAnalyze_NetImpactSign_Positive(t *testing.T) {
	// Stub: all drivers get +10% return.
	// HPG first driver 000300.SS has historicalR=+0.461 → est = +10*0.461 = +4.61%
	// All HPG drivers have positive r — net should be positive.
	closes := make([]float64, 72) // 20*3+12
	for i := range closes {
		closes[i] = 100.0 + float64(i)*0.5 // monotone rising
	}
	hp := stubProvider(closes)
	a := Analyze("HPG", hp)
	if !a.HasMapping {
		t.Fatal("HPG should have mappings")
	}
	if a.NetImpact <= 0 {
		t.Errorf("expected positive net impact for rising HPG drivers, got %.4f", a.NetImpact)
	}
	if a.Summary == "" {
		t.Error("summary should not be empty")
	}
}

func TestAnalyze_SummaryContainsTicker(t *testing.T) {
	closes := make([]float64, 72)
	for i := range closes {
		closes[i] = 100.0 + float64(i)*0.3
	}
	hp := stubProvider(closes)
	a := Analyze("ACB", hp)
	if !a.HasMapping {
		t.Fatal("ACB should have mappings (bank pack)")
	}
	// Summary should mention the ticker
	if len(a.Summary) == 0 {
		t.Error("summary empty")
	}
}

func TestAnalyze_Sides(t *testing.T) {
	closes := make([]float64, 80)
	for i := range closes {
		closes[i] = 50.0 + float64(i)
	}
	hp := stubProvider(closes)
	a := Analyze("HPG", hp)
	// HPG has Input, Peer, Output sides
	if len(a.Input) == 0 {
		t.Error("HPG should have Input drivers")
	}
	if len(a.Peer) == 0 {
		t.Error("HPG should have Peer drivers")
	}
	if len(a.Output) == 0 {
		t.Error("HPG should have Output drivers")
	}
}

func TestFormatForPrompt_Empty(t *testing.T) {
	hp := stubProvider(nil)
	out := FormatForPrompt(nil, hp)
	if out != "" {
		t.Errorf("expected empty string for nil tickers, got %q", out)
	}
}

func TestFormatForPrompt_NoData(t *testing.T) {
	// With empty history, no drivers have return data → no prompt block
	hp := stubProvider(nil)
	out := FormatForPrompt([]string{"HPG"}, hp)
	if out != "" {
		t.Errorf("expected empty string when no history, got %q", out)
	}
}

func TestFormatForPrompt_WithData(t *testing.T) {
	closes := make([]float64, 80)
	for i := range closes {
		closes[i] = 100.0 + float64(i)
	}
	hp := stubProvider(closes)
	out := FormatForPrompt([]string{"HPG", "ACB"}, hp)
	if out == "" {
		t.Error("expected non-empty prompt block when drivers have data")
	}
	// Should contain the section header
	if len(out) < 10 {
		t.Errorf("prompt block too short: %q", out)
	}
}

func TestRegimeChangeFlags_BelowThreshold(t *testing.T) {
	// Flat prices → no regime flags
	closes := make([]float64, 35)
	for i := range closes {
		closes[i] = 100.0
	}
	hp := stubProvider(closes)
	out := RegimeChangeFlags(hp)
	if out != "" {
		t.Errorf("expected empty, got %q", out)
	}
}

func TestRegimeChangeFlags_AboveThreshold(t *testing.T) {
	// Strong rising trend → should produce regime flags
	closes := make([]float64, 35)
	for i := range closes {
		closes[i] = 100.0 + float64(i)*2 // +60% over 30d, >10%
	}
	hp := stubProvider(closes)
	out := RegimeChangeFlags(hp)
	if out == "" {
		t.Error("expected regime flags for large driver moves")
	}
}

func TestBounceSignals_NoOversold(t *testing.T) {
	closes := make([]float64, 15)
	for i := range closes {
		closes[i] = 100.0
	}
	hp := stubProvider(closes)
	rsi := map[string]float64{"HPG": 55, "ACB": 60}
	out := BounceSignals(rsi, hp)
	if out != "" {
		t.Errorf("expected empty for non-oversold tickers, got %q", out)
	}
}

func TestBounceSignals_OversoldPositiveDrivers(t *testing.T) {
	// Rising drivers → net impact positive, RSI < 40 → bounce candidate
	closes := make([]float64, 80)
	for i := range closes {
		closes[i] = 100.0 + float64(i)*0.5
	}
	hp := stubProvider(closes)
	rsi := map[string]float64{"HPG": 32}
	out := BounceSignals(rsi, hp)
	if out == "" {
		t.Error("expected bounce signal for oversold HPG with positive drivers")
	}
}
