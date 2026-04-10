package market

import (
	"fmt"
	"math"
	"sort"

	"github.com/ducnhd/lotusmarket/go/types"
)

type SectorFlow struct {
	Sector     string      `json:"sector"`
	Rank       int         `json:"rank"`
	Signal     string      `json:"signal"`
	NetForeign float64     `json:"net_foreign"`
	Breadth    string      `json:"breadth"`
	TopStocks  []StockFlow `json:"top_stocks"`
	Score      float64     `json:"-"`
}

type StockFlow struct {
	Ticker     string  `json:"ticker"`
	NetForeign float64 `json:"net_foreign"`
	ChangePct  float64 `json:"change_pct"`
}

func RankSectorsByFlow(stocks []types.StockData) []SectorFlow {
	type sectorAcc struct {
		totalForeignValue float64
		advances, total   int
		totalVolume       int64
		stockFlows        []StockFlow
	}
	accMap := make(map[string]*sectorAcc)
	for _, s := range stocks {
		sector := types.GetSector(s.Ticker)
		if sector == "Others" || sector == "" {
			continue
		}
		acc, ok := accMap[sector]
		if !ok {
			acc = &sectorAcc{}
			accMap[sector] = acc
		}
		flowValue := float64(s.ForeignNetVol) * s.Close
		acc.totalForeignValue += flowValue
		acc.total++
		acc.totalVolume += s.Volume
		if s.ChangePercent > 0 {
			acc.advances++
		}
		acc.stockFlows = append(acc.stockFlows, StockFlow{Ticker: s.Ticker, NetForeign: flowValue, ChangePct: s.ChangePercent})
	}
	if len(accMap) == 0 {
		return nil
	}

	allFlows := make([]float64, 0, len(accMap))
	for _, acc := range accMap {
		allFlows = append(allFlows, acc.totalForeignValue)
	}
	minFlow, maxFlow := allFlows[0], allFlows[0]
	for _, v := range allFlows[1:] {
		if v < minFlow {
			minFlow = v
		}
		if v > maxFlow {
			maxFlow = v
		}
	}
	flowRange := maxFlow - minFlow

	allVols := make([]int64, 0, len(accMap))
	for _, acc := range accMap {
		allVols = append(allVols, acc.totalVolume)
	}
	minVol, maxVol := allVols[0], allVols[0]
	for _, v := range allVols[1:] {
		if v < minVol {
			minVol = v
		}
		if v > maxVol {
			maxVol = v
		}
	}
	volRange := float64(maxVol - minVol)

	flows := make([]SectorFlow, 0, len(accMap))
	for sector, acc := range accMap {
		var normFlow float64
		if flowRange == 0 {
			normFlow = 50
		} else {
			normFlow = (acc.totalForeignValue - minFlow) / flowRange * 100
		}
		var breadth float64
		if acc.total > 0 {
			breadth = float64(acc.advances) / float64(acc.total)
		}
		var normVol float64
		if volRange == 0 {
			normVol = 50
		} else {
			normVol = float64(acc.totalVolume-minVol) / volRange * 100
		}
		score := normFlow*0.4 + breadth*100*0.3 + normVol*0.3

		sfCopy := make([]StockFlow, len(acc.stockFlows))
		copy(sfCopy, acc.stockFlows)
		sort.Slice(sfCopy, func(i, j int) bool { return math.Abs(sfCopy[i].NetForeign) > math.Abs(sfCopy[j].NetForeign) })
		top := sfCopy
		if len(top) > 2 {
			top = sfCopy[:2]
		}
		sort.Slice(top, func(i, j int) bool { return top[i].NetForeign > top[j].NetForeign })

		var signal string
		switch {
		case score > 60:
			signal = "inflow"
		case score < 40:
			signal = "outflow"
		default:
			signal = "neutral"
		}
		flows = append(flows, SectorFlow{
			Sector:     sector,
			Signal:     signal,
			NetForeign: acc.totalForeignValue,
			Breadth:    fmt.Sprintf("%d/%d tăng", acc.advances, acc.total),
			TopStocks:  top,
			Score:      score,
		})
	}
	sort.Slice(flows, func(i, j int) bool { return flows[i].Score > flows[j].Score })
	for i := range flows {
		flows[i].Rank = i + 1
	}
	return flows
}
