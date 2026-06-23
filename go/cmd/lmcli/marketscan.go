package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/ducnhd/lotusmarket/go/ai"
	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/technical"
)

const eventCooldownDays = 5

func runMarketScan(ctx context.Context, outDir string) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}
	// Recency for event dedup (5-day window).
	recentVariantKeys = recentContentKeys(outDir, eventCooldownDays)

	snap := buildSnapshot(ctx)
	events := DetectEvents(snap)

	top, ok := TopEvent(events)
	for ok && !notRecentlyPosted(top) {
		// Drop the deduped top event and re-rank the rest.
		filtered := events[:0]
		for _, e := range events {
			if eventSlug(e) != eventSlug(top) {
				filtered = append(filtered, e)
			}
		}
		events = filtered
		top, ok = TopEvent(events)
	}

	if !ok {
		// No event — backfill an evergreen post if the gap is long enough.
		if shouldBackfill(daysSinceLastPost(outDir, time.Now())) {
			log.Printf("[marketscan] no event; backfilling evergreen")
			runAutoBlog(ctx, outDir, "auto")
			return
		}
		log.Printf("[marketscan] no event above threshold and gap < %dd — no post", backfillGapDays)
		return
	}

	log.Printf("[marketscan] top event: %s %s score=%.0f", top.Type, top.Ticker, top.Score)
	writeEventPost(ctx, outDir, top)
}

func writeEventPost(ctx context.Context, outDir string, e Event) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	client, err := ai.New(ai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("ai: %v", err)
	}
	title := eventTitle(e)
	result, err := client.AnalyzeWithContext(ctx, e.DataBlock, eventPromptBase)
	if err != nil {
		log.Fatalf("claude: %v", err)
	}
	date := time.Now().Format("2006-01-02")
	slug := eventSlug(e)
	front := fmt.Sprintf("---\ntitle: %q\ndate: %s\ntopic: %s\n---\n\n", title, date, e.Type)
	full := filepath.Join(outDir, date+"-"+slug+".md")
	if err := os.WriteFile(full, []byte(front+result.Text+"\n"), 0o644); err != nil {
		log.Fatalf("write %s: %v", full, err)
	}
	fmt.Println(full)
}

// eventTitle builds the SEO/Telegram headline (ticker + number + date).
func eventTitle(e Event) string {
	d := time.Now().Format("02/01")
	switch e.Type {
	case "mover":
		dir := "tăng"
		if e.Magnitude < 0 {
			dir = "giảm"
		}
		return fmt.Sprintf("%s %s %.1f%% hôm nay %s — cohort lịch sử nói gì?", e.Ticker, dir, math.Abs(e.Magnitude), d)
	case "volume":
		return fmt.Sprintf("%s volume gấp %.1f lần hôm nay %s — smart money hay noise?", e.Ticker, e.Magnitude, d)
	case "insider":
		return fmt.Sprintf("Nội bộ %s giao dịch lớn (%s) — đọc tín hiệu thế nào?", e.Ticker, d)
	case "flow":
		return fmt.Sprintf("Dòng tiền VN30 %s hôm nay %s — regime nào đang dẫn dắt?", flowDir(e.Magnitude), d)
	case "global":
		return fmt.Sprintf("Sốc global hôm nay %s — VN30 chịu ảnh hưởng ra sao? Cohort lookup", d)
	}
	return fmt.Sprintf("Diễn biến thị trường %s", d)
}

func flowDir(buyPressure float64) string {
	if buyPressure < 50 {
		return "nghiêng bán"
	}
	return "nghiêng mua"
}

// buildSnapshot performs all network I/O and assembles the pure Snapshot.
func buildSnapshot(ctx context.Context) Snapshot {
	universe := marketUniverse()
	quotes, _ := fetchers.VPSMultiple(ctx, universe)

	stocks := make([]StockSnap, 0, len(quotes))
	for _, q := range quotes {
		hist, err := fetchers.EntradeHistory(ctx, q.Ticker, 30)
		avg := 0.0
		rsi := 0.0
		mat := "mixed"
		if err == nil && len(hist) >= 20 {
			var vsum int64
			closes := make([]float64, len(hist))
			for i, h := range hist {
				closes[i] = h.Close
				if i >= len(hist)-20 {
					vsum += h.Volume
				}
			}
			avg = float64(vsum) / 20
			d := technical.Dashboard(closes)
			rsi = d.RSI
			if d.MA200 != nil && d.MA50 != nil {
				if q.Close >= *d.MA200 && *d.MA50 >= *d.MA200 {
					mat = "uptrend"
				} else if q.Close < *d.MA200 && *d.MA50 < *d.MA200 {
					mat = "downtrend"
				}
			}
		}
		stocks = append(stocks, StockSnap{
			Ticker: q.Ticker, Close: q.Close, ChangePct: q.ChangePercent,
			Volume: q.Volume, AvgVol20: avg, RSI: rsi, MATrend: mat,
		})
	}

	globals := []GlobalSnap{}
	for _, g := range fetchers.YahooMultiple(ctx, []string{"^VIX", "^GSPC", "^N225", "^HSI"}) {
		globals = append(globals, GlobalSnap{Symbol: g.Symbol, Price: g.Price, ChangePct: g.ChangePct})
	}

	flow := signalsFlowSafe(quotes) // existing helper in autoblog.go

	insiders := []InsiderSnap{}
	for _, st := range stocks {
		if math.Abs(st.ChangePct) < moverThresholdPct {
			continue // only check insiders for tickers already moving, to bound API calls
		}
		txs, err := fetchers.CafefInsider(ctx, st.Ticker)
		if err != nil {
			continue
		}
		for _, tx := range txs {
			if tx.BuyResult >= insiderShareThreshold {
				insiders = append(insiders, InsiderSnap{Ticker: tx.Ticker, Side: "buy", Shares: int64(tx.BuyResult), InsiderName: tx.InsiderName})
			} else if tx.SellResult >= insiderShareThreshold {
				insiders = append(insiders, InsiderSnap{Ticker: tx.Ticker, Side: "sell", Shares: int64(tx.SellResult), InsiderName: tx.InsiderName})
			}
		}
	}

	return Snapshot{
		Date: time.Now(), Stocks: stocks, Globals: globals,
		Insiders: insiders, FlowBuyPressure: flow.BuyPressure,
	}
}
