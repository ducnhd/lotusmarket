package portfolio

import "math"

type PositionAnalysis struct {
	Action                   string
	TimeframeAlign           string
	RSI                      float64
	MA20, MA50, MA200        *float64
	PriceVsMA20, PriceVsMA50 string
	PE, Beta, DividendYield  float64
	ForeignNetVol            int64
	BuySellRatio             float64
}

type ConfidenceDetail struct {
	Score      int            `json:"score"`
	Components map[string]int `json:"components"`
}

func CalculatePositionConfidence(pa *PositionAnalysis, newsSentiment string) ConfidenceDetail {
	components := make(map[string]int)
	isBuy := pa.Action == "BUY_MORE"
	isSell := pa.Action == "SELL_PARTIAL" || pa.Action == "SELL_ALL" || pa.Action == "CUT_LOSS"

	// 1. Timeframe alignment (25%)
	switch pa.TimeframeAlign {
	case "bullish":
		if isBuy {
			components["timeframe"] = 100
		} else if isSell {
			components["timeframe"] = 30
		} else {
			components["timeframe"] = 60
		}
	case "bearish":
		if isSell {
			components["timeframe"] = 100
		} else if isBuy {
			components["timeframe"] = 30
		} else {
			components["timeframe"] = 60
		}
	default:
		components["timeframe"] = 40
	}

	// 2. RSI zone (20%)
	if isBuy {
		if pa.RSI < 30 {
			components["rsi"] = 90
		} else if pa.RSI < 50 {
			components["rsi"] = 60
		} else {
			components["rsi"] = 30
		}
	} else if isSell {
		if pa.RSI > 70 {
			components["rsi"] = 90
		} else if pa.RSI > 50 {
			components["rsi"] = 60
		} else {
			components["rsi"] = 30
		}
	} else {
		if pa.RSI >= 30 && pa.RSI <= 70 {
			components["rsi"] = 70
		} else {
			components["rsi"] = 40
		}
	}

	// 3. MA trend (20%)
	maScore := 40
	if pa.MA20 != nil && pa.MA50 != nil && pa.MA200 != nil {
		if pa.PriceVsMA20 == "above" && pa.PriceVsMA50 == "above" {
			if isBuy || pa.Action == "HOLD" {
				maScore = 100
			}
		} else if pa.PriceVsMA20 == "below" && pa.PriceVsMA50 == "below" {
			if isSell {
				maScore = 100
			} else if isBuy {
				maScore = 30
			} else {
				maScore = 50
			}
		} else {
			maScore = 50
		}
	}
	components["ma_trend"] = maScore

	// 4. Fundamentals (15%)
	fundScore := 0
	if pa.PE > 0 && pa.PE < 15 {
		fundScore += 50
	} else if pa.PE > 0 {
		fundScore += 20
	}
	if pa.Beta >= 0.5 && pa.Beta <= 1.5 {
		fundScore += 25
	}
	if pa.DividendYield > 0 {
		fundScore += 25
	}
	if fundScore > 100 {
		fundScore = 100
	}
	components["fundamental"] = fundScore

	// 5. Flow (10%)
	flowScore := 50
	if isBuy {
		if pa.ForeignNetVol > 0 {
			flowScore += 25
		} else if pa.ForeignNetVol < 0 {
			flowScore -= 25
		}
		if pa.BuySellRatio > 1 {
			flowScore += 25
		} else if pa.BuySellRatio > 0 && pa.BuySellRatio < 1 {
			flowScore -= 25
		}
	} else if isSell {
		if pa.ForeignNetVol < 0 {
			flowScore += 25
		} else if pa.ForeignNetVol > 0 {
			flowScore -= 25
		}
		if pa.BuySellRatio < 1 && pa.BuySellRatio > 0 {
			flowScore += 25
		} else if pa.BuySellRatio > 1 {
			flowScore -= 25
		}
	}
	if flowScore < 0 {
		flowScore = 0
	}
	if flowScore > 100 {
		flowScore = 100
	}
	components["flow"] = flowScore

	// 6. News (10%)
	if isBuy {
		switch newsSentiment {
		case "positive":
			components["news"] = 90
		case "mixed":
			components["news"] = 50
		default:
			components["news"] = 30
		}
	} else if isSell {
		switch newsSentiment {
		case "negative":
			components["news"] = 90
		case "mixed":
			components["news"] = 50
		default:
			components["news"] = 30
		}
	} else {
		components["news"] = 50
	}

	weights := map[string]int{"timeframe": 25, "rsi": 20, "ma_trend": 20, "fundamental": 15, "flow": 10, "news": 10}
	totalWeight, weightedSum := 0, 0
	for name, w := range weights {
		if score, ok := components[name]; ok {
			totalWeight += w
			weightedSum += score * w
		}
	}
	finalScore := 0
	if totalWeight > 0 {
		finalScore = int(math.Round(float64(weightedSum) / float64(totalWeight)))
	}
	return ConfidenceDetail{Score: finalScore, Components: components}
}
