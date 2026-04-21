// Package market — price driver attribution.
//
// Given a stock's daily move and contextual features (foreign flow, market delta,
// sector delta, news, technical, volume), decompose the move into % contribution
// from each driver. Used to answer "why did this stock move today?" — the result
// can be displayed as progress bars or used as input to retrospective analysis.
//
// Algorithm (fully deterministic, zero external calls):
//  1. Map each raw feature into a signed "impact" number using closed-form functions.
//  2. Aligned drivers (sign matches today's move) get full weight; misaligned get 30%.
//  3. Normalise to sum = 100%. Dominant driver = argmax(aligned weight).
//  4. Residual = 100 - sum — how much of the move we could not attribute.
package market

import "math"

// DriverFeatures is the raw daily feature set for a single ticker.
type DriverFeatures struct {
	PriceChangePct   float64 // today's close change %
	ForeignNetVND    float64 // foreign buy - sell (VND)
	MarketDeltaPct   float64 // market index change % (e.g. VN-Index)
	SectorDeltaPct   float64 // sector average change %
	NewsCount        int
	NewsSentimentAvg float64 // [-1, 1]
	RSI              float64
	MASignal         string  // "BUY" | "SELL" | "HOLD" | ""
	VolumeRatio      float64 // today / MA20
}

// Attribution is the % contribution of each driver to the price move.
// Pcts sum to at most 100; ResidualPct = 100 - sum.
type Attribution struct {
	NewsPct      float64 `json:"news_pct"`
	ForeignPct   float64 `json:"foreign_pct"`
	SectorPct    float64 `json:"sector_pct"`
	MarketPct    float64 `json:"market_pct"`
	TechnicalPct float64 `json:"technical_pct"`
	VolumePct    float64 `json:"volume_pct"`
	ResidualPct  float64 `json:"residual_pct"`
	Dominant     string  `json:"dominant"` // "NEWS"|"FOREIGN"|"SECTOR"|"MARKET"|"TECHNICAL"|"VOLUME"|"UNKNOWN"
}

// AttributeDrivers decomposes a stock's daily move into driver contributions.
// Returns all-zero Attribution{Dominant:"UNKNOWN"} when the move is essentially flat
// (|move| < 0.05%) or no driver has non-zero impact.
func AttributeDrivers(f DriverFeatures) Attribution {
	move := f.PriceChangePct
	if math.Abs(move) < 0.05 {
		return Attribution{Dominant: "UNKNOWN"}
	}

	signs := map[string]float64{
		"FOREIGN":   clampSignal(f.ForeignNetVND/10_000_000_000, 5),
		"MARKET":    f.MarketDeltaPct,
		"SECTOR":    f.SectorDeltaPct - f.MarketDeltaPct, // sector alpha (excess over market)
		"NEWS":      f.NewsSentimentAvg * math.Log1p(float64(f.NewsCount)) * 3,
		"TECHNICAL": technicalSignal(f.RSI, f.MASignal),
		"VOLUME":    volumeSignal(f.VolumeRatio, move),
	}

	moveSign := 1.0
	if move < 0 {
		moveSign = -1.0
	}

	var total float64
	weights := make(map[string]float64)
	for k, v := range signs {
		w := math.Abs(v)
		if w == 0 {
			continue
		}
		if sign(v) != moveSign {
			w *= 0.3 // misaligned driver gets 30% credit
		}
		weights[k] = w
		total += w
	}
	if total < 0.01 {
		return Attribution{Dominant: "UNKNOWN"}
	}

	attr := Attribution{
		ForeignPct:   round1(weights["FOREIGN"] / total * 100),
		MarketPct:    round1(weights["MARKET"] / total * 100),
		SectorPct:    round1(weights["SECTOR"] / total * 100),
		NewsPct:      round1(weights["NEWS"] / total * 100),
		TechnicalPct: round1(weights["TECHNICAL"] / total * 100),
		VolumePct:    round1(weights["VOLUME"] / total * 100),
	}
	attr.ResidualPct = round1(math.Max(0, 100-attr.NewsPct-attr.ForeignPct-attr.SectorPct-attr.MarketPct-attr.TechnicalPct-attr.VolumePct))

	// Dominant = highest aligned contributor
	max := 0.0
	best := "UNKNOWN"
	for k, w := range weights {
		if sign(signs[k]) != moveSign {
			continue
		}
		if w > max {
			max = w
			best = k
		}
	}
	attr.Dominant = best
	return attr
}

func clampSignal(v, limit float64) float64 {
	if v > limit {
		return limit
	}
	if v < -limit {
		return -limit
	}
	return v
}

func technicalSignal(rsi float64, maSignal string) float64 {
	sig := 0.0
	switch maSignal {
	case "BUY":
		sig = 2
	case "SELL":
		sig = -2
	}
	if rsi > 70 {
		sig--
	} else if rsi > 0 && rsi < 30 {
		sig++
	}
	return sig
}

func volumeSignal(ratio, move float64) float64 {
	if ratio < 1.2 {
		return 0
	}
	impact := math.Log(ratio) * 2
	if move < 0 {
		impact = -impact
	}
	return impact
}

func sign(v float64) float64 {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
