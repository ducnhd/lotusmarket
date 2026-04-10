package signals

import "github.com/lotusmarket/lotusmarket-go/types"

type VolumeSurge struct {
	Ticker      string  `json:"ticker"`
	CurrentVol  int64   `json:"current_vol"`
	AvgVol20    float64 `json:"avg_vol_20"`
	StdDevVol20 float64 `json:"stddev_vol_20,omitempty"`
	Ratio       float64 `json:"ratio"`
	ZScore      float64 `json:"z_score,omitempty"`
	PriceChange float64 `json:"price_change"`
	Signal      string  `json:"signal"`
	Intensity   string  `json:"intensity"`
	Label       string  `json:"label"`
}

func ClassifyVolumeSurge(currentVol int64, avgVol20, stdDevVol20, priceChange float64) VolumeSurge {
	s := VolumeSurge{CurrentVol: currentVol, AvgVol20: avgVol20, StdDevVol20: stdDevVol20, PriceChange: priceChange}
	if avgVol20 <= 0 {
		return s
	}
	s.Ratio = float64(currentVol) / avgVol20

	if stdDevVol20 > 0 {
		s.ZScore = (float64(currentVol) - avgVol20) / stdDevVol20
		if s.ZScore < 2.0 {
			return s
		}
		if s.ZScore >= 3.0 {
			s.Intensity = "strong_surge"
		} else {
			s.Intensity = "surge"
		}
	} else {
		if s.Ratio < 2.0 {
			return s
		}
		if s.Ratio >= 3.0 {
			s.Intensity = "strong_surge"
		} else {
			s.Intensity = "surge"
		}
	}

	switch {
	case priceChange >= 0.5:
		s.Signal = "accumulation"
	case priceChange <= -0.5:
		s.Signal = "distribution"
	default:
		s.Signal = "high_activity"
	}
	s.Label = VolumeSurgeLabel(s.Signal, s.Intensity)
	return s
}

func VolumeSurgeLabel(signal, intensity string) string {
	switch {
	case signal == "accumulation" && intensity == "strong_surge":
		return "Tiền đổ vào mạnh"
	case signal == "accumulation" && intensity == "surge":
		return "Tiền đang đổ vào"
	case signal == "distribution" && intensity == "strong_surge":
		return "Bán tháo mạnh"
	case signal == "distribution" && intensity == "surge":
		return "Đang bán tháo"
	case signal == "high_activity":
		return "Giao dịch đột biến"
	default:
		return ""
	}
}

func DetectVolumeSurges(stocks []types.StockData, avgVolumes, stdDevVolumes map[string]float64) []VolumeSurge {
	var surges []VolumeSurge
	for _, s := range stocks {
		avg, ok := avgVolumes[s.Ticker]
		if !ok || avg <= 0 {
			continue
		}
		stddev := 0.0
		if stdDevVolumes != nil {
			stddev = stdDevVolumes[s.Ticker]
		}
		priceChange := 0.0
		if s.RefPrice > 0 {
			priceChange = (s.Close - s.RefPrice) / s.RefPrice * 100
		}
		surge := ClassifyVolumeSurge(s.Volume, avg, stddev, priceChange)
		if surge.Signal != "" {
			surge.Ticker = s.Ticker
			surges = append(surges, surge)
		}
	}
	return surges
}
