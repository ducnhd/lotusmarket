package ratings

import "testing"

func TestComputeBullish(t *testing.T) {
	closes := make([]float64, 30)
	for i := range closes {
		closes[i] = 100 + float64(i)*1.5 // steady uptrend
	}
	volumes := make([]int64, 30)
	for i := range volumes {
		volumes[i] = 1000 + int64(i*50) // rising volume
	}
	ma20 := closes[len(closes)-10]
	ma50 := closes[5]
	ma200 := closes[0]
	r := Compute(closes, volumes, 75, 58, &ma20, &ma50, &ma200)
	if r.OverallVerdict != "Outperform" {
		t.Errorf("bullish path: want Outperform, got %s (gauge=%d)", r.OverallVerdict, r.OverallGauge)
	}
	if r.TrendStrength < 4 {
		t.Errorf("trend stars: want >=4, got %d", r.TrendStrength)
	}
}

func TestComputeBearish(t *testing.T) {
	closes := make([]float64, 30)
	for i := range closes {
		closes[i] = 200 - float64(i)*1.5 // downtrend
	}
	volumes := make([]int64, 30)
	for i := range volumes {
		volumes[i] = 1000
	}
	ma20 := closes[0]
	ma50 := closes[0] + 5
	ma200 := closes[0] + 10
	r := Compute(closes, volumes, 15, 28, &ma20, &ma50, &ma200)
	if r.OverallVerdict == "Outperform" {
		t.Errorf("bearish path: want Neutral/Underperform, got %s (gauge=%d)", r.OverallVerdict, r.OverallGauge)
	}
}

func TestComputeShortHistoryDefaults(t *testing.T) {
	r := Compute([]float64{100, 101, 102}, nil, 50, 50, nil, nil, nil)
	if r.PriceStrength != 3 || r.TrendStrength != 3 {
		t.Errorf("short history should return all 3s; got price=%d trend=%d", r.PriceStrength, r.TrendStrength)
	}
}
