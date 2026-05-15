// Package ratings — deterministic 6-dimensional star ratings for a stock,
// mirroring SSI iResearch's "Đánh giá kỹ thuật" panel. Inputs: price/volume
// history + RSI + MA20/50/200. Outputs: 6 sub-scores (1-5 each) plus an
// overall 0-100 gauge and verdict (Outperform / Neutral / Underperform).
//
// All math is deterministic — no AI, no external data, no DB. Safe to call
// on any closes/volumes slice client-side.
package ratings

import "math"

// Ratings is the per-ticker 6-dim breakdown.
type Ratings struct {
	PriceStrength    int    `json:"price_strength"`    // 1-5: derived from external score bucket
	TrendStrength    int    `json:"trend_strength"`    // 1-5: MA alignment
	ShortTermPos     int    `json:"short_term_pos"`    // 1-5: RSI sweet spot
	MoneyFlow        int    `json:"money_flow"`        // 1-5: volume vs MA20 + 5d direction
	VolatilityRating int    `json:"volatility_rating"` // 1-5: 20d annualized vol bucket
	BaseRange        int    `json:"base_range"`        // 1-5: closeness to MA20
	OverallVerdict   string `json:"overall_verdict"`   // "Outperform" / "Neutral" / "Underperform"
	OverallGauge     int    `json:"overall_gauge"`     // 0-100
}

// Compute turns price history + score + RSI/MA into the 6-dim breakdown.
// `closes` is oldest-first. `volumes` may be nil — MoneyFlow falls back to 3.
// `score` is the existing 0-100 technical signal score (from technical.Score).
// `ma20`/`ma50`/`ma200` may be nil for short histories.
func Compute(closes []float64, volumes []int64, score float64, rsi float64, ma20, ma50, ma200 *float64) Ratings {
	r := Ratings{
		PriceStrength: 3, TrendStrength: 3, ShortTermPos: 3,
		MoneyFlow: 3, VolatilityRating: 3, BaseRange: 3,
	}
	if len(closes) < 20 {
		return r
	}
	current := closes[len(closes)-1]

	r.PriceStrength = priceStrengthStars(score)
	r.TrendStrength = trendStars(current, ma20, ma50, ma200)
	r.ShortTermPos = rsiStars(rsi)
	r.MoneyFlow = moneyFlowStars(closes, volumes)
	r.VolatilityRating = volatilityStars(closes)
	r.BaseRange = baseRangeStars(current, ma20)

	sum := r.PriceStrength + r.TrendStrength + r.ShortTermPos +
		r.MoneyFlow + r.VolatilityRating + r.BaseRange
	r.OverallGauge = sum * 100 / 30
	switch {
	case r.OverallGauge >= 60:
		r.OverallVerdict = "Outperform"
	case r.OverallGauge >= 40:
		r.OverallVerdict = "Neutral"
	default:
		r.OverallVerdict = "Underperform"
	}
	return r
}

func priceStrengthStars(score float64) int {
	switch {
	case score >= 80:
		return 5
	case score >= 60:
		return 4
	case score >= 40:
		return 3
	case score >= 20:
		return 2
	default:
		return 1
	}
}

func trendStars(price float64, ma20, ma50, ma200 *float64) int {
	if ma20 == nil {
		return 3
	}
	above20 := price > *ma20
	above50 := ma50 != nil && price > *ma50
	above200 := ma200 != nil && price > *ma200
	stacked := ma50 != nil && ma200 != nil && *ma20 > *ma50 && *ma50 > *ma200

	switch {
	case above20 && above50 && above200 && stacked:
		return 5
	case above20 && above50 && above200:
		return 4
	case above20 && above50:
		return 3
	case above20:
		return 2
	default:
		return 1
	}
}

func rsiStars(rsi float64) int {
	switch {
	case rsi >= 50 && rsi <= 65:
		return 5
	case (rsi >= 45 && rsi < 50) || (rsi > 65 && rsi <= 70):
		return 4
	case (rsi >= 35 && rsi < 45) || (rsi > 70 && rsi <= 75):
		return 3
	case (rsi >= 25 && rsi < 35) || (rsi > 75 && rsi <= 80):
		return 2
	default:
		return 1
	}
}

func moneyFlowStars(closes []float64, volumes []int64) int {
	if len(volumes) < 21 || len(closes) < 6 {
		return 3
	}
	var sum int64
	start := len(volumes) - 21
	for i := start; i < len(volumes)-1; i++ {
		sum += volumes[i]
	}
	avg20 := float64(sum) / 20
	if avg20 <= 0 {
		return 3
	}
	ratio := float64(volumes[len(volumes)-1]) / avg20
	ret5 := (closes[len(closes)-1] - closes[len(closes)-6]) / closes[len(closes)-6] * 100
	priceUp := ret5 > 0.5
	priceDown := ret5 < -0.5

	switch {
	case ratio >= 1.5 && priceUp:
		return 5
	case ratio >= 1.2 && priceUp:
		return 4
	case ratio >= 0.8 && ratio < 1.2:
		return 3
	case ratio >= 1.2 && priceDown:
		return 1 // distribution
	case ratio < 0.8:
		return 2
	default:
		return 3
	}
}

func volatilityStars(closes []float64) int {
	n := len(closes)
	if n < 21 {
		return 3
	}
	rets := make([]float64, 0, 20)
	for i := n - 20; i < n; i++ {
		if closes[i-1] <= 0 {
			continue
		}
		rets = append(rets, math.Log(closes[i]/closes[i-1]))
	}
	if len(rets) < 5 {
		return 3
	}
	var mean float64
	for _, x := range rets {
		mean += x
	}
	mean /= float64(len(rets))
	var variance float64
	for _, x := range rets {
		variance += (x - mean) * (x - mean)
	}
	variance /= float64(len(rets) - 1)
	annualized := math.Sqrt(variance) * math.Sqrt(252) * 100 // %

	switch {
	case annualized >= 15 && annualized <= 25:
		return 5
	case (annualized >= 10 && annualized < 15) || (annualized > 25 && annualized <= 35):
		return 4
	case annualized > 35 && annualized <= 45:
		return 3
	case annualized > 45 && annualized <= 60:
		return 2
	default:
		return 1
	}
}

func baseRangeStars(current float64, ma20 *float64) int {
	if ma20 == nil || *ma20 <= 0 {
		return 3
	}
	dist := math.Abs(current-*ma20) / *ma20 * 100
	switch {
	case dist <= 2:
		return 5
	case dist <= 4:
		return 4
	case dist <= 7:
		return 3
	case dist <= 12:
		return 2
	default:
		return 1
	}
}
