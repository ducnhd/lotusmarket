package technical

import (
	"time"

	"github.com/lotusmarket/lotusmarket-go/types"
)

func AggregateWeekly(daily []types.StockData) []types.StockData {
	if len(daily) == 0 {
		return nil
	}
	var weekly []types.StockData
	var current *types.StockData

	for _, d := range daily {
		t, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		year, week := t.ISOWeek()
		weekKey := year*100 + week

		if current == nil || getWeekKey(current.Date) != weekKey {
			if current != nil {
				weekly = append(weekly, *current)
			}
			current = &types.StockData{
				Ticker: d.Ticker, Date: d.Date, Open: d.Open,
				High: d.High, Low: d.Low, Close: d.Close, Volume: d.Volume,
			}
		} else {
			if d.High > current.High {
				current.High = d.High
			}
			if d.Low < current.Low {
				current.Low = d.Low
			}
			current.Close = d.Close
			current.Volume += d.Volume
			current.Date = d.Date
		}
	}
	if current != nil {
		weekly = append(weekly, *current)
	}
	return weekly
}

func getWeekKey(dateStr string) int {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0
	}
	year, week := t.ISOWeek()
	return year*100 + week
}
