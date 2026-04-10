package technical

import "math"

type MovingAverages struct {
	MA20  *float64
	MA50  *float64
	MA200 *float64
}

type DashboardMetrics struct {
	RSI      float64
	MA20     *float64
	MA50     *float64
	MA200    *float64
	Momentum float64
	Signal   string
	Score    float64
}

func MA(prices []float64, period int) float64 {
	n := len(prices)
	if n < period {
		return 0
	}
	sum := 0.0
	for i := n - period; i < n; i++ {
		sum += prices[i]
	}
	return sum / float64(period)
}

func CalculateMovingAverages(closes []float64) MovingAverages {
	var ma MovingAverages
	if len(closes) >= 20 {
		v := MA(closes, 20)
		ma.MA20 = &v
	}
	if len(closes) >= 50 {
		v := MA(closes, 50)
		ma.MA50 = &v
	}
	if len(closes) >= 200 {
		v := MA(closes, 200)
		ma.MA200 = &v
	}
	return ma
}

func RSI(closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 50
	}
	var gains, losses float64
	for i := 1; i <= period; i++ {
		delta := closes[i] - closes[i-1]
		if delta > 0 {
			gains += delta
		} else {
			losses -= delta
		}
	}
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	for i := period + 1; i < len(closes); i++ {
		delta := closes[i] - closes[i-1]
		if delta > 0 {
			avgGain = (avgGain*float64(period-1) + delta) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) - delta) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func Momentum(closes []float64, period int) float64 {
	if len(closes) <= period {
		return 0
	}
	current := closes[len(closes)-1]
	past := closes[len(closes)-1-period]
	if past == 0 {
		return 0
	}
	return ((current - past) / past) * 100
}

func Signal(price float64, ma20, ma50 *float64, rsi float64) string {
	if ma20 == nil || ma50 == nil {
		return "HOLD"
	}
	if price > *ma20 && *ma20 > *ma50 && rsi < 70 {
		return "BUY"
	}
	if price < *ma20 && *ma20 < *ma50 && rsi > 30 {
		return "SELL"
	}
	return "HOLD"
}

func Score(buySignal, sellSignal bool, rsi float64) float64 {
	score := 50.0
	if buySignal {
		score += 25
	}
	if sellSignal {
		score -= 25
	}
	if rsi < 30 {
		score += 10
	}
	if rsi > 70 {
		score -= 10
	}
	return math.Max(0, math.Min(100, score))
}

func Dashboard(closes []float64) DashboardMetrics {
	ma := CalculateMovingAverages(closes)
	rsi := RSI(closes, 14)
	mom := Momentum(closes, 20)

	currentPrice := 0.0
	if len(closes) > 0 {
		currentPrice = closes[len(closes)-1]
	}

	sig := Signal(currentPrice, ma.MA20, ma.MA50, rsi)
	buy := sig == "BUY"
	sell := sig == "SELL"

	return DashboardMetrics{
		RSI: rsi, MA20: ma.MA20, MA50: ma.MA50, MA200: ma.MA200,
		Momentum: mom, Signal: sig, Score: Score(buy, sell, rsi),
	}
}
