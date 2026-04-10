package backtest

import (
	"fmt"
	"math"

	"github.com/ducnhd/lotusmarket/go/technical"
	"github.com/ducnhd/lotusmarket/go/types"
)

type Config struct {
	Strategy string  `json:"strategy"` // "rsi", "ma_cross", "combined"
	Capital  float64 `json:"capital"`
	RSIBuy   float64 `json:"rsi_buy"`
	RSISell  float64 `json:"rsi_sell"`
}

type Trade struct {
	Date   string  `json:"date"`
	Action string  `json:"action"`
	Price  float64 `json:"price"`
	Shares float64 `json:"shares"`
	Value  float64 `json:"value"`
	Cash   float64 `json:"cash"`
	PnL    float64 `json:"pnl"`
}

type Result struct {
	Config        Config  `json:"config"`
	Trades        []Trade `json:"trades"`
	TotalReturn   float64 `json:"total_return_pct"`
	MaxDrawdown   float64 `json:"max_drawdown_pct"`
	WinRate       float64 `json:"win_rate_pct"`
	BuyHoldReturn float64 `json:"buy_hold_return_pct"`
	FinalCapital  float64 `json:"final_capital"`
	TradeCount    int     `json:"trade_count"`
}

// Run executes a backtest strategy on the provided price history.
// history must have at least 20 entries with Close prices.
func Run(cfg Config, history []types.StockData) (*Result, error) {
	if cfg.Capital <= 0 {
		cfg.Capital = 100_000_000 // 100M VND
	}
	if cfg.RSIBuy <= 0 {
		cfg.RSIBuy = 30
	}
	if cfg.RSISell <= 0 {
		cfg.RSISell = 70
	}
	if len(history) < 20 {
		return nil, fmt.Errorf("insufficient data: need 20+ days, got %d", len(history))
	}

	result := &Result{Config: cfg}

	closes := make([]float64, len(history))
	for i, h := range history {
		closes[i] = h.Close
	}

	cash := cfg.Capital
	shares := 0.0
	var trades []Trade
	peakValue := cfg.Capital
	maxDrawdown := 0.0
	wins := 0
	losses := 0

	startIdx := 50
	if startIdx >= len(history) {
		startIdx = 20
	}

	for i := startIdx; i < len(history); i++ {
		price := history[i].Close
		date := history[i].Date

		rsi := technical.RSI(closes[:i+1], 14)
		ma := technical.CalculateMovingAverages(closes[:i+1])

		buySignal := false
		sellSignal := false

		switch cfg.Strategy {
		case "rsi":
			buySignal = rsi < cfg.RSIBuy
			sellSignal = rsi > cfg.RSISell
		case "ma_cross":
			if ma.MA20 != nil && ma.MA50 != nil {
				buySignal = *ma.MA20 > *ma.MA50 && price > *ma.MA20
				sellSignal = *ma.MA20 < *ma.MA50 && price < *ma.MA20
			}
		case "combined":
			maBuy := ma.MA20 != nil && ma.MA50 != nil && *ma.MA20 > *ma.MA50
			maSell := ma.MA20 != nil && ma.MA50 != nil && *ma.MA20 < *ma.MA50
			buySignal = rsi < cfg.RSIBuy && maBuy
			sellSignal = rsi > cfg.RSISell || maSell
		default:
			buySignal = rsi < cfg.RSIBuy
			sellSignal = rsi > cfg.RSISell
		}

		if buySignal && shares == 0 && cash > 0 {
			shares = math.Floor(cash / price)
			if shares > 0 {
				cost := shares * price
				cash -= cost
				trades = append(trades, Trade{
					Date: date, Action: "BUY", Price: price,
					Shares: shares, Value: cost, Cash: cash,
				})
			}
		} else if sellSignal && shares > 0 {
			value := shares * price
			pnl := value - (shares * trades[len(trades)-1].Price)
			if pnl > 0 {
				wins++
			} else {
				losses++
			}
			cash += value
			trades = append(trades, Trade{
				Date: date, Action: "SELL", Price: price,
				Shares: shares, Value: value, Cash: cash, PnL: pnl,
			})
			shares = 0
		}

		totalValue := cash + shares*price
		if totalValue > peakValue {
			peakValue = totalValue
		}
		dd := (peakValue - totalValue) / peakValue * 100
		if dd > maxDrawdown {
			maxDrawdown = dd
		}
	}

	lastPrice := history[len(history)-1].Close
	result.FinalCapital = cash + shares*lastPrice
	result.TotalReturn = (result.FinalCapital - cfg.Capital) / cfg.Capital * 100
	result.MaxDrawdown = maxDrawdown
	result.Trades = trades
	result.TradeCount = len(trades)

	totalTrades := wins + losses
	if totalTrades > 0 {
		result.WinRate = float64(wins) / float64(totalTrades) * 100
	}

	firstPrice := history[startIdx].Close
	result.BuyHoldReturn = (lastPrice - firstPrice) / firstPrice * 100

	return result, nil
}
