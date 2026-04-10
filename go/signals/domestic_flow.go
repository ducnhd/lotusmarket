package signals

import (
	"sort"

	"github.com/lotusmarket/lotusmarket-go/types"
)

type DomesticFlow struct {
	TotalVolume    int64    `json:"total_volume"`
	DomesticVolume int64    `json:"domestic_volume"`
	ForeignVolume  int64    `json:"foreign_volume"`
	BuyPressure    float64  `json:"buy_pressure"`
	Signal         string   `json:"signal"`
	TopBuyers      []string `json:"top_buyers,omitempty"`
	TopSellers     []string `json:"top_sellers,omitempty"`
}

func ComputeDomesticFlow(stocks []types.StockData) DomesticFlow {
	var totalVol, foreignVol, totalBid, totalAsk int64

	type stockPressure struct {
		ticker string
		ratio  float64
	}
	var pressures []stockPressure

	for _, s := range stocks {
		totalVol += s.Volume
		fVol := s.ForeignBuyVol + s.ForeignSellVol
		foreignVol += fVol
		totalBid += s.BidVol
		totalAsk += s.AskVol

		if s.BidVol+s.AskVol > 0 {
			ratio := float64(s.BidVol) / float64(s.BidVol+s.AskVol) * 100
			pressures = append(pressures, stockPressure{s.Ticker, ratio})
		}
	}

	buyPressure := 50.0
	if totalBid+totalAsk > 0 {
		buyPressure = float64(totalBid) / float64(totalBid+totalAsk) * 100
	}

	df := DomesticFlow{
		TotalVolume:    totalVol,
		DomesticVolume: totalVol - foreignVol,
		ForeignVolume:  foreignVol,
		BuyPressure:    buyPressure,
		Signal:         classifyDomesticSignal(buyPressure),
	}

	sort.Slice(pressures, func(i, j int) bool { return pressures[i].ratio > pressures[j].ratio })
	for i := 0; i < 3 && i < len(pressures); i++ {
		df.TopBuyers = append(df.TopBuyers, pressures[i].ticker)
	}
	for i := len(pressures) - 1; i >= 0 && len(df.TopSellers) < 3; i-- {
		df.TopSellers = append(df.TopSellers, pressures[i].ticker)
	}

	return df
}

func classifyDomesticSignal(buyPressure float64) string {
	switch {
	case buyPressure >= 60:
		return "mua mạnh"
	case buyPressure >= 53:
		return "nghiêng mua"
	case buyPressure >= 47:
		return "cân bằng"
	case buyPressure >= 40:
		return "nghiêng bán"
	default:
		return "bán mạnh"
	}
}
