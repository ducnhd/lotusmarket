package technical

import (
	"math"
	"testing"

	"github.com/ducnhd/lotusmarket/go/types"
)

func TestMA(t *testing.T) {
	prices := make([]float64, 20)
	for i := range prices {
		prices[i] = float64(i + 1)
	}
	got := MA(prices, 20)
	want := 10.5
	if math.Abs(got-want) > 0.01 {
		t.Errorf("MA(20) = %v, want %v", got, want)
	}
}

func TestMA_InsufficientData(t *testing.T) {
	got := MA([]float64{1, 2, 3}, 20)
	if got != 0 {
		t.Errorf("MA insufficient = %v, want 0", got)
	}
}

func TestRSI_Uptrend(t *testing.T) {
	prices := make([]float64, 30)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}
	rsi := RSI(prices, 14)
	if rsi <= 50 {
		t.Errorf("RSI uptrend = %v, want > 50", rsi)
	}
}

func TestRSI_InsufficientData(t *testing.T) {
	rsi := RSI([]float64{1, 2, 3}, 14)
	if rsi != 50 {
		t.Errorf("RSI insufficient = %v, want 50", rsi)
	}
}

func TestMomentum(t *testing.T) {
	prices := make([]float64, 11)
	prices[0] = 100
	for i := 1; i <= 10; i++ {
		prices[i] = 100 + float64(i)
	}
	got := Momentum(prices, 10)
	if math.Abs(got-10.0) > 0.01 {
		t.Errorf("Momentum = %v, want 10.0", got)
	}
}

func TestSignal_Buy(t *testing.T) {
	ma20, ma50 := 100.0, 90.0
	sig := Signal(110, &ma20, &ma50, 55)
	if sig != "BUY" {
		t.Errorf("Signal = %q, want BUY", sig)
	}
}

func TestSignal_Sell(t *testing.T) {
	ma20, ma50 := 90.0, 100.0
	sig := Signal(80, &ma20, &ma50, 55)
	if sig != "SELL" {
		t.Errorf("Signal = %q, want SELL", sig)
	}
}

func TestSignal_Hold(t *testing.T) {
	ma20, ma50 := 100.0, 100.0
	sig := Signal(100, &ma20, &ma50, 50)
	if sig != "HOLD" {
		t.Errorf("Signal = %q, want HOLD", sig)
	}
}

func TestScore(t *testing.T) {
	score := Score(true, false, 50)
	if score != 75 {
		t.Errorf("Score(buy) = %v, want 75", score)
	}
	score = Score(false, true, 75)
	if score != 15 {
		t.Errorf("Score(sell+highRSI) = %v, want 15", score)
	}
}

func TestDashboard(t *testing.T) {
	prices := make([]float64, 200)
	for i := range prices {
		prices[i] = 20000 + float64(i)*10
	}
	d := Dashboard(prices)
	if d.RSI <= 0 {
		t.Error("Dashboard RSI should be > 0")
	}
	if d.MA20 == nil {
		t.Error("Dashboard MA20 should not be nil")
	}
	if d.Signal == "" {
		t.Error("Dashboard Signal should not be empty")
	}
}

func TestAggregateWeekly_Empty(t *testing.T) {
	result := AggregateWeekly(nil)
	if result != nil {
		t.Errorf("AggregateWeekly(nil) = %v, want nil", result)
	}
}

func TestAggregateWeekly_SameWeek(t *testing.T) {
	daily := []types.StockData{
		{Ticker: "ACB", Date: "2025-03-03", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000},
		{Ticker: "ACB", Date: "2025-03-04", Open: 105, High: 115, Low: 95, Close: 108, Volume: 1500},
		{Ticker: "ACB", Date: "2025-03-05", Open: 108, High: 120, Low: 100, Close: 112, Volume: 2000},
	}
	result := AggregateWeekly(daily)
	if len(result) != 1 {
		t.Fatalf("got %d weeks, want 1", len(result))
	}
	if result[0].Open != 100 {
		t.Errorf("Open = %v, want 100", result[0].Open)
	}
	if result[0].High != 120 {
		t.Errorf("High = %v, want 120", result[0].High)
	}
	if result[0].Low != 90 {
		t.Errorf("Low = %v, want 90", result[0].Low)
	}
	if result[0].Close != 112 {
		t.Errorf("Close = %v, want 112", result[0].Close)
	}
	if result[0].Volume != 4500 {
		t.Errorf("Volume = %v, want 4500", result[0].Volume)
	}
}
