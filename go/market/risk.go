package market

func ScoreVIX(vix float64) int {
	switch {
	case vix < 15:
		return 10
	case vix < 20:
		return 30
	case vix < 25:
		return 50
	case vix < 30:
		return 70
	default:
		return 90
	}
}

func ScoreYieldCurve(spread float64) int {
	switch {
	case spread > 1.0:
		return 10
	case spread > 0.5:
		return 30
	case spread > 0:
		return 50
	default:
		return 90
	}
}

func ScoreNewsSentiment(negative, total int) int {
	if total == 0 {
		return 20
	}
	ratio := float64(negative) / float64(total) * 100
	switch {
	case ratio > 60:
		return 80
	case ratio > 40:
		return 50
	default:
		return 20
	}
}

func ScoreFedTrend(current, previous float64) int {
	diff := current - previous
	switch {
	case diff > 0.05:
		return 70
	case diff < -0.05:
		return 15
	default:
		return 40
	}
}
