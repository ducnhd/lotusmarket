package portfolio

import (
	"math"
	"testing"
)

func TestCompoundReturns(t *testing.T) {
	segments := []ReturnSegment{{StartValue: 100, EndValue: 110}, {StartValue: 110, EndValue: 132}}
	twr := CompoundReturns(segments)
	if math.Abs(twr-0.32) > 0.001 {
		t.Errorf("TWR = %v, want 0.32", twr)
	}
}

func TestCompoundReturns_Empty(t *testing.T) {
	if twr := CompoundReturns(nil); twr != 0 {
		t.Errorf("TWR(nil) = %v, want 0", twr)
	}
}

func TestCompoundReturns_ZeroStart(t *testing.T) {
	segments := []ReturnSegment{{StartValue: 0, EndValue: 100}}
	if twr := CompoundReturns(segments); twr != 0 {
		t.Errorf("TWR zero start = %v, want 0", twr)
	}
}

func TestCalculatePositionConfidence_Buy(t *testing.T) {
	pa := PositionAnalysis{
		Action: "BUY_MORE", TimeframeAlign: "bullish", RSI: 35,
		MA20: ptrFloat(100), MA50: ptrFloat(95), MA200: ptrFloat(90),
		PriceVsMA20: "above", PriceVsMA50: "above",
		PE: 10, Beta: 1.0, DividendYield: 3.0,
		ForeignNetVol: 100000, BuySellRatio: 1.5,
	}
	result := CalculatePositionConfidence(&pa, "positive")
	if result.Score < 60 {
		t.Errorf("Buy confidence = %d, want >= 60", result.Score)
	}
}

func ptrFloat(f float64) *float64 { return &f }
