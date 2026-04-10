package market

func ScoreBreadth(ratio float64) int {
	switch {
	case ratio > 60:
		return 80
	case ratio >= 40:
		return 50
	default:
		return 20
	}
}

func ScoreForeignFlow(netTotal float64) int {
	const threshold = 100e9
	if netTotal > threshold {
		return 90
	}
	if netTotal > 0 {
		ratio := netTotal / threshold
		return 60 + int(ratio*30)
	}
	if netTotal == 0 {
		return 50
	}
	absVal := -netTotal
	if absVal > threshold {
		return 20
	}
	ratio := absVal / threshold
	return 50 - int(ratio*30)
}

func ComputePulseScore(breadthScore, foreignFlowScore, volumeScore, riskScore int) (int, string) {
	score := (breadthScore*30 + foreignFlowScore*25 + volumeScore*20 + riskScore*25) / 100
	var signal string
	switch {
	case score > 60:
		signal = "green"
	case score >= 40:
		signal = "yellow"
	default:
		signal = "red"
	}
	return score, signal
}
