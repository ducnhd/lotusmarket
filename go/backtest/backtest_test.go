package backtest

import (
	"math"
	"testing"

	"github.com/ducnhd/lotusmarket/go/types"
)

func makeHistory(n int, basePrice float64) []types.StockData {
	history := make([]types.StockData, n)
	for i := range history {
		// Oscillating price to trigger RSI signals
		price := basePrice + float64(i%20)*100 - 1000
		if price < 1000 {
			price = 1000
		}
		history[i] = types.StockData{
			Date:  "2025-01-01",
			Close: price,
		}
	}
	return history
}

func TestRun_InsufficientData(t *testing.T) {
	_, err := Run(Config{Strategy: "rsi"}, []types.StockData{{Close: 100}})
	if err == nil {
		t.Error("expected error for insufficient data")
	}
}

func TestRun_RSIStrategy(t *testing.T) {
	history := makeHistory(200, 25000)
	result, err := Run(Config{Strategy: "rsi", Capital: 100_000_000}, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FinalCapital <= 0 {
		t.Errorf("FinalCapital = %v, want > 0", result.FinalCapital)
	}
	if result.MaxDrawdown < 0 || result.MaxDrawdown > 100 {
		t.Errorf("MaxDrawdown = %v, want 0-100", result.MaxDrawdown)
	}
}

func TestRun_MACrossStrategy(t *testing.T) {
	// Create trending data for MA cross signals
	history := make([]types.StockData, 100)
	for i := range history {
		history[i] = types.StockData{
			Date:  "2025-01-01",
			Close: 20000 + float64(i)*50 + float64(i%5)*100,
		}
	}
	result, err := Run(Config{Strategy: "ma_cross", Capital: 100_000_000}, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Config.Strategy != "ma_cross" {
		t.Errorf("Strategy = %q, want ma_cross", result.Config.Strategy)
	}
}

func TestRun_DefaultCapital(t *testing.T) {
	history := makeHistory(100, 25000)
	result, err := Run(Config{Strategy: "rsi"}, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Config.Capital != 100_000_000 {
		t.Errorf("Default capital = %v, want 100M", result.Config.Capital)
	}
}

func TestRun_BuyHoldReturn(t *testing.T) {
	// Steadily rising prices
	history := make([]types.StockData, 100)
	for i := range history {
		history[i] = types.StockData{Date: "2025-01-01", Close: 20000 + float64(i)*100}
	}
	result, err := Run(Config{Strategy: "rsi"}, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Buy & hold on rising prices should be positive
	if result.BuyHoldReturn <= 0 {
		t.Errorf("BuyHoldReturn = %v, want > 0 for rising prices", result.BuyHoldReturn)
	}
}

func TestRun_WinRate(t *testing.T) {
	history := makeHistory(200, 25000)
	result, err := Run(Config{Strategy: "rsi"}, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WinRate < 0 || result.WinRate > 100 {
		t.Errorf("WinRate = %v, want 0-100", result.WinRate)
	}
	_ = math.Abs(0) // use math import
}
