# lotusmarket

Vietnamese Stock Market Toolkit for Go & Python.

Real-time quotes, historical OHLCV, technical analysis, sentiment, money flow signals, and more for HOSE/HNX stocks.

## Install

**Go:**
```bash
go get github.com/lotusmarket/lotusmarket-go
```

**Python:** (coming soon)
```bash
pip install lotusmarket
```

## Quick Start (Go)

```go
package main

import (
    "context"
    "fmt"

    "github.com/lotusmarket/lotusmarket-go/fetchers"
    "github.com/lotusmarket/lotusmarket-go/technical"
    "github.com/lotusmarket/lotusmarket-go/sentiment"
)

func main() {
    ctx := context.Background()

    // Real-time quote (no API key needed)
    stock, _ := fetchers.VPS(ctx, "ACB")
    fmt.Printf("ACB: %v VND, volume: %d\n", stock.Close, stock.Volume)

    // Historical data + technical analysis
    history, _ := fetchers.EntradeHistory(ctx, "ACB", 200)
    closes := make([]float64, len(history))
    for i, h := range history {
        closes[i] = h.Close
    }

    dashboard := technical.Dashboard(closes)
    fmt.Printf("RSI: %.1f, Signal: %s, Score: %.0f\n", dashboard.RSI, dashboard.Signal, dashboard.Score)

    // Vietnamese sentiment
    result := sentiment.Analyze("Ngân hàng tăng mạnh, khối ngoại mua ròng")
    fmt.Printf("Sentiment: %s (%.2f)\n", result.Label, result.Score)
}
```

## Modules

| Module | Description | API Key |
|--------|-------------|---------|
| `technical` | RSI, MA20/50/200, Momentum, Signal, Score, Dashboard, Weekly aggregation | No |
| `sentiment` | Vietnamese financial keyword analysis (70+ keywords), score -1 to +1 | No |
| `signals` | Volume surge detection (z-score), domestic buy/sell pressure | No |
| `market` | Market pulse scoring (0-100), sector flow ranking, risk indicators | No |
| `portfolio` | Time-weighted returns (TWR), 6-factor confidence scoring | No |
| `nlu` | Vietnamese intent classification, ticker extraction | No |
| `fetchers` | VPS (real-time), Entrade (OHLCV), KBS (fundamentals), with fallback | No |
| `types` | Shared types: StockData, KBSQuote, VN30 list, sector mappings | No |

## Data Sources

| Source | Data | Key Required |
|--------|------|--------------|
| VPS | Real-time quotes, foreign flow, bid/ask | No |
| Entrade | Historical daily OHLCV (up to 1000 days) | No |
| KBS | Fundamentals: P/E, P/B, EPS, Beta, Dividend Yield | No |
| SSI | VN30/HNX30 index components | No |

## Support

If you find this useful, consider supporting the project:

- PayPal: (coming soon)
- Telegram Stars: `@lotusmarket_bot /donate`

## License

MIT
