# lotusmarket

Vietnamese Stock Market Toolkit for Go & Python.

Real-time quotes, historical OHLCV, technical analysis, sentiment, money flow signals, and more for HOSE/HNX stocks. All public APIs — no API key required.

[![Go Tests](https://github.com/ducnhd/lotusmarket/actions/workflows/go-test.yml/badge.svg)](https://github.com/ducnhd/lotusmarket/actions/workflows/go-test.yml)
[![Python Tests](https://github.com/ducnhd/lotusmarket/actions/workflows/python-test.yml/badge.svg)](https://github.com/ducnhd/lotusmarket/actions/workflows/python-test.yml)
[![PyPI](https://img.shields.io/pypi/v/lotusmarket)](https://pypi.org/project/lotusmarket/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Install

**Go:**
```bash
go get github.com/ducnhd/lotusmarket/go
```

**Python:**
```bash
pip install lotusmarket

# With fetchers (httpx for API calls):
pip install lotusmarket[fetchers]

# Everything:
pip install lotusmarket[all]
```

## Quick Start

### Go

```go
package main

import (
    "context"
    "fmt"

    "github.com/ducnhd/lotusmarket/go/fetchers"
    "github.com/ducnhd/lotusmarket/go/technical"
    "github.com/ducnhd/lotusmarket/go/sentiment"
    "github.com/ducnhd/lotusmarket/go/signals"
    "github.com/ducnhd/lotusmarket/go/market"
)

func main() {
    ctx := context.Background()

    // 1. Real-time quote
    stock, _ := fetchers.VPS(ctx, "ACB")
    fmt.Printf("ACB: %.0f VND, vol: %d, foreign net: %d\n",
        stock.Close, stock.Volume, stock.ForeignNetVol)

    // 2. Multiple stocks at once
    stocks, _ := fetchers.VPSMultiple(ctx, []string{"ACB", "VNM", "FPT", "HPG"})

    // 3. Historical OHLCV (200 days)
    history, _ := fetchers.EntradeHistory(ctx, "ACB", 200)
    closes := make([]float64, len(history))
    for i, h := range history {
        closes[i] = h.Close
    }

    // 4. Technical analysis — all-in-one
    d := technical.Dashboard(closes)
    fmt.Printf("RSI: %.1f, MA20: %.0f, Signal: %s, Score: %.0f\n",
        d.RSI, *d.MA20, d.Signal, d.Score)

    // 5. Individual indicators
    rsi := technical.RSI(closes, 14)
    ma20 := technical.MA(closes, 20)
    mom := technical.Momentum(closes, 20)

    // 6. Weekly candles from daily
    weekly := technical.AggregateWeekly(history)

    // 7. Vietnamese sentiment
    r := sentiment.Analyze("Ngân hàng tăng mạnh, khối ngoại mua ròng")
    fmt.Printf("Sentiment: %s (%.2f)\n", r.Label, r.Score)

    // 8. Domestic flow — buy/sell pressure
    flow := signals.ComputeDomesticFlow(stocks)
    fmt.Printf("Domestic: %s (buy pressure: %.1f%%)\n", flow.Signal, flow.BuyPressure)

    // 9. Sector flow ranking
    sectors := market.RankSectorsByFlow(stocks)
    for _, s := range sectors {
        fmt.Printf("#%d %s: %s (%.0f)\n", s.Rank, s.Sector, s.Signal, s.Score)
    }

    // 10. Fundamentals
    fund, _ := fetchers.KBS(ctx, "ACB")
    fmt.Printf("P/E: %.1f, P/B: %.1f, EPS: %.0f\n", fund.PE, fund.PB, fund.EPS)

    // 11. With fallback (VPS -> Entrade)
    stock2, _ := fetchers.StockWithFallback(ctx, "VNM")
    _ = stock2
    _ = rsi
    _ = ma20
    _ = mom
    _ = weekly
}
```

### Python

```python
import os
os.environ["LOTUSMARKET_QUIET"] = "1"  # suppress startup banner

import pandas as pd
import lotusmarket as lm

# === 1. Real-time quote ===
from lotusmarket.fetchers import vps, vps_multiple, entrade_history, kbs, stock_with_fallback

stock = vps("ACB")
print(f"ACB: {stock.close:.0f} VND, vol: {stock.volume}, foreign net: {stock.foreign_net_vol}")

# Multiple stocks
stocks = vps_multiple(["ACB", "VNM", "FPT", "HPG"])

# === 2. Historical OHLCV ===
history = entrade_history("ACB", days=200)  # returns list of StockData
closes = pd.Series([h.close for h in history])

# === 3. Technical analysis — all-in-one ===
d = lm.technical.dashboard(closes)
print(f"RSI: {d.rsi:.1f}, MA20: {d.ma20:.0f}, Signal: {d.signal}, Score: {d.score:.0f}")

# Individual indicators (pandas Series in/out)
df = pd.DataFrame({"close": closes})
df["rsi"] = lm.technical.rsi(df["close"], period=14)
df["ma20"] = lm.technical.ma(df["close"], period=20)
df["ma50"] = lm.technical.ma(df["close"], period=50)

# Score and signal
score = lm.technical.score(closes)
signal = lm.technical.signal(closes)

# === 4. Weekly candles ===
daily_df = pd.DataFrame({
    "date": [h.date for h in history],
    "open": [h.open for h in history],
    "high": [h.high for h in history],
    "low": [h.low for h in history],
    "close": [h.close for h in history],
    "volume": [h.volume for h in history],
})
weekly = lm.technical.aggregate_weekly(daily_df)

# === 5. Vietnamese sentiment ===
r = lm.sentiment.analyze("Ngân hàng tăng mạnh, khối ngoại mua ròng")
print(f"Sentiment: {r.label} ({r.score:.2f})")

# === 6. Volume surge detection ===
surge = lm.signals.classify_volume_surge(
    current_vol=5000000, avg_vol_20=2000000, stddev_vol_20=500000, price_change=2.5
)
print(f"Surge: {surge.signal} / {surge.intensity} — {surge.label}")

# === 7. Domestic flow ===
stocks_df = pd.DataFrame({
    "ticker": ["ACB", "VNM", "FPT"],
    "volume": [1000000, 500000, 800000],
    "bid_vol": [600000, 300000, 500000],
    "ask_vol": [400000, 200000, 300000],
    "foreign_buy_vol": [50000, 30000, 40000],
    "foreign_sell_vol": [30000, 40000, 20000],
})
flow = lm.signals.domestic_flow(stocks_df)
print(f"Domestic: {flow.signal} (buy pressure: {flow.buy_pressure:.1f}%)")

# === 8. Sector flow ranking ===
sector_df = pd.DataFrame({
    "ticker": ["ACB", "VNM", "FPT", "HPG"],
    "close": [25000, 70000, 120000, 28000],
    "volume": [1000000, 500000, 800000, 600000],
    "change_percent": [1.5, -0.5, 2.0, -1.0],
    "foreign_net_vol": [500, -200, 300, -100],
})
sectors = lm.market.sector_flow(sector_df)
print(sectors[["sector", "rank", "signal", "net_foreign"]])

# === 9. Market pulse scoring ===
score_val, signal_str = lm.market.pulse_score(
    breadth_score=70, foreign_flow_score=60, volume_score=55, risk_score=40
)
print(f"Pulse: {score_val} ({signal_str})")

# === 10. Risk indicators ===
print(f"VIX risk: {lm.market.score_vix(25)}")
print(f"Yield curve risk: {lm.market.score_yield_curve(-0.3)}")

# === 11. Fundamentals ===
fund = kbs("ACB")
print(f"P/E: {fund.pe:.1f}, P/B: {fund.pb:.1f}, EPS: {fund.eps:.0f}")

# === 12. Portfolio TWR ===
from lotusmarket.portfolio import twr, ReturnSegment
returns = twr([ReturnSegment(100, 110), ReturnSegment(110, 132)])
print(f"TWR: {returns:.2%}")

# === 13. Vietnamese NLU ===
from lotusmarket.nlu import Parser
parser = Parser()
result = parser.parse("giá ACB hiện tại bao nhiêu")
print(f"Intent: {result.intent}, Tickers: {result.tickers}")

# === 14. Backtest ===
from lotusmarket.backtest import run, BacktestConfig
prices = pd.Series([20000 + i * 50 for i in range(200)], dtype=float)
bt = run(BacktestConfig(strategy="rsi"), prices)
print(f"Return: {bt.total_return:.1f}%, Max DD: {bt.max_drawdown:.1f}%, Trades: {bt.trade_count}")
```

## Modules

### Pure computation (no API key, no network)

| Module | Go import | Python import | Description |
|--------|-----------|---------------|-------------|
| **technical** | `technical` | `lotusmarket.technical` | RSI (Wilder's smoothing), MA20/50/200, Momentum, Signal (BUY/SELL/HOLD), Score (0-100), Dashboard, Weekly aggregation |
| **sentiment** | `sentiment` | `lotusmarket.sentiment` | Vietnamese financial keyword analysis — 70+ keywords with weights, score -1.0 to +1.0 |
| **signals** | `signals` | `lotusmarket.signals` | Volume surge detection (z-score or ratio), domestic buy/sell pressure (bid/ask ratio) |
| **market** | `market` | `lotusmarket.market` | Market pulse scoring (0-100, green/yellow/red), sector flow ranking, risk indicators (VIX, yield curve, news sentiment, Fed trend) |
| **portfolio** | `portfolio` | `lotusmarket.portfolio` | Time-weighted returns (TWR), 6-factor confidence scoring (0-100) |
| **nlu** | `nlu` | `lotusmarket.nlu` | Vietnamese intent classification (5 intents), ticker extraction, timeframe detection |
| **backtest** | `backtest` | `lotusmarket.backtest` | Strategy backtesting — RSI, MA cross, combined. Win rate, max drawdown, buy & hold comparison |
| **types** | `types` | `lotusmarket.types` | StockData, KBSQuote, VN30 list, sector mappings, trading fee constants |

### Data fetchers (network calls, no API key needed)

| Fetcher | Go | Python | Data |
|---------|-----|--------|------|
| **VPS** | `fetchers.VPS(ctx, "ACB")` | `fetchers.vps("ACB")` | Real-time: price, volume, change%, foreign buy/sell, bid/ask |
| **Entrade** | `fetchers.EntradeHistory(ctx, "ACB", 200)` | `fetchers.entrade_history("ACB", 200)` | Historical daily OHLCV (up to 1000 days) |
| **KBS** | `fetchers.KBS(ctx, "ACB")` | `fetchers.kbs("ACB")` | Fundamentals: P/E, P/B, EPS, Beta, Dividend Yield, Market Cap |
| **Fallback** | `fetchers.StockWithFallback(ctx, "ACB")` | `fetchers.stock_with_fallback("ACB")` | VPS primary, Entrade backup (3 retries + backoff) |

> Python fetchers require `httpx`: `pip install lotusmarket[fetchers]`

## Data Sources

| Source | URL | Data | Auth |
|--------|-----|------|------|
| VPS | bgapidatafeed.vps.com.vn | Real-time quotes, foreign flow, bid/ask volumes | Public |
| Entrade | services.entrade.com.vn | Historical daily OHLCV candles | Public |
| KBS | kbbuddywts.kbsec.com.vn | Fundamental metrics (P/E, P/B, EPS, Beta) | Public |
| SSI | iboard-query.ssi.com.vn | VN30/HNX30 index component lists | Public |

## Technical Analysis Details

| Indicator | Method | Parameters |
|-----------|--------|------------|
| RSI | Wilder's exponential smoothing | Default period: 14 |
| MA | Simple Moving Average | Periods: 20, 50, 200 |
| Momentum | Rate of change: (current - past) / past * 100 | Default period: 20 |
| Signal | MA crossover + RSI zones | BUY: price > MA20 > MA50, RSI < 70 |
| Score | Composite: signal + RSI zones | Range: 0-100 |
| Weekly | ISO week aggregation from daily OHLCV | O=first, H=max, L=min, C=last, V=sum |

## Vietnamese Sentiment Keywords

The sentiment analyzer includes 70+ Vietnamese financial keywords with weights:

**Positive** (sample): `tăng mạnh` (1.5), `bùng nổ` (2.0), `mua ròng` (1.0), `vượt đỉnh` (1.5), `lợi nhuận` (1.0)

**Negative** (sample): `giảm sâu` (2.0), `bán tháo` (2.0), `lao dốc` (2.0), `nợ xấu` (1.5), `rút vốn` (1.5)

Score formula: `(positive_sum - negative_sum) / total_sum`, clamped to [-1.0, 1.0].

## Domestic Flow Signal

Based on aggregate bid/ask volume ratios across stocks:

| Buy Pressure | Signal |
|-------------|--------|
| >= 60% | mua manh (strong buy) |
| >= 53% | nghieng mua (leaning buy) |
| >= 47% | can bang (balanced) |
| >= 40% | nghieng ban (leaning sell) |
| < 40% | ban manh (strong sell) |

## Volume Surge Detection

Uses z-score (preferred) or ratio fallback:

| Z-Score / Ratio | Intensity |
|-----------------|-----------|
| >= 3.0 / 3x avg | strong_surge |
| >= 2.0 / 2x avg | surge |
| < 2.0 | (no surge) |

Combined with price direction:
- Price up (>= 0.5%) + surge = **accumulation** (tien dang do vao)
- Price down (<= -0.5%) + surge = **distribution** (dang ban thao)
- Flat + surge = **high_activity** (giao dich dot bien)

## Go vs Python API Conventions

| Aspect | Go | Python |
|--------|-----|--------|
| Input/Output | `[]float64`, structs | `pd.Series`, `pd.DataFrame`, dataclasses |
| Async | `context.Context` + goroutines | `asyncio` (fetchers have `_async` variants planned) |
| Error handling | `(result, error)` return | Exceptions (`ValueError`, `ConnectionError`) |
| Config | Option functions | Keyword arguments |
| Dependencies | Zero (stdlib only for pure modules) | `pandas`, `numpy` (core), `httpx` (fetchers) |

## Support

If you find this useful, consider supporting the project:

- [PayPal](https://www.paypal.com/paypalme/ducnhd)
- Telegram Stars: [@vnlotusmarketbot](https://t.me/vnlotusmarketbot) `/donate`

## License

MIT
